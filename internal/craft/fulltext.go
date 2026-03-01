package craft

import (
	"FeedCraft/internal/util"
	"context"
	"github.com/go-shiori/go-readability"
	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
	"net/url"
	"time"
)

var domainLimiter *util.KeyedLimiter

func init() {
	// Initialize domain limiter with default 3 concurrency per domain
	limit := 3
	envClient := util.GetEnvClient()
	if envClient != nil {
		if envLimit := envClient.GetInt("DOMAIN_MAX_CONCURRENCY"); envLimit > 0 {
			limit = envLimit
		}
	}
	domainLimiter = util.NewKeyedLimiter(limit)
	logrus.Infof("Domain KeyedLimiter initialized with max concurrency: %d", limit)
}

type FulltextExtractor func(url string, timeout time.Duration) (string, error)

// ExtractAction 定义了在获取到域名并发许可后，实际执行网页抓取与提取逻辑的动作
type ExtractAction func() (string, error)

// ExecuteWithDomainLimit 统一处理域名级别的并发控制、Context 超时与资源释放
// 参数 targetURL: 目标网址，用于解析 Host
// 参数 timeout: 期望的超时时间，内部会额外增加缓冲时间用于排队等待锁
// 参数 action: 获取到并发锁后，实际执行的抓取逻辑
func ExecuteWithDomainLimit(targetURL string, timeout time.Duration, action ExtractAction) (string, error) {
	parsed, err := url.Parse(targetURL)
	if err != nil {
		return "", err
	}

	// 统一处理超时控制，额外预留 10 秒给可能发生的排队等待
	ctx, cancel := context.WithTimeout(context.Background(), timeout+10*time.Second)
	defer cancel()

	// 尝试获取该域名的并发许可证
	release, err := domainLimiter.Acquire(ctx, parsed.Host)
	if err != nil {
		logrus.Warnf("Failed to acquire permit for domain %s: %v", parsed.Host, err)
		return "", err
	}
	defer release() // 无论 action 是否 panic/error，确保必定释放

	// 执行真正的抓取/提取逻辑
	return action()
}

func TrivialExtractor(targetURL string, timeout time.Duration) (string, error) {
	var action ExtractAction = func() (string, error) {
		article, err := readability.FromURL(targetURL, timeout)
		if err != nil {
			return "", err
		}
		return article.Content, nil
	}

	return ExecuteWithDomainLimit(targetURL, timeout, action)
}

func GetFulltextCraftOptions() []CraftOption {

	transFunc := func(item *feeds.Item) (string, error) {
		link := item.Link.Href
		return TrivialExtractor(link, DefaultExtractFulltextTimeout)
	}
	cachedTransFunc := GetCommonCachedTransformer(cacheKeyForArticleLink, transFunc, "extract fulltext")
	relativeLinkFixOptions := GetRelativeLinkFixCraftOptions()
	var craftOptions []CraftOption
	craftOptions = append(craftOptions, relativeLinkFixOptions...)
	craftOptions = append(craftOptions, OptionTransformFeedItem(GetArticleContentProcessor(cachedTransFunc)))
	return craftOptions
}
