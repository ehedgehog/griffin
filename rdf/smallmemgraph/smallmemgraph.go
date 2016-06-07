package smallmemgraph

import "code.google.com/p/project/trigger/rdf"

type SmallMemGraph map[rdf.Triple]struct{}

var _ rdf.Graph = SmallMemGraph(make(map[rdf.Triple]struct{}))

func NewSmallMemGraph() SmallMemGraph {
    return SmallMemGraph(make(map[rdf.Triple]struct{}))
}

func (s SmallMemGraph) Add(t rdf.Triple) bool {
    _, present := s[t]
    if !present { s[t] = struct{}{} }
    return !present
}

func (s SmallMemGraph) Remove(t rdf.Triple) bool {
    _, present := s[t]
    if present { delete(s, t) }
    return present
}

func (s SmallMemGraph) Size() uint64 {
    return uint64(len(s))
}

func matches(p, t rdf.Triple) bool {
    return match(p.S, t.S) && match(p.P, t.P) && match(p.O, t.O)
}

func match(p, t rdf.Term) bool {
    return p.Type() == rdf.T_VAR || p == t
}

func (s SmallMemGraph) FindAll(t rdf.Triple, f func(rdf.Triple)bool) bool {
    for ot := range s {
        if matches(t, ot) && f(ot) { return false }
    }
    return true
}

func (s SmallMemGraph) Profile() rdf.Profile {
    return nil
}

