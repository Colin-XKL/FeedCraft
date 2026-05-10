package craft

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"FeedCraft/internal/constant"
	"FeedCraft/internal/model"
	"FeedCraft/internal/util"

	"github.com/sirupsen/logrus"
)

type ArticleMutationFunc func(ctx context.Context, article *model.CraftArticle) error
type ArticlePredicateFunc func(ctx context.Context, article *model.CraftArticle) (bool, error)

type ArticleTextTransformProcessor struct {
	CraftName string
	Mutate    ArticleMutationFunc
}

func (p *ArticleTextTransformProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	if feed == nil || p == nil || p.Mutate == nil || len(feed.Articles) == 0 {
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
		if err := p.Mutate(ctx, article); err != nil {
			lastErr = err
			logrus.Warnf("failed to apply craft [%s] for article [%s], err: %v", p.CraftName, article.Title, err)
			continue
		}
		successes += 1
	}

	if attempted > 0 && successes == 0 {
		return nil, fmt.Errorf("all items failed to process. last error: %v", lastErr)
	}
	return cloned, nil
}

type ArticlePredicateProcessor struct {
	CraftName string
	Match     ArticlePredicateFunc
}

func (p *ArticlePredicateProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	if feed == nil || p == nil || p.Match == nil || len(feed.Articles) == 0 {
		return feed, nil
	}

	cloned := cloneCraftFeed(feed)
	filtered := make([]*model.CraftArticle, 0, len(cloned.Articles))
	for _, article := range cloned.Articles {
		if article == nil {
			continue
		}
		matched, err := p.Match(ctx, article)
		if err != nil {
			logrus.Warnf("failed to evaluate craft [%s] for article [%s], err: %v", p.CraftName, article.Title, err)
			filtered = append(filtered, article)
			continue
		}
		if !matched {
			filtered = append(filtered, article)
		}
	}
	cloned.Articles = filtered
	return cloned, nil
}

func CallLLMForArticleTransform(prompt, title, content string, option util.ContentProcessOption) (string, error) {
	contextData := BuildLLMArticlePayload(title, content)
	return llmContextCaller(prompt, contextData, option)
}

func CallLLMForArticlePredicate(prompt, title, content string) (bool, error) {
	return CheckConditionWithLLM(title, content, prompt)
}

func newSummaryProcessor(prompt string) *ArticleTextTransformProcessor {
	finalPrompt := renderTargetLangPrompt(prompt, constant.DefaultPrompts[constant.ProcessorTypeSummary])
	transformer := GetCommonCachedArticleTransformer(
		newArticleTitleContentCacheKeyGenerator(finalPrompt),
		func(ctx context.Context, article *model.CraftArticle) (string, error) {
			original := getPrimaryArticleContent(article)
			if strings.TrimSpace(original) == "" {
				return "", nil
			}
			processed := getArticleContentForPrompt(article, original)
			generated, err := CallLLMForArticleTransform(finalPrompt, article.Title, processed, util.ContentProcessOption{})
			if err != nil {
				return "", err
			}
			return combineArticleHTMLWithGeneratedMarkdown(original, generated), nil
		},
		string(constant.ProcessorTypeSummary),
	)

	return &ArticleTextTransformProcessor{
		CraftName: string(constant.ProcessorTypeSummary),
		Mutate: func(ctx context.Context, article *model.CraftArticle) error {
			original := getPrimaryArticleContent(article)
			if strings.TrimSpace(original) == "" {
				return nil
			}
			transformed, err := transformer(ctx, article)
			if err != nil {
				return err
			}
			article.Content = transformed
			article.Description = transformed
			return nil
		},
	}
}

