package swagger_builder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/goccy/go-yaml"
)

type SwaggerComponent struct {
}

type Path struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`

	Tags []string `json:"tags,omitempty"`

	Responses map[string]SwaggerResponse `json:"responses,omitempty"`

	Parameters []Parameter `json:"parameters,omitempty"`

	RequestBody Body `json:"requestBody,omitempty"`
}

type Items struct {
	Type       string            `json:"type,omitempty"`
	Format     string            `json:"format,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty"`
	Ref        string            `json:"$ref,omitempty"`
}

type Schema struct {
	Type       string            `json:"type,omitempty"`
	Format     string            `json:"format,omitempty"`
	Items      Items             `json:"items,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty"`
	Ref        string            `json:"$ref,omitempty"`
}

func (this Schema) ToItems() Items {
	return Items{
		Type:       this.Type,
		Properties: this.Properties,
		Ref:        this.Ref,
		Format:     this.Format,
	}
}

type MediaTypeObject struct {
	Schema Schema `json:"schema,omitempty"`
	// example, examples, encoding... podem entrar aqui depois
}

type Body struct {
	Description string                     `json:"description"`
	Required    bool                       `json:"required"`
	Content     map[string]MediaTypeObject `json:"content,omitempty"`
}

type SwaggerResponse struct {
	Description string          `json:"description"`
	Content     MediaTypeObject `json:"content"`
}

type Parameter struct {
	Name     string         `json:"name"`
	In       string         `json:"in"`
	Schema   map[string]any `json:"schema"`
	Required bool           `json:"required"`
}

type SwaggerInfo struct {
	Title        string `json:"title"`
	Descriptiotn string `json:"description"`
	Version      string `json:"version"`
}

type SwaggerDocument struct {
	Openapi string `json:"openapi"`

	Info SwaggerInfo `json:"info"`

	Components map[string]SwaggerComponent `json:"components,omitempty"`

	Paths map[string]map[string]Path `json:"paths"`
	//        ^ PATH     ^ METHOD
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
