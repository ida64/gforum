package controller

import "github.com/gin-gonic/gin"

type IndexView struct {
	GlobalView *GlobalView

	IsRestore   bool
	RestorePath string
}

func NewIndexView(c *gin.Context) *IndexView {
	var view = IndexView{
		GlobalView: NewGlobalView(c),
	}

	return &view
}

func NewIndexViewWithRestore(c *gin.Context, path string) *IndexView {
	var view = IndexView{
		GlobalView:  NewGlobalView(c),
		IsRestore:   true,
		RestorePath: path,
	}

	return &view
}

func renderIndexPage(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("layout.html", "root/index.html")

	err := renderLayout(c, content, NewIndexView(c))
	if err != nil {
		c.String(500, "error rendering layout: %s", err)
		return
	}

	c.Status(200)
}

func restorePageContent(c *gin.Context, path string) {
	var content = parseHTMLTemplatesFromResources("layout.html", "root/index.html")

	err := renderLayout(c, content, NewIndexViewWithRestore(c, path))
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
