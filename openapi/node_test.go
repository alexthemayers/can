package openapi

import (
	"github.com/sasswart/gin-in-a-can/test"
	"testing"
)

func TestOpenAPI_Dig(t *testing.T) {
	openapi, _ := LoadOpenAPI(test.AbsOpenAPI)
	endpoint := Dig(openapi, test.Endpoint)
	// TODO check for identity, not just type
	if _, ok := endpoint.(*PathItem); !ok {
		t.Errorf("%#v should have been a %T", endpoint, &PathItem{})
	}

	method := Dig(endpoint, test.Method)
	if _, ok := method.(*Operation); !ok {
		t.Errorf("%#v should have been a %T", method, &Operation{})
	}

	reqBody := Dig(method, test.ReqBody)
	if _, ok := reqBody.(*RequestBody); !ok {
		t.Errorf("%#v should have been a %T", reqBody, &RequestBody{})
	}

	mediaType := Dig(reqBody, test.MediaType)
	if _, ok := mediaType.(*MediaType); !ok {
		t.Errorf("%#v should have been a %T", mediaType, &MediaType{})
	}

	schema := Dig(mediaType, test.Schema)
	if _, ok := schema.(*Schema); !ok {
		t.Errorf("%#v should have been a %T", schema, &Schema{})
	}
}
