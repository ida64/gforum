package main

import "github.com/gin-gonic/gin"

type NavigationView struct {
	GlobalView *GlobalView
}

func renderNavigationComponent(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/navigation.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NavigationView{GlobalView: NewGlobalView(c)})
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderSidebarComponent(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/sidebar.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewGlobalView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}
