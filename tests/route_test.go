package lib_test

import (
	"testing"

	openapi "github.com/LuigiVanin/openapi-builder/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RouteTestSuite struct {
	suite.Suite
}

func (this *RouteTestSuite) SetupTest() {

}

func TestRouteTestSuite(t *testing.T) {
	suite.Run(t, new(RouteTestSuite))
}

func (this *RouteTestSuite) TestRouteBuilderBasic_Success() {
	route := openapi.NewRouteBuilder("/test", "GET")

	assert.NotNil(this.T(), route)

	result := route.Build()

	assert.Contains(this.T(), result, "/test")
	assert.Contains(this.T(), result["/test"], "get")
}

func (this *RouteTestSuite) TestRouteBuilderMetadata_Success() {
	summary := "Route summary"
	description := "Route description"

	route := openapi.NewRouteBuilder("/test", "POST", openapi.Options{
		Summary:     summary,
		Description: description,
	}).AddTag("Teste", "Extra")

	path := route.Build()["/test"]["post"]

	assert.Equal(this.T(), summary, path.Summary)
	assert.Equal(this.T(), description, path.Description)
	assert.Contains(this.T(), path.Tags, "Teste")
	assert.Contains(this.T(), path.Tags, "Extra")
}

func (this *RouteTestSuite) TestRouteBuilderPathParam_Success() {
	route := openapi.NewRouteBuilder("/test/{id}", "GET").
		AddPathParam("id", "string", openapi.Options{Required: true})

	params := route.Build()["/test/{id}"]["get"].Parameters

	assert.NotEmpty(this.T(), params)

	found := false
	for _, param := range params {
		if param.Name == "id" {
			found = true
			assert.Equal(this.T(), "path", param.In)
			assert.Equal(this.T(), "string", param.Schema.Type)
			assert.True(this.T(), param.Required)
		}
	}

	assert.True(this.T(), found)
}

func (this *RouteTestSuite) TestRouteBuilderQueryParam_Success() {
	route := openapi.NewRouteBuilder("/test", "GET").
		AddQueryParam("page", "integer")

	params := route.Build()["/test"]["get"].Parameters

	assert.NotEmpty(this.T(), params)

	found := false
	for _, param := range params {
		if param.Name == "page" {
			found = true
			assert.Equal(this.T(), "query", param.In)
			assert.Equal(this.T(), "integer", param.Schema.Type)
		}
	}

	assert.True(this.T(), found)
}

func (this *RouteTestSuite) TestRouteBuilderHeaderParam_Success() {
	route := openapi.NewRouteBuilder("/test", "GET").
		AddHeaderParam("X-Client-Id", "string")

	params := route.Build()["/test"]["get"].Parameters

	assert.NotEmpty(this.T(), params)

	found := false
	for _, param := range params {
		if param.Name == "X-Client-Id" {
			found = true
			assert.Equal(this.T(), "header", param.In)
		}
	}

	assert.True(this.T(), found)
}

func (this *RouteTestSuite) TestRouteBuilderBody_Success() {
	route := openapi.NewRouteBuilder("/users", "POST").
		AddBody(UserPayload{}, openapi.Options{Required: true, Description: "user body"})

	body := route.Build()["/users"]["post"].RequestBody

	assert.Contains(this.T(), body.Content, "application/json")
	assert.Equal(this.T(), "user body", body.Description)
	assert.True(this.T(), body.Required)

	schema := body.Content["application/json"].Schema

	assert.Equal(this.T(), "object", schema.Type)
	assert.Contains(this.T(), schema.Properties, "id")
	assert.Contains(this.T(), schema.Properties, "name")
}

func (this *RouteTestSuite) TestRouteBuilderResponse_Success() {
	route := openapi.NewRouteBuilder("/users", "POST").
		AddResponse(200, UserPayload{}).
		AddResponse(404, ErrorPayload{}, openapi.Options{Description: "not found"})

	responses := route.Build()["/users"]["post"].Responses

	assert.Contains(this.T(), responses, "200")
	assert.Contains(this.T(), responses, "404")

	okSchema := responses["200"].Content["application/json"].Schema

	assert.Equal(this.T(), "object", okSchema.Type)
	assert.Contains(this.T(), okSchema.Properties, "id")
	assert.Contains(this.T(), okSchema.Properties, "name")

	assert.Equal(this.T(), "not found", responses["404"].Description)
}
