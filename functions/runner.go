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

type Instance struct {
	Name string `json:"Name"`
}

type Instances []Instance

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

func PowerCycle(w http.ResponseWriter, r *http.Request) {
	var projectId string
	config := configs{}
	config.readConfig()
	if config.Projects != nil {
		for _, project := range config.Projects {
			projectId = project.ProjectID
		}
	}

	a := &computeEngine.VmInstances{
		Ctx: context.Background(),
	}
	a.InitVMClient()
	Instances := a.GetInstances(projectId, "australia-southeast1")
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

// func handleRequests() {
// 	http.HandleFunc("/", PowerUtility)
// 	log.Fatal(http.ListenAndServe(":8081", nil))
// }

// func main() {
// 	handleRequests()
// }
