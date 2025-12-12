package recipe

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"time"
)

var Scheduler *util.PreheatingScheduler

var recipeMaxFetchTimeout = time.Minute * 10

func QueryCustomRecipeName(recipeName string) (*dao.CustomRecipe, error) {
	db := util.GetDatabase()
	logrus.Infof("正在查询 recipe，ID: [%s]", recipeName)
	recipe, err := dao.GetCustomRecipeByID(db, recipeName)
	if err != nil {
		logrus.Errorf("查询 recipe 失败，ID: [%s], 错误: %v", recipeName, err)
	} else {
		// Log only the recipe ID to avoid exposing sensitive or lengthy data
		logrus.Infof("成功找到 recipe，ID: [%s]", recipe.ID)
	}
	return recipe, err
}

func CustomRecipe(c *gin.Context) {
	recipeId := c.Param("id")
	logrus.Infof("获取到的 recipe ID: [%s]", recipeId)
	recipe, err := QueryCustomRecipeName(recipeId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Recipe not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	path := GetPathForCustomRecipe(recipe)
	response, err := RetrieveCraftRecipeUsingPath(path)
	logrus.Infof("add preheating task for recipe [%s]", recipeId)
	Scheduler.ScheduleTask(recipeId)

	if err != nil {
		logrus.Errorf("error making request to %s: %s", path, err)
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Internal Server Error"})
		return
	}

	if response.StatusCode() != http.StatusOK {
		logrus.Errorf("request to %s failed with status %d. body: %s", path, response.StatusCode(), response.Body())
		c.Data(response.StatusCode(), "text/plain; charset=utf-8", response.Body())
		return
	}
	c.Data(http.StatusOK, "text/xml; charset=utf-8", response.Body())
}

func GetPathForCustomRecipe(recipe *dao.CustomRecipe) string {
	path := fmt.Sprintf("/craft/%s?input_url=%s", recipe.Craft, url.QueryEscape(recipe.FeedURL))
	return path
}

func RetrieveCraftRecipeUsingPath(path string) (*resty.Response, error) {
	client := newRecipeClient()
	response, err := client.R().Get(path)
	return response, err
}

func newRecipeClient() *resty.Client {
	client := resty.New().SetBaseURL(
		fmt.Sprintf("http://localhost:%d", util.GetLocalPort()),
	).SetTimeout(recipeMaxFetchTimeout)
	return client
}
