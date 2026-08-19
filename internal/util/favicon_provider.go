package util

import (
	"net/url"
	"strconv"
	"strings"
)

type FaviconProvider string

const (
	FaviconProviderGstaticCN  FaviconProvider = "gstatic_cn"
	FaviconProviderGoogle     FaviconProvider = "google"
	FaviconProviderDuckDuckGo FaviconProvider = "duckduckgo"
	FaviconProviderYandex     FaviconProvider = "yandex"
)

// DefaultFaviconProvider is reachable from mainland China.
const DefaultFaviconProvider = FaviconProviderGstaticCN

const defaultFaviconSize = 64

type faviconURLBuilder func(pageURL string) string

var faviconProviders = map[FaviconProvider]faviconURLBuilder{
	FaviconProviderGstaticCN:  buildGstaticCNFaviconURL,
	FaviconProviderGoogle:     buildGoogleFaviconURL,
	FaviconProviderDuckDuckGo: buildDuckDuckGoFaviconURL,
	FaviconProviderYandex:     buildYandexFaviconURL,
}

// BuildFaviconURL returns a favicon URL using the default provider.
func BuildFaviconURL(pageURL string) string {
	return BuildFaviconURLWithProvider(DefaultFaviconProvider, pageURL)
}

// BuildFaviconURLWithProvider returns a favicon URL from the given provider.
// Unknown providers fall back to DefaultFaviconProvider.
func BuildFaviconURLWithProvider(provider FaviconProvider, pageURL string) string {
	builder, ok := faviconProviders[provider]
	if !ok {
		builder = faviconProviders[DefaultFaviconProvider]
	}
	if builder == nil {
		return ""
	}
	return builder(pageURL)
}

func buildGstaticCNFaviconURL(pageURL string) string {
	origin := firstOrigin(pageURL)
	if origin == "" {
		return ""
	}
	return "https://t0.gstatic.cn/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=" +
		url.QueryEscape(origin) + "&size=" + strconv.Itoa(defaultFaviconSize)
}

func buildGoogleFaviconURL(pageURL string) string {
	origin := firstOrigin(pageURL)
	if origin == "" {
		return ""
	}
	return "https://www.google.com/s2/favicons?domain_url=" + url.QueryEscape(origin) + "&sz=" + strconv.Itoa(defaultFaviconSize)
}

func buildDuckDuckGoFaviconURL(pageURL string) string {
	host := hostname(pageURL)
	if host == "" {
		return ""
	}
	return "https://icons.duckduckgo.com/ip3/" + host + ".ico"
}

func buildYandexFaviconURL(pageURL string) string {
	host := hostname(pageURL)
	if host == "" {
		return ""
	}
	return "https://favicon.yandex.net/favicon/" + host
}

func OriginFromURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return ""
	}
	return parsedURL.Scheme + "://" + parsedURL.Host
}

func firstOrigin(rawURLs ...string) string {
	for _, rawURL := range rawURLs {
		if origin := OriginFromURL(rawURL); origin != "" {
			return origin
		}
	}
	return ""
}

func hostname(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(parsedURL.Hostname())
}
