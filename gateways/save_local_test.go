package gateways_test

import (
	"os"
	"testing"

	"github.com/sukimoyoi/chi-go-otoshi/gateways"
)

var repo = &gateways.SaveLocalRepository{}

func TestSaveLocalCreateNumberedFolder(t *testing.T) {
	t.Cleanup(func() {
		os.Remove("./testdata/01")
		os.Remove("./testdata/02")
	})

	folderPath, err := repo.CreateNumberedFolder("./testdata")
	if err != nil {
		t.Errorf("create first numbered folder %s: %w", folderPath, err)
	}
	if _, err := os.Stat(folderPath); err != nil {
		t.Errorf("first numbered folder %s: %w", folderPath, err)
	}

	folderPath, err = repo.CreateNumberedFolder("./testdata")
	if err != nil {
		t.Errorf("create second numbered folder %s: %w", folderPath, err)
	}
	if _, err := os.Stat(folderPath); err != nil {
		t.Errorf("first numbered folder %s: %w", folderPath, err)
	}
}
