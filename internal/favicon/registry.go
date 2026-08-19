package favicon

import (
	"FeedCraft/internal/config"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

const (
	ProviderGstaticCN  = "gstatic_cn"
	ProviderGoogle     = "google"
	ProviderDuckDuckGo = "duckduckgo"
	ProviderYandex     = "yandex"

	defaultIconSize = 64
	maxTemplateSize = 2048
	defaultProvider = ProviderGstaticCN
	sampleTargetURL = "https://example.com/path?q=1"
	customIDPattern = `^[a-z][a-z0-9_-]{0,63}$`
)

type ProviderDescriptor struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URLTemplate string `json:"url_template"`
	Enabled     bool   `json:"enabled"`
	BuiltIn     bool   `json:"built_in"`
}

type providerBuilder func(pageURL string, size int) string

type provider struct {
	descriptor ProviderDescriptor
	build      providerBuilder
}

type snapshot struct {
	settings  config.FaviconSettings
	defaultID string
	providers map[string]provider
}

var (
	activeSnapshot  atomic.Pointer[snapshot]
	providerIDRE    = regexp.MustCompile(customIDPattern)
	placeholderRE   = regexp.MustCompile(`\{[^{}]+\}`)
	warnedFallbacks sync.Map
)

var builtInProviders = []provider{
	{
		descriptor: ProviderDescriptor{
			ID:          ProviderGstaticCN,
			Name:        "Gstatic China",
			URLTemplate: "https://t0.gstatic.cn/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url={origin_query}&size={size}",
			Enabled:     true,
			BuiltIn:     true,
		},
		build: buildGstaticCNURL,
	},
	{
		descriptor: ProviderDescriptor{
			ID:          ProviderGoogle,
			Name:        "Google",
			URLTemplate: "https://www.google.com/s2/favicons?domain_url={origin_query}&sz={size}",
			Enabled:     true,
			BuiltIn:     true,
		},
		build: buildGoogleURL,
	},
	{
		descriptor: ProviderDescriptor{
			ID:          ProviderDuckDuckGo,
			Name:        "DuckDuckGo",
			URLTemplate: "https://icons.duckduckgo.com/ip3/{host}.ico",
			Enabled:     true,
			BuiltIn:     true,
		},
		build: buildDuckDuckGoURL,
	},
	{
		descriptor: ProviderDescriptor{
			ID:          ProviderYandex,
			Name:        "Yandex",
			URLTemplate: "https://favicon.yandex.net/favicon/{host}",
			Enabled:     true,
			BuiltIn:     true,
		},
		build: buildYandexURL,
	},
}

func init() {
	compiled, err := compileSnapshot(DefaultSettings())
	if err != nil {
		panic(err)
	}
	activeSnapshot.Store(compiled)
}

func DefaultSettings() config.FaviconSettings {
	return config.FaviconSettings{DefaultProviderID: defaultProvider}
}

// Replace validates and compiles the complete registry before atomically
// publishing it. Readers always observe either the old or the new snapshot.
func Replace(settings config.FaviconSettings) error {
	compiled, err := compileSnapshot(settings)
	if err != nil {
		return err
	}
	activeSnapshot.Store(compiled)
	warnedFallbacks = sync.Map{}
	return nil
}

func Settings() config.FaviconSettings {
	current := currentSnapshot()
	settings := current.settings
	settings.CustomProviders = append([]config.FaviconProviderConfig(nil), current.settings.CustomProviders...)
	return settings
}

func Providers() []ProviderDescriptor {
	current := currentSnapshot()
	result := make([]ProviderDescriptor, 0, len(current.providers)+len(current.settings.CustomProviders))
	for _, item := range current.providers {
		result = append(result, item.descriptor)
	}
	for _, item := range current.settings.CustomProviders {
		if item.Enabled {
			continue
		}
		result = append(result, ProviderDescriptor{
			ID:          item.ID,
			Name:        item.Name,
			URLTemplate: item.URLTemplate,
			Enabled:     false,
			BuiltIn:     false,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].BuiltIn != result[j].BuiltIn {
			return result[i].BuiltIn
		}
		return result[i].ID < result[j].ID
	})
	return result
}

