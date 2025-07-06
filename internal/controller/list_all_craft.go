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

func ListAllCraft(c *gin.Context) {
	var allCrafts []map[string]interface{}
	db := util.GetDatabase()

	// 获取系统内置 craft atoms
	sysCraftTemplates := craft.GetSysCraftTemplateDict()
	for _, tmpl := range sysCraftTemplates {
		allCrafts = append(allCrafts, map[string]interface{}{
			"type":        SysDefinedCraftAtom,
			"name":        tmpl.Name,
			"description": tmpl.Description,
			"template":    tmpl.Name,
			"params":      map[string]string{}, // 系统内置的 atom 默认没有参数
		})
	}

	// 获取用户自定义 craft atom
	customAtoms, err := dao.GetAllCraftAtoms(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Failed to get user defined craft atoms"})
		return
	}
	for _, atom := range customAtoms {
		allCrafts = append(allCrafts, map[string]interface{}{
			"type":        UserDefinedCraftAtom,
			"name":        atom.Name,
			"description": atom.Description,
			"template":    atom.TemplateName,
			"params":      atom.Params,
		})
	}

	// 获取用户自定义 craft flow
	craftFlows, err := dao.GetAllCraftFlows(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Failed to get craft flows"})
		return
	}
	for _, flow := range craftFlows {
		allCrafts = append(allCrafts, map[string]interface{}{
			"type":              UserDefinedCraftFlow,
			"name":              flow.Name,
			"description":       flow.Description,
			"craft_flow_config": flow.CraftFlowConfig,
		})
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: allCrafts})
}
