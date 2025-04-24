package jujuclient

import (
	jclient "github.com/juju/juju/jujuclient"
)

type Client interface{}

type client struct {
	clientStore jclient.ClientStore
}

func NewClient() (Client, error) {
	clientStore := jclient.NewFileClientStore()
	return &client{clientStore: clientStore}, nil
}
