package controller

import "github.com/gin-gonic/gin"

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
