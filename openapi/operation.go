package openapi

import "fmt"

// Operation is a programmatic representation of the Operation object defined here: https://swagger.io/specification/#operation-object
type Operation struct {
	Tags         []string
	Summary      string
	Description  string
	Parameters   []Parameter
	RequestBody  RequestBody `yaml:"requestBody"`
	Responses    map[string]Response
	OperationId  string `yaml:"operationId"`
	ExternalDocs ExternalDocs
}

// communicate by sharing memory ;)
var _ traversable = &Operation{}

func (o *Operation) Render() error {
	fmt.Println("Rendering API Operation")

	err := o.RequestBody.Render()
	if err != nil {
		return err
	}

	for _, response := range o.Responses {
		err := response.Render()
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *Operation) GetSchemas(name string) (schemas map[string]Schema) {
	schemas = map[string]Schema{}
	if o == nil {
		return schemas
	}
	for _, response := range o.Responses {
		for name, schema := range response.GetSchemas(name) {
			schemas[name] = schema
		}
	}

	return schemas
}

func (o *Operation) getParent() traversable {
	return nil
}

func (o *Operation) getChildren() map[string]traversable {
	return map[string]traversable{}
}

func (o *Operation) setChild(i string, child traversable) {
	// TODO
}
