package feedruntime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/craft"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/engine"
	"FeedCraft/internal/model"
	"FeedCraft/internal/observability"
	"FeedCraft/internal/source"
	"FeedCraft/internal/util"

	"gorm.io/gorm"
)

type InputKind string

const (
	InputKindURI    InputKind = "uri"
	InputKindSource InputKind = "source"
)

const (
	internalScheme             = "feedcraft"
	internalResourceTypeRecipe = "recipe"
	internalResourceTypeTopic  = "topic"
	internalResourceTypeInbox  = "inbox"
)

// InputSpec is the unified runtime input model for RecipeFeed and TopicFeed.
type InputSpec struct {
	Kind         InputKind            `json:"kind"`
	URI          string               `json:"uri,omitempty"`
	SourceConfig *config.SourceConfig `json:"source_config,omitempty"`
}

type Builder struct {
	DB *gorm.DB
}

func NewBuilder(db *gorm.DB) *Builder {
	return &Builder{DB: db}
}

type baseURLProvider interface {
	engine.FeedProvider
	BaseURL() string
}

func BuildProviderFromInput(ctx context.Context, spec InputSpec, stack []string) (engine.FeedProvider, error) {
	return NewBuilder(nil).BuildProviderFromInput(ctx, spec, stack)
}

func BuildTopicProvider(ctx context.Context, topicID string) (engine.FeedProvider, error) {
	return NewBuilder(nil).BuildTopicProvider(ctx, topicID)
}

func BuildTopic(ctx context.Context, topic *dao.TopicFeed, stack []string) (*engine.TopicFeed, error) {
	return NewBuilder(nil).BuildTopic(ctx, topic, stack)
}

func BuildRecipeProvider(ctx context.Context, recipeID string) (engine.FeedProvider, error) {
	return NewBuilder(nil).BuildRecipeProvider(ctx, recipeID)
}

func BuildRecipe(ctx context.Context, recipeData *dao.CustomRecipeV2) (*engine.RecipeFeed, error) {
	return NewBuilder(nil).BuildRecipe(ctx, recipeData)
}

func BuildAggregator(steps []dao.AggregatorStep) (engine.CraftOption, error) {
	return buildAggregator(steps)
}

func (b *Builder) BuildProviderFromInput(ctx context.Context, spec InputSpec, stack []string) (engine.FeedProvider, error) {
	switch spec.Kind {
	case InputKindURI:
		return b.buildProviderFromURI(ctx, spec.URI, stack)
	case InputKindSource:
		if spec.SourceConfig == nil {
			return nil, errors.New("source input requires source_config")
		}
		return newSourceConfigProvider(spec.SourceConfig)
	default:
		return nil, fmt.Errorf("unsupported input kind %q", spec.Kind)
	}
}

func (b *Builder) BuildTopicProvider(ctx context.Context, topicID string) (engine.FeedProvider, error) {
	return b.buildTopicProvider(ctx, topicID, nil)
}

func (b *Builder) BuildRecipeProvider(ctx context.Context, recipeID string) (engine.FeedProvider, error) {
	return b.buildRecipeProvider(ctx, recipeID, observability.TriggerUserRequest)
}

func (b *Builder) buildRecipeProvider(ctx context.Context, recipeID string, trigger string) (engine.FeedProvider, error) {
	recipeData, err := dao.GetCustomRecipeByIDV2(b.db(), recipeID)
	if err != nil {
		return nil, err
	}

	recipeRuntime, err := b.BuildRecipe(ctx, recipeData)
	if err != nil {
		return nil, err
	}

	return &RecipeProvider{
		Recipe:  recipeRuntime,
		Trigger: trigger,
	}, nil
}

func (b *Builder) buildTopicProvider(ctx context.Context, topicID string, stack []string) (engine.FeedProvider, error) {
	topic, err := dao.GetTopicFeedByID(b.db(), topicID)
	if err != nil {
		return nil, err
	}
	return b.BuildTopic(ctx, topic, stack)
}

