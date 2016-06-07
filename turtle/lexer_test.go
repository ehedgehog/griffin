package turtle

import "testing"

// TODO
// check allowed chars in qname/name
// unicode characters
// variables and symbols [trigger extensions]
// tabs and newlines
// newline in thin string body is an error
// spelling of language codes
// error reporting
// rationalise names

func TestNewItemiseFromChecks(t *testing.T) {
	for _, check := range checks {
		context, tokens := lexFromString(check.source)
		_ = context
		token := <-tokens
		if token.k == check.kind && token.spelling == check.spelling {
		} else {
			t.Errorf("%s: expected %v `%s` but got %v `%s` [source: %s]",
				check.tag,
				check.kind, check.spelling, token.k, token.spelling,
				check.source,
			)
		}
		// assertEquals(t, "rest", context.source, check.rest)
	}
}

func TestSemicolonFollowingLanguage(t *testing.T) {
	_, tokens := lexFromString("'hello'@world;")
	ensure(t, <-tokens, tok_STRING)
	ensure(t, <-tokens, tok_LANG)
	ensure(t, <-tokens, tok_SEMI)
	ensure(t, <-tokens, tok_EOF)
}

func ensure(t *testing.T, tok token, kind tokenType) {
	if tok.k != kind {
		t.Fatalf("expected %v, got %v\n", kind, tok.k)
	}
}

var checks = []struct {
	tag      string
	source   string
	kind     tokenType
	spelling string
}{
	{"EOF", "", tok_EOF, ""},
	{"IRI", "<anything>", tok_IRI, "anything"},
	{"IRI", "  <anything>", tok_IRI, "anything"},
	{"IRI", "<another>GUBBINS", tok_IRI, "another"},
	{"integer", "1", tok_INTEGER, "1"},
	{"integer", "  1", tok_INTEGER, "1"},
	{"integer", "17,", tok_INTEGER, "17"},
	{"integer", "+17", tok_INTEGER, "+17"},
	{"integer", "-17", tok_INTEGER, "-17"},
	{"decimal", `1.`, tok_DECIMAL, `1.`},
	{"decimal", `1.ten`, tok_DECIMAL, `1.`},
	{"decimal", `1.2`, tok_DECIMAL, `1.2`},
	{"decimal", `1.234`, tok_DECIMAL, `1.234`},
	{"decimal", `.2x`, tok_DECIMAL, `.2`},
	{"decimal", `.25`, tok_DECIMAL, `.25`},
	{"decimal", `-2.1`, tok_DECIMAL, `-2.1`},
	{"decimal", `+2.1`, tok_DECIMAL, `+2.1`},
	{"decimal", `-.1`, tok_DECIMAL, `-.1`},
	{"decimal", `+.1`, tok_DECIMAL, `+.1`},
	{"double", `1e5`, tok_DOUBLE, `1e5`},
	{"double", `1e512`, tok_DOUBLE, `1e512`},
	{"double", `1E+5`, tok_DOUBLE, `1E+5`},
	{"double", `1e-5`, tok_DOUBLE, `1e-5`},
	{"double", `1e+17`, tok_DOUBLE, `1e+17`},
	{"double", `1.2e3`, tok_DOUBLE, `1.2e3`},
	{"double", `.2e3`, tok_DOUBLE, `.2e3`},
	{"double", `1.e3`, tok_DOUBLE, `1.e3`},
	{"double", `1.2e+3`, tok_DOUBLE, `1.2e+3`},
	{"double", `-1.2e3`, tok_DOUBLE, `-1.2e3`},
	{"integer", "21]", tok_INTEGER, "21"},
	{"dot", ".ZOO", tok_DOT, "."},
	{"dot", "  .OOZE", tok_DOT, "."},
	{"comma", ",", tok_COMMA, ","},
	{"semi", ";", tok_SEMI, ";"},
	{"lbox", "[", tok_LBOX, "["},
	{"rbox", "]", tok_RBOX, "]"},
	{"lpar", "(", tok_LPAR, "("},
	{"rpar", ")", tok_RPAR, ")"},
	{"@prefix", "@prefix", tok_PREFIX, "prefix"},
	{"@prefix", "@prefixed", tok_LANG, "prefixed"},
	{"@base", "@base", tok_BASE, "base"},
	{"@base", "@based", tok_LANG, "based"},
	{"@lang", "@en-uk", tok_LANG, "en-uk"},
	{"@lang", "@en-fr+16", tok_LANG, "en-fr"},
	{"datatype", "^^", tok_DATATYPE, "^^"},
	{"datatype", "^^^", tok_DATATYPE, "^^"},
	{"datatype", "^^<", tok_DATATYPE, "^^"},
	{"boolean", "true", tok_TRUE, "true"},
	{"boolean", "TRUE", tok_TRUE, "true"},
	{"boolean", "TrUe", tok_TRUE, "true"},
	{"boolean", "false", tok_FALSE, "false"},
	{"boolean", "FaLsE", tok_FALSE, "false"},
	{"a", "a", tok_A, "a"},
	{"unboolean", "truely", tok_NAME, "truely"},
	{"una", "and", tok_NAME, "and"},
	{"test", "pre:post", tok_QNAME, "pre:post"},
	{"test", "pre:post@stuff", tok_QNAME, "pre:post"},
	{"test", ":post", tok_QNAME, ":post"},
	{"test", "a:", tok_QNAME, "a:"},
	{"bnode", "_:label", tok_BNODE, "_:label"},
	{"bnode", "_:porkpie,", tok_BNODE, "_:porkpie"},
	{"bnode", "_:B1004", tok_BNODE, "_:B1004"},
	{"unbnode", "_:)", tok_ILL, "_:"},
	{"unbnode", "_spoo", tok_ILL, "_"},
	{"test", "# stuff", tok_COMMENT, "# stuff"},
	{"test", "# staff\n", tok_COMMENT, "# staff\n"},
	{"test", "# stiff\nffits", tok_COMMENT, "# stiff\n"},
	{"test", "\n\nhello\n", tok_NAME, "hello"},
	{"test", `"string"`, tok_STRING, `string`},
	{"test", `'string'`, tok_STRING, `string`},
	{"test", `'s\`, tok_ILL, `eof after backslash`},
	{"test", `"A\nB"`, tok_STRING, "A\nB"},
	{"test", `"A\n"`, tok_STRING, "A\n"},
	{"test", `"\nB"`, tok_STRING, "\nB"},
	{"test", `"\n"`, tok_STRING, "\n"},
	{"test", `'\n\r\\\t\''`, tok_STRING, "\n\r\\\t'"},
	{"test", `'\u10ba!'.`, tok_STRING, "\u10ba!"},
	{"test", `"""AA"BB""CC"""D`, tok_STRING, `AA"BB""CC`},
	{"test", "'''hello\nworld'''", tok_STRING, "hello\nworld"},
}

func assertEquals(t *testing.T, tag string, got, expected interface{}) {
	if got != expected {
		t.Fatalf("%v: expected %v, got %v", tag, expected, got)
	}
}
