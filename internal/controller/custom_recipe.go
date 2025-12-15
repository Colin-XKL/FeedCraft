package controller

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/recipe"
	"FeedCraft/internal/util"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func CreateCustomRecipe(c *gin.Context) {
	var recipeData dao.CustomRecipeV2
	if err := c.ShouldBindJSON(&recipeData); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	if err := dao.CreateCustomRecipeV2(db, &recipeData); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, util.APIResponse[any]{Data: recipeData})
}

func GetCustomRecipe(c *gin.Context) {
	id := c.Param("id")
	db := util.GetDatabase()

	recipeData, err := dao.GetCustomRecipeByIDV2(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Recipe not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: recipeData})
}

// RecipeInfo recipe 的详细信息,包括dao.CustomRecipe的基本信息, 以及一些预热相关的统计信息
type RecipeInfo struct {
	ID             string    `json:"id,omitempty" binding:"required"`
	Description    string    `json:"description,omitempty"`
	Craft          string    `json:"craft" binding:"required"`
	SourceType     string    `json:"source_type"`
	SourceConfig   string    `json:"source_config"`
	IsActive       bool      `json:"is_active" `
	LastAccessedAt time.Time `json:"last_accessed_at"`
}

func ListCustomRecipe(c *gin.Context) {
	db := util.GetDatabase()
	recipeList, err := dao.ListCustomRecipeV2(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	recipeInfoList := lo.Map(recipeList, func(item *dao.CustomRecipeV2, index int) RecipeInfo {
		recipeStatus := recipe.Scheduler.GetContextInfo(item.ID)
		return RecipeInfo{
			ID:             item.ID,
			Description:    item.Description,
			Craft:          item.Craft,
			SourceType:     item.SourceType,
			SourceConfig:   item.SourceConfig,
			IsActive:       recipeStatus.IsActive,
			LastAccessedAt: recipeStatus.LastRequestTime,
		}
	})
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: recipeInfoList})
}

func UpdateCustomRecipe(c *gin.Context) {
	id := c.Param("id")
	var recipeData dao.CustomRecipeV2
	if err := c.ShouldBindJSON(&recipeData); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	_, err := dao.GetCustomRecipeByIDV2(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Recipe not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	if err := dao.UpdateCustomRecipeV2(db, &recipeData); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: recipeData})
}

func DeleteCustomRecipe(c *gin.Context) {
	id := c.Param("id")
	db := util.GetDatabase()

	if err := dao.DeleteCustomRecipeV2(db, id); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{})
}
