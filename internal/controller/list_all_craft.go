package controller

import (
	"FeedCraft/internal/craft"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CraftType string

var (
	SysDefinedCraftAtom  CraftType = "@sys/atom"
	UserDefinedCraftAtom CraftType = "@user/atom"
	UserDefinedCraftFlow CraftType = "@user/flow"
)

type CraftItem struct {
	Type        CraftType `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func ListAllCraft(c *gin.Context) {
	var allCrafts []CraftItem
	db := util.GetDatabase()

	// 获取系统内置 craft atoms (tools)
	sysCraftTemplates := craft.GetSysCraftTemplateDict()
	for _, tmpl := range sysCraftTemplates {
		allCrafts = append(allCrafts, CraftItem{
			Type:        SysDefinedCraftAtom,
			Name:        tmpl.Name,
			Description: tmpl.Description,
		})
	}

	// 获取用户自定义 craft atom (tools)
	customTools, err := dao.GetAllTools(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Failed to get user defined tools"})
		return
	}
	for _, tool := range customTools {
		allCrafts = append(allCrafts, CraftItem{
			Type:        UserDefinedCraftAtom,
			Name:        tool.Name,
			Description: tool.Description,
		})
	}

	// 获取用户自定义 craft flow (blueprints)
	blueprints, err := dao.GetAllBlueprints(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Failed to get blueprints"})
		return
	}
	for _, blueprint := range blueprints {
		allCrafts = append(allCrafts, CraftItem{
			Type:        UserDefinedCraftFlow,
			Name:        blueprint.Name,
			Description: blueprint.Description,
		})
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: allCrafts})
}