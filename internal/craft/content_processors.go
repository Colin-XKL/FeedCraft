package craft

import (
	"context"
	"fmt"
	"strings"
	"time"

	"FeedCraft/internal/model"
	"FeedCraft/internal/util"

	"github.com/sirupsen/logrus"
)

type ArticleCacheKeyGenerator func(article *model.CraftArticle) (string, error)
type ArticleTransformFunc func(ctx context.Context, article *model.CraftArticle) (string, error)
type CleanupTransformFunc func(content string, domain string) (string, error)
type FulltextPlusExtractor func(url string, options util.BrowserlessOptions) (string, error)

var cleanupTransformFunc CleanupTransformFunc = CleanupContent
var fulltextExtractFunc FulltextExtractor = TrivialExtractor
var fulltextPlusExtractFunc FulltextPlusExtractor = getRenderedHTML2

type ArticleContentTransformProcessor struct {
	CraftName string
	Transform ArticleTransformFunc
}

func (p *ArticleContentTransformProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	if feed == nil || p == nil || p.Transform == nil || len(feed.Articles) == 0 {
		return feed, nil
	}

	cloned := cloneCraftFeed(feed)
	var (
		lastErr   error
		successes int
		attempted int
	)

	for _, article := range cloned.Articles {
		if article == nil {
			continue
		}
		attempted += 1
		transformed, err := p.Transform(ctx, article)
		if err != nil {
			lastErr = err
			logrus.Warnf("failed to apply craft [%s] for article [%s], err: %v", p.CraftName, article.Title, err)
			continue
		}
		article.Content = transformed
		article.Description = transformed
		successes += 1
	}

	if attempted > 0 && successes == 0 {
		return nil, fmt.Errorf("all items failed to process. last error: %v", lastErr)
	}
	return cloned, nil
}

type CleanupProcessor struct {
	Processor *ArticleContentTransformProcessor
}

func (p *CleanupProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	if p == nil || p.Processor == nil {
		return feed, nil
	}
	return p.Processor.Process(ctx, feed)
}

type FulltextProcessor struct {
	OriginalFeedURL string
	Processor       *ArticleContentTransformProcessor
}

func (p *FulltextProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	if p == nil || p.Processor == nil {
		return feed, nil
	}
	fixedFeed, err := applyRelativeLinkFix(ctx, feed, p.OriginalFeedURL)
	if err != nil {
		return nil, err
	}
	return p.Processor.Process(ctx, fixedFeed)
}

type FulltextPlusProcessor struct {
	OriginalFeedURL string
	Config          FulltextPlusConfig
	Processor       *ArticleContentTransformProcessor
}

func (p *FulltextPlusProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	if p == nil || p.Processor == nil {
		return feed, nil
	}
	fixedFeed, err := applyRelativeLinkFix(ctx, feed, p.OriginalFeedURL)
	if err != nil {
		return nil, err
	}
	return p.Processor.Process(ctx, fixedFeed)
}

func newCleanupProcessor() *CleanupProcessor {
	transformer := GetCommonCachedArticleTransformer(
		cacheKeyForCraftArticleContent,
		func(ctx context.Context, article *model.CraftArticle) (string, error) {
			content := article.Content
			if strings.TrimSpace(content) == "" {
				content = article.Description
			}
			domain, _ := util.ParseDomainFromUrl(article.Link)
			return cleanupTransformFunc(content, domain)
		},
		"cleanup article content",
	)

	return &CleanupProcessor{
		Processor: &ArticleContentTransformProcessor{
			CraftName: "cleanup article content",
			Transform: transformer,
		},
	}
}

func newFulltextProcessor(originalFeedURL string) *FulltextProcessor {
	transformer := GetCommonCachedArticleTransformer(
		cacheKeyForCraftArticleLink,
		func(ctx context.Context, article *model.CraftArticle) (string, error) {
			return fulltextExtractFunc(article.Link, DefaultExtractFulltextTimeout)
		},
		"extract fulltext",
	)

	return &FulltextProcessor{
		OriginalFeedURL: originalFeedURL,
		Processor: &ArticleContentTransformProcessor{
			CraftName: "extract fulltext",
			Transform: transformer,
		},
	}
}

func newFulltextPlusProcessor(originalFeedURL string, config FulltextPlusConfig) *FulltextPlusProcessor {
	transformer := GetCommonCachedArticleTransformer(
		cacheKeyForCraftArticleLink,
		func(ctx context.Context, article *model.CraftArticle) (string, error) {
			opts := util.BrowserlessOptions{
				Timeout:   DefaultExtractFulltextTimeout,
				WaitUntil: config.Mode,
			}
			if config.Wait > 0 {
				opts.WaitTime = time.Duration(config.Wait) * time.Second
				if opts.WaitTime > opts.Timeout {
					opts.Timeout = opts.WaitTime + 10*time.Second
				}
			}
			return fulltextPlusExtractFunc(article.Link, opts)
		},
		"extract fulltext plus",
	)

	return &FulltextPlusProcessor{
		OriginalFeedURL: originalFeedURL,
		Config:          config,
		Processor: &ArticleContentTransformProcessor{
			CraftName: "extract fulltext plus",
			Transform: transformer,
		},
	}
}

func GetCommonCachedArticleTransformer(cacheKeyGenerator ArticleCacheKeyGenerator, rawTransformer ArticleTransformFunc, craftName string) ArticleTransformFunc {
	return func(ctx context.Context, article *model.CraftArticle) (string, error) {
		hashVal, err := cacheKeyGenerator(article)
		if err != nil {
			return "", err
		}
		cacheKey := getCraftCacheKey(craftName, hashVal)
		return util.CachedFuncWithPreLog(cacheKey, func() (string, error) {
			return rawTransformer(ctx, article)
		}, func(isCached bool) {
			logrus.Infof("applying craft [%s] to article [%s], cached: %v", craftName, article.Title, isCached)
		})
	}
}

func cacheKeyForCraftArticleContent(article *model.CraftArticle) (string, error) {
	content := article.Content
	if strings.TrimSpace(content) == "" {
		content = article.Description
	}
	return util.GetTextContentHash(content), nil
}

func cacheKeyForCraftArticleLink(article *model.CraftArticle) (string, error) {
	uniqLinkStr := article.Title + article.Id + article.Link
	return util.GetTextContentHash(uniqLinkStr), nil
}

func applyRelativeLinkFix(ctx context.Context, feed *model.CraftFeed, originalFeedURL string) (*model.CraftFeed, error) {
	if strings.TrimSpace(originalFeedURL) == "" {
		return feed, nil
	}
	return (&RelativeLinkFixProcessor{OriginalFeedURL: originalFeedURL}).Process(ctx, feed)
}
