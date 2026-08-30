package craft

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"FeedCraft/internal/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckConditionWithLLM_ShortContentSkipsLLM(t *testing.T) {
	called := false
	original := llmContextCaller
	llmContextCaller = func(_ context.Context, prompt, context string, option util.ContentProcessOption) (string, error) {
		called = true
		return "true", nil
	}
	t.Cleanup(func() { llmContextCaller = original })

	result, err := CheckConditionWithLLM(context.Background(), "title", "short", "is this spam?")
	require.NoError(t, err)
	assert.False(t, result)
	assert.False(t, called, "LLM should not be called for content shorter than the minimum length")
}

func TestCheckConditionWithLLM_ParsesTrueFalseResponses(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     bool
	}{
		{name: "lowercase true", response: "true", want: true},
		{name: "uppercase with spaces", response: "  TRUE \n", want: true},
		{name: "lowercase false", response: "false", want: false},
		{name: "mixed case false", response: "False", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			original := llmContextCaller
			llmContextCaller = func(_ context.Context, prompt, context string, option util.ContentProcessOption) (string, error) {
				return tc.response, nil
			}
			t.Cleanup(func() { llmContextCaller = original })

			result, err := CheckConditionWithLLM(context.Background(), "title", strings.Repeat("content ", 10), "is this spam?")
			require.NoError(t, err)
			assert.Equal(t, tc.want, result)
		})
	}
}

func TestCheckConditionWithLLM_UnexpectedResponseReturnsError(t *testing.T) {
	original := llmContextCaller
	llmContextCaller = func(_ context.Context, prompt, context string, option util.ContentProcessOption) (string, error) {
		return "maybe", nil
	}
	t.Cleanup(func() { llmContextCaller = original })

	result, err := CheckConditionWithLLM(context.Background(), "title", strings.Repeat("content ", 10), "is this spam?")
	require.Error(t, err)
	assert.False(t, result)
	assert.Contains(t, err.Error(), "unexpected llm response")
}

func TestCheckConditionWithLLM_PropagatesLLMError(t *testing.T) {
	original := llmContextCaller
	llmContextCaller = func(_ context.Context, prompt, context string, option util.ContentProcessOption) (string, error) {
		return "", fmt.Errorf("upstream timeout")
	}
	t.Cleanup(func() { llmContextCaller = original })

	result, err := CheckConditionWithLLM(context.Background(), "title", strings.Repeat("content ", 10), "is this spam?")
	require.Error(t, err)
	assert.False(t, result)
	assert.Contains(t, err.Error(), "upstream timeout")
}

func TestCheckConditionWithGenericPrompt_WrapsUserCriterion(t *testing.T) {
	var capturedPrompt string
	original := llmContextCaller
	llmContextCaller = func(_ context.Context, prompt, context string, option util.ContentProcessOption) (string, error) {
		capturedPrompt = prompt
		return "true", nil
	}
	t.Cleanup(func() { llmContextCaller = original })

	result, err := CheckConditionWithGenericPrompt(context.Background(), "title", strings.Repeat("content ", 10), "Is this about sports?")
	require.NoError(t, err)
	assert.True(t, result)
	assert.Contains(t, capturedPrompt, "Is this about sports?")
	assert.Contains(t, capturedPrompt, "return 'true'")
	assert.Contains(t, capturedPrompt, "return 'false'")
}

func TestBuildLLMArticlePayload_TitleOnlyAndContentOnly(t *testing.T) {
	titleOnly := BuildLLMArticlePayload("Just A Title", "")
	assert.Contains(t, titleOnly, "Article Title:")
	assert.Contains(t, titleOnly, "Just A Title")
	assert.NotContains(t, titleOnly, "Article Content:")

	contentOnly := BuildLLMArticlePayload("", "Body only content")
	assert.Contains(t, contentOnly, "Article Content:")
	assert.Contains(t, contentOnly, "Body only content")
	assert.NotContains(t, contentOnly, "Article Title:")
}
