package jujuclient

import (
	"context"

	"github.com/juju/errors"
	"github.com/juju/juju/api"
	apiclient "github.com/juju/juju/api/client/client"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/core/logger"
	"github.com/juju/juju/juju"
	jclient "github.com/juju/juju/jujuclient"
	"github.com/rs/zerolog/log"
)

type Client interface {
	GetControllers() (Controllers, error)
	GetModels(controllerName string) (Models, error)
	GetStatus(ctx context.Context, controllerName, modelName string, includeStorage bool) (Status, error)
}

type client struct {
	clientStore jclient.ClientStore
	logger      logger.Logger
}

func NewClient() (Client, error) {
	clientStore := jclient.NewFileClientStore()
	return &client{
		clientStore: clientStore,
		logger:      NewLoggerWrapper(),
	}, nil
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

func (c *client) GetModels(currentController string) (Models, error) {
	if currentController == "" {
		controllers, err := c.GetControllers()
		if err != nil {
			return Models{}, err
		}
		currentController = controllers.Current
	}

	allModels, err := c.clientStore.AllModels(currentController)
	if err != nil {
		return Models{}, err
	}
	currentModel, err := c.clientStore.CurrentModel(currentController)
	if err != nil {
		if errors.Is(err, errors.NotFound) {
			log.Debug().Msg("CurrentModel not found")
		} else {
			return Models{}, err
		}
	}
	return Models{ModelDetails: allModels, Current: currentModel}, nil
}

func (c *client) GetStatus(ctx context.Context, controllerName, modelName string, includeStorage bool) (Status, error) {
	if controllerName == "" {
		currentController, err := c.clientStore.CurrentController()
		if err != nil {
			return Status{}, err
		}
		controllerName = currentController
	}
	if modelName == "" {
		currentModel, err := c.clientStore.CurrentModel(controllerName)
		if err != nil {
			return Status{}, err
		}
		modelName = currentModel
	}

	apiClient, err := c.getAPIClient(ctx, controllerName, modelName)
	if err != nil {
		return Status{}, err
	}
	fullStatus, err := apiClient.Status(ctx, &apiclient.StatusArgs{Patterns: []string{}, IncludeStorage: includeStorage})
	if err != nil {
		return Status{}, err
	}
	return Status{FullStatus: *fullStatus}, nil
}

func (c *client) getAPIClient(ctx context.Context, controllerName string, modelName string) (*apiclient.Client, error) {
	root, err := c.getAPIConn(ctx, controllerName, modelName)
	if err != nil {
		return nil, err
	}
	return apiclient.NewClient(root, c.logger), nil
}

func (c *client) getAPIConn(ctx context.Context, controllerName string, modelName string) (api.Connection, error) {
	accountDetails, err := c.clientStore.AccountDetails(controllerName)
	if err != nil {
		return nil, err
	}
	newAPIConnectionParams, err := c.GetNewAPIConnectionParams(
		c.clientStore,
		controllerName,
		modelName,
		accountDetails,
	)
	if err != nil {
		return nil, err
	}
	apiRoot, err := juju.NewAPIConnection(ctx, newAPIConnectionParams)
	if err != nil {
		return nil, err
	}
	return apiRoot, nil
}

func (c *client) GetNewAPIConnectionParams(
	store jclient.ClientStore,
	controllerName, modelName string,
	accountDetails *jclient.AccountDetails,
) (juju.NewAPIConnectionParams, error) {
	var modelUUID string
	if modelName != "" {
		modelDetails, err := store.ModelByName(controllerName, modelName)
		if err != nil {
			return juju.NewAPIConnectionParams{}, err
		}
		modelUUID = modelDetails.ModelUUID
	}
	dialOpts := api.DefaultDialOpts()
	return juju.NewAPIConnectionParams{
		ControllerStore: store,
		ControllerName:  controllerName,
		AccountDetails:  accountDetails,
		ModelUUID:       modelUUID,
		DialOpts:        dialOpts,
		OpenAPI:         modelcmd.OpenAPIFuncWithMacaroons(api.Open, store, controllerName),
	}, nil
}
