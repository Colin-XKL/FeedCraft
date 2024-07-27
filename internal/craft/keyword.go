package craft

import (
	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
	"strings"
)

type KeywordFilterMode string

var (
	KeywordIncludeMode KeywordFilterMode = "include"
	KeywordExcludeMode KeywordFilterMode = "exclude"
)

type KeywordMatchScope string

var (
	KeywordMatchTitle   KeywordMatchScope = "title"
	KeywordMatchContent KeywordMatchScope = "content"
	KeywordMatchAll     KeywordMatchScope = "all"
)

func optionKeyword(mode KeywordFilterMode, matchScope KeywordMatchScope, keywordList []string) CraftOption {
	searchTitle := matchScope == KeywordMatchAll || matchScope == KeywordMatchTitle
	searchContent := matchScope == KeywordMatchAll || matchScope == KeywordMatchContent

	return func(feed *feeds.Feed) error {
		items := feed.Items
		var filtered []*feeds.Item
		if len(keywordList) == 0 {
			logrus.Warnf("empty keyword list")
			return nil
		}
		for _, feedItem := range items {
			matched := false
			for _, keyword := range keywordList {
				if searchTitle {
					if strings.Contains(feedItem.Title, keyword) {
						matched = true
						break
					}
				}
				if searchContent {
					// todo handle content correctly
					if strings.Contains(feedItem.Content, keyword) || strings.Contains(feedItem.Description, keyword) {
						matched = true
						break
					}
				}
			}
			switch mode {
			case KeywordIncludeMode:
				if matched {
					filtered = append(filtered, feedItem)
				}
				break
			case KeywordExcludeMode:
				if !matched {
					filtered = append(filtered, feedItem)
				}
				break
			default:
				logrus.Warnf("unknown mode %s", mode)
			}
		}
		feed.Items = filtered
		return nil
	}
}

func GetKeywordOption(mode KeywordFilterMode, matchScope KeywordMatchScope, keywordList []string) []CraftOption {
	craftOptions := []CraftOption{
		optionKeyword(mode, matchScope, keywordList),
	}
	return craftOptions
}

var keywordCraftParamTmpl = []ParamTemplate{
	{Key: "mode", Description: "`include` or `exclude`"},
	{Key: "keywords", Description: "keywords that need to match, seperated by comma, example `ad,sell,SALE`"},
	{Key: "scope", Description: "match scope, `title` or `content` or `all`"},
}

func keywordCraftLoadParams(m map[string]string) []CraftOption {
	modeStr, _ := m["mode"]
	mode := KeywordIncludeMode
	if modeStr == string(KeywordIncludeMode) {
		mode = KeywordIncludeMode
	} else if modeStr == string(KeywordExcludeMode) {
		mode = KeywordExcludeMode
	} else {
		logrus.Warnf("unknown mode str %s", modeStr)
	}
	scopeStr, _ := m["scope"]
	scope := KeywordMatchAll
	if scopeStr == string(KeywordMatchTitle) {
		scope = KeywordMatchTitle
	} else if scopeStr == string(KeywordMatchContent) {
		scope = KeywordMatchTitle
	} else if scopeStr == string(KeywordMatchAll) {
		scope = KeywordMatchAll
	} else {
		logrus.Warnf("unknown scope str %s", scopeStr)
	}
	keywordStr, _ := m["keywords"]
	keywordList := strings.Split(keywordStr, ",")
	return GetKeywordOption(mode, scope, keywordList)
}