func newIntroductionProcessor(prompt string) *ArticleTextTransformProcessor {
	finalPrompt := renderTargetLangPrompt(prompt, constant.DefaultPrompts[constant.ProcessorTypeIntroduction])
	transformer := GetCommonCachedArticleTransformer(
		newArticleTitleContentCacheKeyGenerator(finalPrompt),
		func(ctx context.Context, article *model.CraftArticle) (string, error) {
			original := getPrimaryArticleContent(article)
			if strings.TrimSpace(original) == "" {
				return "", nil
			}
			processed := getArticleContentForPrompt(article, original)
			generated, err := CallLLMForArticleTransform(finalPrompt, article.Title, processed, util.ContentProcessOption{})
			if err != nil {
				return "", err
			}
			return combineArticleHTMLWithGeneratedMarkdown(original, generated), nil
		},
		string(constant.ProcessorTypeIntroduction),
	)

	return &ArticleTextTransformProcessor{
		CraftName: string(constant.ProcessorTypeIntroduction),
		Mutate: func(ctx context.Context, article *model.CraftArticle) error {
			original := getPrimaryArticleContent(article)
			if strings.TrimSpace(original) == "" {
				return nil
			}
			transformed, err := transformer(ctx, article)
			if err != nil {
				return err
			}
			article.Content = transformed
			article.Description = transformed
			return nil
		},
	}
}

func newTranslateTitleProcessor(prompt string) *ArticleTextTransformProcessor {
	finalPrompt := renderTargetLangPrompt(prompt, translateArticleTitlePrompt)
	targetLangCode := util.GetDefaultTargetLang()
	transformer := GetCommonCachedArticleTransformer(
		newArticleTitleContentCacheKeyGenerator(finalPrompt),
		func(ctx context.Context, article *model.CraftArticle) (string, error) {
			title := strings.TrimSpace(article.Title)
			if title == "" || util.IsSameLanguage(title, targetLangCode) {
				return title, nil
			}
			return CallLLMForArticleTransform(finalPrompt, "", title, util.ContentProcessOption{})
		},
		"translate title",
	)

	return &ArticleTextTransformProcessor{
		CraftName: "translate title",
		Mutate: func(ctx context.Context, article *model.CraftArticle) error {
			transformed, err := transformer(ctx, article)
			if err != nil {
				return err
			}
			article.Title = transformed
			return nil
		},
	}
}

func newTranslateContentProcessor(prompt string) *ArticleTextTransformProcessor {
	return newArticleContentLLMProcessor("translate article content", prompt, translateArticleContentPrompt)
}

func newTranslateContentImmersiveProcessor(prompt string) *ArticleTextTransformProcessor {
	return newArticleContentLLMProcessor("translate article content immersive", prompt, immersiveTranslatePrompt)
}

func newArticleContentLLMProcessor(craftName, prompt, defaultPrompt string) *ArticleTextTransformProcessor {
	finalPrompt := renderTargetLangPrompt(prompt, defaultPrompt)
	targetLangCode := util.GetDefaultTargetLang()
	transformer := GetCommonCachedArticleTransformer(
		newArticleTitleContentCacheKeyGenerator(finalPrompt),
		func(ctx context.Context, article *model.CraftArticle) (string, error) {
			content := getPrimaryArticleContent(article)
			if strings.TrimSpace(content) == "" || util.IsSameLanguage(content, targetLangCode) {
				return content, nil
			}
			return CallLLMForArticleTransform(finalPrompt, "", content, util.ContentProcessOption{})
		},
		craftName,
	)

	return &ArticleTextTransformProcessor{
		CraftName: craftName,
		Mutate: func(ctx context.Context, article *model.CraftArticle) error {
			transformed, err := transformer(ctx, article)
			if err != nil {
				return err
			}
			article.Content = transformed
			article.Description = transformed
			return nil
		},
	}
}

func newBeautifyContentProcessor(prompt string) *ArticleTextTransformProcessor {
	finalPrompt := strings.TrimSpace(prompt)
	if finalPrompt == "" {
		finalPrompt = beautifyArticleContentPrompt
	}
	transformer := GetCommonCachedArticleTransformer(
		newArticleTitleContentCacheKeyGenerator(finalPrompt),
		func(ctx context.Context, article *model.CraftArticle) (string, error) {
			content := getPrimaryArticleContent(article)
			if strings.TrimSpace(content) == "" {
				return "", nil
			}
			return beautifyArticleContent(content, finalPrompt)
		},
		"beautify article content",
	)

	return &ArticleTextTransformProcessor{
		CraftName: "beautify article content",
		Mutate: func(ctx context.Context, article *model.CraftArticle) error {
			transformed, err := transformer(ctx, article)
			if err != nil {
				return err
			}
			article.Content = transformed
			article.Description = transformed
			return nil
		},
	}
}

