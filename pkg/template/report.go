package template

import (
	"html/template"
	"yaml_comparator/pkg/commons"
	"yaml_comparator/pkg/files"
)

func CreateReport(templatePath string, outputPath string, context any) {

	f, err := files.CreateFile(outputPath)
	commons.Check(err)

	tmpl, err := template.ParseFiles(templatePath)
	commons.Check(err)

	err = tmpl.Execute(f, context)
	commons.Check(err)
}
