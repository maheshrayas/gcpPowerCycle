package computeEngine

import (
	"fmt"
	"sync"
)

func (v *VMInstances) valdiateTags(project string, region string, name string, wg *sync.WaitGroup) {
	stop := false
	//range over the instances and check if the mandatory tags are present, if not stop the instances
	for _, xx := range v.Config.Defaults.Services {
		if xx.Service == "compute" && xx.Active == true {
			for _, val := range xx.Tags {
				// check if compute engine has the tags
				if _, ok := v.instanceDetails[name].Labels[val]; !ok {
					v.instanceDetails[name].State = "stopped"
					stop = true
					fmt.Println("Instance %s doesn't have %s Label", v.instanceDetails[name], val)
				}
				if stop == true {
					fmt.Println("Stopping instance: %s", v.instanceDetails[name])
					v.stopVMInstances(project, region, name)
				}
			}
		}
	}
	wg.Done()
}
