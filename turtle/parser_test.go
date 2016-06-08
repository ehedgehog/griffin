package turtle

import "testing"
import "strconv"
import "fmt"
import "project/trigger/github.com/ehedgehog/griffin/project/trigger/rdf"

// TODO
// names, variables
// error reporting

func TestTurtleParser(t *testing.T) {
	parseTest(t, "")
	parseTest(t, "<spoo> <flarn> <pigsty>", rdf.Triple{rdf.AsIRI("spoo"), rdf.AsIRI("flarn"), rdf.AsIRI("pigsty")})
	A, B, C := rdf.AsIRI("a"), rdf.AsIRI("b"), rdf.AsIRI("c")
	D, E, F := rdf.AsIRI("d"), rdf.AsIRI("e"), rdf.AsIRI("f")
	parseTest(t, "<a> <b> <c>. <d> <e> <f>", rdf.Triple{A, B, C}, rdf.Triple{D, E, F})
	G, H, I := rdf.PlainString("hello"), rdf.PlainString("world"), rdf.PlainString("smudge content")
	_ = I
	P, Q, R := rdf.AsIRI("p"), rdf.AsIRI("q"), rdf.AsIRI("r")
	parseTest(t, `<a> <b> "hello"`, rdf.Triple{A, B, G})
	parseTest(t, `<a> <b> "world"`, rdf.Triple{A, B, H})
	//
	parseTest(t, `<a> <p> <b>; <q> <c>.`, rdf.Triple{A, P, B}, rdf.Triple{A, Q, C})
	parseTest(t, `<a> <p> <b>; <q> <c>; <r> <d>.`, rdf.Triple{A, P, B}, rdf.Triple{A, Q, C}, rdf.Triple{A, R, D})
	//
	parseTest(t, `<a> <p> <b>, <c>.`, rdf.Triple{A, P, B}, rdf.Triple{A, P, C})
	parseTest(t, `<a> <p> <b>; <q> <c>, <d>`, rdf.Triple{A, P, B}, rdf.Triple{A, Q, C}, rdf.Triple{A, Q, D})
}

func TestTurtlePrefixes(t *testing.T) {
	PA, PB, PC := rdf.AsIRI("eh:/PRE/A"), rdf.AsIRI("eh:/PRE/B"), rdf.AsIRI("eh:/PRE/C")
	parseTest(t, `@prefix pre: <eh:/PRE/>. pre:A pre:B pre:C.`, rdf.Triple{PA, PB, PC})
	SA, SB, SC := rdf.AsIRI("eh:/short/A"), rdf.AsIRI("eh:/short/B"), rdf.AsIRI("eh:/short/C")
	parseTest(t, `@prefix : <eh:/short/>. :A :B :C`, rdf.Triple{SA, SB, SC})
}

func TestBase(t *testing.T) {
	S, P, O := rdf.AsIRI("eh:/me/S"), rdf.AsIRI("eh:/me/other/P"), rdf.AsIRI("eh:/me/O")
	parseTest(t, `@base <eh:/me/>. <S> <other/P> <eh:/me/O>`, rdf.Triple{S, P, O})
}

func TestTurtleBNodes(t *testing.T) {
	b := rdf.AsBNode("_:B1001")
	parseTest(t, "[] a <eh:/Type>. ", rdf.Triple{b, rdf.Type, rdf.AsIRI("eh:/Type")})
	parseTest(t, "[<eh:/P> 17] a <eh:/Type>. ", rdf.Triple{b, rdf.AsIRI("eh:/P"), rdf.AsInteger("17")}, rdf.Triple{b, rdf.Type, rdf.AsIRI("eh:/Type")})
}

func TestLiteralBNodes(t *testing.T) {
	b := rdf.AsBNode("_:mehere")
	parseTest(t, "_:mehere a <eh:/Type>. ", rdf.Triple{b, rdf.Type, rdf.AsIRI("eh:/Type")})
}

func TestTurtleLists(t *testing.T) {
	S, P := rdf.AsIRI("eh:/S"), rdf.AsIRI("eh:/P")
	A, B := rdf.AsIRI("eh:/A"), rdf.AsIRI("eh:/B")
	b1, b2 := rdf.AsBNode("_:B1001"), rdf.AsBNode("_:B1002")
	parseTest(t, `<eh:/S> <eh:/P> () `, rdf.Triple{S, P, rdf.Nil})
	parseTest(t, `<eh:/S> <eh:/P> (<eh:/A>) `,
		rdf.Triple{b1, rdf.First, A},
		rdf.Triple{b1, rdf.Rest, rdf.Nil},
		rdf.Triple{S, P, b1})
	parseTest(t, `<eh:/S> <eh:/P> (<eh:/A> <eh:/B>) `,
		rdf.Triple{b1, rdf.First, A},
		rdf.Triple{b1, rdf.Rest, b2},
		rdf.Triple{b2, rdf.First, B},
		rdf.Triple{b2, rdf.Rest, rdf.Nil},
		rdf.Triple{S, P, b1})
}

