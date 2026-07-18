// Package lib is a helper package for the opeapi builder package
package lib

import (
	"math/rand"
	"strings"
)

func GenerateText(length uint) string {
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var result strings.Builder

	for range length {
		index := rand.Intn(len(characters))
		char := characters[index]
		result.WriteString(string(char))
	}

	return result.String()
}
