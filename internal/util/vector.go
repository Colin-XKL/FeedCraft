package util

import "math"

// CosineSimilarity 计算两个向量的余弦相似度
// 返回值范围 [-1, 1]，1 表示完全相同，0 表示正交，-1 表示完全相反
// 边界情况：零向量、维度不匹配等返回 0.0
func CosineSimilarity(a, b []float64) float64 {
	if len(a) == 0 || len(b) == 0 || len(a) != len(b) {
		return 0.0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	// 零向量检查
	if normA == 0 || normB == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
