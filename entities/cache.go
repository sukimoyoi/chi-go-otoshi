package entities

type CacheData struct {
	LastArticles []LastArticle `yaml:"lastArticles"`
}

type LastArticle struct {
	Site string `yaml:"site"`
	URL  string `yaml:"url"`
}

func (data *CacheData) GetSiteUrl(site string) string {
	for _, a := range data.LastArticles {
		if a.Site == site {
			return a.URL
		}
	}
	return ""
}

func (data *CacheData) SetSiteUrl(site, url string) {
	for i, _ := range data.LastArticles {
		if data.LastArticles[i].Site == site {
			data.LastArticles[i].URL = url
		}
	}
}
