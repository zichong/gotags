package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/awalterschulze/gographviz"
)

func Tags2Graph(tags []Tag) {
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}

	fieldNameMaxLen := 0
	for _, tag := range tags {
		if tag.Type == Field {
			if len(tag.Name) > fieldNameMaxLen {
				fieldNameMaxLen = len(tag.Name)
			}
		}
	}

	structMap := make(map[string]map[string][]string)
	for _, tag := range tags {
		if tag.Type == Method || tag.Type == Field {
			structName := tag.Fields[ReceiverType]

			if _, ok := structMap[structName]; ok == false {
				structMap[structName] = make(map[string][]string)
			}
			if _, ok := structMap[structName][string(tag.Type)]; ok == false {
				structMap[structName][string(tag.Type)] = make([]string, 0)
			}

			s := ""
			switch tag.Type {
			case Method:
				s = fmt.Sprintf("- %s%s %s", tag.Name, tag.Fields[Signature], tag.Fields[TypeField])
			case Field:
				s = fmt.Sprintf("%s %s", tag.Name, tag.Fields[TypeField])
			}
			structMap[structName][string(tag.Type)] = append(structMap[structName][string(tag.Type)],
				s)
		}
	}

	graph.AddAttr("G", "rankdir", `"LR"`)
	for s, sInfo := range structMap {
		labelsPart := make([]string, 0)
		labelsPart = append(labelsPart, s)
		for _, members := range sInfo {
			s := ""
			for _, m := range members {
				s += fmt.Sprintf(`%s\l`, m)
			}
			labelsPart = append(labelsPart, s)
		}
		label := strings.Join(labelsPart, "|")
		graph.AddNode("G", s, map[string]string{
			"label": fmt.Sprintf(`"%s"`, label),
			"shape": `"record"`,
		})
	}

	output := graph.String()

	f, err := os.OpenFile("graph.gv", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0776)
	if err != nil {
		panic(err)
	}
	f.WriteString(output)
}
