package rdf

import "testing"

func TestIRIAsTerm(t *testing.T) {
    var u Term = AsIRI("eh:/hello")
    expectEqual(t, "uri.spelling", u.Spelling(), "eh:/hello")
    expectEqual(t, "uri.type", u.Type(), T_IRI)
    expectEqual(t, "uri.language", u.Language(), "")
    expectEqual(t, "uri.value", u.Value(), nil)
}

func TestBNodeAsTerm(t *testing.T) {
    var b Term = AsBNode("_:1829")
    expectEqual(t, "uri.spelling", b.Spelling(), "_:1829")
    expectEqual(t, "uri.type", b.Type(), T_BNODE)
    expectEqual(t, "uri.language", b.Language(), "")
    expectEqual(t, "uri.value", b.Value(), nil)
}

func TestIntegerAsTerm(t *testing.T) {
    var i Term = AsInteger("17")
    expectEqual(t, "integer.spelling", i.Spelling(), "17")
    expectEqual(t, "integer.type", i.Type(), T_LITERAL)
    expectEqual(t, "integer.language", i.Language(), "")
    expectEqual(t, "uri.value", i.Value(), 17)
}

func TestPlainStringAsTerm(t *testing.T) {
    var s Term = PlainString("world")
    expectEqual(t, "string.spelling", s.Spelling(), "world")
    expectEqual(t, "string.type", s.Type(), T_LITERAL)
    expectEqual(t, "string.language", s.Language(), "")
    expectEqual(t, "string.value", s.Value(), "world")
}

func TestLanguagedStringAsTerm(t *testing.T) {
    var l Term = AsLanguagedString("chat", "fr")
    expectEqual(t, "languaged.spelling", l.Spelling(), "chat")
    expectEqual(t, "languaged.type", l.Type(), T_LITERAL)
    expectEqual(t, "languaged.language", l.Language(), "fr")
    expectEqual(t, "languaged.value", l.Value(), "chat")
}

func TestTypedLiteralAsTerm(t *testing.T) {
    var v Term = AsTypedLiteral("spelling", AsIRI("eh:/Type"))
    expectEqual(t, "typed.spelling", v.Spelling(), "spelling")
    expectEqual(t, "typed.type", v.Type(), T_LITERAL)
    expectEqual(t, "typed.language", v.Language(), "")
    expectEqual(t, "typed.type", v.DatatypeIRI(), AsIRI("eh:/Type"))
}

func TestBooleanLiterals(t *testing.T) {
	expectEqual(t, "true.spelling", True.Spelling(), "true")
	expectEqual(t, "true.type", True.Type(), T_LITERAL)
	expectEqual(t, "true.Value", True.Value(), true)
	expectEqual(t, "true.language", True.Language(), "")
	expectEqual(t, "true.datatype", True.DatatypeIRI(), XSD_boolean)
//
    expectEqual(t, "false.spelling", False.Spelling(), "false")
	expectEqual(t, "false.type", False.Type(), T_LITERAL)
	expectEqual(t, "false.Value", False.Value(), false)
	expectEqual(t, "false.language", False.Language(), "")
	expectEqual(t, "false.datatype", False.DatatypeIRI(), XSD_boolean)
}

func TestSomething(t *testing.T) {
}

func expectEqual(t *testing.T, tag string, got, expected interface{}) {
    if got != expected {
        t.Fatalf("%v: expected `%v`, got `%v`", tag, expected, got)
    }
}
