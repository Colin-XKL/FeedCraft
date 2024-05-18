package controller

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func CreateCustomRecipe(c *gin.Context) {
	var recipe dao.CustomRecipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	if err := dao.CreateCustomRecipe(db, &recipe); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, util.APIResponse[any]{Data: recipe})
}

func GetCustomRecipe(c *gin.Context) {
	id := c.Param("id")
	db := util.GetDatabase()

	recipe, err := dao.GetCustomRecipeByID(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Recipe not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: recipe})
}

func ListCustomRecipe(c *gin.Context) {
	db := util.GetDatabase()
	recipeList, err := dao.ListCustomRecipe(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: recipeList})
}

func UpdateCustomRecipe(c *gin.Context) {
	id := c.Param("id")
	var recipe dao.CustomRecipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	_, err := dao.GetCustomRecipeByID(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Recipe not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	if err := dao.UpdateCustomRecipe(db, &recipe); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: recipe})
}

func DeleteCustomRecipe(c *gin.Context) {
	id := c.Param("id")
	db := util.GetDatabase()

	if err := dao.DeleteCustomRecipe(db, id); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{})
}
