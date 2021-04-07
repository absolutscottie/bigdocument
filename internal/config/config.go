package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Name          string `json:"name"`
	BindAddr      string `json:"bind_addr"`
	Port          string `json:"port"`
	DatastoreHost string `json:"datastore_host"`
}

//	Load will open the file identified by filename and Unmarshal the content of
//	the file into Config struct.
//	Returns an error when the file was not found, couldn't be parsed, etc
func Load(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(fileBytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, err
}
