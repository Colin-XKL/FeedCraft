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

type Meta struct {
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	CraftOptionList []CraftOption `json:"craft_option_list"`
}

func GetSysCraftEntries() map[string]Meta {
	craftEntries := make(map[string]Meta)
	craftEntries["proxy"] = Meta{
		Name:            "proxy",
		Description:     "proxy the feed ",
		CraftOptionList: []CraftOption{},
	}
	craftEntries["limit"] = Meta{
		Name:            "limit",
		Description:     "limit the number of entries to a single page",
		CraftOptionList: GetLimitCraftOption(),
	}
	craftEntries["fulltext"] = Meta{
		Name:            "fulltext",
		Description:     "extract fulltext for rss feed",
		CraftOptionList: GetFulltextCraftOptions(),
	}
	craftEntries["fulltext-plus"] = Meta{
		Name:            "fulltext-plus",
		Description:     "emulate the browser to extract fulltext for rss feed",
		CraftOptionList: GetFulltextPlusCraftOptions(),
	}
	craftEntries["introduction"] = Meta{
		Name:            "introduction",
		Description:     "add ai-generated introduction in the beginning of the article",
		CraftOptionList: GetAddIntroductionCraftOptions(),
	}
	craftEntries["ignore-advertorial"] = Meta{
		Name:            "ignore-advertorial",
		Description:     "exclude advertorial article using llm",
		CraftOptionList: GetIgnoreAdvertorialCraftOptions(),
	}
	craftEntries["translate-title"] = Meta{
		Name:            "translate-title",
		Description:     "translate title to Chinese using LLM",
		CraftOptionList: GetTranslateTitleCraftOptions(),
	}
	craftEntries["translate-content"] = Meta{
		Name:            "translate-content",
		Description:     "translate article content to Chinese using LLM",
		CraftOptionList: GetTranslateContentCraftOptions(),
	}
	return craftEntries
}

func Entry(c *gin.Context) {
	craftName := c.Param("craft-name")
	craftAtomDict := GetSysCraftEntries()
	db := util.GetDatabase()

	//TODO IMPLEMENT CUSTOM OPTION PARAMETERS
	craftOptionList, err := inner(db, &craftAtomDict, craftName, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	CommonCraftHandlerUsingCraftOptionList(c, craftOptionList)
	return
}

const MaxCallDepth = 5

// 递归地解出 craft option list
func inner(db *gorm.DB, craftAtomDict *map[string]Meta, craftName string, depthId int) ([]CraftOption, error) {
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

	craftOptionMeta, isKnownCraftAtom := (*craftAtomDict)[craftName]
	if isKnownCraftAtom {
		logrus.Infof("[%s] is known craft atom", craftName)
		return craftOptionMeta.CraftOptionList, nil
	} else {
		craftArr, checkErr := extractCraftArrFromFlow(db, craftName)
		if checkErr != nil {
			// then not a valid  craft name
			return []CraftOption{}, fmt.Errorf("not a valid craft name")
		}
		var retArr []CraftOption
		for _, extractedSubCraftName := range craftArr {
			sub, recurErr := inner(db, craftAtomDict, extractedSubCraftName, depthId+1)
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
