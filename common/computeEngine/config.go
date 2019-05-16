package computeEngine

import (
	"context"

	"google.golang.org/api/compute/v1"
)

//CeInstances Struct to hold the instance details
type CeInstances struct {
	Name   string `json:"name"`
	Labels map[string]string
}

//VMInstances  Intialize VM Instance struct
type VMInstances struct {
	instanceDetails map[string]*CeInstances
	Ctx             context.Context
	computeService  *compute.Service
}
