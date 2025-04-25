package jujuclient

import jclient "github.com/juju/juju/jujuclient"

type Controllers struct {
	ControllerDetails map[string]jclient.ControllerDetails `json:"controllerDetails,omitempty"`
	Current           string                               `json:"current,omitempty"`
}

type Models struct {
	ModelDetails map[string]jclient.ModelDetails `json:"modelDetails,omitempty"`
	Current      string                          `json:"current,omitempty"`
}
