package main

import (
	"fmt"

	sb "github.com/LuigiVanin/swagger-builder/lib"
)

type Parameter struct {
	// Name string `json:"name"`
	Id string `json:"id"`
}

type Query struct {
	Category string `json:"category"`
	Jump     bool
}

type Body struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Query     []Query   `json:"query"`
	Parameter Parameter `json:"parameter"`
	Teste     map[string]string
}

func main() {
	fmt.Println("Hello World!")

	builder := sb.NewSwaggerBuilder(
		"Test API",
		"This is a test API",
		"1.0.0",
	)

	builder.Add(
		builder.Route("POST", "test/").
			AddTag("Teste").
			AddPathParam("id", "integer").
			AddQueryParam("name", "string").
			AddBody(Body{}, sb.Options{Description: "main body hahaha", Required: true}),
	)

	builder = builder.
		AddRoute(sb.Route{
			Method: "GET",
			Path:   "/test",
			Tags:   []string{"Teste"},

			Parameter: Parameter{},
			Query:     Query{},
		})

	builder.AddRoute(sb.Route{
		Path:   "/test",
		Method: "PUT",
	})

	builder.Add(
		builder.Route("POST", "customer"),
	)
	// 	AddRoute(sb.Route{
	// 		Path:      "/test/{id}",
	// 		Method:    "GET",
	// 		Parameter: Parameter{},
	// 	}).
	// 	AddRoute(sb.Route{
	// 		Path:      "/test/{id}",
	// 		Method:    "GET",
	// 		Parameter: Parameter{},
	// 		Body:      Body{},
	// 	})

	document := builder.Build()
	err := document.Write()

	if err != nil {
		fmt.Println(err.Error())
	}

}
