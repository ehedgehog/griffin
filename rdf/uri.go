package rdf

type IRI string

func (u IRI) Type() TermType {
    return T_IRI
}

func (u IRI) Language() string {
    return ""
}

func (u IRI) DatatypeIRI() IRI {
    return IRI("")
}

func (u IRI) Spelling() string {
    return string(u)
}

func (u IRI) Value() interface {} {
    return nil
}

func AsIRI(spelling string) IRI {
    return IRI(spelling)
}