func (b *Builder) BuildTopic(ctx context.Context, topic *dao.TopicFeed, stack []string) (*engine.TopicFeed, error) {
	if topic == nil {
		return nil, errors.New("topic is nil")
	}

	stack, err := pushTopicStack(stack, topic.ID)
	if err != nil {
		return nil, err
	}

	inputs := make([]engine.FeedProvider, 0, len(topic.InputURIs))
	for _, inputURI := range topic.InputURIs {
		spec := InputSpec{
			Kind: InputKindURI,
			URI:  inputURI,
		}
		provider, buildErr := b.BuildProviderFromInput(ctx, spec, stack)
		if buildErr != nil {
			return nil, fmt.Errorf("build topic input %q: %w", inputURI, buildErr)
		}
		inputs = append(inputs, provider)
	}

	aggregator, err := buildAggregator(topic.AggregatorConfig)
	if err != nil {
		return nil, fmt.Errorf("build topic aggregator: %w", err)
	}

	return &engine.TopicFeed{
		ID:          topic.ID,
		Title:       topic.Title,
		Description: topic.Description,
		Inputs:      inputs,
		Aggregator:  aggregator,
	}, nil
}

func (b *Builder) BuildRecipe(ctx context.Context, recipeData *dao.CustomRecipeV2) (*engine.RecipeFeed, error) {
	if recipeData == nil {
		return nil, errors.New("recipe is nil")
	}

	inputSpec, err := buildRecipeInputSpec(recipeData)
	if err != nil {
		return nil, err
	}
	provider, err := b.BuildProviderFromInput(ctx, inputSpec, nil)
	if err != nil {
		return nil, err
	}

	inputProvider, ok := provider.(baseURLProvider)
	if !ok {
		return nil, fmt.Errorf("recipe input provider does not expose base url: %T", provider)
	}

	craftChain, err := craft.BuildOptionChain(b.db(), recipeData.Craft, inputProvider.BaseURL())
	if err != nil {
		return nil, err
	}

	return &engine.RecipeFeed{
		ID:          recipeData.ID,
		Description: recipeData.Description,
		SourceType:  recipeData.SourceType,
		BaseURL:     inputProvider.BaseURL(),
		CraftName:   recipeData.Craft,
		Input:       inputProvider,
		Craft:       craftChain,
	}, nil
}

func buildRecipeInputSpec(recipeData *dao.CustomRecipeV2) (InputSpec, error) {
	var sourceConfig config.SourceConfig
	if err := jsonUnmarshalSourceConfig(recipeData, &sourceConfig); err != nil {
		return InputSpec{}, err
	}

	return InputSpec{
		Kind:         InputKindSource,
		SourceConfig: &sourceConfig,
	}, nil
}

func (b *Builder) buildProviderFromURI(ctx context.Context, rawURI string, stack []string) (engine.FeedProvider, error) {
	if rawURI == "" {
		return nil, errors.New("uri input requires a non-empty uri")
	}

	parsed, err := url.Parse(rawURI)
	if err != nil {
		return nil, fmt.Errorf("invalid uri %q: %w", rawURI, err)
	}

	switch parsed.Scheme {
	case "http", "https":
		return &RawFeedProvider{URL: rawURI}, nil
	case internalScheme:
		resourceType, resourceID, err := parseInternalResourceURI(parsed)
		if err != nil {
			return nil, fmt.Errorf("invalid internal uri %q: %w", rawURI, err)
		}
		switch resourceType {
		case internalResourceTypeRecipe:
			return b.buildRecipeProvider(ctx, resourceID, observability.TriggerTopicAggregation)
		case internalResourceTypeTopic:
			return b.buildTopicProvider(ctx, resourceID, stack)
		case internalResourceTypeInbox:
			return &InboxProvider{DB: b.db(), InboxID: resourceID}, nil
		default:
			return nil, fmt.Errorf("unsupported internal resource type %q", resourceType)
		}
	default:
		return nil, fmt.Errorf("unsupported uri scheme %q", parsed.Scheme)
	}
}

