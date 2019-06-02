package gke

import (
	"fmt"
	"log"
"strconv"
	sch "github.com/maheshrayas/powerCycle/common/schedule"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/compute/v1"
	"regexp"
)

func (k8 *K8Clusters) valdiateTags(nodePool *container.NodePool, clustername string) {
	currentNodeCount:= getCurrentNodeCount(k8, nodePool.InstanceGroupUrls[0],clustername)
	fmt.Println(k8.clusterInstances[clustername].ResourceLabels[nodePool.Name])
		scheduledLabels := &sch.InstaceTimeDetails{InsLabel: k8.clusterInstances[clustername].ResourceLabels[nodePool.Name],
																							 	Localtimezone:k8.Config.Defaults.Timezone,
																								InstanceName :clustername}

		nodesShouldBeRunning:= scheduledLabels.Validate()
		// if there is label with powerCycle
		if status, ok := k8.clusterInstances[clustername].ResourceLabels["status-"+nodePool.Name]; ok {
			fmt.Println(status)
			if status == "stopped" && nodesShouldBeRunning {
				nodeSize,_:=strconv.ParseInt(k8.clusterInstances[clustername].ResourceLabels["nodecount-"+nodePool.Name],10,64)
				k8.StopWorkerNodes(clustername,nodePool.Name, nodeSize, "running", nodeSize)
			}else if status == "running" && !nodesShouldBeRunning{
				k8.StopWorkerNodes(clustername, nodePool.Name,0, "stopped", strconv.Itoa(currentNodeCount))
			}
		}else{
			if !nodesShouldBeRunning {
				k8.StopWorkerNodes(clustername, nodePool.Name, 0, "stopped",strconv.Itoa(currentNodeCount))
			}
		}
}

//Inorder to get the current number of nodes running in node pool, we need to invoke InstanceGroups.ListInstances
// GKE API only provides the Initial node count which is configured at the beginnig of node pool creatig.
func getCurrentNodeCount(k8 *K8Clusters, instanceGroupNameURL string, clustername string)(instanceCount int){
	r := *regexp.MustCompile(`([^/?]*)$`)
	instanceGroup := r.FindString(instanceGroupNameURL)
	rb := &compute.InstanceGroupsListInstancesRequest{
	}
	req := k8.computeService.InstanceGroups.ListInstances(k8.clusterInstances[clustername].ProjectId, k8.clusterInstances[clustername].Zone,
				instanceGroup, rb)
				if err := req.Pages(k8.Ctx, func(page *compute.InstanceGroupsListInstances) error {
					instanceCount = len(page.Items)
					return nil
	});err != nil {
		log.Fatal(err)
	}
	return
}
