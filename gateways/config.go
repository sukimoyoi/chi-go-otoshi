package gateways

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ConfigRepository struct {
	configFilePath string
}

type ConfigOutputData struct {
	Downloader Downloader `yaml:"downloader"`
}

type Downloader struct {
	Sites             []string `yaml:"sites"`
	Titles            []string `yaml:"titles"`
	SaveRootDirectory string   `yaml:"saveRootDirectory"`
}

func NewConfigRepository(configFilePath string) *ConfigRepository {
	return &ConfigRepository{
		configFilePath: configFilePath,
	}
}

func (cr *ConfigRepository) ReadFromFile() (*ConfigOutputData, error) {
	ouputData := &ConfigOutputData{}
	byte, err := ioutil.ReadFile(cr.configFilePath)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}
	err = yaml.UnmarshalStrict(byte, ouputData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return ouputData, nil
}
