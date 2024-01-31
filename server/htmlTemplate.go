package main

import (
	"html/template"
	"sync"

	"github.com/gin-gonic/gin"
)

type LayoutView struct {
	Content *template.Template
	View    interface{}
}

func NewLayoutView(c *gin.Context) *LayoutView {
	var layout LayoutView
	return &layout
}

var (
	htmlTemplateCache     = make(map[string]*template.Template)
	htmlTemplateCacheLock sync.RWMutex
)

/*
* parseHTMLTemplatesFromResources parses the supplied template files from the resources/templates directory.
* It returns a template.Template object.
 */
func parseHTMLTemplatesFromResources(filenames ...string) *template.Template {
	var numFiles = len(filenames)
	var templateFiles = make([]string, numFiles)

	for i := 0; i < numFiles; i++ {
		templateFiles[i] = "resources/templates/" + filenames[i]
	}

	htmlTemplateCacheLock.RLock()
	cachedTemplate, ok := htmlTemplateCache[templateFiles[0]]
	htmlTemplateCacheLock.RUnlock()

	if ok {
		return cachedTemplate
	}

	htmlTemplateCacheLock.Lock()
	defer htmlTemplateCacheLock.Unlock()

	tmpl, err := template.ParseFiles(templateFiles...)
	if err != nil {
		panic(err)
	}

	htmlTemplateCache[templateFiles[0]] = tmpl

	return tmpl
}

/*
* RenderLayoutWithContent renders the layout template with the supplied content template.
* It returns an error if the template could not be rendered.
 */
func renderLayoutWithContent(c *gin.Context, content *template.Template, contentTemplateView interface{}) error {
	var layout = NewLayoutView(c)

	layout.Content = content
	layout.View = contentTemplateView

	err := content.ExecuteTemplate(c.Writer, "layout", layout)
	if err != nil {
		return err
	}

	return nil
}

type Alert struct {
	Message    string
	AutoDelete bool
}

func renderErrorAlert(c *gin.Context, message string) {
	var content = parseHTMLTemplatesFromResources("components/alerts/error.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", &Alert{Message: message})
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderError(c *gin.Context, err error) {
	renderErrorAlert(c, err.Error())
}

func renderSuccessAlert(c *gin.Context, message string, autoDelete bool) {
	var content = parseHTMLTemplatesFromResources("components/alerts/success.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", &Alert{Message: message, AutoDelete: autoDelete})
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}
