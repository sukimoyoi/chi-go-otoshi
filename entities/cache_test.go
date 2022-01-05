package entities_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/sukimoyoi/chi-go-otoshi/entities"
)

func TestGetSiteUrl(t *testing.T) {
	testCases := []struct {
		name  string
		cache entities.CacheData
		site  string
		want  string
	}{
		{
			name:  "empty LastArticles",
			cache: entities.CacheData{},
			site:  "hoge",
			want:  "",
		},
		{
			name: "cache has 'hoge' LastArticles",
			cache: entities.CacheData{
				LastArticles: []entities.LastArticle{
					{
						Site: "hoge",
						URL:  "http://hogehoge.com",
					},
				},
			},
			site: "hoge",
			want: "http://hogehoge.com",
		},
		{
			name: "cache has 'fuga' LastArticles, but not 'hoge' LastArticles",
			cache: entities.CacheData{
				LastArticles: []entities.LastArticle{
					{
						Site: "fuga",
						URL:  "http://fugafuga.com",
					},
				},
			},
			site: "hoge",
			want: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.cache.GetSiteUrl(tc.site)
			if diff := cmp.Diff(tc.want, result); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestSetSiteUrl(t *testing.T) {
	testCases := []struct {
		name  string
		cache entities.CacheData
		site  string
		url   string
		want  string
	}{
		{
			name:  "empty LastArticles",
			cache: entities.CacheData{},
			site:  "hoge",
			url:   "http://hogehoge.com",
			want:  "",
		},
		{
			name: "cache has 'hoge' LastArticles",
			cache: entities.CacheData{
				LastArticles: []entities.LastArticle{
					{
						Site: "hoge",
						URL:  "http://hogehoge.com",
					},
				},
			},
			site: "hoge",
			url:  "http://new.hogehoge.com",
			want: "http://new.hogehoge.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.cache.SetSiteUrl(tc.site, tc.url)
			if diff := cmp.Diff(tc.want, tc.cache.GetSiteUrl(tc.site)); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
