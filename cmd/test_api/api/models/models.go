// GENERATED CODE. DO NOT EDIT

package models

type HeartbeatGetRequest struct {
}

func (model HeartbeatGetRequest) IsValid() error {

	return nil
}

type HeartbeatGetResponse interface {
	ImplementsHeartbeatGetResponse()
	GetStatus() int
}

type HeartbeatGet200Response struct {
}

func (HeartbeatGet200Response) ImplementsHeartbeatGetResponse() {}

func (HeartbeatGet200Response) GetStatus() int {
	return 200
}
