package gke

import (
	"context"

	"github.com/maheshrayas/powerCycle/common/configuration"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
)

type IndividualCluster struct {
	Name             string            `json:"name"`
	Locations        []string          `json:"locations"`
	ResourceLabels   map[string]string `json:"resourceLabels"`
	Status           string            `json:"status"`
	NodePools        []*container.NodePool
	ProjectId        string
	Zone             string
}

//VMInstances  Intialize VM Instance struct
type K8Clusters struct {
	clusterInstances map[string]*IndividualCluster
	Ctx              context.Context
	containerService *container.Service
	Config           *configuration.Configs
	computeService   *compute.Service
}

// Region  struct
type Region struct {
	Zones []string `json:"zones"`
}
