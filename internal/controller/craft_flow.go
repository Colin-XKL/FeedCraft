package controller

import (
	"FeedCraft/internal/craft"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateCraftFlow godoc
// @Summary Create a new CraftFlow
// @Description Create a new CraftFlow
// @Tags CraftFlow
// @Accept json
// @Produce json
// @Param craftFlow body CraftFlow true "CraftFlow data"
// @Success 201 {object} CraftFlow
// @Failure 400 {object} gin.H
// @Router /craft-flows [post]
func CreateCraftFlow(c *gin.Context) {
	var craftFlow dao.CraftFlow
	if err := c.ShouldBindJSON(&craftFlow); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	if err := db.Create(&craftFlow).Error; err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, util.APIResponse[dao.CraftFlow]{Data: craftFlow})
}

// GetCraftFlow godoc
// @Summary Get a CraftFlow by name
// @Description Get a CraftFlow by name
// @Tags CraftFlow
// @Produce json
// @Param name path string true "CraftFlow Name"
// @Success 200 {object} CraftFlow
// @Failure 404 {object} gin.H
// @Router /craft-flows/{name} [get]
func GetCraftFlow(c *gin.Context) {
	name := c.Param("name")
	db := util.GetDatabase()

	craftFlow, err := dao.GetCraftFlowByName(db, name)
	if err != nil {
		c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "CraftFlow not found"})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[dao.CraftFlow]{Data: *craftFlow})
}

// UpdateCraftFlow godoc
// @Summary Update a CraftFlow
// @Description Update a CraftFlow
// @Tags CraftFlow
// @Accept json
// @Produce json
// @Param name path string true "CraftFlow Name"
// @Param craftFlow body CraftFlow true "CraftFlow data"
// @Success 200 {object} CraftFlow
// @Failure 400 {object} gin.H
// @Router /craft-flows/{name} [put]
func UpdateCraftFlow(c *gin.Context) {
	name := c.Param("name")
	var craftFlow dao.CraftFlow
	if err := c.ShouldBindJSON(&craftFlow); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	existingCraftFlow, err := dao.GetCraftFlowByName(db, name)
	if err != nil {
		c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "CraftFlow not found"})
		return
	}

	existingCraftFlow.Description = craftFlow.Description
	existingCraftFlow.CraftFlowConfig = craftFlow.CraftFlowConfig

	if err := db.Save(existingCraftFlow).Error; err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[dao.CraftFlow]{Data: *existingCraftFlow})

}

// DeleteCraftFlow godoc
// @Summary Delete a CraftFlow
// @Description Delete a CraftFlow
// @Tags CraftFlow
// @Produce json
// @Param name path string true "CraftFlow Name"
// @Success 204 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /craft-flows/{name} [delete]
func DeleteCraftFlow(c *gin.Context) {
	name := c.Param("name")
	db := util.GetDatabase()

	craftFlow, err := dao.GetCraftFlowByName(db, name)
	if err != nil {
		c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "CraftFlow not found"})
		return
	}

	if err := db.Delete(craftFlow).Error; err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: nil})
}

// ListCraftFlows godoc
// @Summary List all CraftFlows
// @Description List all CraftFlows
// @Tags CraftFlow
// @Produce json
// @Success 200 {array} CraftFlow
// @Router /craft-flows [get]
func ListCraftFlows(c *gin.Context) {
	var craftFlows []dao.CraftFlow

	db := util.GetDatabase()
	if err := db.Find(&craftFlows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: craftFlows})
}

func ListSysCraftAtoms(c *gin.Context) {
	craftAtoms := craft.GetCraftAtomDict()
	var ret []map[string]string
	for _, meta := range craftAtoms {
		ret = append(ret, map[string]string{
			"name":        meta.Name,
			"description": meta.Description,
		})
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: ret})
}

func ListCraftTemplates(c *gin.Context) {
    values := make([]craft.CraftTemplate, 0, len(craftTemplates))
    for _, template := range craftTemplates {
        values = append(values, template)
    }
}
	craftTemplates := craft.GetSysCraftTemplateDict()
	var ret []craft.CraftTemplate
	for _, template := range craftTemplates {
		ret = append(ret, template)
	}
	c.JSON(http.StatusOK, util.APIResponse[[]craft.CraftTemplate]{Data: ret})
}
