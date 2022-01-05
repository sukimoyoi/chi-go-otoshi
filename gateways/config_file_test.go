package gateways_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sukimoyoi/chi-go-otoshi/entities"
	"github.com/sukimoyoi/chi-go-otoshi/gateways"
)

func TestConfigFileLoad(t *testing.T) {
	repo := &gateways.ConfigFileRepository{
		ConfigFilePath: "./testdata/config.yaml",
	}

	data, err := repo.Load()
	if err != nil {
		t.Errorf("%w", err)
	}

	want := &entities.ConfigData{
		Downloader: entities.Downloader{
			Sites: []string{"gno", "anicobin"},
			Titles: []string{"大正オトメ御伽話",
				"takt op.Destiny",
				"ガルパ☆ピコ"},
			SaveRootDirectory: "./tmp",
		},
	}

	if diff := cmp.Diff(want, data); diff != "" {
		t.Errorf(diff)
	}
}