func jsonUnmarshalSourceConfig(recipeData *dao.CustomRecipeV2, sourceConfig *config.SourceConfig) error {
	if recipeData == nil {
		return errors.New("recipe is nil")
	}
	if sourceConfig == nil {
		return errors.New("source config target is nil")
	}
	if err := json.Unmarshal([]byte(recipeData.SourceConfig), sourceConfig); err != nil {
		return fmt.Errorf("invalid source config: %w", err)
	}
	if sourceConfig.Type == "" && recipeData.SourceType != "" {
		sourceConfig.Type = constant.SourceType(recipeData.SourceType)
	}
	return nil
}

func (b *Builder) db() *gorm.DB {
	if b.DB != nil {
		return b.DB
	}
	return util.GetDatabase()
}

func buildAggregator(steps []dao.AggregatorStep) (engine.CraftOption, error) {
	if len(steps) == 0 {
		return nil, nil
	}

	options := make([]engine.CraftOption, 0, len(steps))
	for idx, step := range steps {
		option, err := buildAggregatorStep(idx, step)
		if err != nil {
			return nil, err
		}
		if option != nil {
			options = append(options, option)
		}
	}
	return composeAggregatorOptions(options...), nil
}

func buildAggregatorStep(index int, step dao.AggregatorStep) (engine.CraftOption, error) {
	stepType := strings.ToLower(strings.TrimSpace(step.Type))
	switch stepType {
	case "deduplicate":
		strategy := strings.ToLower(strings.TrimSpace(step.Option["strategy"]))
		if strategy == "" {
			strategy = "by_link"
		}
		if strategy != "by_link" && strategy != "by_id" {
			return nil, fmt.Errorf("aggregator step %d (%s): invalid strategy %q", index, step.Type, strategy)
		}
		return func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
			_ = ctx
			if feed == nil || len(feed.Articles) == 0 {
				return feed, nil
			}
			cloned := cloneFeedArticles(feed)
			seen := make(map[string]bool)
			uniqueArticles := make([]*model.CraftArticle, 0, len(cloned.Articles))
			for _, article := range cloned.Articles {
				if article == nil {
					continue
				}
				key := article.Link
				if strategy == "by_id" {
					key = article.Id
				}
				if key == "" {
					uniqueArticles = append(uniqueArticles, article)
					continue
				}
				if !seen[key] {
					seen[key] = true
					uniqueArticles = append(uniqueArticles, article)
				}
			}
			cloned.Articles = uniqueArticles
			return cloned, nil
		}, nil
	case "sort":
		sortBy := strings.ToLower(strings.TrimSpace(step.Option["by"]))
		if sortBy == "" {
			sortBy = "date_desc"
		}
		switch sortBy {
		case "date_desc", "date_asc", "quality_desc", "quality_asc":
			return func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
				_ = ctx
				if feed == nil || len(feed.Articles) <= 1 {
					return feed, nil
				}
				cloned := cloneFeedArticles(feed)
				sort.SliceStable(cloned.Articles, func(i, j int) bool {
					a, b := cloned.Articles[i], cloned.Articles[j]
					switch sortBy {
					case "date_asc":
						return a.Updated.Before(b.Updated)
					case "quality_desc":
						return a.QualityScore > b.QualityScore
					case "quality_asc":
						return a.QualityScore < b.QualityScore
					default:
						return a.Updated.After(b.Updated)
					}
				})
				return cloned, nil
			}, nil
		default:
			return nil, fmt.Errorf("aggregator step %d (%s): invalid sort mode %q", index, step.Type, sortBy)
		}
	case "limit":
		rawMax := strings.TrimSpace(step.Option["max"])
		if rawMax == "" {
			return nil, fmt.Errorf("aggregator step %d (%s): option max is required", index, step.Type)
		}
		maxItems, err := strconv.Atoi(rawMax)
		if err != nil || maxItems <= 0 {
			return nil, fmt.Errorf("aggregator step %d (%s): invalid max %q", index, step.Type, rawMax)
		}
		return func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
			_ = ctx
			if feed == nil || len(feed.Articles) == 0 || len(feed.Articles) <= maxItems {
				return feed, nil
			}
			cloned := cloneFeedArticles(feed)
			cloned.Articles = cloned.Articles[:maxItems]
			return cloned, nil
		}, nil
	default:
		return nil, fmt.Errorf("aggregator step %d: unsupported type %q", index, step.Type)
	}
}

