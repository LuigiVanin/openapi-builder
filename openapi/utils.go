package openapi

import (
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"slices"
	"strings"
	"time"
)

var (
	Integer string = "integer"
	String  string = "string"
	Boolean string = "boolean"
	Float   string = "number"
	Array   string = "array"
	Object  string = "object"
)

var (
	timeType       = reflect.TypeOf(time.Time{})
	rawMessageType = reflect.TypeOf(json.RawMessage{})
)

// jsonField é um campo de struct já resolvido do jeito que o encoding/json o vê:
// nome final, tipo, se é opcional e de onde ele veio dentro do struct.
type jsonField struct {
	name      string
	typ       reflect.Type
	index     []int
	omitempty bool
	tagged    bool
}

// embeddedStruct é um struct embedado ainda por percorrer, junto do caminho de
// índices que leva até ele.
type embeddedStruct struct {
	typ   reflect.Type
	index []int
}

// jsonFields devolve os campos que o encoding/json realmente serializa para o
// struct t, com os campos de structs embedados promovidos para o nível de cima.
// A busca é em largura para que a profundidade de cada campo seja conhecida na
// hora de resolver conflito de nome.
func jsonFields(t reflect.Type) []jsonField {
	fields := []jsonField{}

	visited := map[reflect.Type]bool{}
	current := []embeddedStruct{}
	next := []embeddedStruct{{typ: t}}

	for len(next) > 0 {
		current, next = next, current[:0]

		for _, entry := range current {
			if visited[entry.typ] {
				continue
			}

			visited[entry.typ] = true

			for index := range entry.typ.NumField() {
				field := entry.typ.Field(index)

				fieldType := field.Type
				if fieldType.Kind() == reflect.Pointer {
					fieldType = fieldType.Elem()
				}

				if field.Anonymous {
					// Um embedado não exportado ainda entra quando é struct: os
					// campos exportados de dentro dele continuam visíveis.
					if !field.IsExported() && fieldType.Kind() != reflect.Struct {
						continue
					}
				} else if !field.IsExported() {
					continue
				}

				tag := field.Tag.Get("json")

				// `json:"-"` omite o campo; `json:"-,"` nomeia o campo de "-".
				if tag == "-" {
					continue
				}

				name, options, _ := strings.Cut(tag, ",")
				indexPath := append(slices.Clone(entry.index), index)

				// Struct embedado sem nome na tag tem os campos promovidos para o
				// nível de cima em vez de virar uma propriedade própria.
				if name == "" && field.Anonymous && fieldType.Kind() == reflect.Struct {
					next = append(next, embeddedStruct{typ: fieldType, index: indexPath})

					continue
				}

				tagged := name != ""

				if !tagged {
					name = field.Name
				}

				fields = append(fields, jsonField{
					name:      name,
					typ:       field.Type,
					index:     indexPath,
					omitempty: hasTagOption(options, "omitempty"),
					tagged:    tagged,
				})
			}
		}
	}

	return resolveFieldConflicts(fields)
}

// hasTagOption procura uma opção na parte da tag json que vem depois do nome.
func hasTagOption(options string, target string) bool {
	for options != "" {
		var option string
		option, options, _ = strings.Cut(options, ",")

		if option == target {
			return true
		}
	}

	return false
}

// resolveFieldConflicts aplica a mesma regra do encoding/json para nomes
// repetidos: o campo mais raso vence; empatando na profundidade, vence o único
// que tem nome na tag; se nenhum ou mais de um tiver, o nome é descartado.
// No fim a ordem de declaração é restaurada.
func resolveFieldConflicts(fields []jsonField) []jsonField {
	slices.SortFunc(fields, func(a jsonField, b jsonField) int {
		if compare := strings.Compare(a.name, b.name); compare != 0 {
			return compare
		}

		if compare := len(a.index) - len(b.index); compare != 0 {
			return compare
		}

		if a.tagged != b.tagged {
			if a.tagged {
				return -1
			}

			return 1
		}

		return slices.Compare(a.index, b.index)
	})

	resolved := []jsonField{}

	for start := 0; start < len(fields); {
		end := start + 1

		for end < len(fields) && fields[end].name == fields[start].name {
			end++
		}

		if field, ok := dominantField(fields[start:end]); ok {
			resolved = append(resolved, field)
		}

		start = end
	}

	slices.SortFunc(resolved, func(a jsonField, b jsonField) int {
		return slices.Compare(a.index, b.index)
	})

	return resolved
}

