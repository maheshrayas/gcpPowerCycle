package computeEngine

import (
	"context"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

type CeInstance struct {
	Name   string `json:"name"`
	Labels map[string]string
}

type VmInstances struct {
	instanceDetails map[string]*CeInstance
	Ctx             context.Context
	computeService  *compute.Service
}

func (v *VmInstances) InitVMClient() error {
	v.instanceDetails = map[string]*CeInstance{}
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

func (v *VmInstances) GetInstances(project string, region string) []CeInstance {
	Instances := make([]CeInstance, 0)
	req := v.computeService.Instances.List(project, region+"-b")
	if err := req.Pages(v.Ctx, func(page *compute.InstanceList) error {
		for _, instance := range page.Items {
			// TODO: Change code below to process each `instance` resource:
			v.instanceDetails[instance.Name] = &CeInstance{Name: instance.Name,
				Labels: instance.Labels}
			Instances = append(Instances, CeInstance{Name: instance.Name,
				Labels: instance.Labels})
			log.Println("Printing the instance details : ", instance)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	v.stopVMInstances(project, region)
	return Instances
}

func (v *VmInstances) stopVMInstances(project string, region string) {
	for name, _ := range v.instanceDetails {
		log.Println("stopping the instance  : \n", name)
		go func() {
			_, err := v.computeService.Instances.Stop(project, region+"-b", name).Context(v.Ctx).Do()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
}