// BuildURL uses providerID when it resolves to an enabled provider. An empty,
// missing, or disabled provider falls back to the configured default and then
// to the built-in mainland-China provider.
func BuildURL(providerID string, pageURL string, size int) (string, string) {
	return buildURL(currentSnapshot(), providerID, pageURL, size)
}

// BuildURLFromSettings validates and previews a not-yet-published settings
// payload without mutating the active registry.
func BuildURLFromSettings(settings config.FaviconSettings, providerID string, pageURL string, size int) (string, string, error) {
	compiled, err := compileSnapshot(settings)
	if err != nil {
		return "", "", err
	}
	result, resolvedID := buildURL(compiled, providerID, pageURL, size)
	return result, resolvedID, nil
}

func buildURL(current *snapshot, providerID string, pageURL string, size int) (string, string) {
	resolvedID := strings.TrimSpace(providerID)
	item, ok := current.providers[resolvedID]
	if resolvedID == "" || !ok {
		if resolvedID != "" && resolvedID != current.defaultID {
			warnProviderFallback(resolvedID, current.defaultID)
		}
		resolvedID = current.defaultID
		item, ok = current.providers[resolvedID]
	}
	if !ok {
		resolvedID = defaultProvider
		item = current.providers[resolvedID]
	}
	if size <= 0 {
		size = defaultIconSize
	}
	return item.build(pageURL, size), resolvedID
}

func OriginFromURL(rawURL string) string {
	parsedURL := parseTargetURL(rawURL)
	if parsedURL == nil {
		return ""
	}
	return parsedURL.Scheme + "://" + parsedURL.Host
}

func currentSnapshot() *snapshot {
	current := activeSnapshot.Load()
	if current != nil {
		return current
	}
	compiled, err := compileSnapshot(DefaultSettings())
	if err != nil {
		panic(err)
	}
	activeSnapshot.Store(compiled)
	return compiled
}

func compileSnapshot(settings config.FaviconSettings) (*snapshot, error) {
	settings.DefaultProviderID = strings.TrimSpace(settings.DefaultProviderID)
	if settings.DefaultProviderID == "" {
		settings.DefaultProviderID = defaultProvider
	}
	settings.CustomProviders = append([]config.FaviconProviderConfig(nil), settings.CustomProviders...)

	providers := make(map[string]provider, len(builtInProviders)+len(settings.CustomProviders))
	for _, item := range builtInProviders {
		providers[item.descriptor.ID] = item
	}

	seen := make(map[string]struct{}, len(settings.CustomProviders))
	for index := range settings.CustomProviders {
		item := &settings.CustomProviders[index]
		item.ID = strings.TrimSpace(item.ID)
		item.Name = strings.TrimSpace(item.Name)
		item.URLTemplate = strings.TrimSpace(item.URLTemplate)
		if err := validateCustomProvider(*item, providers, seen); err != nil {
			return nil, err
		}
		seen[item.ID] = struct{}{}
		if !item.Enabled {
			continue
		}
		builder, err := compileTemplate(item.URLTemplate)
		if err != nil {
			return nil, fmt.Errorf("favicon provider %q: %w", item.ID, err)
		}
		providers[item.ID] = provider{
			descriptor: ProviderDescriptor{
				ID:          item.ID,
				Name:        item.Name,
				URLTemplate: item.URLTemplate,
				Enabled:     true,
				BuiltIn:     false,
			},
			build: builder,
		}
	}

	if _, ok := providers[settings.DefaultProviderID]; !ok {
		return nil, fmt.Errorf("default favicon provider %q is missing or disabled", settings.DefaultProviderID)
	}

	return &snapshot{
		settings:  settings,
		defaultID: settings.DefaultProviderID,
		providers: providers,
	}, nil
}

