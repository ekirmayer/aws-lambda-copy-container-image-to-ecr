package main

import (
	"encoding/json"
	"fmt"
	"os"
	"log"
)

// Define a struct for DockerConfig
type DockerConfig struct {
	Auths map[string]map[string]string `json:"auths"`
}

// Function to add a new registry
func (dc *DockerConfig) AddRegistry(registry string, auth string) {
	dc.Auths[registry] = map[string]string{"auth": auth}
}

var dockerConfig DockerConfig

func init() {
	// Create an instance of the DockerConfig struct
	dockerConfig = DockerConfig{
		Auths: map[string]map[string]string{},
	}
}

func GenerateAuth() (string, error) {
	path := os.TempDir() + "/docker"
	os.MkdirAll(path, 0777)
	file, err := os.Create(path + "/config.json")
	if err != nil {
		fmt.Println(err)
		return path, err
	}
	err = os.Chmod(file.Name(), 0700)
    if err != nil {
        log.Fatal(err)
    }
	defer file.Close()

	// Write the struct into the JSON file
	encoder := json.NewEncoder(file)
	err = encoder.Encode(dockerConfig)
	if err != nil {
		fmt.Println(err)
	}
	
	return path, nil
}
