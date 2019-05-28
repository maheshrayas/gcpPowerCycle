package main

import (
	"context"
// "time"
	"fmt"

	// "github.com/maheshrayas/powerCycle/common/computeEngine"
	"github.com/maheshrayas/powerCycle/common/configuration"
	"github.com/maheshrayas/powerCycle/common/gke"
)

//PowerCycle Entry point for the cloud functions
func main() {
	var projectID string
	config := &configuration.Configs{}
	config.ReadConfig()
	if config.Projects != nil {
		for _, project := range config.Projects {
			projectID = project.ProjectID
		}
	}
	// a := &computeEngine.VMInstances{
	// 	Ctx: context.Background(),
	// 	Config: config,
	// }
	// a.InitVMClient()
	a := &gke.K8Clusters{
		Ctx: context.Background(),
		Config: config,
	}
	a.InitContainerClient()
	// fmt.Println(projectID)
	Instances := a.GetClusters(projectID)
	fmt.Println(Instances)
	// json.NewEncoder(w).Encode(Instances)

	// jw := writers.NewMessageWriter(Instances)
	// jsonString, err := Instances.JSONString()
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	log.Println(err.Error())
	// 	return
	// }
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(jsonString))
}