// dominantField escolhe o campo que fica com o nome. O grupo já chega ordenado
// por profundidade e, dentro dela, com os que têm tag na frente.
func dominantField(group []jsonField) (jsonField, bool) {
	if len(group) == 1 {
		return group[0], true
	}

	// O primeiro está sozinho no nível mais raso.
	if len(group[0].index) != len(group[1].index) {
		return group[0], true
	}

	// Empate de profundidade: só vale se exatamente um tiver nome na tag.
	if group[0].tagged && !group[1].tagged {
		return group[0], true
	}

	return jsonField{}, false
}

func TypeToSwagger(kind reflect.Kind) string {
	switch kind {
	case reflect.String:
		return String
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Integer
	case reflect.Bool:
		return Boolean
	case reflect.Float32, reflect.Float64:
		return Float
	case reflect.Struct, reflect.Map:
		return Object
	case reflect.Array, reflect.Slice:
		return Array
	default:
		return String
	}
}

// TypeToSchema converte um tipo Go no Schema OpenAPI equivalente, recompondo
// structs, slices, arrays e ponteiros recursivamente.
func TypeToSchema(t reflect.Type) Schema {
	return typeToSchema(t, map[reflect.Type]bool{})
}

// seen guarda os structs abertos no caminho atual da recursão. Sem ele um tipo
// que se referencia — um campo `Parent *Node` dentro de `Node` — recursaria para
// sempre, já que ponteiro passou a ser recomposto. O `defer delete` deixa a
// proteção por caminho, e não global: o mesmo tipo pode aparecer em dois ramos
// irmãos, o que é legítimo.
func typeToSchema(t reflect.Type, seen map[reflect.Type]bool) Schema {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	switch t {
	case timeType:
		return Schema{Type: String, Format: "date-time"}
	case rawMessageType:
		return Schema{Type: Object}
	}

	// []byte sai em base64 no encoding/json, não como lista de números.
	if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
		return Schema{Type: String, Format: "byte"}
	}

	schema := Schema{Type: TypeToSwagger(t.Kind())}

	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		schema.Items = typeToSchema(t.Elem(), seen).ToItems()

		return schema

	case reflect.Struct:
		if seen[t] {
			// O tipo já está aberto acima no caminho: corta o ciclo e documenta
			// como um objeto sem propriedades.
			return schema
		}

		seen[t] = true
		defer delete(seen, t)

		schema.Properties = map[string]Schema{}

		for _, field := range jsonFields(t) {
			schema.Properties[field.name] = typeToSchema(field.typ, seen)
		}

		return schema

	default:
		return schema
	}
}

// TypeToParam converte cada campo exportado e serializável de um struct em um
// Parameter. O `In` fica vazio aqui — quem chama é que decide se é path, query
// ou header.
func TypeToParam(t reflect.Type) []Parameter {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return []Parameter{}
	}

	parameters := []Parameter{}

	for _, field := range jsonFields(t) {
		parameter := Parameter{
			// In:       in,
			Name: field.name,
			// Um campo com `,omitempty` é, por definição, opcional.
			Required: !field.omitempty,
			Schema:   TypeToSchema(field.typ),
		}
		parameters = append(parameters, parameter)
	}

	return parameters
}

func FormatRoutePath(endpoint string) string {
	endpoint = path.Clean(endpoint)

	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	return endpoint
}

func Merge[T any](base T, override T) T {
	baseVal := reflect.ValueOf(base)
	overrideVal := reflect.ValueOf(override)

	if baseVal.Kind() != reflect.Struct && baseVal.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("Merge exige um struct, recebeu %T", base))
	}

	// Se for ponteiro, derreferencia
	if baseVal.Kind() == reflect.Ptr {
		baseVal = baseVal.Elem()
		overrideVal = overrideVal.Elem()
	}

	// Cria cópia do base
	result := reflect.New(baseVal.Type()).Elem()
	result.Set(baseVal)

	// Aplica overrides
	for i := range baseVal.NumField() {
		field := overrideVal.Field(i)
		if !field.IsZero() {
			result.Field(i).Set(field)
		}
	}

	return result.Interface().(T)
}
