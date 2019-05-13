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

func GetInstances(project string, region string) []Instance {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}
	Instances := make([]Instance, 0)
	req := computeService.Instances.List(project, region+"-b")
	if err := req.Pages(ctx, func(page *compute.InstanceList) error {
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
