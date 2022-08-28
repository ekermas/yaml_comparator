package files

import (
	"fmt"
	"os"
	"yaml_comparator/pkg/commons"
)

func ReadFile(src string) string {
	res, err := os.ReadFile(src)
	commons.Check(err)
	fmt.Printf("readFile(%s) dump:\n%s\n", src, string(res))
	return string(res)
}
