package feedruntime

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/engine"
	"FeedCraft/internal/model"
	"FeedCraft/internal/recipe"
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
)

// InputSpec is the unified runtime input model for RecipeFeed and TopicFeed.
type InputSpec struct {
	Kind         InputKind            `json:"kind"`
	URI          string               `json:"uri,omitempty"`
	SourceConfig *config.SourceConfig `json:"source_config,omitempty"`
}

// Builder compiles persisted feed configs into executable runtime providers.
type Builder struct {
	DB *gorm.DB
}

func NewBuilder(db *gorm.DB) *Builder {
	return &Builder{DB: db}
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

func BuildAggregator(steps []dao.AggregatorStep) (engine.FeedProcessor, error) {
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
		return &SourceConfigProvider{SourceConfig: spec.SourceConfig}, nil
	default:
		return nil, fmt.Errorf("unsupported input kind %q", spec.Kind)
	}
}

func (b *Builder) BuildTopicProvider(ctx context.Context, topicID string) (engine.FeedProvider, error) {
	return b.buildTopicProvider(ctx, topicID, nil)
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
			if _, err := dao.GetCustomRecipeByIDV2(b.db(), resourceID); err != nil {
				return nil, err
			}
			return &RecipeProvider{RecipeID: resourceID}, nil
		case internalResourceTypeTopic:
			return b.buildTopicProvider(ctx, resourceID, stack)
		default:
			return nil, fmt.Errorf("unsupported internal resource type %q", resourceType)
		}
	default:
		return nil, fmt.Errorf("unsupported uri scheme %q", parsed.Scheme)
	}
}

func (b *Builder) db() *gorm.DB {
	if b.DB != nil {
		return b.DB
	}
	return util.GetDatabase()
}

func buildAggregator(steps []dao.AggregatorStep) (engine.FeedProcessor, error) {
	if len(steps) == 0 {
		return nil, nil
	}

	processors := make([]engine.FeedProcessor, 0, len(steps))
	for idx, step := range steps {
		processor, err := buildAggregatorStep(idx, step)
		if err != nil {
			return nil, err
		}
		processors = append(processors, processor)
	}

	return &engine.FlowCraftProcessor{Processors: processors}, nil
}

func buildAggregatorStep(index int, step dao.AggregatorStep) (engine.FeedProcessor, error) {
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
		return &engine.DeduplicateProcessor{Strategy: strategy}, nil
	case "sort":
		sortBy := strings.ToLower(strings.TrimSpace(step.Option["by"]))
		if sortBy == "" {
			sortBy = "date_desc"
		}
		switch sortBy {
		case "date_desc", "date_asc", "quality_desc", "quality_asc":
			return &engine.SortProcessor{SortBy: sortBy}, nil
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
		return &engine.LimitProcessor{MaxItems: maxItems}, nil
	default:
		return nil, fmt.Errorf("aggregator step %d: unsupported type %q", index, step.Type)
	}
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

// RecipeProvider adapts an existing RecipeFeed into the FeedProvider interface.
type RecipeProvider struct {
	RecipeID string
}

func (p *RecipeProvider) Fetch(ctx context.Context) (*model.CraftFeed, error) {
	feed, err := recipe.ProcessRecipeByID(ctx, p.RecipeID)
	if err != nil {
		return nil, err
	}
	return model.FromFeedsFeed(feed), nil
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
	return (&SourceConfigProvider{SourceConfig: sourceConfig}).Fetch(ctx)
}

// SourceConfigProvider adapts a full SourceConfig into the FeedProvider interface.
type SourceConfigProvider struct {
	SourceConfig *config.SourceConfig
}

func (p *SourceConfigProvider) Fetch(ctx context.Context) (*model.CraftFeed, error) {
	if p.SourceConfig == nil {
		return nil, errors.New("source config is nil")
	}

	factory, err := source.Get(p.SourceConfig.Type)
	if err != nil {
		return nil, err
	}

	src, err := factory(p.SourceConfig)
	if err != nil {
		return nil, err
	}

	return (&source.LegacySourceAdapter{LegacySource: src}).Fetch(ctx)
}
