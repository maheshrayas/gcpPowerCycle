package gke

import (
	"log"
	"google.golang.org/api/container/v1"
	"fmt"
	"time"

)
const(
	statusRunning     = "RUNNING"
)

func (k8 *K8Clusters) StopWorkerNodes(clustername string, nodePoolId string, nodeSize int64, status string, currentNodeCount string) {
	fmt.Println("HIHI",currentNodeCount)
	rb := &container.SetNodePoolSizeRequest{
		NodeCount: nodeSize,
	}
	if err := k8.waitForCluster(clustername); err != nil {
		log.Fatal("Contains error", err)
	}
	k8.addResourceLabels(clustername, status, nodePoolId, currentNodeCount)
	resp, err := k8.containerService.Projects.Zones.Clusters.NodePools.SetSize(k8.clusterInstances[clustername].ProjectId,
		k8.clusterInstances[clustername].Zone,
		clustername,
		nodePoolId,
		rb).Context(k8.Ctx).Do()
	if err != nil {
		log.Fatal("Contains error", err)
	}
	fmt.Println(resp)
}


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

func (k8 *K8Clusters)addResourceLabels(clustername string, status string, nodePoolId string, currentNodeCount string) {
	k8.clusterInstances[clustername].ResourceLabels["status-"+nodePoolId] = status
	k8.clusterInstances[clustername].ResourceLabels["nodecount-"+nodePoolId] = currentNodeCount
	rb := &container.SetLabelsRequest{
		ResourceLabels: k8.clusterInstances[clustername].ResourceLabels,
	}
	_, err := k8.containerService.Projects.Zones.Clusters.ResourceLabels(k8.clusterInstances[clustername].ProjectId , k8.clusterInstances[clustername].Zone, clustername, rb).Context(k8.Ctx).Do()
	if err != nil {
		log.Fatal(err)
	}
}