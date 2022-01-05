package usecases_test

import (
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/sukimoyoi/chi-go-otoshi/entities"
	"github.com/sukimoyoi/chi-go-otoshi/usecases"
	"github.com/sukimoyoi/chi-go-otoshi/usecases/mock_usecases"
)

func TestRegularlyDownload(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	testHtmlRoot := "./testdata/test_anicobin.html"
	testHtmlRootBytes, err := ioutil.ReadFile(testHtmlRoot)
	if err != nil {
		t.Fatalf("read %s: %s", testHtmlRoot, err)
	}
	testHtmlArticle := "./testdata/test_article.html"
	testHtmlArticleBytes, err := ioutil.ReadFile(testHtmlArticle)
	if err != nil {
		t.Fatalf("read %s: %s", testHtmlArticle, err)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://anicobin.ldblog.jp/",
		httpmock.NewStringResponder(200, string(testHtmlRootBytes)))
	httpmock.RegisterResponder("GET", "http://anicobin.ldblog.jp/archives/58864759.html",
		httpmock.NewStringResponder(200, string(testHtmlArticleBytes)))
	httpmock.RegisterResponder("GET", `=~^https://livedoor\.blogimg\.jp/.*`,
		httpmock.NewStringResponder(200, "test desu"))

	dl := &usecases.Download{
		RootPage:                 entities.NewWebPage("anicobin", "あにこ便", "http://anicobin.ldblog.jp/"),
		ArticleSelectorFromRoot:  "div.autopagerize_page_element > div.top-article-outer > div.top-right > h2.top-article-title > a",
		ImageSelectorFromArticle: "div.tw_matome > a",
	}

	testCases := []struct {
		name      string
		titles    []string
		setupMock func(*mock_usecases.MockSaveRepository, *mock_usecases.MockCacheRepository)
	}{
		{
			name:   "regularly download one title",
			titles: []string{"大正オトメ御伽話"},
			setupMock: func(msr *mock_usecases.MockSaveRepository, mcr *mock_usecases.MockCacheRepository) {
				mcr.EXPECT().Load().Return(&entities.CacheData{
					LastArticles: []entities.LastArticle{
						{
							Site: "anicobin",
							URL:  "http://anicobin.ldblog.jp/",
						},
					},
				}, nil)

				msr.EXPECT().CreateFolder("大正オトメ御伽話")
				msr.EXPECT().CreateNumberedFolder("大正オトメ御伽話")
				msr.EXPECT().SaveFileReader(gomock.Any(), gomock.Any()).AnyTimes()
				mcr.EXPECT().Save(&entities.CacheData{
					LastArticles: []entities.LastArticle{
						{
							Site: "anicobin",
							URL:  "http://anicobin.ldblog.jp/archives/58878526.html",
						},
					},
				})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			saveRepo := mock_usecases.NewMockSaveRepository(c)
			cacheRepo := mock_usecases.NewMockCacheRepository(c)
			tc.setupMock(saveRepo, cacheRepo)

			dl.SaveRepository = saveRepo
			dl.CacheRespository = cacheRepo
			if err := dl.Regularly(tc.titles, ""); err != nil {
				t.Errorf("%w", err)
			}
		})
	}
}
