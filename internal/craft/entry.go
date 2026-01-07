package craft

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

func GetSysCraftTemplateDict() map[string]CraftTemplate {
	sysCraftTempList := make(map[string]CraftTemplate)
	sysCraftTempList["proxy"] = CraftTemplate{
		Name:                "proxy",
		Description:         "代理订阅源",
		ParamTemplateDefine: []ParamTemplate{},
		OptionFunc: func(m map[string]string) []CraftOption {
			return []CraftOption{}
		},
	}
	sysCraftTempList["limit"] = CraftTemplate{
		Name:                "limit",
		Description:         "限制单页条目数量",
		ParamTemplateDefine: limitCraftParamTmpl,
		OptionFunc:          limitCraftLoadParams,
	}
	sysCraftTempList["time-limit"] = CraftTemplate{
		Name:                "time-limit",
		Description:         "根据时间限制文章保留天数",
		ParamTemplateDefine: timeLimitCraftParamTmpl,
		OptionFunc:          timeLimitCraftLoadParams,
	}
	sysCraftTempList["keyword"] = CraftTemplate{
		Name:                "keyword",
		Description:         "关键词过滤",
		ParamTemplateDefine: keywordCraftParamTmpl,
		OptionFunc:          keywordCraftLoadParams,
	}
	sysCraftTempList["guid-fix"] = CraftTemplate{
		Name:                "guid-fix",
		Description:         "修复 RSS GUID。使用文章内容 MD5 作为唯一 ID。",
		ParamTemplateDefine: []ParamTemplate{},
		OptionFunc: func(m map[string]string) []CraftOption {
			return GetGuidCraftOptions()
		},
	}
	sysCraftTempList["relative-link-fix"] = CraftTemplate{
		Name:                "relative-link-fix",
		Description:         "修复文章链接,确保是绝对url. 这样可以保证在获取全文等场景时可以跳转到正确的网页",
		ParamTemplateDefine: []ParamTemplate{},
		OptionFunc: func(m map[string]string) []CraftOption {
			return GetRelativeLinkFixCraftOptions()
		},
	}
	sysCraftTempList["fulltext"] = CraftTemplate{
		Name:                "fulltext",
		Description:         "提取 RSS 订阅源的全文",
		ParamTemplateDefine: []ParamTemplate{},
		OptionFunc: func(m map[string]string) []CraftOption {
			return GetFulltextCraftOptions()
		},
	}
	sysCraftTempList["fulltext-plus"] = CraftTemplate{
		Name:                "fulltext-plus",
		Description:         "模拟浏览器提取 RSS 订阅源的全文",
		ParamTemplateDefine: []ParamTemplate{},
		OptionFunc: func(m map[string]string) []CraftOption {
			return GetFulltextPlusCraftOptions()
		},
	}
	sysCraftTempList["cleanup"] = CraftTemplate{
		Name:                "cleanup",
		Description:         "清理文章HTML内容，保留核心内容",
		ParamTemplateDefine: []ParamTemplate{},
		OptionFunc: func(m map[string]string) []CraftOption {
			return GetCleanupCraftOptions()
		},
	}
	sysCraftTempList["introduction"] = CraftTemplate{
		Name:                "introduction",
		Description:         "在文章开头添加 AI 生成的导读",
		ParamTemplateDefine: introCraftParamTmpl,
		OptionFunc:          introCraftLoadParam,
	}
	sysCraftTempList["summary"] = CraftTemplate{
		Name:                "summary",
		Description:         "让AI总结文章主要内容,并附在原文开头",
		ParamTemplateDefine: summaryCraftParamTmpl,
		OptionFunc:          summaryCraftLoadParam,
	}
	sysCraftTempList["ignore-advertorial"] = CraftTemplate{
		Name:                "ignore-advertorial",
		Description:         "使用 LLM 排除广告文章",
		ParamTemplateDefine: llmFilterCraftParamTmpl,
		OptionFunc:          llmFilterCraftLoadParam,
	}
	sysCraftTempList["llm-filter"] = CraftTemplate{
		Name:                "llm-filter",
		Description:         "使用 LLM 根据自定义条件过滤文章 (如果满足条件则排除)",
		ParamTemplateDefine: llmFilterGenericParamTmpl,
		OptionFunc:          llmFilterGenericLoadParam,
	}
	sysCraftTempList["translate-title"] = CraftTemplate{
		Name:                "translate-title",
		Description:         "使用 LLM 将标题翻译为中文",
		ParamTemplateDefine: transTitleParamTmpl,
		OptionFunc:          transTitleCraftLoadParam,
	}
	sysCraftTempList["translate-content"] = CraftTemplate{
		Name:                "translate-content",
		Description:         "使用 LLM 将文章内容翻译为中文, 替换原文内容, 只输出翻译后文章",
		ParamTemplateDefine: transContentParamTmpl,
		OptionFunc:          transContentCraftLoadParam,
	}
	sysCraftTempList["translate-content-immersive"] = CraftTemplate{
		Name:                "translate-content-immersive",
		Description:         "使用 LLM 将文章内容翻译为中文, 沉浸式翻译模式, 每个原文段落后面添加翻译后内容",
		ParamTemplateDefine: immersiveTranslateParamTmpl,
		OptionFunc:          immersiveTranslateLoadParam,
	}
	sysCraftTempList["beautify-content"] = CraftTemplate{
		Name:                "beautify-content",
		Description:         "使用 LLM 美化文章排版，去除广告和无关信息",
		ParamTemplateDefine: beautifyContentParamTmpl,
		OptionFunc:          beautifyContentCraftLoadParam,
	}
	return sysCraftTempList
}

