package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/flowchartsman/swaggerui"
)

func main() {
	fmt.Println("Hello World!")

	spec, err := os.ReadFile("swagger/index.yaml")

	if err != nil {
		fmt.Println("Erro: ", err.Error())
	}

	http.Handle("/docs/", http.StripPrefix("/docs", swaggerui.Handler(spec)))
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Erro: ", err.Error())
	}
}
