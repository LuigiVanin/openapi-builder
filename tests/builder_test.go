package lib_test

import (
	"reflect"
	"testing"

	openapi "github.com/LuigiVanin/openapi-builder/openapi"
	lib "github.com/LuigiVanin/openapi-builder/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BuilderTestSuite struct {
	suite.Suite
}

type Parameter struct {
	// Name string `json:"name"`
	Id string `json:"id"`
}

type Query struct {
	Category string `json:"category"`
	Jump     bool
}

type Header struct {
	ClientId string `json:"client-id"`
}

type UserPayload struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ErrorPayload struct {
	Error  bool   `json:"error"`
	Reason string `json:"reason"`
}

func (this *BuilderTestSuite) SetupTest() {

}

func TestBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}

func (this *BuilderTestSuite) TestCreateBuilder_Success() {
	title := lib.GenerateText(10)
	description := lib.GenerateText(10)

	builder := openapi.NewBuilder(title, description, "1.0.0")

	assert.NotNil(this.T(), builder)

	var document *openapi.Document

	assert.NotPanics(this.T(), func() {
		builder.Build()
	})

	document = builder.Build()

	assert.NotNil(this.T(), document)

	assert.NotPanics(this.T(), func() {
		document.Output("json") // nolint:errcheck
	})

	assert.NotPanics(this.T(), func() {
		document.Output("yaml") // nolint:errcheck
	})

	str, err := document.Output("json")

	assert.Nil(this.T(), err)

	json := string(str)

	assert.Contains(this.T(), json, title)
	assert.Contains(this.T(), json, description)
}

func (this *BuilderTestSuite) TestBuilderRoute_Success() {
	title := lib.GenerateText(10)
	description := lib.GenerateText(10)

	builder := openapi.NewBuilder(title, description, "1.0.0")

	builder.AddRoute(openapi.Route{
		Method: "GET",
		Path:   "/test",
		Tags:   []string{"Teste"},

		Parameter: Parameter{},
		Query:     Query{},
		Header:    Header{},
	})

	document := builder.Build()

	assert.NotNil(this.T(), document)

	assert.Equal(this.T(), document.Info.Title, title)
	assert.Equal(this.T(), document.Info.Description, description)

	assert.NotNil(this.T(), document.Paths)
	assert.NotEmpty(this.T(), document.Paths)

	params := document.Paths["/test"]["get"].Parameters

	assert.NotNil(this.T(), params)
	assert.NotEmpty(this.T(), params)
}

func (this *BuilderTestSuite) TestBuilderRouteParams_Success() {
	title := lib.GenerateText(10)
	description := lib.GenerateText(10)

	builder := openapi.NewBuilder(title, description, "1.0.0")

	builder.AddRoute(openapi.Route{
		Method: "GET",
		Path:   "/test",
		Tags:   []string{"Teste"},

		Parameter: Parameter{},
		Query:     Query{},
		Header:    Header{},
	})

	document := builder.Build()

	assert.NotNil(this.T(), document)

	assert.Equal(this.T(), document.Info.Title, title)
	assert.Equal(this.T(), document.Info.Description, description)

	assert.NotNil(this.T(), document.Paths)
	assert.NotEmpty(this.T(), document.Paths)

	params := document.Paths["/test"]["get"].Parameters

	assert.NotNil(this.T(), params)
	assert.NotEmpty(this.T(), params)

	for _, expectedParam := range openapi.TypeToParam(reflect.TypeOf(Query{})) {
		for _, param := range params {
			if expectedParam.Name == param.Name {
				assert.Equal(this.T(), param.In, "query")
				assert.Equal(this.T(), param.Name, expectedParam.Name)
			}
		}
	}

	for _, expectedParam := range openapi.TypeToParam(reflect.TypeOf(Parameter{})) {
		for _, param := range params {
			if expectedParam.Name == param.Name {
				assert.Equal(this.T(), param.In, "path")
				assert.Equal(this.T(), param.Name, expectedParam.Name)
			}
		}
	}

	for _, expectedParam := range openapi.TypeToParam(reflect.TypeOf(Header{})) {
		for _, param := range params {
			if expectedParam.Name == param.Name {
				assert.Equal(this.T(), param.In, "header")
				assert.Equal(this.T(), param.Name, expectedParam.Name)
			}
		}
	}
}

func (this *BuilderTestSuite) TestBuilderRouteBody_Success() {
	title := lib.GenerateText(10)
	description := lib.GenerateText(10)

	builder := openapi.NewBuilder(title, description, "1.0.0")

	builder.AddRoute(openapi.Route{
		Method: "POST",
		Path:   "/test",
		Tags:   []string{"Teste"},

		Body: UserPayload{},
	})

	document := builder.Build()

	assert.NotNil(this.T(), document)

	assert.NotNil(this.T(), document.Paths)
	assert.NotEmpty(this.T(), document.Paths)

	body := document.Paths["/test"]["post"].RequestBody

	assert.NotNil(this.T(), body.Content)
	assert.Contains(this.T(), body.Content, "application/json")
	assert.True(this.T(), body.Required)

	schema := body.Content["application/json"].Schema

	assert.Equal(this.T(), "object", schema.Type)
	assert.Contains(this.T(), schema.Properties, "id")
	assert.Contains(this.T(), schema.Properties, "name")
}

func (this *BuilderTestSuite) TestBuilderRouteResponse_Success() {
	title := lib.GenerateText(10)
	description := lib.GenerateText(10)

	builder := openapi.NewBuilder(title, description, "1.0.0")

	builder.AddRoute(openapi.Route{
		Method: "POST",
		Path:   "/test",
		Tags:   []string{"Teste"},

		Responses: map[string]any{
			"200": UserPayload{},
			"400": ErrorPayload{},
		},
	})

	document := builder.Build()

	assert.NotNil(this.T(), document)

	responses := document.Paths["/test"]["post"].Responses

	assert.NotEmpty(this.T(), responses)
	assert.Contains(this.T(), responses, "200")
	assert.Contains(this.T(), responses, "400")

	okResponse := responses["200"]
	assert.Contains(this.T(), okResponse.Content, "application/json")

	okSchema := okResponse.Content["application/json"].Schema
	assert.Equal(this.T(), "object", okSchema.Type)
	assert.Contains(this.T(), okSchema.Properties, "id")
	assert.Contains(this.T(), okSchema.Properties, "name")

	errSchema := responses["400"].Content["application/json"].Schema
	assert.Contains(this.T(), errSchema.Properties, "error")
	assert.Contains(this.T(), errSchema.Properties, "reason")
}
