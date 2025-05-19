package productstc

import tclient "go.temporal.io/sdk/client"

type ClientImpl struct {
	client tclient.Client
}

func NewClientImpl(client tclient.Client) *ClientImpl {
	return &ClientImpl{
		client: client,
	}
}
