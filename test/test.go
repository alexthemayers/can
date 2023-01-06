package test

import (
	"path"
	"path/filepath"
)

// constants and other information used for unit testing.
// This file serves as a single source of truth for data drawn from in multiple places during testing
const (
	Endpoint                 = "/endpoint"
	Method                   = "post"
	ReqBody                  = "RequestBody"
	EmptyParamName           = "Param"
	MediaType                = "application/json"
	OpenAPIFile              = "../openapi/fixtures/validation.yaml"
	Schema                   = "Model" // the Dig() key used to access any schema held within a MediaType
	Pattern                  = "^([a-zA-Z0-9])+([-_ @\\.]([a-zA-Z0-9])+)*$"
	GinRenderedPathItemName  = "EndpointValidationFixture"
	GinRenderedResponseName  = "PostEndpointValidationFixture201Response"
	GinRenderedMediaItemName = "PostEndpointValidationFixtureRequestbody"
	GinRenderedOpenAPIName   = "ValidationFixture"
)

var AbsOpenAPI, _ = filepath.Abs(OpenAPIFile)

var BasePath = path.Dir(AbsOpenAPI)

// These functions are used purely for testing purposes.
// If a function finds use outside of testing it should be moved out of this file.

/*
	Test Config

	return
*/