func composeAggregatorOptions(options ...engine.CraftOption) engine.CraftOption {
	return func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
		currentFeed := feed
		var err error
		for _, option := range options {
			if option == nil {
				continue
			}
			currentFeed, err = option(ctx, currentFeed)
			if err != nil {
				return nil, err
			}
		}
		return currentFeed, nil
	}
}

func cloneFeedArticles(feed *model.CraftFeed) *model.CraftFeed {
	if feed == nil {
		return nil
	}
	cloned := *feed
	cloned.Articles = make([]*model.CraftArticle, 0, len(feed.Articles))
	for _, article := range feed.Articles {
		if article == nil {
			cloned.Articles = append(cloned.Articles, nil)
			continue
		}
		articleCopy := *article
		cloned.Articles = append(cloned.Articles, &articleCopy)
	}
	return &cloned
}

func pushTopicStack(stack []string, topicID string) ([]string, error) {
	for idx, existing := range stack {
		if existing == topicID {
			cycle := append(append([]string{}, stack[idx:]...), topicID)
			return nil, fmt.Errorf("topic dependency cycle detected: %s", strings.Join(cycle, " -> "))
		}
	}
	next := append([]string{}, stack...)
	next = append(next, topicID)
	return next, nil
}

func parseInternalResourceURI(parsed *url.URL) (string, string, error) {
	resourceType := strings.TrimSpace(parsed.Host)
	resourceID := strings.Trim(strings.TrimSpace(parsed.Path), "/")
	if resourceType == "" {
		return "", "", errors.New("missing resource type")
	}
	if resourceID == "" {
		return "", "", errors.New("missing resource id")
	}
	if strings.Contains(resourceID, "/") {
		return "", "", errors.New("resource id must be a single path segment")
	}
	return resourceType, resourceID, nil
}

// InboxProvider adapts inbox items into a feed without requiring an intermediate recipe.
type InboxProvider struct {
	DB      *gorm.DB
	InboxID string
}

func (p *InboxProvider) Fetch(ctx context.Context) (*model.CraftFeed, error) {
	_ = ctx
	db := p.DB
	if db == nil {
		db = util.GetDatabase()
	}

	inbox, err := dao.GetInboxByID(db, p.InboxID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inbox %s: %w", p.InboxID, err)
	}

	items, err := dao.ListInboxItems(db, p.InboxID)
	if err != nil {
		return nil, fmt.Errorf("failed to list inbox items: %w", err)
	}

	feed := &model.CraftFeed{
		Title:       inbox.Title,
		Description: inbox.Description,
		Id:          fmt.Sprintf("inbox:%s", p.InboxID),
		Link:        p.BaseURL(),
		Articles:    make([]*model.CraftArticle, 0, len(items)),
	}
	for _, item := range items {
		feed.Articles = append(feed.Articles, &model.CraftArticle{
			Title:       item.Title,
			Link:        item.URL,
			Content:     item.Content,
			Description: item.Summary,
			Id:          item.ItemID,
			AuthorName:  item.Author,
			Created:     item.PublishedAt,
			Updated:     item.PublishedAt,
		})
	}

	return feed, nil
}

func (p *InboxProvider) BaseURL() string {
	return fmt.Sprintf("feedcraft://inbox/%s", p.InboxID)
}

// RecipeProvider adapts a runtime RecipeFeed and adds execution metadata.
type RecipeProvider struct {
	Recipe  *engine.RecipeFeed
	Trigger string
}

