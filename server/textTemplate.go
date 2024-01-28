package main

import (
	"text/template"
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

	return template.Must(template.ParseFiles(templateFiles...))
}
