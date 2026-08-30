package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTruncateHeadTail_ShortInputUnchanged(t *testing.T) {
	input := "hello 世界"
	require.Equal(t, input, TruncateHeadTail(input, 100))
}

func TestTruncateHeadTail_KeepsHeadAndTail(t *testing.T) {
	input := strings.Repeat("H", 40) + strings.Repeat("M", 40) + strings.Repeat("T", 40)
	got := TruncateHeadTail(input, 50)

	require.LessOrEqual(t, len([]rune(got)), 50)
	require.True(t, strings.HasPrefix(got, "HHHH"))
	require.True(t, strings.HasSuffix(got, "TTTT"))
	require.Contains(t, got, "[...truncated...]")
	require.NotContains(t, got, "MMMM")
}

func TestLLMPromptMaxChars_DefaultAndEnv(t *testing.T) {
	t.Setenv("FC_LLM_PROMPT_MAX_CHARS", "")
	require.Equal(t, defaultLLMPromptMaxChars, LLMPromptMaxChars())

	t.Setenv("FC_LLM_PROMPT_MAX_CHARS", "8000")
	require.Equal(t, 8000, LLMPromptMaxChars())
}
