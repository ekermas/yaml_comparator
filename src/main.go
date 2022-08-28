package main

import (
	"flag"
	"fmt"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
	"html/template"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile(src string) string {
	res, err := os.ReadFile(src)
	check(err)
	fmt.Printf("readFile(%s) dump:\n%s\n", src, string(res))
	return string(res)
}

func mapSlice(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func mapSliceAsFreeInterface(vs []string, f func(string) interface{}) []interface{} {
	vsm := make([]interface{}, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func mapSliceAsNodeList(vs []interface{}, f func(interface{}) []Node) [][]Node {
	vsm := make([][]Node, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

type files []string

func (i *files) String() string {
	return "my string representation"
}

func (i *files) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {

	var fileNames files
	flag.Var(&fileNames, "file", "Some description for this param.")
	flag.Parse()

	fileContents := mapSlice(fileNames, func(item string) string { return readFile(item) })

	parsedFiles := mapSliceAsFreeInterface(fileContents, func(item string) interface{} {
		var freeForm interface{}
		err := yaml.Unmarshal([]byte(item), &freeForm)
		check(err)
		fmt.Printf("---> freeForm:\n%+v\n\n", freeForm)
		return freeForm
	})

	linearFiles := mapSliceAsNodeList(parsedFiles, func(item interface{}) []Node {
		res := checkNode(item, Node{Path: "root"}, []Node{})
		fmt.Printf("---> Nodes:\n%+v\n\n", res)
		return res
	})

	fmt.Printf("---> linearFiles:\n%+v\n\n", linearFiles)

	joined := make([][]Link, 0)
	for index := range linearFiles {
		if index+1 < len(linearFiles) {
			left := linearFiles[index]
			right := linearFiles[index+1]
			joined = append(joined, makeLinks(left, right))
		}
	}
	fmt.Printf("---> joined:\n%+v\n\n", joined)

	table := writer(fileNames, joined)

	fmt.Printf("---> table:\n%+v\n\n", table)
	fmt.Printf("---> table.len:\n%+v\n\n", len(table))

	f, err := os.Create("./data/index.html")
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	tmpl, _ := template.ParseFiles("./data/template.html")
	tmpl.Execute(f, table)

	open("C:\\workspace\\mlr\\yaml_comparator\\data\\index.html")

}

func writer(fileNames []string, files [][]Link) map[string][]string {

	columnsCount := 2*len(fileNames) - 1

	var rows = make(map[string][]string)
	rows["fileName"] = make([]string, columnsCount)
	for _, links := range files {
		for _, link := range links {
			if link.Master != (Node{}) {
				rows[link.Master.Path] = make([]string, columnsCount)
			}
			if link.Slave != (Node{}) {
				rows[link.Slave.Path] = make([]string, columnsCount)
			}
		}
	}

	for fileIdx, _ := range files {
		baseIdx := 2 * fileIdx
		rows["fileName"][baseIdx+0] = fileNames[fileIdx]
		rows["fileName"][baseIdx+1] = "link"
		rows["fileName"][baseIdx+2] = fileNames[fileIdx+1]
	}

	for fileIdx, links := range files {
		baseIdx := 2 * fileIdx
		for _, link := range links {
			rows[link.Path][baseIdx+0] = link.Master.Value
			rows[link.Path][baseIdx+1] = fmt.Sprintf("%s / %s", link.GetValueStatus().ToString(), link.GetLinkStatus().ToString())
			rows[link.Path][baseIdx+2] = link.Slave.Value
		}
	}

	return rows

}

func makeLinks(left []Node, right []Node) []Link {

	joined := make([]Link, 0)

	for li := range left {
		lItem := left[li]
		joined = append(joined, Link{
			Path:   lItem.Path,
			Master: lItem,
		})
	}

	for ri := range right {
		rItem := right[ri]
		idx := slices.IndexFunc(joined, func(c Link) bool {
			return c.Path == rItem.Path
		})
		if idx == -1 {
			joined = append(joined, Link{
				Path:  rItem.Path,
				Slave: rItem,
			})
		} else {
			joined[idx].Slave = rItem
		}
	}

	return joined
}

type Node struct {
	Path  string
	Value string
	Type  string
}

type ValueStatus int

const (
	Eq ValueStatus = iota
	Ne
	Wt
)

func (src ValueStatus) ToString() string {
	switch src {
	case Eq:
		return "Eq"
	case Ne:
		return "Ne"
	case Wt:
		return "WrongType"
	default:
		panic(fmt.Sprintf("Unknown ValueStatus type %v", src))
	}
}

type LinkStatus int

const (
	Present LinkStatus = iota
	NoLeft
	NoRight
)

func (src LinkStatus) ToString() string {
	switch src {
	case Present:
		return "Present"
	case NoLeft:
		return "NoLeft"
	case NoRight:
		return "NoRight"
	default:
		panic(fmt.Sprintf("Unknown LinkStatus type %v", src))
	}
}

type Link struct {
	Path   string
	Master Node
	Slave  Node
}

func (src Link) GetValueStatus() ValueStatus {
	if src.Master.Type != src.Slave.Type {
		return Wt
	}
	if src.Master.Value != src.Slave.Value {
		return Ne
	}
	return Eq
}

func (src Link) GetLinkStatus() LinkStatus {
	if src.Master == (Node{}) {
		return NoLeft
	}
	if src.Slave == (Node{}) {
		return NoRight
	}
	return Present
}

func checkNode(src interface{}, currentNode Node, resAcc []Node) []Node {

	addToCurrent := func(c Node, value string, tpe string) Node {
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
			resAcc = append(resAcc, checkNode(item, copyNode, []Node{})...)
		}
	case map[string]interface{}:
		for key := range v {
			copyNode := currentNode
			copyNode.Path = fmt.Sprintf("%s.%s", currentNode.Path, key)
			fmt.Printf("%v -> map %v\n", copyNode.Path, v)
			resAcc = append(resAcc, checkNode(v[key], copyNode, []Node{})...)
		}
	default:
		panic(fmt.Sprintf("I don't know about type %T!\n", v))
	}

	return resAcc
}
