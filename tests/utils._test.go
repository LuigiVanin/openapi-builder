package lib_test

import (
	"reflect"
	"testing"

	openapi "github.com/LuigiVanin/openapi-builder/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
}

func (this *UtilsTestSuite) SetupTest() {

}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (this *UtilsTestSuite) TestTypeToSwagger_Success() {
	assert.Equal(this.T(), "string", openapi.TypeToSwagger(reflect.String))
	assert.Equal(this.T(), "integer", openapi.TypeToSwagger(reflect.Int))
	assert.Equal(this.T(), "integer", openapi.TypeToSwagger(reflect.Int64))
	assert.Equal(this.T(), "boolean", openapi.TypeToSwagger(reflect.Bool))
	assert.Equal(this.T(), "number", openapi.TypeToSwagger(reflect.Float32))
	assert.Equal(this.T(), "number", openapi.TypeToSwagger(reflect.Float64))
	assert.Equal(this.T(), "object", openapi.TypeToSwagger(reflect.Struct))
	assert.Equal(this.T(), "object", openapi.TypeToSwagger(reflect.Map))
	assert.Equal(this.T(), "array", openapi.TypeToSwagger(reflect.Slice))
	assert.Equal(this.T(), "array", openapi.TypeToSwagger(reflect.Array))
}

func (this *UtilsTestSuite) TestTypeToSwaggerDefault_Success() {
	assert.Equal(this.T(), "string", openapi.TypeToSwagger(reflect.Chan))
}

func (this *UtilsTestSuite) TestTypeToSchemaPrimitive_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf("some string"))

	assert.Equal(this.T(), "string", schema.Type)
	assert.Empty(this.T(), schema.Properties)
}

func (this *UtilsTestSuite) TestTypeToSchemaStruct_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(UserPayload{}))

	assert.Equal(this.T(), "object", schema.Type)
	assert.Contains(this.T(), schema.Properties, "id")
	assert.Contains(this.T(), schema.Properties, "name")
	assert.Equal(this.T(), "string", schema.Properties["id"].Type)
}

func (this *UtilsTestSuite) TestTypeToSchemaSlice_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf([]UserPayload{}))

	assert.Equal(this.T(), "array", schema.Type)
	assert.Equal(this.T(), "object", schema.Items.Type)
	assert.Contains(this.T(), schema.Items.Properties, "id")
	assert.Contains(this.T(), schema.Items.Properties, "name")
}

func (this *UtilsTestSuite) TestTypeToSchemaPointer_Success() {
	schema := openapi.TypeToSchema(reflect.TypeOf(&UserPayload{}))

	assert.Equal(this.T(), "object", schema.Type)
	assert.Contains(this.T(), schema.Properties, "id")
}

func (this *UtilsTestSuite) TestTypeToParam_Success() {
	params := openapi.TypeToParam(reflect.TypeOf(Query{}))

	assert.NotEmpty(this.T(), params)

	names := map[string]bool{}
	for _, param := range params {
		names[param.Name] = true
		assert.True(this.T(), param.Required)
	}

	assert.True(this.T(), names["category"])
	assert.True(this.T(), names["Jump"])
}

func (this *UtilsTestSuite) TestTypeToParamNonStruct_Success() {
	params := openapi.TypeToParam(reflect.TypeOf(""))

	assert.Empty(this.T(), params)
}

func (this *UtilsTestSuite) TestFormatRoutePath_Success() {
	assert.Equal(this.T(), "/test", openapi.FormatRoutePath("test"))
	assert.Equal(this.T(), "/test", openapi.FormatRoutePath("/test"))
	assert.Equal(this.T(), "/test/id", openapi.FormatRoutePath("test/id"))
	assert.Equal(this.T(), "/test", openapi.FormatRoutePath("/test/"))
}

func (this *UtilsTestSuite) TestMerge_Success() {
	base := openapi.WriteOptions{
		Formats:    []string{"yaml", "json"},
		FolderPath: "specs",
		FileName:   "index",
	}
	override := openapi.WriteOptions{
		FileName: "custom",
	}

	result := openapi.Merge(base, override)

	assert.Equal(this.T(), "custom", result.FileName)
	assert.Equal(this.T(), "specs", result.FolderPath)
	assert.Equal(this.T(), []string{"yaml", "json"}, result.Formats)
}
