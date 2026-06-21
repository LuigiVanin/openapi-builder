package lib_test

import (
	sb "github.com/LuigiVanin/swagger-builder/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DocumentTestSuite struct {
	suite.Suite
}

func (this *DocumentTestSuite) SetupTest() {

}

func (this *DocumentTestSuite) CreatingDocumentPaht_Success() {
	document := sb.SwaggerDocument{
		Openapi: "3.0.0",
		Info: sb.SwaggerInfo{
			Title:       "Swagger Test",
			Description: "Description Ha Ha!",
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
