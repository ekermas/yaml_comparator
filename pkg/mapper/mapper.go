package mapper

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strconv"
	"yaml_comparator/pkg/domain"
)

func mapSliceAsNodeList(vs []interface{}, f func(interface{}) []domain.Node) [][]domain.Node {
	vsm := make([][]domain.Node, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func CreateLinksModel(src []interface{}) [][]domain.Link {
	linearFiles := mapSliceAsNodeList(src, func(item interface{}) []domain.Node {
		res := checkNode(item, domain.Node{Path: "root"}, []domain.Node{})
		fmt.Printf("---> Nodes:\n%+v\n\n", res)
		return res
	})

	fmt.Printf("---> linearFiles:\n%+v\n\n", linearFiles)

	joined := make([][]domain.Link, 0)
	for index := range linearFiles {
		if index+1 < len(linearFiles) {
			left := linearFiles[index]
			right := linearFiles[index+1]
			joined = append(joined, makeLinks(left, right))
		}
	}
	fmt.Printf("---> joined:\n%+v\n\n", joined)

	return joined
}

func makeLinks(left []domain.Node, right []domain.Node) []domain.Link {

	joined := make([]domain.Link, 0)

	for li := range left {
		lItem := left[li]
		joined = append(joined, domain.Link{
			Path: lItem.Path,
			Left: lItem,
		})
	}

	for ri := range right {
		rItem := right[ri]
		idx := slices.IndexFunc(joined, func(c domain.Link) bool {
			return c.Path == rItem.Path
		})
		if idx == -1 {
			joined = append(joined, domain.Link{
				Path:  rItem.Path,
				Right: rItem,
			})
		} else {
			joined[idx].Right = rItem
		}
	}

	return joined
}

func checkNode(src interface{}, currentNode domain.Node, resAcc []domain.Node) []domain.Node {

	addToCurrent := func(c domain.Node, value string, tpe string) domain.Node {
		c.Value = value
		c.Type = tpe
		return c
	}

	switch v := src.(type) {
	case int:
		fmt.Printf("%v -> int %v\n", currentNode.Path, v)
		c := addToCurrent(currentNode, strconv.Itoa(v), "int")
		resAcc = append(resAcc, c)
	case string:
		fmt.Printf("%v -> string %v\n", currentNode.Path, v)
		c := addToCurrent(currentNode, v, "string")
		resAcc = append(resAcc, c)
	case float64:
		fmt.Printf("%v -> float %v\n", currentNode.Path, v)
		c := addToCurrent(currentNode, fmt.Sprintf("%f", v), "float")
		resAcc = append(resAcc, c)
	case nil:
		fmt.Printf("%v -> nil %v\n", currentNode.Path, v)
		c := addToCurrent(currentNode, "null", "null")
		resAcc = append(resAcc, c)
	case []interface{}:
		for index, item := range v {
			copyNode := currentNode
			copyNode.Path = fmt.Sprintf("%s.%d", currentNode.Path, index)
			fmt.Printf("%v -> array %v\n", item, v)
			resAcc = append(resAcc, checkNode(item, copyNode, []domain.Node{})...)
		}
	case map[string]interface{}:
		for key := range v {
			copyNode := currentNode
			copyNode.Path = fmt.Sprintf("%s.%s", currentNode.Path, key)
			fmt.Printf("%v -> map %v\n", copyNode.Path, v)
			resAcc = append(resAcc, checkNode(v[key], copyNode, []domain.Node{})...)
		}
	default:
		panic(fmt.Sprintf("I don't know about type %T!\n", v))
	}

	return resAcc
}
