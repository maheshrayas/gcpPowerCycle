package configuration

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Configs struct {
	Defaults struct {
		Regions  []string `yaml:"regions"`
		Timezone string   `yaml:"timezone"`
		Services []struct {
			Service string   `yaml:"service"`
			Active  bool     `yaml:"active"`
			Tags    []string `yaml:"tags"`
			Action  string   `yaml:"action"`
		} `yaml:"services"`
	} `yaml:"defaults"`
	Projects []struct {
		ProjectID string `yaml:"project_id"`
	} `yaml:"projects"`
}

func (config *Configs) ReadConfig() {
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

type Region struct {
	Zones []string `json:"zones"`
}
