package gateways

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sukimoyoi/chi-go-otoshi/entities"
	"gopkg.in/yaml.v2"
)

var CachePath string = ".chi-go-otocshi.cache.yaml"

type CacheFileRepository struct {
	CacheFilePath string
}

func (cr *CacheFileRepository) Load() (*entities.CacheData, error) {
	ouputData := &entities.CacheData{}
	byte, err := ioutil.ReadFile(cr.CacheFilePath)
	if err != nil {
		return nil, fmt.Errorf("read cache file '%s': %w", cr.CacheFilePath, err)
	}
	err = yaml.UnmarshalStrict(byte, ouputData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal cache: %w", err)
	}
	return ouputData, nil
}

func (cr *CacheFileRepository) Save(data *entities.CacheData) error {
	buf, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal cache: %w", err)
	}

	err = ioutil.WriteFile(cr.CacheFilePath, buf, os.ModePerm)
	if err != nil {
		return fmt.Errorf("write cache file '%s': %w", cr.CacheFilePath, err)
	}
	return nil
}
