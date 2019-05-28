package gke

import (
	"fmt"
	"sync"
"log"
	sch "github.com/maheshrayas/powerCycle/common/schedule"
	"google.golang.org/api/container/v1"
)

func (k8 *K8Clusters) valdiateTags(project string, region string, name string, wg *sync.WaitGroup) {
	stop := false
	//range over the instances and check if the mandatory tags are present, if not stop the instances
	for _, xx := range k8.Config.Defaults.Services {
		if xx.Service == "gke" && xx.Active == true {
			for _, val := range xx.Tags {
				// check if gke cluster has the tags
				if _, ok := k8.clusterInstances[name].ResourceLabels[val]; !ok {
					stop = true
					fmt.Println("Cluster %s doesn't have %s Label %s", k8.clusterInstances[name], val,stop)
				}
			}
		}
	}
	// stop the instance if the mandatory tags are not present
	if stop == true && k8.clusterInstances[name].Status == "RUNNING" {
		k8.clusterInstances[name].Status = "STOPPED"
		fmt.Println("Stopping instance in the cluster: %s", k8.clusterInstances[name])
	}
	// now check for the scheduled up tim
		scheduledLabels := &sch.InstaceTimeDetails{InsLabel: k8.clusterInstances[name].ResourceLabels["schedule"],
																							 	Localtimezone:k8.Config.Defaults.Timezone,
																								 InstanceName : name,}

		fmt.Println(k8.clusterInstances[name].NodePools[0].Name)
		nodesShouldBeRunning:= scheduledLabels.Validate()
		fmt.Println(nodesShouldBeRunning)
		// if there is label with powerCycle
		if status, ok := k8.clusterInstances[name].ResourceLabels["node-status"]; ok {
			fmt.Println(status)
			if status == "stopped" && nodesShouldBeRunning {
				fmt.Println(k8.clusterInstances[name].NodePools[0].InitialNodeCount)
				k8.addResourceLabels( project, region, name, "running")
				k8.StopWorkerNodes(project, region, name,k8.clusterInstances[name].NodePools[0].Name,k8.clusterInstances[name].NodePools[0].InitialNodeCount)
			}else if status == "running" && !nodesShouldBeRunning{
				fmt.Println("I am here 2")
				k8.addResourceLabels( project, region, name, "stopped")
				k8.StopWorkerNodes(project, region, name,k8.clusterInstances[name].NodePools[0].Name,0)
			}
		}else{
			if stop == false && !nodesShouldBeRunning {
				fmt.Println("I am here 3")
				k8.addResourceLabels( project, region, name, "stopped")
				k8.StopWorkerNodes(project, region, name,k8.clusterInstances[name].NodePools[0].Name,0)
			}
		}
	wg.Done()
}

func (k8 *K8Clusters)addResourceLabels(projectId string, zone string, name string, status string) {
	k8.clusterInstances[name].ResourceLabels["node-status"] = status
	rb := &container.SetLabelsRequest{
		ResourceLabels: k8.clusterInstances[name].ResourceLabels,
	}
	_, err := k8.containerService.Projects.Zones.Clusters.ResourceLabels(projectId, zone, k8.clusterInstances[name].Name, rb).Context(k8.Ctx).Do()
	if err != nil {
		log.Fatal(err)
	}
}

