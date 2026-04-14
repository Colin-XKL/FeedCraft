package craft

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/util"
	"context"
	"strconv"
	"strings"

	"github.com/gorilla/feeds"
	"github.com/samber/lo"
	"github.com/samber/lo/parallel"
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

// OptionEmbeddingFilter 创建 Embedding 主题过滤器的 CraftOption
// anchors: 锚点文本列表
// threshold: 余弦相似度阈值，≥ 阈值即保留
// maxContentLen: 文章正文最大截取长度
// instruction: 传递给 Embedding 模型的 instruction 参数
func OptionEmbeddingFilter(anchors []string, threshold float64, maxContentLen int, instruction string) CraftOption {
	return func(feed *feeds.Feed, payload ExtraPayload) error {
		items := feed.Items
		if len(items) == 0 {
			return nil
		}

		// 1. 校验锚点
		if len(anchors) == 0 {
			logrus.Warn("[embedding-filter] anchors list is empty, skipping filter (returning all items)")
			return nil
		}

		ctx := context.Background()

		// 2. 获取或计算锚点向量（带内存缓存）
		anchorVectors, err := adapter.GetOrComputeAnchorVectors(ctx, anchors, instruction)
		if err != nil {
			logrus.Errorf("[embedding-filter] failed to compute anchor vectors: %v, returning all items", err)
			return nil
		}

		// 3. 并发计算每篇文章的 Embedding 向量
		type articleResult struct {
			vector []float64
			err    error
		}

		results := parallel.Map(items, func(item *feeds.Item, _ int) articleResult {
			// 拼接标题 + 正文
			text := buildArticleText(item, maxContentLen)
			if len(strings.TrimSpace(text)) == 0 {
				return articleResult{err: nil, vector: nil}
			}

			vectors, embedErr := adapter.EmbedTexts(ctx, []string{text}, instruction)
			if embedErr != nil {
				logrus.Warnf("[embedding-filter] failed to embed article [%s]: %v, keeping article", item.Title, embedErr)
				return articleResult{err: embedErr, vector: nil}
			}
			if len(vectors) == 0 {
				logrus.Warnf("[embedding-filter] empty embedding result for article [%s], keeping article", item.Title)
				return articleResult{err: nil, vector: nil}
			}
			return articleResult{err: nil, vector: vectors[0]}
		})

		// 4. 检查是否全部失败
		allFailed := lo.EveryBy(results, func(r articleResult) bool {
			return r.err != nil
		})
		if allFailed && len(results) > 0 {
			logrus.Error("[embedding-filter] all article embeddings failed, returning original feed without filtering")
			return nil
		}

		// 5. 根据相似度过滤
		feed.Items = lo.Filter(items, func(item *feeds.Item, index int) bool {
			r := results[index]

			// Embedding 失败的文章保留（保守策略）
			if r.err != nil || r.vector == nil {
				return true
			}

			// 与所有锚点计算相似度，任一超过阈值即保留
			maxSim := 0.0
			for _, anchorVec := range anchorVectors {
				sim := util.CosineSimilarity(r.vector, anchorVec)
				if sim > maxSim {
					maxSim = sim
				}
			}

			matched := maxSim >= threshold
			if matched {
				logrus.Infof("[embedding-filter] article [%s] MATCHED (max similarity: %.4f >= %.4f)", item.Title, maxSim, threshold)
			} else {
				logrus.Infof("[embedding-filter] article [%s] DROPPED (max similarity: %.4f < %.4f)", item.Title, maxSim, threshold)
			}
			return matched
		})

		return nil
	}
}

// buildArticleText 拼接文章标题和正文，正文超过 maxLen 时截取
func buildArticleText(item *feeds.Item, maxLen int) string {
	content := item.Content
	if len(content) == 0 {
		content = item.Description
	}

	// 截取正文
	if len(content) > maxLen {
		content = content[:maxLen]
	}

	if item.Title != "" && content != "" {
		return item.Title + "\n" + content
	}
	if item.Title != "" {
		return item.Title
	}
	return content
}

// GetEmbeddingFilterOptions 返回 Embedding 过滤器的 CraftOption 列表
func GetEmbeddingFilterOptions(anchors []string, threshold float64, maxContentLen int, instruction string) []CraftOption {
	return []CraftOption{
		OptionEmbeddingFilter(anchors, threshold, maxContentLen, instruction),
	}
}

// --- CraftTemplate 参数定义 ---

var embeddingFilterParamTmpl = []ParamTemplate{
	{
		Key:         "anchors",
		Description: "自然语言描述的主题锚点文本，每行一条。文章与任一锚点相似度超过阈值即被保留。",
		Default:     "",
	},
	{
		Key:         "threshold",
		Description: "余弦相似度阈值（0-1），越高越严格。默认 0.6。",
		Default:     "0.6",
	},
	{
		Key:         "max_content_length",
		Description: "文章正文截取的最大字符数，用于控制 Embedding 输入长度。默认 2000。",
		Default:     "2000",
	},
	{
		Key:         "instruction",
		Description: "传递给 Embedding 模型的 instruction 参数（如果模型支持）。留空则使用全局配置。",
		Default:     "",
	},
}

// embeddingFilterLoadParam 从参数 map 加载 Embedding 过滤器配置
func embeddingFilterLoadParam(m map[string]string) []CraftOption {
	// 解析锚点文本（换行分隔）
	anchorsStr := m["anchors"]
	if anchorsStr == "" {
		logrus.Warn("[embedding-filter] anchors parameter is empty")
		return []CraftOption{}
	}
	rawAnchors := strings.Split(anchorsStr, "\n")
	var anchors []string
	for _, a := range rawAnchors {
		a = strings.TrimSpace(a)
		if a != "" {
			anchors = append(anchors, a)
		}
	}
	if len(anchors) == 0 {
		logrus.Warn("[embedding-filter] no valid anchors after parsing")
		return []CraftOption{}
	}

	// 解析阈值
	threshold := defaultEmbeddingThreshold
	if thresholdStr, ok := m["threshold"]; ok && thresholdStr != "" {
		parsed, err := strconv.ParseFloat(thresholdStr, 64)
		if err != nil || parsed < 0 || parsed > 1 {
			logrus.Warnf("[embedding-filter] invalid threshold value [%s], using default %.2f", thresholdStr, defaultEmbeddingThreshold)
		} else {
			threshold = parsed
		}
	}

	// 解析最大内容长度
	maxContentLen := defaultMaxContentLength
	if maxLenStr, ok := m["max_content_length"]; ok && maxLenStr != "" {
		parsed, err := strconv.Atoi(maxLenStr)
		if err != nil || parsed <= 0 {
			logrus.Warnf("[embedding-filter] invalid max_content_length value [%s], using default %d", maxLenStr, defaultMaxContentLength)
		} else {
			maxContentLen = parsed
		}
	}

	// 解析 instruction
	instruction := m["instruction"]

	return GetEmbeddingFilterOptions(anchors, threshold, maxContentLen, instruction)
}
