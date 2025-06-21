package craft

import (
	"FeedCraft/internal/adapter"
	"FeedCraft/internal/util"
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
)

// TextProcessor defines interface for processing text content
type TextProcessor interface {
	Process(original string) (string, error)
	Combine(original, processed string) string
}

// LLMTextProcessor implements TextProcessor using LLM
type LLMTextProcessor struct {
	prompt string
}

func (p *LLMTextProcessor) Process(original string) (string, error) {
	return adapter.CallLLMUsingContext(p.prompt, original)
}

func (p *LLMTextProcessor) Combine(original, processed string) string {
	processedHTML := util.Markdown2HTML(processed)
	return fmt.Sprintf(`<div>%s<div><hr/><br/>%s</div>`, processedHTML, original)
}

// ProcessorType defines supported processing types
type ProcessorType string

const (
	ProcessorTypeIntroduction ProcessorType = "introduction"
	// Add more processor types here
)

// DefaultPrompts contains default prompts for different processing types
var DefaultPrompts = map[ProcessorType]string{
	ProcessorTypeIntroduction: `
你是一位专业的文章导读撰写专家...`, // Keep original prompt content
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
	
	var processed string
	var err error
	
	if len(cleanedContent) > 0 {
		processed, err = processor.Process(cleanedContent)
	} else {
		processed, err = processor.Process(originalContent)
	}
	
	if err != nil {
		errMsg := fmt.Sprintf("process article content failed. err: %v", err)
		logrus.Warnf(errMsg)
		processed = errMsg
	}

	return processor.Combine(originalContent, processed)
}

func GetTextProcessingCraftOptions(processor TextProcessor, cacheKeyPrefix string) []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		ret := processItemContent(item, processor)
		return ret, nil
	}
	
	cachedTransformer := GetCommonCachedTransformer(
		func(item *feeds.Item) (string, error) {
			return cacheKeyPrefix + util.GetMD5Hash(item.Title+item.Link.Href), nil
		}, 
		transFunc, 
		"text processing")
		
	return []CraftOption{
		OptionTransformFeedItem(GetArticleContentProcessor(cachedTransformer)),
	}
}

func NewLLMTextProcessor(processorType ProcessorType, customPrompt string) TextProcessor {
	prompt := customPrompt
	if prompt == "" {
		prompt = DefaultPrompts[processorType]
	}
	return &LLMTextProcessor{prompt: prompt}
}

func textProcessingCraftLoadParam(processorType ProcessorType, m map[string]string) []CraftOption {
	prompt := m["prompt"]
	processor := NewLLMTextProcessor(processorType, prompt)
	return GetTextProcessingCraftOptions(processor, string(processorType))
}

func GetIntroductionCraftOptions() []CraftOption {
	return textProcessingCraftLoadParam(ProcessorTypeIntroduction, map[string]string{})
}

var introCraftParamTmpl = []ParamTemplate{
	{Key: "prompt", Description: "the llm prompt for generate summary", Default: DefaultPrompts[ProcessorTypeIntroduction]},
}
