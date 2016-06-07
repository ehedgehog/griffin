package main

import "strconv"
import "fmt"

import "code.google.com/p/project/trigger/rdf"
import "code.google.com/p/project/trigger/turtle"

type ToGraph struct {
	g          rdf.Graph
	prefixes   map[string]string
	bnodeCount int
}

func (tg *ToGraph) Add(t rdf.Triple) bool {
	return tg.g.Add(t)
}

func (tg *ToGraph) SetBase(base rdf.IRI) {
}

func (tg *ToGraph) SetPrefix(prefix string, ns rdf.IRI) {
	tg.prefixes[prefix] = ns.Spelling()
}

func (tg *ToGraph) NewBNode() rdf.Term {
	s := rdf.AsBNode("_:B" + strconv.Itoa(tg.bnodeCount))
	tg.bnodeCount += 1
	return s
}

func (tg *ToGraph) Report(message string, where turtle.Location) {
	fmt.Printf("! %s %v\n", message, where)
}
