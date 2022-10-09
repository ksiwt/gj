package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pohedev/gj.git/token"
)

const (
	nullValue      = "null"
	boolTrueValue  = "true"
	boolFalseValue = "false"
)

const (
	nullValueLen      = 4
	boolTrueValueLen  = 4
	boolFalseValueLen = 5
)

const eof = -1

// Item represents a Token returned from the scanner.
type Item struct {
	Token token.Token // The Token of this Item.
	Pos   int         // The starting position, in bytes, of this Item in the input string.
	Val   string      // The value of this Item.
}

func (i Item) String() string {
	switch i.Token {
	case token.EOF:
		return "EOF"
	case token.Error:
		return i.Val
	}
	if len(i.Val) > 10 {
		return fmt.Sprintf("%.10q...", i.Val)
	}
	return fmt.Sprintf("%q", i.Val)
}

// Lexer holds the state of the scanner.
type Lexer struct {
	input string    // the string being scanned.
	start int       // start position of this Item.
	pos   int       // current position in the input.
	width int       // width of last rune read from input.
	items chan Item // channel of scanned items.
}

// Lex creates a new lexer.
func Lex(input string) *Lexer {
	l := &Lexer{
		input: input,
		items: make(chan Item),
	}
	go l.run() // concurrently run state machine.
	return l
}

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// run lexer the input by executing state functions until
// the state is nil.
func (l *Lexer) run() {
	for state := lexToken; state != nil; {
		state = state(l)

	}
	close(l.items) // No more tokens will be delivered.
}

// emit passes an Item back to the client.
func (l *Lexer) emit(t token.Token) {
	l.items <- Item{
		Token: t,
		Pos:   l.start,
		Val:   l.input[l.start:l.pos],
	}
	l.start = l.pos
}

// next returns the next rune in the input.
func (l *Lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// backup steps back one rune.
// can be called only once per call of next.
func (l *Lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume
// the next rune in the input.
func (l *Lexer) peek() int32 {
	rune := l.next()
	l.backup()
	return rune
}

// accept consumes the next rune
// if it's from the valid set.
func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *Lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// error returns an error Token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- Item{token.Error, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// NextItem returns the next Item from the input. The Lexer has to be
// drained (all items received until itemEOF or itemError) - otherwise
// the Lexer goroutine will leak.
func (l *Lexer) NextItem() Item {
	return <-l.items
}

// lexToken scans current char and creates a new Token.
func lexToken(l *Lexer) stateFn {
	for {
		r := l.next()
		if r == eof {
			break
		}

		switch {
		case isSpace(r):
			l.ignore()
		case r == '{':
			l.emit(token.LeftBrace)
			return lexToken
		case r == '}':
			l.emit(token.RightBrace)
			return lexToken
		case r == '[':
			l.emit(token.LeftBracket)
			return lexToken
		case r == ']':
			l.emit(token.RightBracket)
			return lexToken
		case r == ':':
			l.emit(token.Colon)
			return lexToken
		case r == ',':
			l.emit(token.Comma)
			return lexToken
		case r == '"':
			return lexQuote
		case isNumber(r):
			l.backup()
			return lexNumber
		case r == 'n':
			l.backup()
			return lexNull
		case r == 't' || r == 'f':
			l.backup()
			return lexBool
		default:
			l.emit(token.Unknown)
		}
	}
	// Correctly reached EOF.
	l.emit(token.EOF)
	return nil // Stop the run loop.
}

// lexQuote scans a run of quoted string.
func lexQuote(l *Lexer) stateFn {
	for {
		switch l.next() {
		case '\\':
			if r := l.next(); r != eof && r != '\n' {
				break
			}
		case eof, '\n':
			return l.errorf("unterminated quoted string")
		case '"':
			l.emit(token.String)
			return lexToken
		}
	}
}

// lexNumber scans a run of number.
func lexNumber(l *Lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(token.Number)
	return lexToken
}

func (l *Lexer) scanNumber() bool {
	// Optional leading sign.
	l.accept("+-")

	digits := "0123456789_"
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}

	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}

	return true
}

// lexNull scans a run of null.
func lexNull(l *Lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], nullValue) {
		for i := 0; i < nullValueLen; i++ {
			l.next()
		}
		l.emit(token.Null)
	}
	return lexToken
}

// lexBool scans a run of boolean.
func lexBool(l *Lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], boolTrueValue) {
		for i := 0; i < boolTrueValueLen; i++ {
			l.next()
		}
		l.emit(token.True)
	} else if strings.HasPrefix(l.input[l.pos:], boolFalseValue) {
		for i := 0; i < boolFalseValueLen; i++ {
			l.next()
		}
		l.emit(token.False)
	}
	return lexToken
}

// isSpace reports whether rune is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

// isNumber reports whether rune is a number.
func isNumber(r rune) bool {
	return r == '+' || r == '-' || ('0' <= r && r <= '9')
}

// isAlphaNumeric reports whether rune is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
