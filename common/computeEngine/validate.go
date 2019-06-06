package computeEngine

import (
	"fmt"
	"log"
	"sync"
	sch "github.com/maheshrayas/powerCycle/common/schedule"
)

func (v *VMInstances) valdiateTags(project string, region string, name string, wg *sync.WaitGroup) {
	// check only if the instances have the labe 'Schedule
	if scheduledTime, ok := v.instanceDetails[name].Labels["schedule"]; ok {
		if v.instanceDetails[name].State == "RUNNING" {
			scheduledLabels := &sch.InstaceTimeDetails{InsLabel: scheduledTime,
				Localtimezone: v.Config.Defaults.Timezone,
				InstanceName:  name}
			if !scheduledLabels.Validate() {
				v.stopVMInstances(project, region, name)
			}
		}
		if v.instanceDetails[name].State == "TERMINATED" {
			scheduledLabels := &sch.InstaceTimeDetails{InsLabel: scheduledTime,
				Localtimezone: v.Config.Defaults.Timezone,
				InstanceName:  name}
			if scheduledLabels.Validate() {
				v.StartVMInstances(project, region, name)
			}
		}
	}
	wg.Done()
}
//stopVMInstances : Stop the instances
func (v *VMInstances) stopVMInstances(project string, zone string, name string) {
	_, err := v.computeService.Instances.Stop(project, zone, name).Context(v.Ctx).Do()
	if err != nil {
		log.Fatal(err)
	}
}
//StartVMInstances : Start the instances
func (v *VMInstances) StartVMInstances(project string, zone string, name string) {
	_, err := v.computeService.Instances.Start(project, zone, name).Context(v.Ctx).Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Started instance %s", name)
}
