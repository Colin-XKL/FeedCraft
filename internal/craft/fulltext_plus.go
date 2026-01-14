package craft

import (
	"FeedCraft/internal/util"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
)

func getRenderedHTML2(websiteUrl string, options util.BrowserlessOptions) (string, error) {
	content, err := util.GetBrowserlessContent(websiteUrl, options)
	if err != nil {
		return "", err
	}

	parseUrl, err := url.Parse(websiteUrl)
	if err != nil {
		logrus.Errorf("parse url failed: %v", err)
		return "", err
	}

	article, err := readability.FromReader(strings.NewReader(content), parseUrl)
	if err != nil {
		return "", err
	}
	return article.Content, err
}

type FulltextPlusConfig struct {
	Wait int    // seconds
	Mode string // networkidle2, etc.
}

func GetFulltextPlusCraftOptions(config FulltextPlusConfig) []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		link := item.Link.Href
		opts := util.BrowserlessOptions{
			Timeout:   DefaultExtractFulltextTimeout,
			WaitUntil: config.Mode,
		}
		if config.Wait > 0 {
			opts.WaitTime = time.Duration(config.Wait) * time.Second
			// Increase total timeout if explicit wait is longer
			if opts.WaitTime > opts.Timeout {
				opts.Timeout = opts.WaitTime + 10*time.Second
			}
		}

		return getRenderedHTML2(link, opts)
	}

	cachedTransFunc := GetCommonCachedTransformer(cacheKeyForArticleLink, transFunc, "extract fulltext plus")

	relativeLinkFixOptions := GetRelativeLinkFixCraftOptions()

	var craftOptions []CraftOption
	craftOptions = append(craftOptions, relativeLinkFixOptions...)
	craftOptions = append(craftOptions, OptionTransformFeedItem(GetArticleContentProcessor(cachedTransFunc)))
	return craftOptions
}

func fulltextPlusLoadParam(m map[string]string) []CraftOption {
	config := FulltextPlusConfig{
		Wait: 0,
		Mode: "networkidle2",
	}

	if val, ok := m["wait"]; ok {
		if v, err := strconv.Atoi(val); err == nil {
			config.Wait = v
		}
	}
	if val, ok := m["mode"]; ok && val != "" {
		config.Mode = val
	}

	return GetFulltextPlusCraftOptions(config)
}

var fulltextPlusParamTmpl = []ParamTemplate{
	{
		Key:         "wait",
		Description: "Wait time in seconds (0 to disable)",
		Default:     "0",
	},
	{
		Key:         "mode",
		Description: "Page load wait mode (load, domcontentloaded, networkidle0, networkidle2)",
		Default:     "networkidle2",
	},
}
