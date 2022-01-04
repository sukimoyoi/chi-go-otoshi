package usecases

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sukimoyoi/chi-go-otoshi/entities"
	"github.com/sukimoyoi/chi-go-otoshi/gateways"
)

type DownloadInputData struct {
	Site              string
	Titles            []string
	SaveRootDirectory string
}

func RegularlyDownloadBus(inputData *DownloadInputData) error {
	var downloadTask DownloadInterface

	switch inputData.Site {
	case Anicobin.RootPage.Identifier:
		downloadTask = Anicobin
	case Gno.RootPage.Identifier:
		downloadTask = Gno
	default:
		return fmt.Errorf("unsupported site '%s'", inputData.Site)
	}
	return downloadTask.Regularly(inputData.Titles, inputData.SaveRootDirectory)
}

type DownloadInterface interface {
	Singularly(page entities.WebPage, saveRootDirectory string) error
	Regularly(targetTitles []string, saveRootDirectory string) error
}

type Download struct {
	RootPage                                entities.WebPage
	ArticleSelectorFromRoot                 string
	ImageSelectorFromArticle                string
	ImageSelectorFromArticleExterUrlPattern string
	DownloadInterface
	SaveRepository   SaveRepository
	CacheRespository CacheRepository
}

type SaveRepository interface {
	CreateTitleFolder(folderPath string) error
	CreateNumberedFolder(folderParentPath string) (folderPath string, err error)
	SaveFileReader(filePath string, r io.Reader) error
}

type CacheRepository interface {
	Load() (*entities.CacheData, error)
	Save(*entities.CacheData) error
}

var (
	Anicobin = &Download{
		RootPage:       entities.NewWebPage("anicobin", "あにこ便", "http://anicobin.ldblog.jp/"),
		SaveRepository: &gateways.SaveLocalRepository{},
		CacheRespository: &gateways.CacheFileRepository{
			CacheFilePath: gateways.CachePath,
		},
		ArticleSelectorFromRoot:  "div.autopagerize_page_element > div.top-article-outer > div.top-right > h2.top-article-title > a",
		ImageSelectorFromArticle: "div.tw_matome > a",
	}
	Gno = &Download{
		RootPage:       entities.NewWebPage("gno", "アニメと漫画と 連邦 こっそり日記", "https://gno.blog.jp/"),
		SaveRepository: &gateways.SaveLocalRepository{},
		CacheRespository: &gateways.CacheFileRepository{
			CacheFilePath: gateways.CachePath,
		},
		ArticleSelectorFromRoot:                 "h1.article-index-title > a",
		ImageSelectorFromArticle:                "div.article-body-more > a",
		ImageSelectorFromArticleExterUrlPattern: "livedoor.blogimg",
	}
)

func (d *Download) Singularly(page entities.WebPage, saveRootDirectory string) error {
	return nil
}

func (d *Download) Regularly(targetTitles []string, saveRootDirectory string) error {
	d.CommonRegularly(targetTitles, saveRootDirectory)
	return nil
}

