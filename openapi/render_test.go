package openapi

import (
	"github.com/sasswart/gin-in-a-can/test"
	"testing"
)

func TestRender_Render(t *testing.T) {
	var (
		got      []byte
		expected = "package \n\n// GENERATED MODEL. DO NOT EDIT\n\ntype RequestBodyModel struct {\n\tDescription string\n\tEnabled bool\n\tId string\n\tName string\n}\n"
		err      error

		schema           Traversable
		configData       = test.NewTestConfig()
		templateFileName = "schema.tmpl"
	)

	openAPI, _ := LoadOpenAPI(test.AbsOpenAPI)
	SetRenderer(openAPI, GinRenderer{})
	schema = Dig(openAPI, test.Endpoint, test.Method, test.ReqBody, test.MediaType, test.Schema)

	got, err = Render(configData, schema, templateFileName)
	if err != nil {
		t.Errorf("Test error: %s\n", err.Error())
	}
	if string(got) != expected {
		t.Errorf("got %s, expected %s\n", got, expected)
	}
}
