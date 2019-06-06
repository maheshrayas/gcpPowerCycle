package computeEngine

import (
	"fmt"
	"log"
	"sync"
	confg "github.com/maheshrayas/powerCycle/common/configuration"
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

//GetInstances Lists all the running instances
func (v *VMInstances) GetInstances(project string)  {
	var wg sync.WaitGroup
	for _, region := range v.Config.Defaults.Regions {
		fmt.Println("Checking for reqion", region)
		for _, zone := range confg.GetZones(v.Ctx, v.computeService, project, region) {
			fmt.Println("Checking for zone", zone)
			req := v.computeService.Instances.List(project, zone)
			if err := req.Pages(v.Ctx, func(page *compute.InstanceList) error {
				for _, instance := range page.Items {
					wg.Add(1)
					v.instanceDetails[instance.Name] = &CeInstances{Name: instance.Name,
						Labels: instance.Labels, Zone: zone, State: instance.Status}
					go v.valdiateTags(project, zone, instance.Name, &wg)
				}
				return nil
			}); err != nil {
				log.Fatal(err)
			}
		}
	}
	wg.Wait()
}
