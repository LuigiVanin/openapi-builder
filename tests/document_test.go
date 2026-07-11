package lib_test

import (
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
