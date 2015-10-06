package main

import (
	"log"
	"strings"
)

// An lexed item
type item struct {
	typ int
	val string
}

// A state function returns the next state function where lexing should continue
type stateFn func(*lexer) stateFn

// Removes comments from the input.
func filterComments(i string) string {
	out := make([]rune, 0)

	inString := false
	inComment := false
	for _, c := range i {
		switch string(c) {
			case "'":
				if !inComment {
					inString = !inString
				}
			case "#":
				if !inString {
					inComment = true
				}
				continue
			case "\n":
				if inComment {
					inComment = false
					continue
				}
		}
		if !inComment {
			out = append(out, c)
		}
	}

	return string(out)
}

type lexer struct {
	input    []string
	start    int
	pos      int
	items    chan item
	lastitem item
	debug		bool
}

// Create a new lexer which lexes the given input
func newLexer(input string, debug bool) *lexer {
	l := new(lexer)
	l.items = make(chan item, 0)
	l.input = strings.Fields(filterComments(input))
	l.debug = debug

	if debug {
		log.Println("Schema input with comments removed.")
		log.Println(filterComments(input))
	}

	// Start the lexing asynchronly
	go l.run()
	return l
}

func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
}

// Method to fulfill the yacc interface
func (l *lexer) Lex(lval *yySymType) int {
	item := <-l.items
	l.lastitem = item
	lval.val = item.val

	if l.debug {
		printToken(l.lastitem)
	}

	return item.typ
}

// dito
func (l *lexer) Error(e string) {
	printToken(l.lastitem)
	log.Println(e)
}

// print the token type and the token value for debugging
func printToken(i item) {
	typ := i.typ
	val := i.val
	switch typ {
	case ATTRIBUTETYPE:
		log.Print("attributetype ", val)
	case OBJECTCLASS:
		log.Print("objectclass ", val)
	case NAME:
		log.Print("name ", val)
	case DESC:
		log.Print("desc ", val)
	case OBSOLETE:
		log.Print("obsolete ", val)
	case SUP:
		log.Print("sup ", val)
	case EQUALITY:
		log.Print("equality ", val)
	case ORDERING:
		log.Print("ordering ", val)
	case SUBSTR:
		log.Print("substr ", val)
	case SYNTAX:
		log.Print("syntax ", val)
	case SINGLEVALUE:
		log.Print("singlevalue ", val)
	case COLLECTIVE:
		log.Print("collective ", val)
	case NOUSERMODIFICATION:
		log.Print("nousermodification ", val)
	case ABSTRACT:
		log.Print("abstract ", val)
	case STRUCTURAL:
		log.Print("structural ", val)
	case AUXILIARY:
		log.Print("auxiliary ", val)
	case MUST:
		log.Print("must ", val)
	case MAY:
		log.Print("may ", val)
	case DOLLAR:
		log.Print("dollar ", val)
	case LBRACE:
		log.Print("lbrace ", val)
	case RBRACE:
		log.Print("rbrace ", val)
	case USAGE:
		log.Print("usage ", val)
	case EOF:
		log.Print("eof ", val)
	case STRING:
		log.Print("string ", val)
	case OID:
		log.Print("oid ", val)
	case NUMERICOID:
		log.Print("noid ", val)
	}
}

// Write a lexed token to the channel of lexed items
func (l *lexer) emit(typ int, val string) {
	switch typ {
	case STRING:
		fallthrough
	case OID:
		fallthrough
	case NUMERICOID:
		l.items <- item{typ: typ, val: val}
	default:
		l.items <- item{typ: typ}
	}
	l.pos++
	l.start = l.pos
}

// Emits an EOF token and closes the items channel
func (l *lexer) eof(e string) {
	l.emit(EOF, "")
	close(l.items)
}

// Lexes raw text until an attributetype or objectclass definition
func lexText(l *lexer) stateFn {
	for {
		if l.start == len(l.input) {
			l.eof("lexText")
			return nil
		}
		if strings.ToLower(l.input[l.start]) == "attributetype" {
			return lexAttributeType
		}
		if strings.ToLower(l.input[l.start]) == "objectclass" {
			return lexObjectType
		}
		l.pos++
		l.start = l.pos
	}
}

// Lexes an attributetype token
func lexAttributeType(l *lexer) stateFn {
	l.emit(ATTRIBUTETYPE, "")
	return lexLeftBrace
}

// Lexes an objectclass token
func lexObjectType(l *lexer) stateFn {
	l.emit(OBJECTCLASS, "")
	return lexLeftBrace
}

// Lexes a left brace (
func lexLeftBrace(l *lexer) stateFn {
	l.emit(LBRACE, "")
	return lexNumericOid
}

// Lexes a numericOid like 1.2.3.4.1545612.1
func lexNumericOid(l *lexer) stateFn {
	l.emit(NUMERICOID, l.input[l.start])
	return lexAttributes
}

