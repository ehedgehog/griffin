package turtle

// TODO
//
// lexing from streams
// error streams and/or error tokens?
// tracking line/col position

// import "unicode/utf8"
import "io"
import "strings"
import "bufio"

const eof rune = 0

type tokenType int

const (
	tok_ILL tokenType = iota
	tok_EOF
	tok_IRI
	tok_INTEGER
	tok_DECIMAL
	tok_DOUBLE
	tok_STRING
	tok_BNODE
	tok_DOT
	tok_COMMA
	tok_SEMI
	tok_LBOX
	tok_RBOX
	tok_LPAR
	tok_RPAR
	tok_DATATYPE
	tok_PREFIX
	tok_BASE
	tok_LANG
	tok_NAME
	tok_A
	tok_TRUE
	tok_FALSE
	tok_QNAME
	tok_COMMENT
)

func (k tokenType) String() string {
	switch k {
	case tok_ILL:
		return "illegal"
	case tok_TRUE:
		return "true"
	case tok_FALSE:
		return "false"
	case tok_A:
		return "A"
	case tok_NAME:
		return "NAME"
	case tok_EOF:
		return "EOF"
	case tok_DOT:
		return "DOT"
	case tok_COMMA:
		return "COMMA"
	case tok_SEMI:
		return "SEMI"
	case tok_LPAR:
		return "LPAR"
	case tok_RPAR:
		return "RPAR"
	case tok_LBOX:
		return "LBOX"
	case tok_RBOX:
		return "RBOX"
	case tok_QNAME:
		return "QNAME"
	case tok_COMMENT:
		return "#COMMENT"
	case tok_IRI:
		return "IRI"
	case tok_BASE:
		return "@BASE"
	case tok_PREFIX:
		return "@PREFIX"
	case tok_LANG:
		return "LANG"
	case tok_DATATYPE:
		return "DATATYPE"
	case tok_INTEGER:
		return "integer"
	case tok_DECIMAL:
		return "decimal"
	case tok_DOUBLE:
		return "double"
	case tok_STRING:
		return "string"
	case tok_BNODE:
		return "bnode"
	default:
		return "!unknown"
	}
	panic("unreachable")
}

type token struct {
	k        tokenType
	spelling string
	where    Location
}

type context struct {
	source       *bufio.Reader
	sink         chan token
	accumulating []rune
	latest       rune
	backedup     bool
	quote        rune
	where        Location
}

type stateFn func(*context) stateFn

func (c *context) run() {
	for state := initialState; state != nil; state = state(c) {
	}
	close(c.sink)
}

func (c *context) backup() {
	c.backedup = true
}

func (c *context) clearSpelling() {
	c.accumulating = c.accumulating[0:0]
}

func (c *context) append(r rune) {
	c.accumulating = append(c.accumulating, r)
}

func (c *context) next() rune {
	if c.backedup {
		c.backedup = false
		return c.latest
	}

	r, width, err := c.source.ReadRune()

	_ = width

	if err == io.EOF {
		c.latest = eof
		return eof
	}

	// r, width := utf8.DecodeRuneInString(c.source)
	// c.source = c.source[width:]

	c.latest = r
	c.backedup = false
	if r == '\n' {
		c.where.line += 1
		c.where.col = 1
	} else {
		c.where.col += 1
	}
	return r
}

func lexFromString(source string) (*context, chan token) {
	return lexFromReader(strings.NewReader(source))
}

func lexFromReader(source io.Reader) (*context, chan token) {
	c := &context{source: bufio.NewReader(source),
		sink:  make(chan token, 2),
		where: Location{1, 0}}
	go c.run()
	return c, c.sink
}

func (c *context) send(t tokenType) {
	c.sink <- token{t, string(c.accumulating), c.where}
}

func (c *context) sendWith(t tokenType, spelling string) {
	c.sink <- token{t, spelling, c.where}
}