func TestNonplainLiterals(t *testing.T) {
	S, P := rdf.AsIRI("eh:/S"), rdf.AsIRI("eh:/P")
	parseTest(t, `<eh:/S> <eh:/P> "chat"@fr`, rdf.Triple{S, P, rdf.AsLanguagedString("chat", "fr")})
}

func TestTypedLiterals(t *testing.T) {
	S, P := rdf.AsIRI("eh:/S"), rdf.AsIRI("eh:/P")
	parseTest(t, `<eh:/S> <eh:/P> "chat"^^<eh:/datatype>`, rdf.Triple{S, P, rdf.AsTypedLiteral("chat", rdf.AsIRI("eh:/datatype"))})
	parseTest(t, `<eh:/S> <eh:/P> true.`, rdf.Triple{S, P, rdf.True})
	parseTest(t, `<eh:/S> <eh:/P> false.`, rdf.Triple{S, P, rdf.False})
}

func TestNumericLiterals(t *testing.T) {
	S, P := rdf.AsIRI("eh:/S"), rdf.AsIRI("eh:/P")
	parseTest(t, `<eh:/S> <eh:/P> 1.`, rdf.Triple{S, P, rdf.AsDecimal("1.")})
	parseTest(t, `<eh:/S> <eh:/P> .1`, rdf.Triple{S, P, rdf.AsDecimal(".1")})
	parseTest(t, `<eh:/S> <eh:/P> 1.2`, rdf.Triple{S, P, rdf.AsDecimal("1.2")})
	parseTest(t, `<eh:/S> <eh:/P> 1.e5`, rdf.Triple{S, P, rdf.AsDouble("1.e5")})
	parseTest(t, `<eh:/S> <eh:/P> .1e5`, rdf.Triple{S, P, rdf.AsDouble(".1e5")})
}

func TestLocationalThingies(t *testing.T) {
	reportTest(t, "@broken", Location{1, 7})
	reportTest(t, "\n@broken", Location{2, 8})
	reportTest(t, "\n\n@broken", Location{3, 8})
	reportTest(t, "@broken", Location{1, 7})
	reportTest(t, "<eh:/P> ).", Location{1, 9})
}

type Report struct {
	message string
	where   Location
}

type Triples struct {
	them       []rdf.Triple
	prefixes   map[string]string
	bnodeCount int
	base       rdf.IRI
	reports    []Report
}

func (ts *Triples) Add(t rdf.Triple) bool {
	ts.them = append(ts.them, t)
	return true
}

func (ts *Triples) SetPrefix(prefix string, ns rdf.IRI) {
	ts.prefixes[prefix] = ns.Spelling()
}

func (ts *Triples) SetBase(base rdf.IRI) {
	ts.base = base
}

func (ts *Triples) NewBNode() rdf.Term {
	ts.bnodeCount += 1
	s := rdf.AsBNode("_:B" + strconv.Itoa(ts.bnodeCount))
	return s
}

func (ts *Triples) Report(message string, where Location) {
	ts.reports = append(ts.reports, Report{message, where})
}

func reportTest(t *testing.T, source string, expect Location) {
	them := &Triples{[]rdf.Triple{}, map[string]string{}, 1000, rdf.AsIRI(""), []Report{}}
	ParseFromString(source, them)
	assert_equals(t, "one report", len(them.reports), 1)
	if false {
		fmt.Println("reports:", them.reports)
	}
	assert_equals(t, "the location", them.reports[0].where, expect)
}

func parseTest(t *testing.T, source string, triples ...rdf.Triple) {
	them := &Triples{[]rdf.Triple{}, map[string]string{}, 1000, rdf.AsIRI(""), []Report{}}
	ParseFromString(source, them)
	assertSameTriples(t, "triples", them.them, triples)
}

func assertSameTriples(t *testing.T, tag string, got, expected []rdf.Triple) {
	assert_equals(t, "#triples", len(got), len(expected))
	for i, g := range got {
		assert_equals(t, "a triple", g, expected[i])
	}
}

func assert_equals(t *testing.T, tag string, got, expected interface{}) {
	if got != expected {
		t.Fatalf("%v: expected %v, got %v", tag, expected, got)
	}
}
