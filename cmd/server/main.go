package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/flowchartsman/swaggerui"
)

func main() {
	fmt.Println("Hello World!")

	spec, err := os.ReadFile("../specs/index.json")

	if err != nil {
		fmt.Println("Erro: ", err.Error())

		spec, _ = os.ReadFile("./specs/index.json")
	}

	http.Handle("/docs/", http.StripPrefix("/docs", swaggerui.Handler(spec)))
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Erro: ", err.Error())
	}
}
