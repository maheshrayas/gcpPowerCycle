package computeEngine

import (
	"encoding/json"
	"log"

	"sync"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

//InitVMClient   Initialize
func (v *VMInstances) InitVMClient() error {
	v.instanceDetails = map[string]*CeInstances{}
	c, err := google.DefaultClient(v.Ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}
	v.computeService, err = compute.New(c)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (v *VMInstances) getZones(project string, region string) []string {
	resp, err := v.computeService.Regions.Get(project, region).Context(v.Ctx).Do()
	if err != nil {
		log.Fatal(err)
	}
	jsonRegions, _ := json.Marshal(resp)
	var regions Region
	json.Unmarshal(jsonRegions, &regions)
	return parseRegion(&regions.Zones)
}

//GetInstances Lists all the running instances
func (v *VMInstances) GetInstances(project string, region string) []CeInstances {
	Instances := make([]CeInstances, 0)
	var wg sync.WaitGroup
	for _, zone := range v.getZones(project, region) {
		wg.Add(1)
		req := v.computeService.Instances.List(project, zone)
		if err := req.Pages(v.Ctx, func(page *compute.InstanceList) error {
			for _, instance := range page.Items {
				// TODO: Change code below to process each `instance` resource:
				v.instanceDetails[instance.Name] = &CeInstances{Name: instance.Name,
					Labels: instance.Labels, Zone: zone}
				Instances = append(Instances, CeInstances{Name: instance.Name,
					Labels: instance.Labels, Zone: zone})
				go v.valdiateTags(project, zone, instance.Name, &wg)
				// log.Println("Printing the instance details : ", instance)
			}
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}
	wg.Wait()
	return Instances
}
