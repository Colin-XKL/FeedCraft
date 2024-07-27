package controller

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateCraftAtom godoc
// @Summary Create a new CraftAtom
// @Description Create a new CraftAtom
// @Tags CraftAtom
// @Accept json
// @Produce json
// @Param craftAtom body dao.CraftAtom true "CraftAtom data"
// @Success 201 {object} dao.CraftAtom
// @Failure 400 {object} gin.H
// @Router /craft-atoms [post]
func CreateCraftAtom(c *gin.Context) {
	var craftAtom dao.CraftAtom
	if err := c.ShouldBindJSON(&craftAtom); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	if err := dao.CreateCraftAtom(db, &craftAtom); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, util.APIResponse[dao.CraftAtom]{Data: craftAtom})
}

// GetCraftAtom godoc
// @Summary Get a CraftAtom by name
// @Description Get a CraftAtom by name
// @Tags CraftAtom
// @Produce json
// @Param name path string true "CraftAtom Name"
// @Success 200 {object} dao.CraftAtom
// @Failure 404 {object} gin.H
// @Router /craft-atoms/{name} [get]
func GetCraftAtom(c *gin.Context) {
	name := c.Param("name")
	db := util.GetDatabase()

	craftAtom, err := dao.GetCraftAtomByName(db, name)
	if err != nil {
		c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "CraftAtom not found"})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[dao.CraftAtom]{Data: *craftAtom})
}

// UpdateCraftAtom godoc
// @Summary Update a CraftAtom
// @Description Update a CraftAtom
// @Tags CraftAtom
// @Accept json
// @Produce json
// @Param name path string true "CraftAtom Name"
// @Param craftAtom body dao.CraftAtom true "CraftAtom data"
// @Success 200 {object} dao.CraftAtom
// @Failure 400 {object} gin.H
// @Router /craft-atoms/{name} [put]
func UpdateCraftAtom(c *gin.Context) {
	name := c.Param("name")
	var craftAtom dao.CraftAtom
	if err := c.ShouldBindJSON(&craftAtom); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	existingCraftAtom, err := dao.GetCraftAtomByName(db, name)
	if err != nil {
		c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "CraftAtom not found"})
		return
	}

	existingCraftAtom.Description = craftAtom.Description
	existingCraftAtom.TemplateName = craftAtom.TemplateName
	existingCraftAtom.Params = craftAtom.Params

	if err := dao.UpdateCraftAtom(db, existingCraftAtom); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[dao.CraftAtom]{Data: *existingCraftAtom})
}

// DeleteCraftAtom godoc
// @Summary Delete a CraftAtom
// @Description Delete a CraftAtom
// @Tags CraftAtom
// @Produce json
// @Param name path string true "CraftAtom Name"
// @Success 204 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /craft-atoms/{name} [delete]
func DeleteCraftAtom(c *gin.Context) {
	name := c.Param("name")
	db := util.GetDatabase()

	if err := dao.DeleteCraftAtom(db, name); err != nil {
		c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "CraftAtom not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListSysCraftAtoms godoc
// @Summary List all CraftAtoms
// @Description List all CraftAtoms
// @Tags CraftAtom
// @Produce json
// @Success 200 {array} dao.CraftAtom
// @Router /craft-atoms [get]
func ListSysCraftAtoms(c *gin.Context) {
	db := util.GetDatabase()

	craftAtoms, err := dao.GetAllCraftAtoms(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: craftAtoms})
}
