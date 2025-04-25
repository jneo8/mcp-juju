package jujuclient

import (
	"github.com/juju/errors"
	jclient "github.com/juju/juju/jujuclient"
	"github.com/rs/zerolog/log"
)

type Client interface {
	GetControllers() (Controllers, error)
	GetModels() (Models, error)
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

func (c *client) GetModels() (Models, error) {
	controllers, err := c.GetControllers()
	if err != nil {
		return Models{}, err
	}
	allModels, err := c.clientStore.AllModels(controllers.Current)
	if err != nil {
		return Models{}, err
	}
	currentModel, err := c.clientStore.CurrentModel(controllers.Current)
	if err != nil {
		if errors.Is(err, errors.NotFound) {
			log.Debug().Msg("CurrentModel not found")
		} else {
			return Models{}, err
		}
	}
	return Models{ModelDetails: allModels, Current: currentModel}, nil
}
