package controller

import (
	"html/template"
	"sync"
)

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
