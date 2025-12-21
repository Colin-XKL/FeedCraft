package craft

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"fmt"
	"net/http"

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

func GetCraftAtomDict() map[string]dao.CraftAtom {
	tmplDict := GetSysCraftTemplateDict()
	craftAtomDict := make(map[string]dao.CraftAtom)
	for name, craftTemplate := range tmplDict {
		item := dao.CraftAtom{
			Name:         craftTemplate.Name, // 默认会有个跟template 同名的craft atom
			Description:  craftTemplate.Description,
			TemplateName: craftTemplate.Name,
			Params:       map[string]string{},
		}
		craftAtomDict[name] = item
	}

	db := util.GetDatabase()
	craftAtomList, err := dao.GetAllCraftAtoms(db)
	if err != nil {
		logrus.Errorf("read craft atom list from db error. only built-in atom will work now. err: %s", err)
	} else {
		for _, atom := range craftAtomList {
			craftAtomDict[atom.Name] = atom
		}
	}

	return craftAtomDict
}

func Entry(c *gin.Context) {
	craftName := c.Param("craft-name")
	db := util.GetDatabase()

	craftOptionList, err := getCraftOptions(db, craftName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	CommonCraftHandlerUsingCraftOptionList(c, craftOptionList)
}

func ProcessFeed(feed *gofeed.Feed, feedURL string, craftName string) (*feeds.Feed, error) {
	db := util.GetDatabase()
	craftOptionList, err := getCraftOptions(db, craftName)
	if err != nil {
		return nil, err
	}

	craftedFeed, err := NewCraftedFeedFromGofeed(feed, feedURL, craftOptionList...)
	if err != nil {
		return nil, err
	}

	return craftedFeed.OutputFeed, nil
}

func getCraftOptions(db *gorm.DB, craftName string) ([]CraftOption, error) {
	craftAtomDict := GetCraftAtomDict()
	craftTmplDict := GetSysCraftTemplateDict()
	return inner(db, &craftAtomDict, &craftTmplDict, craftName, 0)
}

const MaxCallDepth = 5

// 递归地解出 craft option list
func inner(db *gorm.DB, craftAtomDict *map[string]dao.CraftAtom, craftTmplDict *map[string]CraftTemplate, craftName string, depthId int) ([]CraftOption, error) {
	if depthId+1 > MaxCallDepth {
		return []CraftOption{}, fmt.Errorf("max call depth hit")
	}
	logrus.Infof("checking %s", craftName)

	// 对于每一个 craft name in flow
	// 1. 判断是否为 built-in craft atom
	// 1.1 如果是, 直接返回对应的 craft option
	// 2. 如果不是, 判断是否为flow
	// 2.1 如果是flow, 调用 entry check
	// 如果不是, 说明craft name 无效, 返回error

	craftAtom, isKnownCraftAtom := (*craftAtomDict)[craftName]
	if isKnownCraftAtom {
		logrus.Infof("[%s] is known craft atom", craftName)
		tmplContent, tmplValid := (*craftTmplDict)[craftAtom.TemplateName]
		if !tmplValid {
			return []CraftOption{}, fmt.Errorf("invalid tmpl name [%s] for craft atom [%s]", craftAtom.TemplateName, craftAtom.Name)
		}
		return tmplContent.GetOptions(craftAtom.Params), nil
	} else {
		craftArr, checkErr := extractCraftArrFromFlow(db, craftName)
		if checkErr != nil {
			// then not a valid  craft name
			return []CraftOption{}, fmt.Errorf("not a valid craft name")
		}
		var retArr []CraftOption
		for _, extractedSubCraftName := range craftArr {
			sub, recurErr := inner(db, craftAtomDict, craftTmplDict, extractedSubCraftName, depthId+1)
			if recurErr != nil {
				return []CraftOption{}, recurErr
			}
			retArr = append(retArr, sub...)
		}
		return retArr, nil
	}
}

// 给定flowName,查询其详情, 获取 craft array
// 只输出craft name 的array, 不做任何检查和转换
func extractCraftArrFromFlow(db *gorm.DB, flowName string) ([]string, error) {
	flowContent, err := dao.GetCraftFlowByName(db, flowName)
	if err != nil {
		return []string{}, fmt.Errorf("craft flow name [%s] not found", flowName)
	}
	craftNameList := lo.Map(flowContent.CraftFlowConfig, func(item dao.CraftFlowItem, index int) string {
		return item.CraftName
	})
	return craftNameList, nil
}