func (p *RecipeProvider) Fetch(ctx context.Context) (*model.CraftFeed, error) {
	if p.Recipe == nil {
		return nil, errors.New("recipe runtime is nil")
	}

	startedAt := time.Now()
	feed, err := p.Recipe.Fetch(ctx)
	if err != nil {
		reportRecipeRuntimeFailure(ctx, p.Recipe, p.Trigger, startedAt, err)
		return nil, err
	}

	observability.Report(observability.ExecutionEvent{
		ResourceType: dao.ResourceTypeRecipe,
		ResourceID:   p.Recipe.ID,
		ResourceName: p.Recipe.ID,
		Trigger:      p.triggerOrDefault(),
		Status:       dao.ExecutionStatusSuccess,
		Message:      fmt.Sprintf("recipe executed successfully with %d items", len(feed.Articles)),
		Details: map[string]any{
			"source_type": p.Recipe.SourceType,
			"base_url":    p.Recipe.BaseURL,
			"item_count":  len(feed.Articles),
		},
		RequestID: observability.RequestIDFromContext(ctx),
		Duration:  time.Since(startedAt),
	})

	return feed, nil
}

func (p *RecipeProvider) triggerOrDefault() string {
	if strings.TrimSpace(p.Trigger) != "" {
		return p.Trigger
	}
	return observability.TriggerUserRequest
}

func reportRecipeRuntimeFailure(ctx context.Context, recipeRuntime *engine.RecipeFeed, trigger string, startedAt time.Time, err error) {
	if recipeRuntime == nil {
		return
	}
	observability.Report(observability.ExecutionEvent{
		ResourceType: dao.ResourceTypeRecipe,
		ResourceID:   recipeRuntime.ID,
		ResourceName: recipeRuntime.ID,
		Trigger:      trigger,
		Status:       dao.ExecutionStatusFailure,
		ErrorKind:    observability.ClassifyError(err),
		Message:      err.Error(),
		Details: map[string]any{
			"source_type": recipeRuntime.SourceType,
			"base_url":    recipeRuntime.BaseURL,
			"craft":       recipeRuntime.CraftName,
		},
		RequestID: observability.RequestIDFromContext(ctx),
		Duration:  time.Since(startedAt),
	})
}

// RawFeedProvider fetches a third-party URL using the minimal raw-feed semantics.
type RawFeedProvider struct {
	URL string
}

func (p *RawFeedProvider) Fetch(ctx context.Context) (*model.CraftFeed, error) {
	sourceConfig := &config.SourceConfig{
		Type: constant.SourceRSS,
		HttpFetcher: &config.HttpFetcherConfig{
			URL: p.URL,
		},
	}
	provider, err := newSourceConfigProvider(sourceConfig)
	if err != nil {
		return nil, err
	}
	return provider.Fetch(ctx)
}

// SourceConfigProvider adapts a full SourceConfig into the FeedProvider interface.
type SourceConfigProvider struct {
	SourceConfig *config.SourceConfig
	Source       source.Source
}

func newSourceConfigProvider(sourceConfig *config.SourceConfig) (*SourceConfigProvider, error) {
	if sourceConfig == nil {
		return nil, errors.New("source config is nil")
	}

	factory, err := source.Get(sourceConfig.Type)
	if err != nil {
		return nil, err
	}

	src, err := factory(sourceConfig)
	if err != nil {
		return nil, err
	}

	return &SourceConfigProvider{
		SourceConfig: sourceConfig,
		Source:       src,
	}, nil
}

func (p *SourceConfigProvider) BaseURL() string {
	if p == nil || p.Source == nil {
		return ""
	}
	return p.Source.BaseURL()
}

func (p *SourceConfigProvider) Fetch(ctx context.Context) (*model.CraftFeed, error) {
	if p == nil || p.Source == nil {
		return nil, errors.New("source provider is not initialized")
	}
	return p.Source.Fetch(ctx)
}
