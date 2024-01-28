package main

import "github.com/gin-gonic/gin"

type NavigationView struct {
	GenericView    *GenericView
	CategoryModels []CategoryModel
}

func renderNavigationComponent(c *gin.Context) {
	var categories []CategoryModel
	err := database.Find(&categories).Error
	if err != nil {
		renderErrorAlert(c, "error fetching categories")
		return
	}

	var content = parseTmplFromResources("components/navigation.html")

	err = content.ExecuteTemplate(c.Writer, "componentBody", NavigationView{GenericView: NewGenericView(c), CategoryModels: categories})
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderSidebarComponent(c *gin.Context) {
	var content = parseTmplFromResources("components/sidebar.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewGenericView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}
