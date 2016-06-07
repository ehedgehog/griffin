package rdf

import "strconv"

func AsInteger(spelling string) Integer {
    return Integer{SpellingTerm{BaseTerm{}, spelling}}
}

func AsLanguagedString(spelling, language string) LanguagedString {
    return LanguagedString{spelling, language}
}

func AsTypedLiteral(spelling string, datatype IRI) TypedLiteral {
    return TypedLiteral{spelling, datatype}
}

func AsDecimal(spelling string) TypedLiteral {
		return TypedLiteral{spelling, XSD_decimal}
}

func AsDouble(spelling string) TypedLiteral {
	return TypedLiteral{spelling, XSD_double}
}

// ---------------------------------------------------------

type TypedLiteral struct { spelling string; dataType IRI }

func (t TypedLiteral) Spelling() string {
    return t.spelling
}

func (t TypedLiteral) Value() interface {} {
    return "TO BE DONE"
}

func (t TypedLiteral) Language() string {
    return ""
}

func (t TypedLiteral) DatatypeIRI() IRI {
    return t.dataType
}

func (it TypedLiteral) Type() TermType {
    return T_LITERAL
}

// ------------------------------------------------------------


type LanguagedString struct { spelling string; language string }

func (l LanguagedString) Spelling() string {
    return l.spelling
}

func (t LanguagedString) DatatypeIRI() IRI {
    return AsIRI("")
    }

func (l LanguagedString) Type() TermType {
    return T_LITERAL
}

func (l LanguagedString) Language() string {
    return l.language
}

func (l LanguagedString) Value() interface{} {
    return l.spelling
}

// -------------------------------------------------------------

type PlainString string 

func (p PlainString) Spelling() string {
    return string(p)
}

func (p PlainString) Type() TermType {
    return T_LITERAL
}

func (p PlainString) Value() interface{} {
    return string(p)
}

func (p PlainString) DatatypeIRI() IRI {
    return AsIRI("")
}

func (p PlainString) Language() string {
    return ""
}

// ---------------------------------------------------------------

type Integer struct { SpellingTerm }

func (i Integer) Type() TermType {
    return T_LITERAL
}

func (i Integer) Value() interface {} {
    value, err := strconv.Atoi(i.spelling)
    _ = err // TODO
    return value
}

func (i Integer) DatatypeIRI() IRI {
    return XSD_integer
}

// -----------------------------------------------------------------

type Boolean bool

func (b Boolean) Type() TermType {
	return T_LITERAL
}

func (b Boolean) DatatypeIRI() IRI {
	return XSD_boolean
}

func (b Boolean) Value() interface{} {
	return bool(b)
}

func (b Boolean) Spelling() string {
    if bool(b) { return "true" }
    return "false"
}

func (b Boolean) Language() string {
    return ""
}

