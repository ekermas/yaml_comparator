package template

import (
	"fmt"
	"yaml_comparator/pkg/domain"
)

func PrepareData(fileNames []string, files [][]domain.Link) map[string][]string {
	columnsCount := 2*len(fileNames) - 1

	var rows = make(map[string][]string)
	rows["fileName"] = make([]string, columnsCount)
	for _, links := range files {
		for _, link := range links {
			if link.Left != (domain.Node{}) {
				rows[link.Left.Path] = make([]string, columnsCount)
			}
			if link.Right != (domain.Node{}) {
				rows[link.Right.Path] = make([]string, columnsCount)
			}
		}
	}

	for fileIdx := range files {
		baseIdx := 2 * fileIdx
		rows["fileName"][baseIdx+0] = fileNames[fileIdx]
		rows["fileName"][baseIdx+1] = "link"
		rows["fileName"][baseIdx+2] = fileNames[fileIdx+1]
	}

	for fileIdx, links := range files {
		baseIdx := 2 * fileIdx
		for _, link := range links {
			rows[link.Path][baseIdx+0] = link.Left.Value
			rows[link.Path][baseIdx+1] = fmt.Sprintf("%s / %s", link.GetValueStatus().ToString(), link.GetLinkStatus().ToString())
			rows[link.Path][baseIdx+2] = link.Right.Value
		}
	}

	fmt.Printf("---> table:\n%+v\n\n", rows)
	fmt.Printf("---> table.len:\n%+v\n\n", len(rows))

	return rows
}
