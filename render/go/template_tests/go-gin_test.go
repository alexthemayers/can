package template_tests

import "path/filepath"

const basePath = "../../../templates/go-gin/"

var (
	// openapiTemplate's file contains the template for the base api. A simple interface and a function to
	// register the api with your program
	openapiTemplate, _ = filepath.Abs(filepath.Join(basePath, "openapi.tmpl"))

	// operationTemplate's file contains a template for handling the operation's parameters and their validation
	operationTemplate, _ = filepath.Abs(filepath.Join(basePath, "operation.tmpl"))

	// pathItemTemplate's file contains a template for the generation of a particular path's http handlers and an
	// an interface for their associated methods
	pathItemTemplate, _ = filepath.Abs(filepath.Join(basePath, "path_item.tmpl"))

	// schemaTemplate's file contains a template for the generation of an openAPI-3's schema object and the associated
	// validation logic it contains.
	schemaTemplate, _ = filepath.Abs(filepath.Join(basePath, "schema.tmpl"))
)
