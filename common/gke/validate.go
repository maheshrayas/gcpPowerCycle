package gke

import (
	"log"
	"strconv"
	sch "github.com/maheshrayas/powerCycle/common/schedule"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/compute/v1"
	"regexp"
)

func (k8 *K8Clusters) valdiateTags(nodePool *container.NodePool, clustername string) {
	mNodePool := map[string]string{"Clustername":clustername,
																 "NodePoolName":nodePool.Name,
																 "Status": "stopped",
																 "NodeSize": "0",
																}
	currentNodeCount := strconv.Itoa(getCurrentNodeCount(k8, nodePool.InstanceGroupUrls[0], clustername))
	scheduledLabels := &sch.InstaceTimeDetails{InsLabel: k8.clusterInstances[clustername].ResourceLabels[nodePool.Name],
		Localtimezone: k8.Config.Defaults.Timezone,
		InstanceName:  clustername}

	nodesShouldBeRunning := scheduledLabels.Validate()
	// if there is label with status-nodepoolname
	// strconv.ParseInt(
	if status, ok := k8.clusterInstances[clustername].ResourceLabels["status-"+nodePool.Name]; ok {
		if status == "stopped" && nodesShouldBeRunning {
			mNodePool["NodeSize"]= k8.clusterInstances[clustername].ResourceLabels["nodecount-"+nodePool.Name]
			mNodePool["CurrentNodeCount"] = mNodePool["NodeSize"]
			mNodePool["Status"] = "running"
			k8.StopWorkerNodes(mNodePool)
		} else if status == "running" && !nodesShouldBeRunning {
			mNodePool["CurrentNodeCount"] = currentNodeCount
			k8.StopWorkerNodes(mNodePool)
		}
	} else {
		if !nodesShouldBeRunning {
			mNodePool["CurrentNodeCount"] = currentNodeCount
			k8.StopWorkerNodes(mNodePool)
		}
	}
}

//Inorder to get the current number of nodes running in node pool, we need to invoke InstanceGroups.ListInstances
// GKE API only provides the Initial node count which is configured at the beginnig of node pool creatig.
func getCurrentNodeCount(k8 *K8Clusters, instanceGroupNameURL string, clustername string) (instanceCount int) {
	r := *regexp.MustCompile(`([^/?]*)$`)
	instanceGroup := r.FindString(instanceGroupNameURL)
	rb := &compute.InstanceGroupsListInstancesRequest{}
	req := k8.computeService.InstanceGroups.ListInstances(k8.clusterInstances[clustername].ProjectId, k8.clusterInstances[clustername].Zone,
		instanceGroup, rb)
	if err := req.Pages(k8.Ctx, func(page *compute.InstanceGroupsListInstances) error {
		instanceCount = len(page.Items)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	return
}
