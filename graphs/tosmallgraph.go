package graphs

import "strconv"
import "fmt"

import "github.com/ehedgehog/griffin/rdf"
import "github.com/ehedgehog/griffin/turtle"

type ToGraph struct {
	G          rdf.Graph
	Prefixes   map[string]string
	BnodeCount int
}

func (tg *ToGraph) Add(t rdf.Triple) bool {
	return tg.G.Add(t)
}

func (tg *ToGraph) SetBase(base rdf.IRI) {
}

func (tg *ToGraph) SetPrefix(prefix string, ns rdf.IRI) {
	tg.Prefixes[prefix] = ns.Spelling()
}

func (tg *ToGraph) NewBNode() rdf.Term {
	s := rdf.AsBNode("_:B" + strconv.Itoa(tg.BnodeCount))
	tg.BnodeCount += 1
	return s
}

func (tg *ToGraph) Report(message string, where turtle.Location) {
	fmt.Printf("! %s %v\n", message, where)
}
