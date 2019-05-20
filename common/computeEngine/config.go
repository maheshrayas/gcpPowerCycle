package computeEngine

import (
	"context"

	"github.com/maheshrayas/powerCycle/common/configuration"
	"google.golang.org/api/compute/v1"
)

//CeInstances Struct to hold the instance details
type CeInstances struct {
	Name   string `json:"name"`
	Labels map[string]string
	State  string
	Zone   string
}

//VMInstances  Intialize VM Instance struct
type VMInstances struct {
	instanceDetails map[string]*CeInstances
	Ctx             context.Context
	computeService  *compute.Service
	Config          *configuration.Configs
}

// Region  struct
type Region struct {
	Zones []string `json:"zones"`
}
