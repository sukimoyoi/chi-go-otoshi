package main

import (
	"log"

	"github.com/sukimoyoi/chi-go-otoshi/gateways"
	"github.com/sukimoyoi/chi-go-otoshi/usecases"
)

func main() {

	cr := &gateways.ConfigFileRepository{
		ConfigFilePath: "./config.yaml",
	}
	config, err := cr.Load()
	if err != nil {
		log.Fatalln(err)
	}

	inputData := &usecases.DownloadInputData{
		Site:              config.Downloader.Sites[0],
		Titles:            config.Downloader.Titles,
		SaveRootDirectory: config.Downloader.SaveRootDirectory,
	}

	if err := usecases.RegularlyDownloadBus(inputData); err != nil {
		log.Fatalln(err)
	}
}