func GetToolDict() map[string]dao.Tool {
	tmplDict := GetSysCraftTemplateDict()
	toolDict := make(map[string]dao.Tool)
	for name, craftTemplate := range tmplDict {
		item := dao.Tool{
			Name:         craftTemplate.Name, // 默认会有个跟template 同名的 tool
			Description:  craftTemplate.Description,
			TemplateName: craftTemplate.Name,
			Params:       map[string]string{},
		}
		toolDict[name] = item
	}

	db := util.GetDatabase()
	toolList, err := dao.GetAllTools(db)
	if err != nil {
		logrus.Errorf("read tool list from db error. only built-in tool will work now. err: %s", err)
	} else {
		for _, tool := range toolList {
			toolDict[tool.Name] = tool
		}
	}

	return toolDict
}

func Entry(c *gin.Context) {
	processorName := c.Param("craft-name") // Keep param name for now (Phase 2)
	db := util.GetDatabase()

	craftOptionList, err := getCraftOptions(db, processorName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	CommonCraftHandlerUsingCraftOptionList(c, craftOptionList)
}

func ProcessFeed(feed *gofeed.Feed, feedURL string, processorName string) (*feeds.Feed, error) {
	db := util.GetDatabase()
	craftOptionList, err := getCraftOptions(db, processorName)
	if err != nil {
		return nil, err
	}

	craftedFeed, err := NewCraftedFeedFromGofeed(feed, feedURL, craftOptionList...)
	if err != nil {
		return nil, err
	}

	return craftedFeed.OutputFeed, nil
}

func getCraftOptions(db *gorm.DB, processorName string) ([]CraftOption, error) {
	toolDict := GetToolDict()
	craftTmplDict := GetSysCraftTemplateDict()

	if strings.Contains(processorName, ",") {
		var allOptions []CraftOption
		parts := strings.Split(processorName, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			options, err := inner(db, &toolDict, &craftTmplDict, part, 0)
			if err != nil {
				return nil, err
			}
			allOptions = append(allOptions, options...)
		}
		return allOptions, nil
	}

	return inner(db, &toolDict, &craftTmplDict, processorName, 0)
}

const MaxCallDepth = 5

// 递归地解出 craft option list
func inner(db *gorm.DB, toolDict *map[string]dao.Tool, craftTmplDict *map[string]CraftTemplate, processorName string, depthId int) ([]CraftOption, error) {
	if depthId+1 > MaxCallDepth {
		return []CraftOption{}, fmt.Errorf("max call depth hit")
	}
	logrus.Infof("checking %s", processorName)

	// 对于每一个 processor name in flow
	// 1. 判断是否为 built-in tool (formerly craft atom)
	// 1.1 如果是, 直接返回对应的 craft option
	// 2. 如果不是, 判断是否为 blueprint (formerly flow)
	// 2.1 如果是 blueprint, 调用 entry check
	// 如果不是, 说明 processor name 无效, 返回error

	tool, isKnownTool := (*toolDict)[processorName]
	if isKnownTool {
		logrus.Infof("[%s] is known tool", processorName)
		tmplContent, tmplValid := (*craftTmplDict)[tool.TemplateName]
		if !tmplValid {
			return []CraftOption{}, fmt.Errorf("invalid tmpl name [%s] for tool [%s]", tool.TemplateName, tool.Name)
		}
		return tmplContent.GetOptions(tool.Params), nil
	} else {
		processorArr, checkErr := extractProcessorArrFromBlueprint(db, processorName)
		if checkErr != nil {
			// then not a valid processor name
			return []CraftOption{}, fmt.Errorf("not a valid processor name")
		}
		var retArr []CraftOption
		for _, extractedSubProcessorName := range processorArr {
			sub, recurErr := inner(db, toolDict, craftTmplDict, extractedSubProcessorName, depthId+1)
			if recurErr != nil {
				return []CraftOption{}, recurErr
			}
			retArr = append(retArr, sub...)
		}
		return retArr, nil
	}
}

// 给定 blueprintName, 查询其详情, 获取 processor array
// 只输出 processor name 的 array, 不做任何检查和转换
func extractProcessorArrFromBlueprint(db *gorm.DB, blueprintName string) ([]string, error) {
	blueprint, err := dao.GetBlueprintByName(db, blueprintName)
	if err != nil {
		return []string{}, fmt.Errorf("blueprint name [%s] not found", blueprintName)
	}
	processorNameList := lo.Map(blueprint.BlueprintConfig, func(item dao.BlueprintItem, index int) string {
		return item.ProcessorName
	})
	return processorNameList, nil
}