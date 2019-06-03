package gke

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	confg "github.com/maheshrayas/powerCycle/common/configuration"
	// "time"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
)

//InitVMClient   Initialize
func (v *K8Clusters) InitContainerClient() error {
	v.clusterInstances = map[string]*IndividualCluster{}
	c, err := google.DefaultClient(v.Ctx, container.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}
	v.containerService, err = container.New(c)
	if err != nil {
		log.Fatal(err)
	}
	comp, err := google.DefaultClient(v.Ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}
	v.computeService, err = compute.New(comp)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (v *K8Clusters) getZones(project string, region string) []string {
	resp, err := v.computeService.Regions.Get(project, region).Context(v.Ctx).Do()
	if err != nil {
		log.Fatal(err)
	}
	jsonRegions, _ := json.Marshal(resp)
	var regions Region
	json.Unmarshal(jsonRegions, &regions)
	return confg.ParseRegion(&regions.Zones)
}

func (k8 *K8Clusters) GetClusters(project string) []string {
	var wg sync.WaitGroup
	for _, region := range k8.Config.Defaults.Regions {
		fmt.Println("Checking for region", region)
		for _, zone := range k8.getZones(project, region) {
			fmt.Println("Checking for zone", zone)
			clusters, err := k8.containerService.Projects.Zones.Clusters.List(project, zone).Context(k8.Ctx).Do()
			if err != nil {
				log.Fatal("Bomb")
			}
			for _, cl := range clusters.Clusters {
				wg.Add(1)
				go k8.powerCycleGKE(cl, project, &wg)
			}
		}
	}
	wg.Wait()
	return []string{}
}

func (k8 *K8Clusters) powerCycleGKE(cl *container.Cluster, project string, wg *sync.WaitGroup ){
	k8.clusterInstances[cl.Name] = &IndividualCluster{
		Name:           cl.Name,
		ResourceLabels: cl.ResourceLabels,
		Status:         cl.Status,
		NodePools:      cl.NodePools,
		ProjectId:      project,
		Zone:           cl.Zone,
	}
	for _, nodePool := range cl.NodePools {
		k8.valdiateTags(nodePool, cl.Name)
	}
	wg.Done()
}

