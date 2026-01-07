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
func CreateBlueprint(c *gin.Context) {
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

// GetBlueprint godoc
// @Summary Get a Blueprint by name
// @Description Get a Blueprint by name
// @Tags Blueprint
// @Produce json
// @Param name path string true "Blueprint Name"
// @Success 200 {object} dao.Blueprint
// @Failure 404 {object} gin.H
// @Router /api/admin/blueprints/{name} [get]
func GetBlueprint(c *gin.Context) {
	name := c.Param("name")
	db := util.GetDatabase()

	blueprint, err := dao.GetBlueprintByName(db, name)
	if err != nil {
		c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Blueprint not found"})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[dao.Blueprint]{Data: *blueprint})
}

// UpdateBlueprint godoc
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
func UpdateBlueprint(c *gin.Context) {
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

// DeleteBlueprint godoc
// @Summary Delete a Blueprint
// @Description Delete a Blueprint
// @Tags Blueprint
// @Produce json
// @Param name path string true "Blueprint Name"
// @Success 204 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/admin/blueprints/{name} [delete]
func DeleteBlueprint(c *gin.Context) {
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

// ListBlueprints godoc
// @Summary List all Blueprints
// @Description List all Blueprints
// @Tags Blueprint
// @Produce json
// @Success 200 {array} dao.Blueprint
// @Router /api/admin/blueprints [get]
func ListBlueprints(c *gin.Context) {
	var blueprints []dao.Blueprint

	db := util.GetDatabase()
	if err := db.Find(&blueprints).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: blueprints})
}

func ListSysTools(c *gin.Context) {
	tools := craft.GetToolDict()
	var ret []map[string]string
	for _, meta := range tools {
		ret = append(ret, map[string]string{
			"name":        meta.Name,
			"description": meta.Description,
		})
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: ret})
}

func ListToolTemplates(c *gin.Context) {
	toolTemplates := craft.GetSysCraftTemplateDict()
	ret := make([]craft.CraftTemplate, len(toolTemplates))
	i := 0
	for _, template := range toolTemplates {
		ret[i] = template
		i++
	}
	c.JSON(http.StatusOK, util.APIResponse[[]craft.CraftTemplate]{Data: ret})
}
