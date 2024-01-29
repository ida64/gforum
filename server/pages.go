package main

import "github.com/gin-gonic/gin"

type RootIndexView struct {
	GenericView *GenericView

	IsRestore   bool
	RestorePath string
}

func NewRootIndexView(c *gin.Context) *RootIndexView {
	var view = RootIndexView{
		GenericView: NewGenericView(c),
	}

	return &view
}

func NewRootIndexViewWithRestore(c *gin.Context, path string) *RootIndexView {
	var view = RootIndexView{
		GenericView: NewGenericView(c),
		IsRestore:   true,
		RestorePath: path,
	}

	return &view
}

func renderRootIndexPage(c *gin.Context) {
	var content = parseTmplFromResources("layout.html", "root/index.html")

	err := renderLayoutWithContent(c, content, RootIndexView{GenericView: NewGenericView(c)})
	if err != nil {
		c.String(500, "error rendering layout: %s", err)
		return
	}

	c.Status(200)
}

func restorePageContent(c *gin.Context, path string) {
	var content = parseTmplFromResources("layout.html", "root/index.html")

	err := renderLayoutWithContent(c, content, NewRootIndexViewWithRestore(c, path))
	if err != nil {
		c.String(500, "error rendering layout: %s", err)
		return
	}

	c.Status(200)
}

func restorePageMiddleware(c *gin.Context) {
	if c.GetHeader("Hx-Request") == "" {
		restorePageContent(c, c.Request.URL.Path)
		c.Abort()
	}
	c.Next()
}
