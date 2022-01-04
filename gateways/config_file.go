package gateways

import (
	"fmt"
	"io/ioutil"

	"github.com/sukimoyoi/chi-go-otoshi/entities"
	"gopkg.in/yaml.v2"
)

type ConfigFileRepository struct {
	ConfigFilePath string
}

func (cr *ConfigFileRepository) Load() (*entities.ConfigData, error) {
	ouputData := &entities.ConfigData{}
	byte, err := ioutil.ReadFile(cr.ConfigFilePath)
	if err != nil {
		return nil, fmt.Errorf("read config file '%s': %w", cr.ConfigFilePath, err)
	}
	err = yaml.UnmarshalStrict(byte, ouputData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return ouputData, nil
}
