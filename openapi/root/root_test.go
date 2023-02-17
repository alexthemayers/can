package root_test

import (
	"github.com/sasswart/gin-in-a-can/openapi/path"
	"github.com/sasswart/gin-in-a-can/openapi/root"
	"github.com/sasswart/gin-in-a-can/openapi/test"
	"github.com/sasswart/gin-in-a-can/tree"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestOpenAPI_LoadOpenAPI(t *testing.T) {
	apiSpec, err := root.LoadAPISpec(test.AbsOpenAPI)
	if err != nil {
		t.Errorf("could not load file %s:%s", test.OpenapiFile, err.Error())
	}
	if apiSpec == nil {
		t.Errorf("could not load file %s:%s", test.OpenapiFile, err.Error())
	}
}

func TestOpenAPI_SetRenderer(t *testing.T) {
	//apiSpec, _ := LoadAPISpec(test.AbsOpenAPI)
	//want := render.GinRenderer{}
	//_ = SetRenderer(apiSpec, want)
	//leaf := test.Dig(apiSpec, test.Endpoint, test.Method, test.ReqBody, test.MediaType, test.Schema, "name")
	//
	//got := leaf.GetRenderer()
	//if !reflect.DeepEqual(got, want) {
	//	t.Fatalf("new metadata did not populate to root of tree. wanted %v, got %v", want, got)
	//}
}

func TestOpenAPI_GetBasePath(t *testing.T) {
	apiSpec, _ := root.LoadAPISpec(test.AbsOpenAPI)
	wanted := filepath.Dir(test.AbsOpenAPI)
	if apiSpec.GetBasePath() != wanted {
		t.Errorf("could not get basePath %s, got %s", wanted, apiSpec.GetBasePath())
	}
}

func TestOpenAPI_GetParent(t *testing.T) {
	apiSpec, _ := root.LoadAPISpec(test.AbsOpenAPI)
	p := apiSpec.GetParent()
	if p != nil {
		t.Errorf("the root openapi file found a parent: %v", p)
	}
}

func TestGetOpenAPI_GetChildren(t *testing.T) {
	apiSpec, _ := root.LoadAPISpec(test.AbsOpenAPI)
	paths := apiSpec.GetChildren()
	if len(paths) == 0 {
		t.Errorf("error occured or test yaml file has no paths to get")
	}
	for k, v := range paths {
		if k == test.Endpoint {
			if _, ok := v.(*path.Item); ok {
				return // test that it's a *PathItem
			}
		}
	}
	t.Errorf("could not find a valid child in openapi file")
}

func TestOpenAPI_SetChild(t *testing.T) {
	apiSpec, _ := root.LoadAPISpec(test.AbsOpenAPI)
	pathKey := "new"
	p := path.Item{
		Description: "new path item",
	}
	apiSpec.SetChild(pathKey, &p)

	paths := apiSpec.GetChildren()
	for k, v := range paths {
		if k == pathKey {
			path, ok := v.(*path.Item) // test that it's a *PathItem
			if !ok {
				t.Errorf("Non-valid pathItem found")
			}
			if !reflect.DeepEqual(*path, p) {
				t.Errorf("path item set is not equivalent to path item retrieved")
			}
		}
	}
}

func TestOpenAPI_GetName(t *testing.T) {
	apiSpec, _ := root.LoadAPISpec(test.AbsOpenAPI)
	//_ = SetRenderer(apiSpec, render.GinRenderer{})
	name := apiSpec.GetName()
	if name != test.GinRenderedOpenAPIName {
		t.Errorf("wanted %s, got %s", test.GinRenderedOpenAPIName, name)
	}
}

func TestOpenAPI_ResolveRefs(t *testing.T) {
	apiSpec := root.Root{
		Node: tree.Node{
			//basePath: filepath.Dir(test.AbsOpenAPI),
		},
		//Components: components.Components{},
		Paths: map[string]*path.Item{},
	}
	content, _ := os.ReadFile(test.AbsOpenAPI)
	_ = yaml.Unmarshal(content, &apiSpec)

	newApi, err := tree.Traverse(&apiSpec, root.ResolveRefs)

	if err != nil {
		t.Errorf(err.Error()) // just bubbling up is enough here
	}
	if newApi == nil {
		t.Errorf("%s resulted in a nil api tree", test.OpenapiFile)
	}
}
func TestOpenAPI_readRef(t *testing.T) {
	// ReadRef takes a reference and attempts to unmarshal it's content into the struct being passed as `i`.
	// As it happens, this ref is contained within the struct that is being unmarshalled into.
	//func readRef(absFilename string, n tree.NodeTraverser) error {

	t.Errorf("TODO")
}

func TestOpenAPI_SetMetadata(t *testing.T) {
	apiSpec, _ := root.LoadAPISpec(test.AbsOpenAPI)
	data := map[string]string{"this": "is", "some": "metadata"}
	traversable := test.Dig(apiSpec, test.Endpoint, test.Method, test.ReqBody, test.MediaType, test.Schema, "name")

	traversable.SetMetadata(data)
	if !reflect.DeepEqual(apiSpec.GetMetadata(), data) {
		t.Fatalf("new metadata did not populate to root of tree. wanted %v, got %v", data, apiSpec.GetMetadata())
	}
}
