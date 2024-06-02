package craft

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
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
	entries := GetSysCraftEntries()
	craftOptionMeta, exist := entries[craftName]
	if exist {
		logrus.Infof("requesting craft [%s]", craftName)
		CommonCraftHandlerUsingCraftOptionList(c, craftOptionMeta.CraftOptionList)
		return
	}

	db := util.GetDatabase()
	flowByName, err := dao.GetCraftFlowByName(db, craftName)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: fmt.Sprintf("craft flow name [%s] not found", craftName)})
		logrus.Infof("invalid craft flow [%s]", craftName)
		return
	}
	logrus.Infof("requesting craft flow [%s]", craftName)

	craftNameList := lo.Map(flowByName.CraftFlowConfig, func(item dao.CraftFlowItem, index int) string {
		return item.CraftName // TODO IMPLEMENT CUSTOM OPTION PARAMETERS
	})
	var craftListForFlow []CraftOption
	for _, craftNameInFlow := range craftNameList {
		craftOptionMeta, exist := entries[craftNameInFlow]
		if !exist {
			c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: fmt.Sprintf("craft name [%s] in flow [%s] not found", craftNameInFlow, flowByName.Name)})
			return
		}
		craftListForFlow = append(craftListForFlow, craftOptionMeta.CraftOptionList...)
	}

	CommonCraftHandlerUsingCraftOptionList(c, craftListForFlow)
	return

}
