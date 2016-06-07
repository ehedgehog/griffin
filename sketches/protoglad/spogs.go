package main

import "code.google.com/p/project/trigger/rdf"

func SubjectsWithProperty(g rdf.Graph, P, O rdf.Term) []rdf.Term {
	terms := []rdf.Term{}
	g.FindAll(rdf.Triple{rdf.ANY, P, O}, func(t rdf.Triple) bool {
		terms = append(terms, t.S)
		return false
	})
	return terms
}

func ObjectsOf(g rdf.Graph, S, P rdf.Term) []rdf.Term {
	result := []rdf.Term{}
	g.FindAll(rdf.Triple{S, P, rdf.ANY}, func(t rdf.Triple) bool {
		result = append(result, t.O)
		return false
	})
	return result
}
