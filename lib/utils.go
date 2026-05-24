package swagger_builder

import (
	"fmt"
	"path"
	"reflect"
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
	default:
		return "string"
	}
}

func FormatRoutePath(endpoint string) string {
	endpoint = path.Clean(endpoint)

	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	return endpoint
}

type Struct interface{}

func Merge[T Struct](base T, override T) T {
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
