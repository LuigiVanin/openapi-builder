package lib_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	openapi "github.com/LuigiVanin/openapi-builder/openapi"
	lib "github.com/LuigiVanin/openapi-builder/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DocumentTestSuite struct {
	suite.Suite
}

func (this *DocumentTestSuite) SetupTest() {

}

func TestDocumentTestSuite(t *testing.T) {
	suite.Run(t, new(DocumentTestSuite))
}

func (this *DocumentTestSuite) TestCreatingDocument_Success() {
	document := openapi.Document{
		Openapi: "3.0.0",
		Info: openapi.Info{
			Title:       "Swagger Test",
			Description: "Generic Description!",
			Version:     "3.0.0",
		},
	}

	assert.NotNil(this.T(), document)

	assert.NotPanics(this.T(), func() {
		document.Output("json") // nolint:errcheck
	})

	assert.NotPanics(this.T(), func() {
		document.Output("yaml") // nolint:errcheck
	})

}

func (this *DocumentTestSuite) TestCreatingDocumentJson_Success() {
	title := lib.GenerateText(10)
	description := lib.GenerateText(20)

	document := openapi.Document{
		Openapi: "3.0.0",
		Info: openapi.Info{
			Title:       title,
			Description: description,
			Version:     "3.0.0",
		},
	}

	assert.NotNil(this.T(), document)

	jsonBytes, err := document.Output("json")

	assert.Nil(this.T(), err)

	json := string(jsonBytes)

	assert.Contains(this.T(), json, title)
	assert.Contains(this.T(), json, description)

}

func (this *DocumentTestSuite) TestCreatingDocumentYaml_Success() {
	title := lib.GenerateText(10)
	description := lib.GenerateText(20)

	document := openapi.Document{
		Openapi: "3.0.0",
		Info: openapi.Info{
			Title:       title,
			Description: description,
			Version:     "3.0.0",
		},
	}

	assert.NotNil(this.T(), document)

	yamlBytes, err := document.Output("yaml")

	assert.Nil(this.T(), err)

	yaml := string(yamlBytes)

	assert.Contains(this.T(), yaml, title)
	assert.Contains(this.T(), yaml, description)

}

func (this *DocumentTestSuite) TestDocumentWithPath_Success() {
	summary := lib.GenerateText(10)
	description := lib.GenerateText(20)
	tags := []string{lib.GenerateText(4)}

	document := openapi.Document{
		Openapi: "3.0.0",
		Info: openapi.Info{
			Title:       "Swagger Test",
			Description: "Generic Description!",
			Version:     "3.0.0",
		},

		Paths: map[string]map[string]openapi.Path{
			"/test": {
				"POST": {
					Summary:     summary,
					Description: description,
					Tags:        tags,
				},
			},
		},
	}

	assert.NotNil(this.T(), document)

	jsonBytes, err := document.Output("json")

	assert.Nil(this.T(), err)

	json := string(jsonBytes)

	assert.Contains(this.T(), json, summary)
	assert.Contains(this.T(), json, description)
	assert.Contains(this.T(), json, tags[0])

}

func (this *DocumentTestSuite) TestScalarSchemaDoesNotEmitEmptyItems_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf("texto"))

	raw, err := json.Marshal(schema)

	assert.NoError(this.T(), err)
	assert.JSONEq(this.T(), `{"type":"string"}`, string(raw))
}

func (this *DocumentTestSuite) TestNestedItemsSerializeInJson_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf([][]int{}))

	raw, err := json.Marshal(schema)

	assert.NoError(this.T(), err)
	assert.JSONEq(
		this.T(),
		`{"type":"array","items":{"type":"array","items":{"type":"integer"}}}`,
		string(raw),
	)
}

func (this *DocumentTestSuite) TestNestedItemsSerializeInYaml_Success() {
	builder := openapi.NewBuilder(lib.GenerateText(10), lib.GenerateText(10), "1.0.0")

	builder.AddRoute(openapi.Route{
		Method: "POST",
		Path:   "/matrix",
		Body:   [][]string{},
	})

	output, err := builder.Build().Output("yaml")

	assert.NoError(this.T(), err)
	assert.Equal(this.T(), 2, strings.Count(string(output), "items:"))
}
