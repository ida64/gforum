package controller

import (
	utils "gforum/utils"

	"github.com/gin-gonic/gin"
)

type NavigationView struct {
	GlobalView *GlobalView
}

func NewNavigationView(c *gin.Context) *NavigationView {
	var view = NavigationView{
		GlobalView: NewGlobalView(c),
	}

	return &view
}

func renderNavigation(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/navigation.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewNavigationView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

type FeedSelectorView struct {
	GlobalView *GlobalView
	FeedID     int
}

func NewFeedSelectorView(c *gin.Context) *FeedSelectorView {
	var id int = utils.GetParamterInt(c, "feed_id")

	var view = FeedSelectorView{
		GlobalView: NewGlobalView(c),
		FeedID:     id,
	}

	return &view
}

func renderFeedSelector(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/feedSelector.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewFeedSelectorView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}
