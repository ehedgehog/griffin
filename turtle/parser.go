package turtle

// TODO
// The name `Reader` is wrong. Fix it. The comment phrasing
// is off, too.
//
// decent error messages (incl line/col number)

import "code.google.com/p/project/trigger/rdf"
import "fmt"
import "strings"
import "io"

type Reporter interface {
	Report(message string, where Location)
}

type Location struct {
	line, col int
}

// TODO rename
// A Reader provides the output context for a Turtle
// reader.
type Consumer interface {
	Reporter
	// Supply a parsed triple to the requesting envionment
	Add(rdf.Triple) bool
	// Set the current base for relative IRIs.
	SetBase(base rdf.IRI)
	// Set the given prefix to be bound to the given IRI
	SetPrefix(prefix string, ns rdf.IRI)
	// Deliver a new bnode different from any other this
	// Reader instance has delivered.
	NewBNode() rdf.Term
}

type parsingState struct {
	items    chan (token)
	kind     tokenType
	spelling string
	where    Location
	prefixes map[string]string
	base     rdf.IRI
}

func (l *parsingState) SetBase(base rdf.IRI) {
	l.base = base
}

func (l *parsingState) Resolve(spelling string) rdf.IRI {
	if isAbsolute(spelling) {
		return rdf.AsIRI(spelling)
	}
	return rdf.AsIRI(string(l.base) + spelling)
}

func isAbsolute(spelling string) bool {
	for _, ch := range spelling {
		if ch == ':' {
			return true
		}
		if 'a' <= ch && ch <= 'z' {
			continue
		}
		if 'A' <= ch && ch <= 'Z' {
			continue
		}
		if '0' <= ch && ch <= '9' {
			continue
		}
		if ch == '.' || ch == '-' || ch == '_' {
			continue
		}
		return false
	}
	return false
}

func (l *parsingState) Expand(qname string) (spelling string, e error) {
	colon := strings.IndexRune(qname, ':')
	prefix := qname[:colon+1]
	ns := l.prefixes[prefix]
	return ns + qname[colon+1:], nil
}

func (l *parsingState) advance() {
	for {
		item := <-l.items
		_ = item
		l.kind, l.spelling, l.where = item.k, item.spelling, item.where
		if false {
			fmt.Println(l.kind, l.spelling)
		}
		if l.kind != tok_COMMENT {
			break
		}
	}
}

// Parse parses the source string using the sink Reader
// to deliver its triples and prefixes.
func ParseFromString(source string, sink Consumer) {
	_, items := lexFromString(source)
	ParseFromLexer(items, sink)
}

func ParseFromReader(source io.Reader, sink Consumer) {
	_, items := lexFromReader(source)
	ParseFromLexer(items, sink)
}

func ParseFromLexer(items chan token, sink Consumer) {
	x := &parsingState{items: items, prefixes: make(map[string]string)}
	x.advance()
	parseIt(x, sink)
}

func parseIt(l *parsingState, sink Consumer) {
	for startsTerm(l.kind) || l.kind == tok_PREFIX || l.kind == tok_BASE {
		if l.kind == tok_PREFIX {
			// TODO the checking
			l.advance()
			prefix := l.spelling
			l.advance()
			term := parseNode(l, sink)
			if l.kind == tok_DOT {
				l.advance()
			} else {
				panic("expected tok_DOT for @prefix.")
			}
			l.prefixes[prefix] = term.Spelling()
			sink.SetPrefix(prefix, rdf.AsIRI(term.Spelling()))
		} else if l.kind == tok_BASE {
			l.advance()
			l.SetBase(rdf.AsIRI(l.spelling))
			sink.SetBase(rdf.AsIRI(l.spelling)) // TODO check it's an IRI!
			l.advance()
			if l.kind == tok_DOT {
				l.advance()
			} else {
				panic("expected tok_DOT for @base.")
			}
		} else {
			s := parseNode(l, sink)
			parsePredications(l, s, sink)
			if l.kind == tok_DOT {
				l.advance()
			}
		}
	}
	if l.kind != tok_EOF {
		sink.Report("should have exhausted input: "+l.kind.String()+": "+l.spelling, l.where)
	}
}

