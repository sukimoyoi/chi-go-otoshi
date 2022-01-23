package main

import (
	"flag"
	"log"

	"github.com/sukimoyoi/chi-go-otoshi/gateways"
	"github.com/sukimoyoi/chi-go-otoshi/usecases"
)

var (
	runType = flag.String("type", "from-root-page", "download images [from-root-page, single-page]")
)

func main() {
	flag.Parse()

	cr := &gateways.ConfigFileRepository{
		ConfigFilePath: "./config.yaml",
	}
	config, err := cr.Load()
	if err != nil {
		log.Fatalln(err)
	}

	switch *runType {
	case "from-root-page":
		for _, site := range config.Downloader.Sites {
			inputData := &usecases.DownloadInputData{
				Site:              site,
				Titles:            config.Downloader.Titles,
				SaveRootDirectory: config.Downloader.SaveRootDirectory,
			}

			if err := usecases.FromRootPageDownloadBus(inputData); err != nil {
				log.Println(err)
			}
		}
	case "single-page":
	default:
		log.Fatalf("Unsupport download type '%s'\n", *runType)
	}

}