func validateCustomProvider(item config.FaviconProviderConfig, builtIns map[string]provider, seen map[string]struct{}) error {
	if !providerIDRE.MatchString(item.ID) {
		return fmt.Errorf("invalid favicon provider ID %q", item.ID)
	}
	if _, reserved := builtIns[item.ID]; reserved {
		return fmt.Errorf("favicon provider ID %q is reserved", item.ID)
	}
	if _, duplicate := seen[item.ID]; duplicate {
		return fmt.Errorf("duplicate favicon provider ID %q", item.ID)
	}
	if item.Name == "" {
		return fmt.Errorf("favicon provider %q name is required", item.ID)
	}
	if item.URLTemplate == "" {
		return fmt.Errorf("favicon provider %q URL template is required", item.ID)
	}
	if len(item.URLTemplate) > maxTemplateSize {
		return fmt.Errorf("favicon provider %q URL template exceeds %d characters", item.ID, maxTemplateSize)
	}
	_, err := compileTemplate(item.URLTemplate)
	if err != nil {
		return fmt.Errorf("favicon provider %q: %w", item.ID, err)
	}
	return nil
}

func compileTemplate(template string) (providerBuilder, error) {
	allowed := map[string]bool{
		"{host}":         true,
		"{origin}":       true,
		"{origin_query}": true,
		"{url_query}":    true,
		"{size}":         true,
	}
	for _, placeholder := range placeholderRE.FindAllString(template, -1) {
		if !allowed[placeholder] {
			return nil, fmt.Errorf("unsupported placeholder %q", placeholder)
		}
	}
	if !strings.Contains(template, "{host}") &&
		!strings.Contains(template, "{origin}") &&
		!strings.Contains(template, "{origin_query}") &&
		!strings.Contains(template, "{url_query}") {
		return nil, fmt.Errorf("URL template must include a target placeholder")
	}

	builder := func(pageURL string, size int) string {
		parsedURL := parseTargetURL(pageURL)
		if parsedURL == nil {
			return ""
		}
		origin := parsedURL.Scheme + "://" + parsedURL.Host
		return strings.NewReplacer(
			"{host}", parsedURL.Hostname(),
			"{origin}", origin,
			"{origin_query}", url.QueryEscape(origin),
			"{url_query}", url.QueryEscape(parsedURL.String()),
			"{size}", strconv.Itoa(size),
		).Replace(template)
	}

	sample := builder(sampleTargetURL, defaultIconSize)
	parsedTemplateURL, err := url.Parse(sample)
	if err != nil || parsedTemplateURL.Scheme != "https" || parsedTemplateURL.Host == "" || parsedTemplateURL.User != nil {
		return nil, fmt.Errorf("URL template must render an absolute HTTPS URL without credentials")
	}
	return builder, nil
}

func warnProviderFallback(providerID string, defaultID string) {
	key := providerID + "\x00" + defaultID
	if _, loaded := warnedFallbacks.LoadOrStore(key, struct{}{}); loaded {
		return
	}
	logrus.Warnf("favicon provider %q is missing or disabled; falling back to %q", providerID, defaultID)
}

func buildGstaticCNURL(pageURL string, size int) string {
	origin := OriginFromURL(pageURL)
	if origin == "" {
		return ""
	}
	return "https://t0.gstatic.cn/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=" +
		url.QueryEscape(origin) + "&size=" + strconv.Itoa(size)
}

func buildGoogleURL(pageURL string, size int) string {
	origin := OriginFromURL(pageURL)
	if origin == "" {
		return ""
	}
	return "https://www.google.com/s2/favicons?domain_url=" + url.QueryEscape(origin) + "&sz=" + strconv.Itoa(size)
}

func buildDuckDuckGoURL(pageURL string, _ int) string {
	host := hostname(pageURL)
	if host == "" {
		return ""
	}
	return "https://icons.duckduckgo.com/ip3/" + host + ".ico"
}

func buildYandexURL(pageURL string, _ int) string {
	host := hostname(pageURL)
	if host == "" {
		return ""
	}
	return "https://favicon.yandex.net/favicon/" + host
}

func hostname(rawURL string) string {
	parsedURL := parseTargetURL(rawURL)
	if parsedURL == nil {
		return ""
	}
	return strings.TrimSpace(parsedURL.Hostname())
}

func parseTargetURL(rawURL string) *url.URL {
	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL.Host == "" || parsedURL.User != nil {
		return nil
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil
	}
	return parsedURL
}
