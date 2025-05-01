package jujuclient

import (
	jclient "github.com/juju/juju/jujuclient"
	"github.com/juju/juju/rpc/params"
)

type Controllers struct {
	ControllerDetails map[string]jclient.ControllerDetails `json:"controllerDetails,omitempty"`
	Current           string                               `json:"current,omitempty"`
}

type Models struct {
	ModelDetails map[string]jclient.ModelDetails `json:"modelDetails,omitempty"`
	Current      string                          `json:"current,omitempty"`
}

type Status struct {
	FullStatus params.FullStatus `json:"fullStatus,omitempty"`
}

type ApplicationConfig struct {
	ApplicationConfig map[string]interface{} `json:"application-config,omitempty"`
	Application       string                 `json:"application,omitempty"`
	Charm             string                 `json:"charm,omitempty"`
	Settings          map[string]interface{} `json:"settings,omitempty"`
}
