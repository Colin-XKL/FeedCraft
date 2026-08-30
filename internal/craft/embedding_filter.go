package craft

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/model"
	"FeedCraft/internal/util"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
)

/*
 * Embedding 零样本分类主题过滤器
 * 通过将文章和预设的"课题锚点"分别编码为向量，利用余弦相似度进行语义级别的主题匹配。
 * 文章与任一锚点的相似度超过阈值即被保留，否则丢弃。
 */

const (
	defaultEmbeddingThreshold = 0.6
	defaultMaxContentLength   = 2000
)

var errEmbeddingFilterAnchorsRequired = errors.New("[embedding-filter] anchors parameter is required")

// EmbeddingFilterMode 定义 Embedding 过滤器的工作模式
type EmbeddingFilterMode string

var (
	EmbeddingIncludeMode EmbeddingFilterMode = "include" // 匹配锚点的文章被保留（默认）
	EmbeddingExcludeMode EmbeddingFilterMode = "exclude" // 匹配锚点的文章被移除（反选）
)

type EmbeddingFilterProcessor struct {
	anchors       []string
	threshold     float64
	maxContentLen int
	instruction   string
	mode          EmbeddingFilterMode
}

func NewEmbeddingFilterProcessor(anchors []string, threshold float64, maxContentLen int, instruction string, mode EmbeddingFilterMode) *EmbeddingFilterProcessor {
	return &EmbeddingFilterProcessor{
		anchors:       anchors,
		threshold:     threshold,
		maxContentLen: maxContentLen,
		instruction:   instruction,
		mode:          mode,
	}
}

func (p *EmbeddingFilterProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	if feed == nil {
		return feed, nil
	}
	if len(feed.Articles) == 0 {
		return feed, nil
	}

	threshold := p.threshold
	if threshold < 0 {
		logrus.Warnf("[embedding-filter] threshold %.4f is below 0, clamping to 0", threshold)
		threshold = 0
	}
	if threshold > 1 {
		logrus.Warnf("[embedding-filter] threshold %.4f is above 1, clamping to 1", threshold)
		threshold = 1
	}

	if len(p.anchors) == 0 {
		return nil, errEmbeddingFilterAnchorsRequired
	}

	anchorVectors, err := adapter.GetOrComputeAnchorVectors(ctx, p.anchors, p.instruction)
	if err != nil {
		return nil, fmt.Errorf("[embedding-filter] failed to compute anchor vectors: %w", err)
	}

	articles := feed.Articles
	texts := make([]string, len(articles))
	for i, article := range articles {
		texts[i] = buildArticleTextFromArticle(article, p.maxContentLen)
	}

	var validTexts []string
	var validIndices []int
	emptyTextSet := make(map[int]bool)
	for i, text := range texts {
		if len(strings.TrimSpace(text)) > 0 {
			validTexts = append(validTexts, text)
			validIndices = append(validIndices, i)
		} else {
			emptyTextSet[i] = true
		}
	}

	articleVectors := make([][]float64, len(articles))
	var embedErr error
	if len(validTexts) > 0 {
		vectors, batchErr := adapter.EmbedTexts(ctx, validTexts, p.instruction)
		if batchErr != nil {
			embedErr = batchErr
			logrus.Errorf("[embedding-filter] batch embedding failed: %v", batchErr)
		} else {
			if len(vectors) < len(validTexts) {
				logrus.Warnf("[embedding-filter] embedding returned %d vectors for %d texts, some articles may not be properly filtered", len(vectors), len(validTexts))
			}
			for j, idx := range validIndices {
				if j < len(vectors) {
					articleVectors[idx] = vectors[j]
				}
			}
		}
	}

	if embedErr != nil {
		return nil, fmt.Errorf("[embedding-filter] all article embeddings failed: %w", embedErr)
	}

	cloned := cloneCraftFeed(feed)
	filtered := make([]*model.CraftArticle, 0, len(cloned.Articles))
	totalCount := len(cloned.Articles)

	for index, article := range cloned.Articles {
		if article == nil {
			continue
		}
		vec := articleVectors[index]

		var keep bool
		if vec == nil {
			if emptyTextSet[index] {
				keep = true
			} else {
				logrus.Warnf("[embedding-filter] article [%s] has valid text but nil vector (embedding failed), keeping by default", article.Title)
				keep = true
			}
		} else {
			maxSim := -1.0
			for _, anchorVec := range anchorVectors {
				sim := util.CosineSimilarity(vec, anchorVec)
				if sim > maxSim {
					maxSim = sim
				}
			}
			matched := maxSim >= threshold
			switch p.mode {
			case EmbeddingExcludeMode:
				keep = !matched
			default:
				keep = matched
			}
			if keep {
				logrus.Debugf("[embedding-filter] article [%s] KEPT (max similarity: %.4f, threshold: %.4f, mode: %s)", article.Title, maxSim, threshold, p.mode)
			} else {
				logrus.Debugf("[embedding-filter] article [%s] DROPPED (max similarity: %.4f, threshold: %.4f, mode: %s)", article.Title, maxSim, threshold, p.mode)
			}
		}

		if keep {
			filtered = append(filtered, article)
		}
	}

	cloned.Articles = filtered
	keptCount := len(cloned.Articles)
	droppedCount := totalCount - keptCount
	logrus.Infof("[embedding-filter] filtering complete: %d total, %d kept, %d dropped (threshold: %.4f, mode: %s)", totalCount, keptCount, droppedCount, threshold, p.mode)

	return cloned, nil
}

func OptionEmbeddingFilter(anchors []string, threshold float64, maxContentLen int, instruction string, mode EmbeddingFilterMode) LegacyCraftOption {
	return OptionEmbeddingFilterWithContext(context.Background(), anchors, threshold, maxContentLen, instruction, mode)
}

