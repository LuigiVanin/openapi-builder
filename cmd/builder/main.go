package main

import (
	"fmt"

	oas "github.com/LuigiVanin/openapi-builder/openapi"
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

	builder := oas.NewBuilder(
		"Test API",
		"This is a test API",
		"1.0.0",
	)

	builder.Add(
		builder.Route("POST", "test/", oas.Options{Summary: "Resumo do teste", Description: "Descrição do teste"}).
			AddTag("Teste").
			AddPathParam("id", oas.Integer).
			AddQueryParam("name", oas.Integer).
			AddBody(Body{}, oas.Options{Description: "main body hahaha", Required: true}),
	)

	builder = builder.
		AddRoute(oas.Route{
			Method: "GET",
			Path:   "/test",
			Tags:   []string{"Teste"},

			Parameter: Parameter{},
			Query:     Query{},
		})

	builder.AddRoute(oas.Route{
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
			AddHeaderParam("x-string", oas.String, oas.Options{Description: "Api key"}).
			AddResponse(200, Response{}).
			AddResponse(404, Error{}),
	)

	document := builder.Build()
	err := document.Write()

	if err != nil {
		fmt.Println(err.Error())
	}

}
