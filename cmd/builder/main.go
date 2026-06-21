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

type Response struct {
	Success   string    `json:"success"`
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Query     []Query   `json:"query"`
	Parameter Parameter `json:"parameter"`

	Body Body `json:"body"`

	Teste map[string]string
}

type Error struct {
	Error  bool   `json:"error"`
	Reason string `json:"reason"`
}

func main() {
	fmt.Println("Hello World!")

	builder := sb.NewSwaggerBuilder(
		"Test API",
		"This is a test API",
		"1.0.0",
	)

	builder.Add(
		builder.Route("POST", "test/", sb.Options{Summary: "Resumo do teste", Description: "Descrição do teste"}).
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
		Body:   Body{},
		Responses: map[string]any{
			"200": Response{},
			"404": Error{},
		},
	})

	builder.Add(
		builder.Route("POST", "/customer").
			AddResponse(200, Response{}).
			AddResponse(404, Error{}),
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