func initialState(c *context) stateFn {
	ch := c.next()
	switch ch {
	case '<':
		return lexIRI
	case '_':
		c.clearSpelling()
		c.append(ch)
		ch = c.next()
		if ch == ':' {
			c.append(ch)
			ch = c.next()
			for IsBnodeLabelChar(ch) {
				c.append(ch)
				ch = c.next()
			}
		}
		if len(c.accumulating) < 3 {
			c.send(tok_ILL)
		} else {
			c.send(tok_BNODE)
		}
		c.backup()
	case '+', '-':
		c.clearSpelling()
		c.append(ch)
		return lexNumber
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		c.clearSpelling()
		c.backup()
		return lexNumber
	case '#':
		return lexComment
	case '@':
		return lexAt
	case eof:
		c.sendWith(tok_EOF, "")
		return nil
	case ':':
		return lexColon
	case '^':
		return lexHat
	case '.':
		c.clearSpelling()
		c.append(ch)
		ch = c.next()
		count := 0
		for '0' <= ch && ch <= '9' {
			count += 1
			c.append(ch)
			ch = c.next()
		}
		if count == 0 {
			c.sendWith(tok_DOT, ".")
		} else if ch == 'e' || ch == 'E' {
			c.append(ch)
			return lexExponent
		} else {
			c.send(tok_DECIMAL)
		}
		c.backup()
	case ' ', '\t', '\n':
		// TODO track character positions
	//	if ch == '\n' {
	//		c.where.line += 1
	//	}
	case '"', '\'':
		c.quote = ch
		return lexString
	case '[':
		c.sendWith(tok_LBOX, "[")
	case ']':
		c.sendWith(tok_RBOX, "]")
	case '(':
		c.sendWith(tok_LPAR, "(")
	case ')':
		c.sendWith(tok_RPAR, ")")
	case ';':
		c.sendWith(tok_SEMI, ";")
	case ',':
		c.sendWith(tok_COMMA, ",")
	default:
		if 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' {
			c.backup()
			return lexName
		} else {
			c.sendWith(tok_ILL, string([]rune{ch}))
		}
	}
	return initialState
}

func IsBnodeLabelChar(ch rune) bool {
	if 'a' <= ch && ch <= 'z' {
		return true
	}
	if 'A' <= ch && ch <= 'Z' {
		return true
	}
	if '0' <= ch && ch <= '9' {
		return true
	}
	//
	if ch == '_' {
		return true
	}
	// TODO: . except at beginning or end;  -, \uB7, \u300 to \u36F and \u203F to 2040 not at start.
	//
	return false
}

func lexNumber(c *context) stateFn {
	for {
		ch := c.next()
		if '0' <= ch && ch <= '9' {
			c.append(ch)
		} else if ch == '.' {
			c.append(ch)
			ch = c.next()
			for '0' <= ch && ch <= '9' {
				c.append(ch)
				ch = c.next()
			}
			if ch == 'e' || ch == 'E' {
				c.append(ch)
				return lexExponent
			} else {
				c.send(tok_DECIMAL)
				c.backup()
				return initialState
			}
		} else if ch == 'e' || ch == 'E' {
			c.append(ch)
			return lexExponent
		} else {
			c.backup()
			c.send(tok_INTEGER)
			return initialState
		}
	}
	panic("unreachable")
}

func lexExponent(c *context) stateFn {
	ch := c.next()
	if ch == '+' || ch == '-' {
		c.append(ch)
		ch = c.next()
	}
	count := 0
	for '0' <= ch && ch <= '9' {
		count += 1
		c.append(ch)
		ch = c.next()
	}
	if count == 0 {
		c.send(tok_ILL)
	} else {
		c.send(tok_DOUBLE)
	}
	c.backup()
	return initialState
}

func lexComment(c *context) stateFn {
	c.clearSpelling()
	c.append('#')
	for {
		ch := c.next()
		if ch != eof {
			c.append(ch)
		}
		if ch == eof || ch == '\n' {
			break
		}
	}
	c.send(tok_COMMENT)
	return initialState
}

func lexColon(c *context) stateFn {
	c.clearSpelling()
	c.append(':')
	for {
		ch := c.next()
		if 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || '0' <= ch && ch <= '9' || ch == '_' || ch == '-' {
			c.append(ch)
		} else {
			break
		}
	}
	// c.backup()
	c.send(tok_QNAME)
	return initialState
}

func lexIRI(c *context) stateFn {
	c.accumulating = c.accumulating[0:0]
	for {
		r := c.next()
		if r == '>' || r == eof {
			break
		}
		c.accumulating = append(c.accumulating, r)
	}
	c.send(tok_IRI)
	return initialState
}