func startsTerm(k tokenType) bool {
	switch k {
	case tok_IRI, tok_QNAME, tok_INTEGER, tok_STRING, tok_TRUE, tok_FALSE,
		tok_LPAR, tok_LBOX, tok_A, tok_DECIMAL, tok_DOUBLE, tok_BNODE:
		return true
	}
	return false
}

func parsePredications(l *parsingState, s rdf.Term, sink Consumer) {
	for {
		if !startsTerm(l.kind) {
			return
		}
		p := parseNode(l, sink)
		parseObjects(l, s, p, sink)
		if l.kind != tok_SEMI {
			return
		}
		l.advance()
	}
}

func parseObjects(l *parsingState, s, p rdf.Term, sink Consumer) {
	for {
		sink.Add(rdf.Triple{s, p, parseNode(l, sink)})
		if l.kind != tok_COMMA {
			return
		}
		l.advance()
	}
}

func parseNode(l *parsingState, sink Consumer) rdf.Term {
	switch l.kind {
	case tok_IRI:
		result := l.Resolve(l.spelling)
		l.advance()
		return result

	case tok_LPAR:
		l.advance()
		if l.kind == tok_RPAR {
			l.advance()
			return rdf.Nil
		} else {
			result := sink.NewBNode()
			here := result
			for {
				this := parseNode(l, sink)
				sink.Add(rdf.Triple{here, rdf.First, this})
				if !startsTerm(l.kind) {
					break
				}
				next := sink.NewBNode()
				sink.Add(rdf.Triple{here, rdf.Rest, next})
				here = next
			}
			if l.kind != tok_RPAR {
				sink.Report("missing )", l.where)
			} else {
				l.advance()
			}
			sink.Add(rdf.Triple{here, rdf.Rest, rdf.Nil})
			return result
		}

	case tok_LBOX:
		l.advance()
		s := sink.NewBNode()
		parsePredications(l, s, sink)
		if l.kind == tok_RBOX {
			l.advance()
		} else {
			sink.Report("missing ]", l.where)
		}
		return s

	case tok_QNAME:
		resolved, err := l.Expand(l.spelling)
		if err != nil {
			sink.Report("Could not expand "+l.spelling, l.where)
		} // TODO
		result := rdf.AsIRI(resolved)
		l.advance()
		return result

	case tok_A:
		// TODO
		result := rdf.Type
		l.advance()
		return result

	case tok_INTEGER:
		result := rdf.AsInteger(l.spelling)
		l.advance()
		return result

	case tok_STRING:
		spelling := l.spelling
		var result rdf.Term = rdf.PlainString(spelling)
		l.advance()
		if l.kind == tok_LANG {
			result = rdf.AsLanguagedString(spelling, l.spelling)
			l.advance()
		} else if l.kind == tok_DATATYPE {
			l.advance()
			// TODO only parse qname or iri here, and avoid
			// type assertion.
			t := parseNode(l, sink)
			return rdf.AsTypedLiteral(spelling, t.(rdf.IRI))
		}
		return result

	case tok_TRUE:
		l.advance()
		return rdf.True

	case tok_FALSE:
		l.advance()
		return rdf.False

	case tok_DECIMAL:
		spelling := l.spelling
		l.advance()
		return rdf.AsDecimal(spelling)

	case tok_DOUBLE:
		spelling := l.spelling
		l.advance()
		return rdf.AsDouble(spelling)

	case tok_BNODE:
		result := rdf.AsBNode(l.spelling)
		l.advance()
		return result

	default:
		cc := fmt.Sprintf("%v", l.where)
		panic("oops: can't handle this token: " + l.kind.String() + " " + l.spelling + " on line " + cc)
	}
	return nil
}
