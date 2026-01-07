package controller

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateTool godoc (formerly CreateCraftAtom)
// @Summary Create a new Tool
// @Description Create a new Tool
// @Tags Tool
// @Accept json
// @Produce json
// @Param tool body dao.Tool true "Tool data"
// @Success 201 {object} dao.Tool
// @Failure 400 {object} gin.H
// @Router /api/admin/tools [post]
func CreateTool(c *gin.Context) {
	var tool dao.Tool
	if err := c.ShouldBindJSON(&tool); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	if err := dao.CreateTool(db, &tool); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, util.APIResponse[dao.Tool]{Data: tool})
}

// GetTool godoc (formerly GetCraftAtom)
// @Summary Get a Tool by name
// @Description Get a Tool by name
// @Tags Tool
// @Produce json
// @Param name path string true "Tool Name"
// @Success 200 {object} dao.Tool
// @Failure 404 {object} gin.H
// @Router /api/admin/tools/{name} [get]
func GetTool(c *gin.Context) {
	name := c.Param("name")
	db := util.GetDatabase()

	tool, err := dao.GetToolByName(db, name)
	if err != nil {
		c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Tool not found"})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[dao.Tool]{Data: *tool})
}

// UpdateTool godoc (formerly UpdateCraftAtom)
// @Summary Update a Tool
// @Description Update a Tool
// @Tags Tool
// @Accept json
// @Produce json
// @Param name path string true "Tool Name"
// @Param tool body dao.Tool true "Tool data"
// @Success 200 {object} dao.Tool
// @Failure 400 {object} gin.H
// @Router /api/admin/tools/{name} [put]
func UpdateTool(c *gin.Context) {
	name := c.Param("name")
	var tool dao.Tool
	if err := c.ShouldBindJSON(&tool); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	existingTool, err := dao.GetToolByName(db, name)
	if err != nil {
		c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Tool not found"})
		return
	}

	existingTool.Description = tool.Description
	existingTool.TemplateName = tool.TemplateName
	existingTool.Params = tool.Params

	if err := dao.UpdateTool(db, existingTool); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[dao.Tool]{Data: *existingTool})
}

// DeleteTool godoc (formerly DeleteCraftAtom)
// @Summary Delete a Tool
// @Description Delete a Tool
// @Tags Tool
// @Produce json
// @Param name path string true "Tool Name"
// @Success 204 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/admin/tools/{name} [delete]
func DeleteTool(c *gin.Context) {
	name := c.Param("name")
	db := util.GetDatabase()

	if err := dao.DeleteTool(db, name); err != nil {
		c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Tool not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListTools godoc (formerly ListCraftAtoms)
// @Summary List all Tools
// @Description List all Tools
// @Tags Tool
// @Produce json
// @Success 200 {array} dao.Tool
// @Router /api/admin/tools [get]
func ListTools(c *gin.Context) {
	db := util.GetDatabase()

	tools, err := dao.GetAllTools(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: tools})
}
