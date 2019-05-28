package gke

import (
	"log"
	"google.golang.org/api/container/v1"
	"fmt"
)


func (k8 *K8Clusters) StopWorkerNodes(projectId string, zone string, clusterId string, nodePoolId string, nodeSize int64) {
	rb := &container.SetNodePoolSizeRequest{
		NodeCount:nodeSize,
	}
	// create a new label powerCycle and set it as stopped
	resp, err := k8.containerService.Projects.Zones.Clusters.NodePools.SetSize(projectId, zone, clusterId, nodePoolId, rb).Context(k8.Ctx).Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}