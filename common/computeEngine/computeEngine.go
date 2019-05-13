package computeEngine

import (
	"context"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

type Instance struct {
	Name   string `json:"name"`
	Labels map[string]string
}

type VmInstances struct {
	instanceName   map[string]*Instance
	Ctx            context.Context
	computeService *compute.Service
}

func (v *VmInstances) InitVMClient() error {
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

func (v *VmInstances) GetInstances(project string, region string) []Instance {
	Instances := make([]Instance, 0)
	req := v.computeService.Instances.List(project, region+"-b")
	if err := req.Pages(v.Ctx, func(page *compute.InstanceList) error {
		for _, instance := range page.Items {
			// TODO: Change code below to process each `instance` resource:
			Instances = append(Instances, Instance{Name: instance.Name,
				Labels: instance.Labels})

			log.Println("Printing the instance details : ", instance)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	return Instances
}

// func stopVMInstances() {
// 	resp, err := computeService.Instances.Stop(project, zone, instance).Context(ctx).Do()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// }
