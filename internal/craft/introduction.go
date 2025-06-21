package craft

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/util"
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
)

// TextProcessor defines interface for processing text content
type TextProcessor interface {
	Process(original string) (string, error)
	Combine(original, processed string) string
	GetName() string
}

// LLMTextProcessor 调用LLM处理文章内容,然后把处理结果放到文章开头
type LLMTextProcessor struct {
	prompt string
	name   string
}

func (p *LLMTextProcessor) Process(original string) (string, error) {
	return adapter.CallLLMUsingContext(p.prompt, original)
}
func (p *LLMTextProcessor) Combine(original, processed string) string {
	processedHTML := util.Markdown2HTML(processed)
	return fmt.Sprintf(`<div><div>%s</div><hr/><br/><div>%s</div></div>`, processedHTML, original)
}
func (p *LLMTextProcessor) GetName() string {
	return p.name
}

func processItemContent(item *feeds.Item, processor TextProcessor) string {
	originalContent := item.Content
	originalTitle := item.Title

	if len(originalContent) == 0 {
		if len(item.Description) == 0 {
			logrus.Warnf("empty content, both content and description fields are empty. title [%s]", originalTitle)
		}
		originalContent = item.Description
		logrus.Warnf("empty content, using description field as fallback")
	}

	domain, _ := util.ParseDomainFromUrl(item.Link.Href)
	cleanedContent := util.Html2Markdown(originalContent, &domain)

	var processedContent string
	var err error

	if len(cleanedContent) > 0 {
		processedContent, err = processor.Process(cleanedContent)
	} else {
		processedContent, err = processor.Process(originalContent)
	}

	if err != nil {
		errMsg := fmt.Sprintf("process article content using processsor [%s] failed. err: %v", processor.GetName(), err)
		logrus.Warnf(errMsg)
		processedContent = errMsg
	}

	return processor.Combine(originalContent, processedContent)
}

func NewLLMTextProcessor(processorType constant.ProcessorType, customPrompt string) TextProcessor {
	prompt := customPrompt
	if prompt == "" {
		prompt = constant.DefaultPrompts[processorType]
	}
	return &LLMTextProcessor{prompt: prompt, name: string(processorType)}
}

func GetAddIntroductionCraftOptions(prompt string) []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		processorType := constant.ProcessorTypeIntroduction
		processor := NewLLMTextProcessor(processorType, prompt)
		ret := processItemContent(item, processor)
		return ret, nil
	}
	cachedTransformer := GetCommonCachedTransformer(cacheKeyForArticleTitle, transFunc, string(constant.ProcessorTypeIntroduction))
	craftOption := []CraftOption{
		OptionTransformFeedItem(GetArticleContentProcessor(cachedTransformer)),
	}
	return craftOption
}

func introCraftLoadParam(m map[string]string) []CraftOption {
	prompt, exist := m["prompt"]
	if !exist || len(prompt) == 0 {
		prompt = constant.DefaultPrompts[constant.ProcessorTypeIntroduction]
	}
	return GetAddIntroductionCraftOptions(prompt)
}

var introCraftParamTmpl = []ParamTemplate{
	{
		Key:         "prompt",
		Description: "the llm prompt for generate introduction",
		Default:     constant.DefaultPrompts[constant.ProcessorTypeIntroduction],
	},
}
