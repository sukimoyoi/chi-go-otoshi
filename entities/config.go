package entities

type ConfigData struct {
	Downloader Downloader `yaml:"downloader"`
}

type Downloader struct {
	Sites             []string `yaml:"sites"`
	Titles            []string `yaml:"titles"`
	SaveRootDirectory string   `yaml:"saveRootDirectory"`
}
