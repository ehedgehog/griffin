package smallmemgraph

import "project/trigger/github.com/ehedgehog/griffin/project/trigger/rdf"
import "testing"

func TestSomething(t *testing.T) {
    g := NewSmallMemGraph()
    if g.Size() != 0 { t.Fatalf( "empty graph show have size 0" ) }
    S, P, O := rdf.IRI("http://example.com/S"), rdf.IRI("http://example.com/P"), rdf.IRI("http://example.com/O")
    tr := rdf.Triple{S, P, O}
    added := g.Add(tr)
    if g.Size() != 1 { t.Fatalf( "After Add, graph shouwld have size 1") }
    if !added { t.Fatalf("Adding fresh triple should return true:") }
    added2 := g.Add(tr)
    if g.Size() != 1 { t.Fatalf("Re-adding triple should not increase size.")}
    if added2 { t.Fatalf( "Adding existing triple should return false.")}
}

