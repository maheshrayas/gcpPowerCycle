package computeEngine

import (
	"log"

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
func (v *VMInstances) GetInstances(project string, region string) []CeInstances {
	Instances := make([]CeInstances, 0)
	req := v.computeService.Instances.List(project, region+"-b")
	if err := req.Pages(v.Ctx, func(page *compute.InstanceList) error {
		for _, instance := range page.Items {
			// TODO: Change code below to process each `instance` resource:
			v.instanceDetails[instance.Name] = &CeInstances{Name: instance.Name,
				Labels: instance.Labels}
			Instances = append(Instances, CeInstances{Name: instance.Name,
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

func (v *VMInstances) stopVMInstances(project string, region string) {
	for name := range v.instanceDetails {
		log.Println("stopping the instance  : \n", name)
		go func(instanceName string) {
			_, err := v.computeService.Instances.Stop(project, region+"-b", instanceName).Context(v.Ctx).Do()
			if err != nil {
				log.Fatal(err)
			}
		}(name)
	}
}
