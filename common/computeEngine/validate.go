package computeEngine

import (
	"fmt"
	"sync"

	sch "github.com/maheshrayas/powerCycle/common/schedule"
)

func (v *VMInstances) valdiateTags(project string, region string, name string, wg *sync.WaitGroup) {
	stop := false
	//range over the instances and check if the mandatory tags are present, if not stop the instances
	for _, xx := range v.Config.Defaults.Services {
		if xx.Service == "compute" && xx.Active == true {
			for _, val := range xx.Tags {
				// check if compute engine has the tags
				if _, ok := v.instanceDetails[name].Labels[val]; !ok {
					stop = true
					fmt.Printf("Instance %s doesn't have %s Label", v.instanceDetails[name], val)
				}
			}
		}
	}
	// stop the instance if the mandatory tags are not present
	if stop == true && v.instanceDetails[name].State == "RUNNING" {
		v.instanceDetails[name].State = "TERMINATED"
		fmt.Printf("Stopping instance: %s", v.instanceDetails[name])
		v.stopVMInstances(project, region, name)
	}
	// now check for the scheduled up tim
	if stop == false && v.instanceDetails[name].State == "RUNNING" {
		scheduledLabels := &sch.InstaceTimeDetails{InsLabel: v.instanceDetails[name].Labels["schedule"],
			Localtimezone: v.Config.Defaults.Timezone,
			InstanceName:  name}
		if !scheduledLabels.Validate() {
			v.stopVMInstances(project, region, name)
		}
	}
	if stop == false && v.instanceDetails[name].State == "TERMINATED" {
		scheduledLabels := &sch.InstaceTimeDetails{InsLabel: v.instanceDetails[name].Labels["schedule"],
			Localtimezone: v.Config.Defaults.Timezone,
			InstanceName:  name}
		if scheduledLabels.Validate() {
			v.StartVMInstances(project, region, name)
		}
	}
	wg.Done()
}
