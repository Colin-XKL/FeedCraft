package craft

import (
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
)

const browserFailCircuitThreshold = 3

func getRenderedHTML2(websiteUrl string, options util.BrowserlessOptions) (string, error) {
	parseUrl, err := url.Parse(websiteUrl)
	if err != nil {
		logrus.Errorf("parse url failed: %v", err)
		return "", err
	}
	host := parseUrl.Hostname()
	if host == "" {
		return "", fmt.Errorf("empty hostname in URL: %s", websiteUrl)
	}

	// Acquire concurrency permit for this domain
	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout+10*time.Second)
	defer cancel()

	release, err := domainLimiter.Acquire(ctx, host)
	if err != nil {
		logrus.Warnf("Failed to acquire permit for domain %s: %v", host, err)
		return "", err
	}
	defer release()

	content, err := util.GetBrowserlessContent(websiteUrl, options)
	if err != nil {
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
	Mode string // load, networkidle2, etc.
}

type browserFailBudget struct {
	mu    sync.Mutex
	fails int
}

func (b *browserFailBudget) exhausted() bool {
	if b == nil {
		return false
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.fails >= browserFailCircuitThreshold
}

func (b *browserFailBudget) success() {
	if b == nil {
		return
	}
	b.mu.Lock()
	b.fails = 0
	b.mu.Unlock()
}

func (b *browserFailBudget) failure() {
	if b == nil {
		return
	}
	b.mu.Lock()
	b.fails++
	b.mu.Unlock()
}

func buildFulltextPlusBrowserOptions(config FulltextPlusConfig) util.BrowserlessOptions {
	opts := util.BrowserlessOptions{
		Timeout:   util.ResolveBrowserRenderTimeout(),
		WaitUntil: config.Mode,
	}
	if config.Wait > 0 {
		opts.WaitTime = time.Duration(config.Wait) * time.Second
		if opts.WaitTime > opts.Timeout {
			opts.Timeout = opts.WaitTime + 10*time.Second
		}
	}
	return opts
}

func runFulltextPlusExtract(link string, opts util.BrowserlessOptions, budget *browserFailBudget) (string, error) {
	if budget.exhausted() {
		logrus.Warnf("browser render circuit open, using HTTP extract for %s", link)
		return TrivialExtractor(link, DefaultExtractFulltextTimeout)
	}

	content, err := fulltextPlusExtractFunc(link, opts)
	if err == nil {
		budget.success()
		return content, nil
	}

	budget.failure()
	logrus.Warnf("browser fulltext failed for %s, falling back to HTTP extract: %v", link, err)
	fallback, ferr := TrivialExtractor(link, DefaultExtractFulltextTimeout)
	if ferr != nil {
		return "", fmt.Errorf("%w; http fallback failed: %v", err, ferr)
	}
	return fallback, nil
}

func GetFulltextPlusCraftOptions(config FulltextPlusConfig) []LegacyCraftOption {
	budget := &browserFailBudget{}
	transFunc := func(item *feeds.Item) (string, error) {
		return runFulltextPlusExtract(item.Link.Href, buildFulltextPlusBrowserOptions(config), budget)
	}

	cachedTransFunc := GetCommonCachedTransformer(cacheKeyForArticleLink, transFunc, "extract fulltext plus")

	relativeLinkFixOptions := GetRelativeLinkFixCraftOptions()

	var craftOptions []LegacyCraftOption
	craftOptions = append(craftOptions, relativeLinkFixOptions...)
	craftOptions = append(craftOptions, OptionTransformFeedItem(GetArticleContentProcessor(cachedTransFunc)))
	return craftOptions
}

func fulltextPlusLoadParam(m map[string]string) []LegacyCraftOption {
	return GetFulltextPlusCraftOptions(parseFulltextPlusConfig(m))
}

func parseFulltextPlusConfig(m map[string]string) FulltextPlusConfig {
	config := FulltextPlusConfig{
		Wait: 0,
		Mode: "load",
	}

	if val, ok := m["wait"]; ok {
		if v, err := strconv.Atoi(val); err == nil {
			config.Wait = v
		}
	}
	if val, ok := m["mode"]; ok && val != "" {
		config.Mode = val
	}

	return config
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
		Default:     "load",
	},
}
