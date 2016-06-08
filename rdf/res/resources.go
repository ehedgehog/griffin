package res

import "project/trigger/github.com/ehedgehog/griffin/project/trigger/rdf"

type Resource struct {
	graph rdf.Graph
	term  rdf.Term
}

type Terms struct {
	graph rdf.Graph
	terms []rdf.Term
}

func (r Resource) Objects(predicate rdf.Term) Terms {
	terms := []rdf.Term{}
	pt := rdf.Triple{r.term, predicate, rdf.ANY}
	r.graph.FindAll(pt, func(t rdf.Triple) bool { terms = append(terms, t.O); return false })
	return Terms{r.graph, terms}
}

func (t Terms) Objects(predicate rdf.Term) Terms {
	terms := []rdf.Term{}
	for _, tt := range t.terms {
		pt := rdf.Triple{tt, predicate, rdf.ANY}
		t.graph.FindAll(pt, func(t rdf.Triple) bool { terms = append(terms, t.O); return false })
	}
	return Terms{t.graph, terms}
}
