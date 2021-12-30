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

func RegularyDownloadBus(inputData *DownloadInputData) error {
	var downloadTask Download

	switch inputData.Site {
	case Anicobin.RootPage.Identifier:
		downloadTask = Anicobin
	default:
		return fmt.Errorf("unsupported site '%s'", inputData.Site)
	}
	return downloadTask.Regulary(inputData.Titles, inputData.SaveRootDirectory)
}

type Download interface {
	Singulary(targetUrl, saveRootDirectory string) error
	Regulary(targetTitles []string, saveRootDirectory string) error
}

type AnicobinDownload struct {
	Download
	RootPage   entities.WebPage
	Repository SaveRepository
}

type SaveRepository interface {
	CreateFolder(folderPath string) error
	CreateNumberedFolder(folderParentPath string) (folderPath string, err error)
	Save(filePath string, r io.Reader) error
}

var (
	Anicobin = &AnicobinDownload{
		RootPage:   entities.NewWebPage("anicobin", "あにこ便", "http://anicobin.ldblog.jp/"),
		Repository: &gateways.SaveLocalRepository{},
	}
)

func (a *AnicobinDownload) Singulary(targetUrl, saveRootDirectory string) error {
	return nil
}

func (a *AnicobinDownload) Regulary(targetTitles []string, saveRootDirectory string) error {

	log.Printf("start regulary donwnload images on '%s'\n", a.RootPage.PrintName)

	baseUrl, err := url.Parse(a.RootPage.URL)
	if err != nil {
		return fmt.Errorf("pase url '%s': %w", baseUrl, err)
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
	doc.Find("div.autopagerize_page_element > div.top-article-outer > div.top-right > h2.top-article-title > a").Each(func(_ int, s *goquery.Selection) {
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

		href, exists := s.Attr("href")
		if !exists {
			return
		}

		aUrl, err := baseUrl.Parse(href)
		if err == nil {
			articlePages = append(articlePages, entities.NewWebPage(hitTitle, s.Text(), aUrl.String()))
		}
	})

	for _, ap := range articlePages {

		log.Printf("download images from '%s' (%s)\n", ap.PrintName, ap.URL)
		resp, err = http.Get(ap.URL)
		if err != nil {
			return fmt.Errorf("access '%s': %w", ap.URL, err)
		}
		defer resp.Body.Close()

		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return fmt.Errorf("read html: %w", err)
		}

		saveUrls := []string{}
		doc.Find("div.tw_matome > a").Each(func(_ int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists {
				saveUrl, err := baseUrl.Parse(href)
				if err == nil {
					saveUrls = append(saveUrls, saveUrl.String())
				}
			}
		})

		titleFolderPath := filepath.Join(saveRootDirectory, ap.Identifier)
		if err := a.Repository.CreateFolder(titleFolderPath); err != nil {
			return fmt.Errorf("create folder '%s': %w", titleFolderPath, err)
		}

		fmt.Println(titleFolderPath)

		numberedFolderPath, err := a.Repository.CreateNumberedFolder(titleFolderPath)
		if err != nil {
			return fmt.Errorf("create numbered folder '%s': %w", numberedFolderPath, err)
		}

		number := 1
		for _, su := range saveUrls {
			resp, err = http.Get(su)
			if err != nil {
				log.Printf("download image '%s': %s\n", su, err)
			}
			defer resp.Body.Close()

			// TODO: .ついてない可能性もある
			splittedSu := strings.Split(su, ".")
			filePath := filepath.Join(numberedFolderPath, fmt.Sprintf("%s_%05d.%s", a.RootPage.Identifier, number, splittedSu[len(splittedSu)-1]))

			if err = a.Repository.Save(filePath, resp.Body); err != nil {
				return fmt.Errorf("save image '%s': %w", filePath, err)
			}

			number += 1
		}

	}

	return err
}
