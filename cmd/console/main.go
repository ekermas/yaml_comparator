package main

import (
	"flag"
	"yaml_comparator/pkg/app"
)

type Files []string

func (i *Files) String() string {
	return "my string representation"
}

func (i *Files) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {

	var fileNames Files
	flag.Var(&fileNames, "file", "Some description for this param.")
	flag.Parse()

	app.Process(fileNames)

}
