package rdf

type BNode string

func (b BNode) Spelling() string {
    return string(b)
}

func (b BNode) Value() interface {} {
    return nil
}

func (b BNode) DatatypeIRI() IRI {
    return IRI("")
}

func (b BNode) Language() string {
    return ""
}

func (b BNode) Type() TermType {
    return T_BNODE
}

func AsBNode(spelling string) BNode {
    return BNode(spelling)
}