// Lexes the attributes of an attributetype or objectclass
func lexAttributes(l *lexer) stateFn {
	switch l.input[l.start] {
	case "NAME":
		l.emit(NAME, "")
		return lexName
	case "DESC":
		l.emit(DESC, "")
		return lexQuotedString
	case "OBSOLETE":
		l.emit(OBSOLETE, "")
		return lexAttributes
	case "SUP":
		l.emit(SUP, "")
		return lexOids
	case "EQUALITY":
		l.emit(EQUALITY, "")
		return lexOid
	case "ORDERING":
		l.emit(ORDERING, "")
		return lexOid
	case "SUBSTR":
		l.emit(SUBSTR, "")
		return lexOid
	case "SYNTAX":
		l.emit(SYNTAX, "")
		return lexNoidLength
	case "SINGLE-VALUE":
		l.emit(SINGLEVALUE, "")
		return lexAttributes
	case "COLLECTIVE":
		l.emit(COLLECTIVE, "")
		return lexAttributes
	case "NO-USER-MODIFICATION":
		l.emit(NOUSERMODIFICATION, "")
		return lexAttributes
	case "USAGE":
		l.emit(USAGE, "")
		return lexUsage
	case "ABSTRACT":
		l.emit(ABSTRACT, "")
		return lexAttributes
	case "STRUCTURAL":
		l.emit(STRUCTURAL, "")
		return lexAttributes
	case "AUXILIARY":
		l.emit(AUXILIARY, "")
		return lexAttributes
	case "MUST":
		l.emit(MUST, "")
		return lexOids
	case "MAY":
		l.emit(MAY, "")
		return lexOids
	case ")":
		l.emit(RBRACE, "")
		return lexText
	}
	l.eof("Attributes")
	return nil
}

// Lexes a single name in single quotes or multiple names in braces
// Names can't have whitespaces in them.
// Examples are: 'foobar' or ( 'foo' 'bar' )
func lexName(l *lexer) stateFn {
	if l.input[l.pos] == "(" {
		for {
			if l.start == len(l.input) {
				l.eof("Name")
			}

			switch l.input[l.start] {
			case "(":
				l.emit(LBRACE, "")
			case ")":
				l.emit(RBRACE, "")
				return lexAttributes
			default:
				l.emit(STRING, strings.TrimPrefix(strings.TrimRight(l.input[l.start], "'"), "'"))
			}
			l.start = l.pos
		}
	} else {
		l.emit(STRING, strings.TrimPrefix(strings.TrimRight(l.input[l.start], "'"), "'"))
		return lexAttributes
	}
}

// Lexes a string in single quotes, whitespaces in the string are permitted.
// Example: 'foo bar is great'
func lexQuotedString(l *lexer) stateFn {
	// the string only consists of the current fild
	if strings.HasPrefix(l.input[l.start], "'") && strings.HasSuffix(l.input[l.start], "'") {
		out := strings.TrimPrefix(l.input[l.start], "'")
		out = strings.TrimRight(out, "'")
		l.emit(STRING, out)
		return lexAttributes
	}

	// the string consists of multiple fields
	out := make([]string, 0)
	out = append(out, strings.TrimPrefix(l.input[l.start], "'"))
	l.pos++
	for {
		if l.pos == len(l.input) {
			l.eof("QuotedString")
		}

		if strings.HasSuffix(l.input[l.pos], "'") {
			out = append(out, strings.TrimRight(l.input[l.pos], "'"))
			l.emit(STRING, strings.Join(out, " "))
			return lexAttributes
		}

		out = append(out, l.input[l.pos])

		l.pos++
	}
}

// Lexes an Oid (string without enclosing braces and no whitespaces)
func lexOid(l *lexer) stateFn {
	l.emit(OID, l.input[l.start])
	return lexAttributes
}

// Lexes a numeric Oid, with an optional length specification in curly braces
// The length is currently ignored and not returned.
// Example: 1.2.3.4.5.6.7.8.9.0.1.2.3{32}
func lexNoidLength(l *lexer) stateFn {
	// ignore the length for now
	oid := strings.SplitN(l.input[l.start], "{", 2)
	l.emit(NUMERICOID, oid[0])
	return lexAttributes
}

// Lexes an usage string, but currently drops it silently.
func lexUsage(l *lexer) stateFn {
	l.emit(STRING, "")
	return lexAttributes
}

// Lexes a single oid, or a list of oids in braces seperated by dollar signs ( oid1 $ oid2 $ oid3 )
func lexOids(l *lexer) stateFn {
	if l.input[l.pos] == "(" {
		for {
			if l.start == len(l.input) {
				l.eof("Oids")
			}

			switch l.input[l.start] {
			case "(":
				l.emit(LBRACE, "")
			case ")":
				l.emit(RBRACE, "")
				return lexAttributes
			case "$":
				l.emit(DOLLAR, "")
			default:
				l.emit(OID, l.input[l.start])
			}
			l.start = l.pos
		}
	} else {
		l.emit(OID, l.input[l.start])
		return lexAttributes
	}
}
