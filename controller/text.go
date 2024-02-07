package controller

import (
	"sync"
	"text/template"
)

var (
	textTemplateCache     = make(map[string]*template.Template)
	textTemplateCacheLock sync.RWMutex
)

/*
* parseTextTemplatesFromResources parses the supplied template files from the resources/templates directory.
* It returns a template.Template object.
 */
func parseTextTemplatesFromResources(filenames ...string) *template.Template {
	var numFiles = len(filenames)
	var templateFiles = make([]string, numFiles)

	for i := 0; i < numFiles; i++ {
		templateFiles[i] = "resources/templates/" + filenames[i]
	}

	textTemplateCacheLock.RLock()
	cachedTemplate, ok := textTemplateCache[templateFiles[0]]
	textTemplateCacheLock.RUnlock()

	if ok {
		return cachedTemplate
	}

	textTemplateCacheLock.Lock()
	defer textTemplateCacheLock.Unlock()

	tmpl, err := template.ParseFiles(templateFiles...)
	if err != nil {
		panic(err)
	}

	textTemplateCache[templateFiles[0]] = tmpl

	return tmpl

}
