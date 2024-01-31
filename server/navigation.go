package main

import "github.com/gin-gonic/gin"

type NavigationView struct {
	GenericView *GenericView
}

func renderNavigationComponent(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/navigation.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NavigationView{GenericView: NewGenericView(c)})
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderSidebarComponent(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/sidebar.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewGenericView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}
