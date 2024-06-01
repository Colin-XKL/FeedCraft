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
)

func CustomRecipe(c *gin.Context) {
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
	client := newRecipeClient()
	path := fmt.Sprintf("/craft/%s?input_url=%s", recipe.Craft, url.QueryEscape(recipe.FeedURL))
	response, err := client.R().Get(path)

	if err != nil {
		logrus.Errorf("error making request to %s: %s", path, err)
		return
	}
	c.Data(http.StatusOK, "text/xml; charset=utf-8", response.Body())
}
func newRecipeClient() *resty.Client {
	client := resty.New().SetBaseURL(fmt.Sprintf("http://localhost:%d", util.GetLocalPort()))
	return client
}
