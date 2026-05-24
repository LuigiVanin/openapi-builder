package main

import (
	"fmt"

	swg "github.com/LuigiVanin/swagger-builder/lib"
)

type Parameter struct {
	// Name string `json:"name"`
	Id string `json:"id"`
}

type Query struct {
	Category string `json:"category"`
	Jump     bool
}

func main() {
	fmt.Println("Hello World!")

	builder := swg.NewSwaggerBuilder(
		"Test API",
		"This is a test API",
		"1.0.0",
	)

	builder = builder.
		AddRoute(swg.Route{
			Path:      "/test",
			Method:    "POST",
			Parameter: Parameter{},
			Query:     Query{},
		}).
		AddRoute(swg.Route{
			Path:   "/test",
			Method: "PUT",
		}).
		AddRoute(swg.Route{
			Path:      "/test/{id}",
			Method:    "GET",
			Parameter: Parameter{},
		})

	document := builder.Build()
	err := document.Write()

	if err != nil {
		fmt.Println(err.Error())
	}

}
