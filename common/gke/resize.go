package gke

import (
	"log"
	"google.golang.org/api/container/v1"
	"fmt"
	"time"
	"strconv"

)
const(
	statusRunning = "RUNNING"
)
func (k8 *K8Clusters) StopWorkerNodes(mNodePool map[string]string) {
	nodeCount, _ := strconv.ParseInt(mNodePool["NodeSize"], 10, 64)
	rb := &container.SetNodePoolSizeRequest{
		NodeCount: nodeCount,
	}
	if err := k8.waitForCluster(mNodePool["Clustername"]); err != nil {
		log.Fatal("Contains error", err)
	}
	k8.addResourceLabels(mNodePool)
	resp, err := k8.containerService.Projects.Zones.Clusters.NodePools.SetSize(k8.clusterInstances[mNodePool["Clustername"]].ProjectId,
		k8.clusterInstances[mNodePool["Clustername"]].Zone,
		mNodePool["Clustername"],
		mNodePool["NodePoolName"],
		rb).Context(k8.Ctx).Do()
	if err != nil {
		log.Fatal("Contains error", err)
	}
	fmt.Println(resp)
}

//Only one nodepool can be updated at a given point of time
//When the nodepool size is changed, the cluster is status RECONCILING
// so the nodepool size change must be sequential and before the next change is done, the cluster need to be in the status=RUNNING
func (k8 *K8Clusters) waitForCluster(clustername string) error {
	message := ""
	for {
		cluster, err := k8.containerService.Projects.Zones.Clusters.Get(k8.clusterInstances[clustername].ProjectId,
			k8.clusterInstances[clustername].Zone,
			clustername).Context(k8.Ctx).Do()
		if err != nil {
			return err
		}
		if cluster.Status == statusRunning {
			log.Printf("Cluster %v is running", clustername)
			return nil
		}
		if cluster.Status != message {
			log.Printf("%v cluster %v", string(cluster.Status), clustername)
			message = cluster.Status
		}
		time.Sleep(time.Second * 5)
	}
}

// add the 2 new resource label for GKE cluster
// status-nodepoolid : running or stopped
// nodecount-nodepoolid : set the current running number of instances so that when the node size is set to this count
func (k8 *K8Clusters)addResourceLabels(mNodePool map[string]string) {
	clustername:=mNodePool["Clustername"]
	k8.clusterInstances[clustername].ResourceLabels["status-"+ mNodePool["NodePoolName"]] = mNodePool["Status"]
	k8.clusterInstances[clustername].ResourceLabels["nodecount-"+ mNodePool["NodePoolName"]] = mNodePool["CurrentNodeCount"]
	rb := &container.SetLabelsRequest{
		ResourceLabels: k8.clusterInstances[clustername].ResourceLabels,
	}
	_, err := k8.containerService.Projects.Zones.Clusters.ResourceLabels(k8.clusterInstances[clustername].ProjectId , k8.clusterInstances[clustername].Zone, clustername, rb).Context(k8.Ctx).Do()
	if err != nil {
		log.Fatal(err)
	}
}