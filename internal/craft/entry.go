package craft

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

func GetSysCraftTemplateDict() map[string]CraftTemplate {
	sysCraftTempList := make(map[string]CraftTemplate)
	sysCraftTempList["proxy"] = CraftTemplate{
		Name:                "proxy",
		Description:         "proxy the feed",
		ParamTemplateDefine: []ParamTemplate{},
		OptionFunc: func(m map[string]string) []CraftOption {
			return []CraftOption{}
		},
	}
	sysCraftTempList["limit"] = CraftTemplate{
		Name:                "limit",
		Description:         "limit the number of entries to a single page",
		ParamTemplateDefine: limitCraftParamTmpl,
		OptionFunc:          limitCraftLoadParams,
	}
	sysCraftTempList["fulltext"] = CraftTemplate{
		Name:                "fulltext",
		Description:         "extract fulltext for rss feed",
		ParamTemplateDefine: []ParamTemplate{},
		OptionFunc: func(m map[string]string) []CraftOption {
			return GetFulltextCraftOptions()
		},
	}
	sysCraftTempList["fulltext-plus"] = CraftTemplate{
		Name:                "fulltext-plus",
		Description:         "emulate the browser to extract fulltext for rss feed",
		ParamTemplateDefine: []ParamTemplate{},
		OptionFunc: func(m map[string]string) []CraftOption {
			return GetFulltextPlusCraftOptions()
		},
	}
	sysCraftTempList["introduction"] = CraftTemplate{
		Name:                "introduction",
		Description:         "add ai-generated introduction in the beginning of the article",
		ParamTemplateDefine: introCraftParamTmpl,
		OptionFunc:          introCraftLoadParam,
	}
	sysCraftTempList["ignore-advertorial"] = CraftTemplate{
		Name:                "ignore-advertorial",
		Description:         "exclude advertorial article using llm",
		ParamTemplateDefine: llmFilterCraftParamTmpl,
		OptionFunc:          llmFilterCraftLoadParam,
	}
	sysCraftTempList["translate-title"] = CraftTemplate{
		Name:                "translate-title",
		Description:         "translate title to Chinese using LLM",
		ParamTemplateDefine: transTitleParamTmpl,
		OptionFunc:          transTitleCraftLoadParam,
	}
	sysCraftTempList["translate-content"] = CraftTemplate{
		Name:                "translate-content",
		Description:         "translate article content to Chinese using LLM",
		ParamTemplateDefine: transContentParamTmpl,
		OptionFunc:          transContentCraftLoadParam,
	}
	return sysCraftTempList
}

func GetCraftAtomDict() map[string]CraftAtom {
	tmplDict := GetSysCraftTemplateDict()
	craftAtomDict := make(map[string]CraftAtom)
	for name, craftTemplate := range tmplDict {
		item := CraftAtom{
			Name:         craftTemplate.Name, // 默认会有个跟template 同名的craft atom
			Description:  fmt.Sprintf("(sys predefined)%s", craftTemplate.Description),
			TemplateName: craftTemplate.Name,
			Params:       map[string]string{},
		}
		craftAtomDict[name] = item
	}
	return craftAtomDict
}

func Entry(c *gin.Context) {
	craftName := c.Param("craft-name")
	craftAtomDict := GetCraftAtomDict()
	craftTmplDict := GetSysCraftTemplateDict()
	db := util.GetDatabase()

	//TODO IMPLEMENT CUSTOM OPTION PARAMETERS
	craftOptionList, err := inner(db, &craftAtomDict, &craftTmplDict, craftName, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	CommonCraftHandlerUsingCraftOptionList(c, craftOptionList)
	return
}

const MaxCallDepth = 5

// 递归地解出 craft option list
func inner(db *gorm.DB, craftAtomDict *map[string]CraftAtom, craftTmplDict *map[string]CraftTemplate, craftName string, depthId int) ([]CraftOption, error) {
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
