package lib_test

import (
	openapi "github.com/LuigiVanin/openapi-builder/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DocumentTestSuite struct {
	suite.Suite
}

func (this *DocumentTestSuite) SetupTest() {

}

func (this *DocumentTestSuite) CreatingDocumentPath_Success() {
	document := openapi.Document{
		Openapi: "3.0.0",
		Info: openapi.Info{
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
