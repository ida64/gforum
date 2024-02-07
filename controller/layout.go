package controller

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

type LayoutView struct {
	Content *template.Template
}

func NewLayoutView(content *template.Template) *LayoutView {
	var layout LayoutView
	layout.Content = content
	return &layout
}

func renderLayout(c *gin.Context, content *template.Template, data interface{}) error {
	err := content.ExecuteTemplate(c.Writer, "layout", data)
	if err != nil {
		return err
	}

	c.Status(200)

	return nil
}
