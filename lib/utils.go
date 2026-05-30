package swagger_builder

import (
	"fmt"
	"path"
	"reflect"
	"slices"
	"strings"
)

func TypeToSwagger(kind reflect.Kind) string {
	switch kind {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int64:
		return "integer"
	case reflect.Bool:
		return "boolean"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Struct, reflect.Map:
		return "object"
	case reflect.Array, reflect.Slice:
		return "array"
	default:
		return "string"
	}
}

func StructToSchema(t reflect.Type) Schema {

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
		s := StructToSchema(t.Elem())
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
			schema.Properties[fieldName] = StructToSchema(field.Type)
		}

	}

	return schema

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
