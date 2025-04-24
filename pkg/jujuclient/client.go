package jujuclient

import (
	"github.com/juju/errors"
	jclient "github.com/juju/juju/jujuclient"
	"github.com/rs/zerolog/log"
)

type Client interface {
	GetControllers() (Controllers, error)
}

type client struct {
	clientStore jclient.ClientStore
}

func NewClient() (Client, error) {
	clientStore := jclient.NewFileClientStore()
	return &client{clientStore: clientStore}, nil
}

func (c *client) GetControllers() (Controllers, error) {
	allControllers, err := c.clientStore.AllControllers()
	if err != nil {
		return Controllers{}, err
	}
	currentController, err := c.clientStore.CurrentController()
	if err != nil {
		if errors.Is(err, errors.NotFound) {
			log.Debug().Msg("CurrentController not found")
		} else {
			return Controllers{}, err
		}
	}
	return Controllers{ControllerDetails: allControllers, Current: currentController}, nil
}

type Controllers struct {
	ControllerDetails map[string]jclient.ControllerDetails `json:"controllerDetails,omitempty"`
	Current           string                               `json:"current,omitempty"`
}