func lexAt(c *context) stateFn {
	c.clearSpelling()
	ch := c.next()
	for 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '-' {
		c.append(ch)
		ch = c.next()
	}
	c.backup()
	spelling := string(c.accumulating)
	switch spelling {
	case "prefix":
		c.sendWith(tok_PREFIX, "prefix")
	case "base":
		c.sendWith(tok_BASE, "base")
	default:
		c.sendWith(tok_LANG, spelling)
	}
	return initialState
}

func lexHat(c *context) stateFn {
	ch := c.next()
	if ch == '^' {
		c.sendWith(tok_DATATYPE, "^^")
	} else {
		c.sendWith(tok_ILL, "^")
		c.backup()
	}
	return initialState
}

func lexString(c *context) stateFn {
	// We have just eaten the initial c.quote
	c.clearSpelling()
	ch2 := c.next()
	if ch2 == c.quote {
		// two quotes together. A third makes it a fat string,
		// otherwise it's an empty string.
		ch3 := c.next()
		if ch3 == c.quote {
			_ = fatStringBody(c)
			return initialState
		} else {
			c.backup()
			c.sendWith(tok_STRING, "")
			return initialState
		}
	} else {
		// quote not followed by quote: thin string
		c.backup()
		_ = stringBody(c)
		c.send(tok_STRING)
		return initialState
	}
	panic("notreachable")
}

func fatStringBody(c *context) (ok bool) {
	for {
		ok := stringBody(c)
		if !ok {
			return false
		}
		// We've just hit a c.quote, so we may be ending the string
		ch2 := c.next()
		if ch2 == c.quote {
			ch3 := c.next()
			if ch3 == c.quote {
				c.send(tok_STRING)
				return true
			} else {
				c.append(c.quote)
				c.append(c.quote)
				c.append(ch3)
			}
		} else {
			c.append(c.quote)
			c.append(ch2)
		}
	}
	panic("unreachable")
}

func stringBody(c *context) (ok bool) {
	for {
		ch := c.next()
		if ch == eof {
			c.sendWith(tok_ILL, "eof in string")
			return false
		}
		if ch == c.quote {
			return true
		}
		if ch == '\\' {
			ch2 := c.next()
			switch ch2 {
			case eof:
				c.sendWith(tok_ILL, "eof after backslash")
				return false
			// TODO b, f
			case 'n':
				c.append('\n')
			case 'r':
				c.append('\r')
			case '\\':
				c.append('\\')
			case '"':
				c.append('"')
			case 't':
				c.append('\t')
			case '\'':
				c.append('\'')
			case 'u':
				c.append(lexU(c, 4))
			case 'U':
				c.append(lexU(c, 8))
			default:
				c.sendWith(tok_ILL, "illegal escape after backslash")
				return false
			}
		} else {
			c.append(ch)
		}
	}
	panic("unreachable")
}

func lexU(c *context, count int) rune {
	r := 0
	for i := 0; i < count; i += 1 {
		r = r<<4 + unhex(c.next())
	}
	return rune(r)
}

func unhex(ch rune) int {
	if '0' <= ch && ch <= '9' {
		return int(ch) - '0'
	}
	if 'a' <= ch && ch <= 'z' {
		return int(ch) - 'a' + 10
	}
	if 'A' <= ch && ch <= 'Z' {
		return int(ch) - 'A' + 10
	}
	panic("oops, non-hex digit")
}

func lexPostColon(c *context) stateFn {
	for {
		ch := c.next()
		if 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || '0' <= ch && ch <= '9' || ch == '_' || ch == '-' {
			c.append(ch)
		} else {
			c.backup()
			c.send(tok_QNAME)
			return initialState
		}
	}
	panic("unreachable")
}

func lexName(c *context) stateFn {
	c.clearSpelling()
	for {
		ch := c.next()
		if 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || '0' <= ch && ch <= '9' || ch == '_' || ch == '-' {
			c.append(ch)
		} else if ch == ':' {
			c.append(ch)
			return lexPostColon
		} else {
			spelling := string(c.accumulating)
			kind := tok_NAME
			if strings.EqualFold(spelling, "true") {
				kind = tok_TRUE
				spelling = "true"
			} else if strings.EqualFold(spelling, "false") {
				kind = tok_FALSE
				spelling = "false"
			} else if strings.EqualFold(spelling, "a") {
				kind = tok_A
				spelling = "a"
			}
			c.sendWith(kind, spelling)
			c.backup()
			return initialState
		}
	}
	return initialState
}
