package computeEngine

import (
	"log"
)

func (v *VMInstances) stopVMInstances(project string, zone string, name string) {
	_, err := v.computeService.Instances.Stop(project, zone, name).Context(v.Ctx).Do()
	if err != nil {
		log.Fatal(err)
	}
}