func newLLMFilterProcessor(condition string) *ArticlePredicateProcessor {
	condition = strings.TrimSpace(condition)
	if condition == "" {
		condition = "Is this content spam or low quality?"
	}
	matcher := GetCommonCachedArticlePredicate(
		newArticleTitleContentCacheKeyGenerator(condition),
		func(ctx context.Context, article *model.CraftArticle) (bool, error) {
			content := getPrimaryArticleContent(article)
			return CheckConditionWithGenericPrompt(article.Title, content, condition)
		},
		"llm filter",
	)
	return &ArticlePredicateProcessor{
		CraftName: "llm filter",
		Match:     matcher,
	}
}

func newIgnoreAdvertorialProcessor(prompt string) *ArticlePredicateProcessor {
	finalPrompt := strings.TrimSpace(prompt)
	if finalPrompt == "" {
		finalPrompt = fmt.Sprintf("%s\n%s\n", judgePrompt, promptCheckIfAdvertorial)
	}
	matcher := GetCommonCachedArticlePredicate(
		newArticleTitleContentCacheKeyGenerator(finalPrompt),
		func(ctx context.Context, article *model.CraftArticle) (bool, error) {
			content := getPrimaryArticleContent(article)
			return CallLLMForArticlePredicate(finalPrompt, article.Title, content)
		},
		"ignore advertorial",
	)
	return &ArticlePredicateProcessor{
		CraftName: "ignore advertorial",
		Match:     matcher,
	}
}

func GetCommonCachedArticlePredicate(cacheKeyGenerator ArticleCacheKeyGenerator, rawPredicate ArticlePredicateFunc, craftName string) ArticlePredicateFunc {
	return func(ctx context.Context, article *model.CraftArticle) (bool, error) {
		hashVal, err := cacheKeyGenerator(article)
		if err != nil {
			return false, err
		}
		cacheKey := getCraftCacheKey(craftName, hashVal)
		value, err := util.CachedFunc(cacheKey, func() (string, error) {
			matched, callErr := rawPredicate(ctx, article)
			if callErr != nil {
				return "", callErr
			}
			if matched {
				return "true", nil
			}
			return "false", nil
		})
		if err != nil {
			return false, err
		}
		return strings.EqualFold(strings.TrimSpace(value), "true"), nil
	}
}

func newArticleTitleContentCacheKeyGenerator(prompt string) ArticleCacheKeyGenerator {
	promptHash := util.GetTextContentHash(prompt)
	return func(article *model.CraftArticle) (string, error) {
		payloadHash := util.GetTextContentHash(strings.Join([]string{
			promptHash,
			strings.TrimSpace(article.Title),
			strings.TrimSpace(getPrimaryArticleContent(article)),
		}, "|"))
		return payloadHash, nil
	}
}

func renderTargetLangPrompt(prompt string, defaultPrompt string) string {
	finalPrompt := strings.TrimSpace(prompt)
	if finalPrompt == "" {
		finalPrompt = defaultPrompt
	}

	tmpl, err := template.New("prompt").Parse(finalPrompt)
	if err != nil {
		logrus.Debugf("failed to parse llm prompt template: %v", err)
		return finalPrompt
	}

	data := struct {
		TargetLang string
	}{
		TargetLang: util.GetLanguageName(util.GetDefaultTargetLang()),
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		logrus.Debugf("failed to execute llm prompt template: %v", err)
		return finalPrompt
	}
	return buf.String()
}

func getPrimaryArticleContent(article *model.CraftArticle) string {
	if article == nil {
		return ""
	}
	content := article.Content
	if strings.TrimSpace(content) == "" {
		content = article.Description
	}
	return content
}

func getArticleContentForPrompt(article *model.CraftArticle, original string) string {
	domain, _ := util.ParseDomainFromUrl(article.Link)
	cleaned := util.Html2Markdown(original, &domain)
	if strings.TrimSpace(cleaned) != "" {
		return cleaned
	}
	return original
}

func combineArticleHTMLWithGeneratedMarkdown(originalHTML string, generatedMarkdown string) string {
	processedHTML := util.Markdown2HTML(generatedMarkdown)
	return fmt.Sprintf(`<div><div>%s</div><hr/><br/><div>%s</div></div>`, processedHTML, originalHTML)
}
