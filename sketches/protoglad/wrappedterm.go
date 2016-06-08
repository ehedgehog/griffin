package main

import (
	"sort"
	"strings"

	"project/trigger/github.com/ehedgehog/griffin/project/trigger/rdf"
)

type WrappedTerm struct {
	Term  rdf.Term
	Graph rdf.Graph
}

func (wt WrappedTerm) Label() string {
	obs := ObjectsOf(wt.Graph, wt.Term, rdf.Label)
	if len(obs) > 0 {
		return obs[0].Spelling()
	}
	return wt.Term.Spelling()
}

func (wt WrappedTerm) Properties() []WrappedTerm {
	properties := []WrappedTerm{}
	propertySet := map[rdf.Term]bool{}
	wt.Graph.FindAll(rdf.Triple{wt.Term, rdf.ANY, rdf.ANY}, func(t rdf.Triple) bool {
		_, present := propertySet[t.P]
		if !present {
			propertySet[t.P] = true
			properties = append(properties, WrappedTerm{t.P, wt.Graph})
		}
		return false
	})
	sort.Sort(Sortable(properties))
	return properties
}

type Sortable []WrappedTerm

func (s Sortable) Len() int {
	return len(s)
}

func (s Sortable) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Sortable) Less(i, j int) bool {
	return s[i].Term.Spelling() < s[j].Term.Spelling()
}

func (wt WrappedTerm) Values(P WrappedTerm) []WrappedTerm {
	objects := []WrappedTerm{}
	wt.Graph.FindAll(rdf.Triple{wt.Term, P.Term, rdf.ANY}, func(t rdf.Triple) bool {
		objects = append(objects, WrappedTerm{t.O, wt.Graph})
		return false
	})
	return objects
}

// should check that is URI node
func (wt WrappedTerm) URI() string {
	return wt.Term.Spelling()
}

// should look up shortname, or label, whatever
// will use local name for now
func (wt WrappedTerm) Display() string {
	spelling := wt.Term.Spelling()
	cut := strings.LastIndexAny(spelling, "#/")
	return spelling[cut+1:]
}

func (wt WrappedTerm) HasValue() bool {
	result := false
	wt.Graph.FindAll(rdf.Triple{wt.Term, rdf.ANY, rdf.ANY},
		func(t rdf.Triple) bool {
			result = true
			return false
		})
	return result
}

func (wt WrappedTerm) Prettily() string {
	if wt.Term.Type() == rdf.T_IRI {
		return wt.Display()
	}
	return wt.Term.Spelling()
}

func (wt WrappedTerm) String() string {
	return wt.Term.Spelling()
}
