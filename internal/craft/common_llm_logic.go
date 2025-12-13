package craft

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/util"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// CheckConditionWithLLM 检查文章内容是否符合特定条件
// prompt: 用户提供的 prompt 模板，必须明确要求返回 'true' 或 'false'
// content: 文章内容
// 返回值: (bool, error)
// true: 表示符合条件
// false: 表示不符合条件或无法判断
func CheckConditionWithLLM(content string, conditionPrompt string) (bool, error) {
	const MinContentLength = 20
	if len(strings.TrimSpace(content)) < MinContentLength {
		return false, nil
	}
	option := util.ContentProcessOption{
		RemoveImage: true,
		ConvertToMd: true,
	}

	// 如果 prompt 没有明确包含 true/false 的指令，建议调用者在传入前拼接好。
	// 但为了通用性，我们这里假设传入的 prompt 已经包含了判断逻辑。

	result, err := adapter.CallLLMUsingContext(conditionPrompt, content, option)
	if err != nil {
		logrus.Errorf("Error checking condition with LLM: %v", err)
		return false, err
	}

	logrus.Infof("LLM check result: [%s] for prompt prefix: [%.20s...]", result, conditionPrompt)

	// 简单的处理: 只要包含 true 就算 true (不区分大小写)
	// 或者严格匹配 "true"
	// 原有 ignore-advertorial 逻辑是: strings.TrimSpace(result) == "true"

	trimmedResult := strings.ToLower(strings.TrimSpace(result))
	return trimmedResult == "true", nil
}

// CheckConditionWithGenericPrompt 使用通用模板构造 prompt 并调用 LLM
// userPrompt: 用户提供的判断标准，例如 "Is this regarding politics?"
func CheckConditionWithGenericPrompt(content string, userPrompt string) (bool, error) {
	// 构造强制要求 true/false 的完整 prompt
	fullPrompt := fmt.Sprintf(`
Evaluate the following content based on this criterion:
"%s"

If the content matches the criterion, return 'true'.
Otherwise return 'false'.
Do not include any other text.
`, userPrompt)

	return CheckConditionWithLLM(content, fullPrompt)
}