func (d *Download) CommonRegularly(targetTitles []string, saveRootDirectory string) error {
	log.Printf("start regularly donwnload images on '%s'\n", d.RootPage.PrintName)

	lastArticleUrl := ""
	cache, err := d.CacheRespository.Load()
	if err != nil {
		log.Printf("load cache: '%s'\n", err)
	}
	if cache != nil {
		lastArticleUrl = cache.GetSiteUrl(d.RootPage.Identifier)
		log.Printf("last article is '%s'\n", lastArticleUrl)
	}

	baseUrl, err := url.Parse(d.RootPage.URL)
	if err != nil {
		return fmt.Errorf("parse url '%s': %w", baseUrl, err)
	}
	resp, err := http.Get(baseUrl.String())
	if err != nil {
		return fmt.Errorf("access '%s': %w", baseUrl, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("read html: %w", err)
	}

	articlePages := []entities.WebPage{}
	oldSkip := false
	newLastArticleUrl := ""
	doc.Find(d.ArticleSelectorFromRoot).Each(func(_ int, s *goquery.Selection) {
		if oldSkip {
			return
		}

		href, exists := s.Attr("href")
		if !exists {
			return
		}
		aUrl, err := baseUrl.Parse(href)
		if err != nil {
			return
		}
		if newLastArticleUrl == "" {
			newLastArticleUrl = aUrl.String()
		}

		if aUrl.String() == lastArticleUrl {
			oldSkip = true
			return
		}

		hitTitle := ""
		for _, t := range targetTitles {
			if strings.Contains(s.Text(), t) {
				hitTitle = t
				break
			}
		}
		if hitTitle == "" {
			return
		}
		articlePages = append(articlePages, entities.NewWebPage(hitTitle, s.Text(), aUrl.String()))
	})

	for _, ap := range articlePages {
		if err := d.CommonSingularly(ap, saveRootDirectory); err != nil {
			return fmt.Errorf("download singularly: %w", err)
		}
	}

	// update cache
	if cache == nil {
		cache = &entities.CacheData{
			LastArticles: []entities.LastArticle{
				{
					Site: d.RootPage.Identifier,
				},
			},
		}
	} else {
		if cache.GetSiteUrl(d.RootPage.Identifier) == "" {
			cache.LastArticles = append(cache.LastArticles, entities.LastArticle{
				Site: d.RootPage.Identifier,
			})
		}
	}
	cache.SetSiteUrl(d.RootPage.Identifier, newLastArticleUrl)
	if err := d.CacheRespository.Save(cache); err != nil {
		return fmt.Errorf("save cache: %w", err)
	}

	return nil
}

func (d *Download) CommonSingularly(page entities.WebPage, saveRootDirectory string) error {
	log.Printf("download images from '%s' (%s)\n", page.PrintName, page.URL)

	baseUrl, err := url.Parse(page.URL)
	if err != nil {
		return fmt.Errorf("parse url '%s': %w", baseUrl, err)
	}

	resp, err := http.Get(page.URL)
	if err != nil {
		return fmt.Errorf("access '%s': %w", page.URL, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("read html: %w", err)
	}

	saveUrls := []string{}
	doc.Find(d.ImageSelectorFromArticle).Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		saveUrl, err := baseUrl.Parse(href)
		if err != nil {
			return
		}
		if !strings.Contains(saveUrl.String(), d.ImageSelectorFromArticleExterUrlPattern) {
			return
		}
		saveUrls = append(saveUrls, saveUrl.String())

	})

	titleFolderPath := filepath.Join(saveRootDirectory, page.Identifier)
	if err := d.SaveRepository.CreateTitleFolder(titleFolderPath); err != nil {
		return fmt.Errorf("create folder '%s': %w", titleFolderPath, err)
	}

	numberedFolderPath, err := d.SaveRepository.CreateNumberedFolder(titleFolderPath)
	if err != nil {
		return fmt.Errorf("create numbered folder '%s': %w", numberedFolderPath, err)
	}

	log.Printf("save files to '%s'\n", numberedFolderPath)

	number := 1
	for _, su := range saveUrls {
		resp, err = http.Get(su)
		if err != nil {
			log.Printf("download image '%s': %s\n", su, err)
		}
		defer resp.Body.Close()

		filePath := filepath.Join(numberedFolderPath, fmt.Sprintf("%s_%05d%s", d.RootPage.Identifier, number, GetFileExtension(su)))

		if err = d.SaveRepository.SaveFileReader(filePath, resp.Body); err != nil {
			return fmt.Errorf("save image '%s': %w", filePath, err)
		}

		number += 1
	}
	return nil
}

func Contains(s []string, elem string) bool {
	for _, v := range s {
		if v == elem {
			return true
		}
	}
	return false
}

// I'm too lazy to determine the file type, so I do it as follows
func GetFileExtension(url string) string {
	splitted := strings.Split(url, "/")
	if len(splitted) == 1 {
		return ""
	}

	splitted2 := strings.Split(splitted[len(splitted)-1], ".")
	// no extension
	if len(splitted2) == 1 {
		return ""
	}
	return "." + splitted2[len(splitted2)-1]
}
