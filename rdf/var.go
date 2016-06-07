package rdf

type Var string

func (v Var) Type() TermType { return T_VAR }

func (v Var) DatatypeIRI() IRI { return IRI("") }

func (v Var) Value() interface{} { return "" }

func (v Var) Spelling() string { return string(v) }

func (v Var) Language() string { return "" }

