package main

import (
	"flag"
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type Inventory struct {
	All *All `yaml:"all"`
}

type All struct {
	Hosts    *Hosts    `yaml:"hosts,omitempty"`
	Children *Children `yaml:"children"`
}

type Children struct {
	Child map[string]*Hosts `yaml:",inline"`
}

type Hosts struct {
	Host map[string]string `yaml:"hosts,omitempty"`
}

type Outputs struct {
	Group map[string]*TFEntry `json:"outputs"`
}

type TFEntry struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

func main() {

	pathTFState := flag.String("pathTFState", "./terraform.tfstate", "path to terraform.tfstate")
	pathAnsible := flag.String("pathAnsible", "./inventory.yaml", "path to inventory.yaml")
	flag.Parse()

	terraformState, err := convertFromTFState(*pathTFState)
	if err != nil {
		log.Errorf("failed Converting Terraform State file: %v", err)
	}
	err = saveAnsibleInventory(terraformState, *pathAnsible)
	if err != nil {
		log.Errorf("failed Saving Result: %v", err)
	}
}

func saveAnsibleInventory(data []byte, pathAnsible string) error {
	os.Remove(pathAnsible)
	err := ioutil.WriteFile(pathAnsible, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed writing ansible inventory: %v", err)
	}
	return nil
}

func convertFromTFState(pathTFState string) ([]byte, error) {
	tfStateFile, err := ioutil.ReadFile(pathTFState)
	if err != nil {
		return nil, fmt.Errorf("failed Reading Terraform State file: %v", err)
	}

	tfState := &Outputs{}

	err = json.Unmarshal(tfStateFile, tfState)
	if err != nil {
		return nil, fmt.Errorf("failed Unmarshaling Terraform State File: %v", err)
	}

	inventory := &Inventory{
		All: &All{
			Hosts: nil,
			Children: &Children{
				Child: map[string]*Hosts{},
			},
		},
	}
	for group, addresses := range tfState.Group {
		inventory.All.Children.Child[group] = &Hosts{Host: arrayToMapKeys(addresses.Value)}
	}
	res, err := yaml.Marshal(inventory)
	if err != nil {
		return nil, fmt.Errorf("failed to Marshal yaml: %v", err)
	}
	result := strings.ReplaceAll(string(res), "\"\"", "")
	return []byte(result), nil
}

func arrayToMapKeys(array string) map[string]string {
	source := strings.Split(array, ", ")
	result := map[string]string{}
	for _, i := range source {
		result[i] = ""
	}
	return result
}
