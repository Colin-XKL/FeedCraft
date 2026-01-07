package controller

import (
	"FeedCraft/internal/craft"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateBlueprint godoc
// @Summary Create a new Blueprint
// @Description Create a new Blueprint
// @Tags Blueprint
// @Accept json
// @Produce json
// @Param blueprint body dao.Blueprint true "Blueprint data"
// @Success 201 {object} dao.Blueprint
// @Failure 400 {object} gin.H
// @Router /api/admin/blueprints [post]
func CreateCraftFlow(c *gin.Context) {
	var blueprint dao.Blueprint
	if err := c.ShouldBindJSON(&blueprint); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	if err := db.Create(&blueprint).Error; err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, util.APIResponse[dao.Blueprint]{Data: blueprint})
}

// GetCraftFlow godoc
// @Summary Get a Blueprint by name
// @Description Get a Blueprint by name
// @Tags Blueprint
// @Produce json
// @Param name path string true "Blueprint Name"
// @Success 200 {object} dao.Blueprint
// @Failure 404 {object} gin.H
// @Router /api/admin/blueprints/{name} [get]
func GetCraftFlow(c *gin.Context) {
	name := c.Param("name")
	db := util.GetDatabase()

	blueprint, err := dao.GetBlueprintByName(db, name)
	if err != nil {
		c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Blueprint not found"})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[dao.Blueprint]{Data: *blueprint})
}

// UpdateCraftFlow godoc
// @Summary Update a Blueprint
// @Description Update a Blueprint
// @Tags Blueprint
// @Accept json
// @Produce json
// @Param name path string true "Blueprint Name"
// @Param blueprint body dao.Blueprint true "Blueprint data"
// @Success 200 {object} dao.Blueprint
// @Failure 400 {object} gin.H
// @Router /api/admin/blueprints/{name} [put]
func UpdateCraftFlow(c *gin.Context) {
	name := c.Param("name")
	var blueprint dao.Blueprint
	if err := c.ShouldBindJSON(&blueprint); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	existingBlueprint, err := dao.GetBlueprintByName(db, name)
	if err != nil {
		c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Blueprint not found"})
		return
	}

	existingBlueprint.Description = blueprint.Description
	existingBlueprint.BlueprintConfig = blueprint.BlueprintConfig

	if err := db.Save(existingBlueprint).Error; err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[dao.Blueprint]{Data: *existingBlueprint})

}

// DeleteCraftFlow godoc
// @Summary Delete a Blueprint
// @Description Delete a Blueprint
// @Tags Blueprint
// @Produce json
// @Param name path string true "Blueprint Name"
// @Success 204 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/admin/blueprints/{name} [delete]
func DeleteCraftFlow(c *gin.Context) {
	name := c.Param("name")
	db := util.GetDatabase()

	blueprint, err := dao.GetBlueprintByName(db, name)
	if err != nil {
		c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Blueprint not found"})
		return
	}

	if err := db.Delete(blueprint).Error; err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: nil})
}

// ListCraftFlows godoc
// @Summary List all Blueprints
// @Description List all Blueprints
// @Tags Blueprint
// @Produce json
// @Success 200 {array} dao.Blueprint
// @Router /api/admin/blueprints [get]
func ListCraftFlows(c *gin.Context) {
	var blueprints []dao.Blueprint

	db := util.GetDatabase()
	if err := db.Find(&blueprints).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: blueprints})
}

func ListSysCraftAtoms(c *gin.Context) {
	craftAtoms := craft.GetToolDict()
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
	craftTemplates := craft.GetSysCraftTemplateDict()
	ret := make([]craft.CraftTemplate, len(craftTemplates))
	i := 0
	for _, template := range craftTemplates {
		ret[i] = template
		i++
	}
	c.JSON(http.StatusOK, util.APIResponse[[]craft.CraftTemplate]{Data: ret})
}