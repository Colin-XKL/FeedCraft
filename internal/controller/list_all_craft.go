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

	// 获取系统内置 craft atoms
	sysCraftTemplates := craft.GetSysCraftTemplateDict()
	for _, tmpl := range sysCraftTemplates {
		allCrafts = append(allCrafts, CraftItem{
			Type:        SysDefinedCraftAtom,
			Name:        tmpl.Name,
			Description: tmpl.Description,
		})
	}

	// 获取用户自定义 craft atom
	customAtoms, err := dao.GetAllCraftAtoms(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Failed to get user defined craft atoms"})
		return
	}
	for _, atom := range customAtoms {
		allCrafts = append(allCrafts, CraftItem{
			Type:        UserDefinedCraftAtom,
			Name:        atom.Name,
			Description: atom.Description,
		})
	}

	// 获取用户自定义 craft flow
	craftFlows, err := dao.GetAllCraftFlows(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Failed to get craft flows"})
		return
	}
	for _, flow := range craftFlows {
		allCrafts = append(allCrafts, CraftItem{
			Type:        UserDefinedCraftFlow,
			Name:        flow.Name,
			Description: flow.Description,
		})
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: allCrafts})
}
