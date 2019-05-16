package functions

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/maheshrayas/powerCycle/common/computeEngine"

	"gopkg.in/yaml.v2"
)

type configs struct {
	Defaults struct {
		Region []string `yaml:"region"`
	} `yaml:"defaults"`
	Projects []struct {
		ProjectID string `yaml:"project_id"`
	} `yaml:"projects"`
}

func (config *configs) readConfig() {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	source := []byte(data)
	err1 := yaml.Unmarshal(source, &config)
	if err1 != nil {
		log.Fatalf("error: %v", err)
	}
}

//PowerCycle Entry point for the cloud functions
func PowerCycle(w http.ResponseWriter, r *http.Request) {
	var projectID string
	config := configs{}
	config.readConfig()
	if config.Projects != nil {
		for _, project := range config.Projects {
			projectID = project.ProjectID
		}
	}

	a := &computeEngine.VMInstances{
		Ctx: context.Background(),
	}
	a.InitVMClient()
	Instances := a.GetInstances(projectID, "australia-southeast1")
	json.NewEncoder(w).Encode(Instances)
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
