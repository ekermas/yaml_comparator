package app

import (
	"yaml_comparator/pkg/commons"
	"yaml_comparator/pkg/files"
	"yaml_comparator/pkg/mapper"
	"yaml_comparator/pkg/template"
)

func mapSlice(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func Process(fileNames []string) {

	fileContents :=
		mapSlice(fileNames, func(item string) string { return files.ReadFile(item) })

	parsedFiles :=
		mapper.ParseYamlAsFreeForm(fileContents)

	joinedLinks := mapper.CreateLinksModel(parsedFiles)

	table := template.PrepareData(fileNames, joinedLinks)

	template.CreateReport("./data/template.html", "./data/index.html", table)

	err := commons.OpenBrowser("C:\\workspace\\mlr\\yaml_comparator\\data\\index.html")
	commons.Check(err)

}
