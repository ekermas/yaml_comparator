package mapper

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"yaml_comparator/pkg/commons"
)

func mapSliceAsFreeInterface(vs []string, f func(string) interface{}) []interface{} {
	vsm := make([]interface{}, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func ParseYamlAsFreeForm(fileContents []string) []interface{} {
	return mapSliceAsFreeInterface(fileContents, func(item string) interface{} {
		var freeForm interface{}
		err := yaml.Unmarshal([]byte(item), &freeForm)
		commons.Check(err)
		fmt.Printf("---> freeForm:\n%+v\n\n", freeForm)
		return freeForm
	})
}
