package gateways_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sukimoyoi/chi-go-otoshi/entities"
	"github.com/sukimoyoi/chi-go-otoshi/gateways"
)

func TestCacheFileLoad(t *testing.T) {
	repo := &gateways.CacheFileRepository{
		CacheFilePath: "./testdata/cache.yaml",
	}

	data, err := repo.Load()
	if err != nil {
		t.Errorf("%w", err)
	}

	if diff := cmp.Diff("https://hogehoge.com/45665223.html", data.GetSiteUrl("hoge")); diff != "" {
		t.Errorf(diff)
	}
	if diff := cmp.Diff("http://fugafuga.com/58875694.html", data.GetSiteUrl("fuga")); diff != "" {
		t.Errorf(diff)
	}
}

func TestCacheFileSaveAndLoad(t *testing.T) {
	data := &entities.CacheData{
		LastArticles: []entities.LastArticle{
			{
				Site: "nyan",
				URL:  "http://nyannyan.com/9999999.html",
			},
		},
	}

	repo := gateways.CacheFileRepository{
		CacheFilePath: "./testdata/cache_tmp.yaml",
	}
	if err := repo.Save(data); err != nil {
		t.Errorf("%w", err)
	}

	data, err := repo.Load()
	if err != nil {
		t.Errorf("%w", err)
	}
	if diff := cmp.Diff("http://nyannyan.com/9999999.html", data.GetSiteUrl("nyan")); diff != "" {
		t.Errorf(diff)
	}
}
