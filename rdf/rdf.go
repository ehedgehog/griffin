package rdf

// TODO
// implementations and Symbol
// checking of IRI spelling
// protection of Triple fields
// Consider *Triple rather than Triple
// BNode channel for safe global bnode allocation?

const XSD_integer = IRI("http://www.w3.org/2001/XMLSchema#integer")
const XSD_boolean = IRI("http://www.w3.org/2001/XMLSchema#boolean")
const XSD_decimal = IRI("http://www.w3.org/2001/XMLSchema#decimal")
const XSD_double = IRI("http://www.w3.org/2001/XMLSchema#double")

const True = Boolean(true)
const False = Boolean(false)

const ANY = Var("??")

var EVERY = Triple{ANY, ANY, ANY}

type TermType int

const (T_NONE TermType = iota; T_IRI; T_BNODE; T_LITERAL; T_VAR; T_SYMBOL)

type Term interface {
    Spelling() string
    Language() string
    DatatypeIRI() IRI
    Value() interface{}
    Type() TermType
}

type Triple struct {S, P, O Term}

type Graph interface {
    Add(t Triple) bool
    Remove(t Triple) bool
    Size() uint64
    FindAll(t Triple, f func (t Triple) bool) bool
    Profile() Profile
}

type Profile interface {
}

type Prefixes interface {
    Shorten(u IRI) string
    SetPrefix(prefix string, u IRI)
    SetAll(other Prefixes, overwriteExisting bool)
}