func OptionEmbeddingFilterWithContext(ctx context.Context, anchors []string, threshold float64, maxContentLen int, instruction string, mode EmbeddingFilterMode) LegacyCraftOption {
	processor := NewEmbeddingFilterProcessor(anchors, threshold, maxContentLen, instruction, mode)
	return func(feed *feeds.Feed, payload ExtraPayload) error {
		_ = payload
		return applyLocalProcessorToLegacyFeed(ctx, processor, feed)
	}
}

func buildArticleText(item *feeds.Item, maxLen int) string {
	if item == nil {
		return ""
	}
	return buildArticleTextParts(item.Title, item.Content, item.Description, maxLen)
}

func buildArticleTextFromArticle(article *model.CraftArticle, maxLen int) string {
	if article == nil {
		return ""
	}
	return buildArticleTextParts(article.Title, article.Content, article.Description, maxLen)
}

func buildArticleTextParts(title, content, description string, maxLen int) string {
	if len(content) == 0 {
		content = description
	}

	if utf8.RuneCountInString(content) > maxLen {
		runes := []rune(content)
		content = string(runes[:maxLen])
	}

	if title != "" && content != "" {
		return title + "\n" + content
	}
	if title != "" {
		return title
	}
	return content
}

func GetEmbeddingFilterOptions(anchors []string, threshold float64, maxContentLen int, instruction string, mode EmbeddingFilterMode) []LegacyCraftOption {
	return []LegacyCraftOption{
		OptionEmbeddingFilter(anchors, threshold, maxContentLen, instruction, mode),
	}
}

var embeddingFilterParamTmpl = []ParamTemplate{
	{
		Key:         "anchors",
		Description: "自然语言描述的主题锚点文本，每行一条。文章与任一锚点相似度超过阈值即被视为匹配。",
		Default:     "",
	},
	{
		Key:         "threshold",
		Description: "余弦相似度阈值（0-1），越高越严格。默认 0.6。",
		Default:     "0.6",
	},
	{
		Key:         "mode",
		Description: "过滤模式：include（保留匹配项，默认）或 exclude（移除匹配项，反选）。",
		Default:     "include",
	},
	{
		Key:         "max_content_length",
		Description: "文章正文截取的最大字符数。最终发送给 Embedding 服务的单条输入还会受 FC_EMBEDDING_MAX_INPUT_CHARS 保护。默认 2000。",
		Default:     "2000",
	},
	{
		Key:         "instruction",
		Description: "作为文本前缀拼接到每条 Embedding 输入前。留空则使用全局配置。",
		Default:     "",
	},
}

type embeddingFilterParsedConfig struct {
	anchors       []string
	threshold     float64
	maxContentLen int
	instruction   string
	mode          EmbeddingFilterMode
}

func parseEmbeddingFilterParams(m map[string]string) (embeddingFilterParsedConfig, error) {
	cfg := embeddingFilterParsedConfig{
		threshold:     defaultEmbeddingThreshold,
		maxContentLen: defaultMaxContentLength,
		mode:          EmbeddingIncludeMode,
		instruction:   m["instruction"],
	}

	anchorsStr := m["anchors"]
	if anchorsStr == "" {
		logrus.Warn("[embedding-filter] anchors parameter is empty")
		return cfg, errEmbeddingFilterAnchorsRequired
	}
	rawAnchors := strings.Split(anchorsStr, "\n")
	for _, a := range rawAnchors {
		a = strings.TrimSpace(a)
		if a != "" {
			cfg.anchors = append(cfg.anchors, a)
		}
	}
	if len(cfg.anchors) == 0 {
		logrus.Warn("[embedding-filter] no valid anchors after parsing")
		return cfg, errEmbeddingFilterAnchorsRequired
	}

	if thresholdStr, ok := m["threshold"]; ok && thresholdStr != "" {
		parsed, err := strconv.ParseFloat(thresholdStr, 64)
		if err != nil || parsed < 0 || parsed > 1 {
			logrus.Warnf("[embedding-filter] invalid threshold value [%s], using default %.2f", thresholdStr, defaultEmbeddingThreshold)
		} else {
			cfg.threshold = parsed
		}
	}

	if maxLenStr, ok := m["max_content_length"]; ok && maxLenStr != "" {
		parsed, err := strconv.Atoi(maxLenStr)
		if err != nil || parsed <= 0 {
			logrus.Warnf("[embedding-filter] invalid max_content_length value [%s], using default %d", maxLenStr, defaultMaxContentLength)
		} else {
			cfg.maxContentLen = parsed
		}
	}

	if modeStr, ok := m["mode"]; ok && modeStr != "" {
		switch strings.ToLower(strings.TrimSpace(modeStr)) {
		case "exclude":
			cfg.mode = EmbeddingExcludeMode
		case "include":
			cfg.mode = EmbeddingIncludeMode
		default:
			logrus.Warnf("[embedding-filter] unknown mode [%s], using default 'include'", modeStr)
		}
	}

	return cfg, nil
}

func embeddingFilterLoadParam(m map[string]string) []LegacyCraftOption {
	cfg, err := parseEmbeddingFilterParams(m)
	if err != nil {
		return []LegacyCraftOption{embeddingFilterConfigError(err)}
	}
	return GetEmbeddingFilterOptions(cfg.anchors, cfg.threshold, cfg.maxContentLen, cfg.instruction, cfg.mode)
}

func embeddingFilterConfigError(err error) LegacyCraftOption {
	return func(feed *feeds.Feed, payload ExtraPayload) error {
		return err
	}
}
