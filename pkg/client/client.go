package client

import (
	"github.com/weaveworks/flintlock/api/services/microvm/v1alpha1"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

// You will add lots of directives like these in the same package...
//counterfeiter:generate -o fakeclient/ github.com/weaveworks/flintlock/api/services/microvm/v1alpha1.MicroVMClient

type client struct {
	flClient v1alpha1.MicroVMClient
}

func New(flClient v1alpha1.MicroVMClient) *client {
	return &client{
		flClient: flClient,
	}
}

func (c *client) Create() (*v1alpha1.CreateMicroVMResponse, error) {
	return nil, nil
}
