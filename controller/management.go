package controller

import "github.com/gin-gonic/gin"

type ManagementView struct {
	GlobalView *GlobalView
}

func NewManagementView(globalView *GlobalView) *ManagementView {
	return &ManagementView{
		GlobalView: globalView,
	}
}

func renderManagement(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/management/main.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewManagementView(NewGlobalView(c)))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderManagementNavigation(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/management/navigation.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewManagementView(NewGlobalView(c)))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderManagementCategories(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/management/categories.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewManagementView(NewGlobalView(c)))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}
