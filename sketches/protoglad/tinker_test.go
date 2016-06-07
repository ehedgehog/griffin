package main

import "fmt"
import "testing"
import "code.google.com/p/project/trigger/rdf/smallmemgraph"
import "code.google.com/p/project/trigger/rdf"

func TestSomethingOrOther(t *testing.T) {
	g := smallmemgraph.NewSmallMemGraph()
	S := rdf.IRI("http://example.com/example/nwas")
	P := rdf.Label
	O := rdf.PlainString("Now We Are Six")
	O2 := rdf.PlainString("Then We Were Five")

	g.Add(rdf.Triple{S, P, O})
	g.Add(rdf.Triple{rdf.IRI("eh:/S"), rdf.Label, rdf.PlainString("Alternative Label")})
	g.Add(rdf.Triple{S, P, O2})
	fmt.Println("S labels:", ObjectsOf(g, S, P))
	fmt.Println("eh:/S labels:", ObjectsOf(g, rdf.IRI("eh:/S"), rdf.Label))

	t.Fatalf("BOOM today.", g)
}
