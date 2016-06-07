package rdf

type BaseTerm struct {}

func (b BaseTerm) Spelling() string {
    return ""
}

func (b BaseTerm) Language() string {
    return ""
}

func (b BaseTerm) DatatypeIRI() IRI {
    return AsIRI("")
}

func (b BaseTerm) Value() interface{} {
    return nil 
}

func (b BaseTerm) Type() TermType {
    return T_NONE
}

type SpellingTerm struct { BaseTerm; spelling string }

func (s SpellingTerm) Spelling () string {
    return s.spelling
}


