package swagger_builder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/goccy/go-yaml"
)

type Method = string
type PathName = string
type MediaType = string

type SwaggerComponent struct {
}

type Content struct {
	Schema map[string]any `json:"schema"`
}

type SwaggerResponse struct {
	Description string                `json:"description"`
	Content     map[MediaType]Content `json:"content"`
}

type SwaggerParameter struct {
	Name     string         `json:"name"`
	In       string         `json:"in"`
	Schema   map[string]any `json:"schema"`
	Required bool           `json:"required"`
}

type SwaggerPath struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`

	Tags []string `json:"tags,omitempty"`

	Responses map[string]SwaggerResponse `json:"responses,omitempty"`

	Parameters []SwaggerParameter `json:"parameters,omitempty"`
}

type SwaggerInfo struct {
	Title        string `json:"title"`
	Descriptiotn string `json:"descriptiotn"`
	Version      string `json:"version"`
}

type SwaggerDocument struct {
	Openapi string `json:"openapi"`

	Info SwaggerInfo `json:"info"`

	Components map[string]SwaggerComponent `json:"components,omitempty"`

	Paths map[PathName]map[Method]SwaggerPath `json:"paths"`
}

type WriteOptions struct {
	Formats    []string
	FolderPath string
	FileName   string
}

func (this *SwaggerDocument) Write(options ...WriteOptions) error {
	option := WriteOptions{
		Formats:    []string{"yaml", "json"},
		FolderPath: "swagger",
		FileName:   "index",
	} // default options

	if len(options) > 0 {
		overrideOpts := options[0]
		option = Merge(option, overrideOpts)
	}

	for _, format := range []string{"yaml", "json"} {

		if !slices.Contains(option.Formats, format) {
			continue
		}

		var output []byte
		var err error

		switch format {
		case "yaml":
			output, err = yaml.Marshal(this)
		case "json":
			output, err = json.MarshalIndent(this, "", "\t")
		}

		if err != nil {
			fmt.Println("Erro ao converter para Yaml")
			return err
		}

		if option.FolderPath != "" {
			err := os.MkdirAll(option.FolderPath, 0755)

			if err != nil {
				return err
			}
		}

		fileName := filepath.Join(option.FolderPath, option.FileName+"."+format)

		file, err := os.Create(fileName)

		if err != nil {
			return err
		}

		defer file.Close()

		_, err = file.Write(output)

		if err != nil {
			fmt.Println("Erro ao tentar escrever no arquivo")
			return err
		}
	}

	return nil
}
