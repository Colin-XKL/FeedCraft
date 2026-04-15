package craft

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gorilla/feeds"
	"github.com/samber/lo"
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

		// 0. 阈值 clamp：将 threshold 钳制到 [0, 1] 范围
		if threshold < 0 {
			logrus.Warnf("[embedding-filter] threshold %.4f is below 0, clamping to 0", threshold)
			threshold = 0
		}
		if threshold > 1 {
			logrus.Warnf("[embedding-filter] threshold %.4f is above 1, clamping to 1", threshold)
			threshold = 1
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
			return fmt.Errorf("[embedding-filter] failed to compute anchor vectors: %w", err)
		}

		// 3. 收集所有文章文本，批量调用 Embedding
		texts := make([]string, len(items))
		for i, item := range items {
			texts[i] = buildArticleText(item, maxContentLen)
		}

		// 过滤掉空文本的索引，记录有效文本的映射
		var validTexts []string
		var validIndices []int
		// emptyTextSet 记录空文本文章的索引，用于区分"空文本保留"和"计算失败保留"
		emptyTextSet := make(map[int]bool)
		for i, text := range texts {
			if len(strings.TrimSpace(text)) > 0 {
				validTexts = append(validTexts, text)
				validIndices = append(validIndices, i)
			} else {
				emptyTextSet[i] = true
			}
		}

		// 批量计算文章向量
		articleVectors := make([][]float64, len(items)) // nil 表示未计算
		var embedErr error
		if len(validTexts) > 0 {
			vectors, batchErr := adapter.EmbedTexts(ctx, validTexts, instruction)
			if batchErr != nil {
				embedErr = batchErr
				logrus.Errorf("[embedding-filter] batch embedding failed: %v", batchErr)
			} else {
				// 检查返回的向量数量是否与请求数量一致
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

		// 4. 如果批量 Embedding 全部失败，返回错误
		if embedErr != nil {
			return fmt.Errorf("[embedding-filter] all article embeddings failed: %w", embedErr)
		}

		// 5. 根据相似度过滤
		totalCount := len(items)
		feed.Items = lo.Filter(items, func(item *feeds.Item, index int) bool {
			vec := articleVectors[index]

			// 向量为 nil 时区分原因
			if vec == nil {
				if emptyTextSet[index] {
					// 空文本文章：保留（行为不变）
					return true
				}
				// 有效文本但向量计算失败：保留但记录警告
				logrus.Warnf("[embedding-filter] article [%s] has valid text but nil vector (embedding failed), keeping by default", item.Title)
				return true
			}

			// 与所有锚点计算相似度，任一超过阈值即保留
			maxSim := -1.0
			for _, anchorVec := range anchorVectors {
				sim := util.CosineSimilarity(vec, anchorVec)
				if sim > maxSim {
					maxSim = sim
				}
			}

			matched := maxSim >= threshold
			if matched {
				logrus.Debugf("[embedding-filter] article [%s] MATCHED (max similarity: %.4f >= %.4f)", item.Title, maxSim, threshold)
			} else {
				logrus.Debugf("[embedding-filter] article [%s] DROPPED (max similarity: %.4f < %.4f)", item.Title, maxSim, threshold)
			}
			return matched
		})

		keptCount := len(feed.Items)
		droppedCount := totalCount - keptCount
		logrus.Infof("[embedding-filter] filtering complete: %d total, %d kept, %d dropped (threshold: %.4f)", totalCount, keptCount, droppedCount, threshold)

		return nil
	}
}

// buildArticleText 拼接文章标题和正文，正文超过 maxLen 时按 Unicode 字符截取
func buildArticleText(item *feeds.Item, maxLen int) string {
	content := item.Content
	if len(content) == 0 {
		content = item.Description
	}

	// 按 Unicode 字符（rune）截取正文，避免切断多字节字符
	if utf8.RuneCountInString(content) > maxLen {
		runes := []rune(content)
		content = string(runes[:maxLen])
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
