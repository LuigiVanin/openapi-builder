package openapi

import (
	"fmt"
	"path"
	"reflect"
	"slices"
	"strings"
)

var (
	Integer string = "integer"
	String  string = "string"
	Boolean string = "boolean"
	Float   string = "number"
	Array   string = "array"
	Object  string = "object"
)

func TypeToSwagger(kind reflect.Kind) string {
	switch kind {
	case reflect.String:
		return String
	case reflect.Int, reflect.Int64:
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

func TypeToSchema(t reflect.Type) Schema {

	schema := Schema{
		Properties: map[string]Schema{},
	}

	// Tipos que serão recompostos recursivamente
	compositeTypes := []reflect.Kind{reflect.Struct, reflect.Array, reflect.Slice}

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	schema.Type = TypeToSwagger(t.Kind())

	if !slices.Contains(compositeTypes, t.Kind()) {
		return schema
	}

	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		s := TypeToSchema(t.Elem())
		schema.Items = s.ToItems()

		return schema
	}

	for index := range t.NumField() {
		field := t.Field(index)

		fieldKind := field.Type.Kind()
		fieldName := field.Tag.Get("json")

		if fieldName == "" {
			fieldName = field.Name
		}

		schema.Properties[fieldName] = Schema{
			Type: TypeToSwagger(fieldKind),
		}

		if slices.Contains(compositeTypes, fieldKind) {
			schema.Properties[fieldName] = TypeToSchema(field.Type)
		}

	}

	return schema

}

func TypeToParam(t reflect.Type) []Parameter {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return []Parameter{}
	}

	parameters := []Parameter{}

	for index := range t.NumField() {
		field := t.Field(index)

		name := field.Tag.Get("json")

		if name == "" {
			name = field.Name
		}

		parameter := Parameter{
			// In:       in,
			Name:     name,
			Required: true,
			Schema: Schema{
				Type: TypeToSwagger(field.Type.Kind()),
			},
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
