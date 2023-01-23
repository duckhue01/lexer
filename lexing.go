package lexing

import (
	"errors"
	"strings"
	"unicode/utf8"
)

type L struct {
	source       string
	start, pos   int
	tokens       chan Token
	initState    StateFunc
	Err          error
	ErrorHandler func(e string)
	rewind       runeStack
}

type TokenType int

type Token struct {
	Typ TokenType
	Val string
}

type StateFunc func(*L) StateFunc

const (
	EOFRune    rune      = -1
	EmptyToken TokenType = 0
)

func New(src string, initState StateFunc) *L {
	return &L{
		source:    src,
		start:     0,
		pos:       0,
		tokens:    make(chan Token),
		initState: initState,
		rewind:    newRuneStack(),
	}

}

func (l *L) Lex() {
	go l.run()
}

func (l *L) run() {
	state := l.initState
	for state != nil {
		state = state(l)
	}

	close(l.tokens)
}

// Current returns the value being being analyzed at this moment.
func (l *L) Current() string {
	return l.source[l.start:l.pos]
}

// Emit will receive a token type and push a new token with the current analyzed
// value into the tokens channel.
func (l *L) Emit(t TokenType) {
	tok := Token{
		Typ: t,
		Val: l.Current(),
	}
	l.tokens <- tok
	l.start = l.pos
	l.rewind.clear()
}

// Ignore clears the rewind stack and then sets the current beginning pos
// to the current pos in the source which effectively ignores the section
// of the source being analyzed.
func (l *L) Ignore() {
	l.rewind.clear()
	l.start = l.pos
}

// Peek performs a Next operation immediately followed by a Rewind returning the
// peeked rune.
func (l *L) Peek() rune {
	r := l.Next()
	l.Rewind()

	return r
}

// Rewind will take the last rune read (if any) and rewind back. Rewinds can
// occur more than once per call to Next but you can never rewind past the
// last point a token was emitted.
func (l *L) Rewind() {
	r := l.rewind.pop()
	if r > EOFRune {
		size := utf8.RuneLen(r)
		l.pos -= size
		if l.pos < l.start {
			l.pos = l.start
		}
	}
}

// Next pulls the next rune from the Lexer and returns it, moving the pos
// forward in the source.
func (l *L) Next() rune {
	var (
		r rune
		s int
	)
	str := l.source[l.pos:]
	if len(str) == 0 {
		r, s = EOFRune, 0
	} else {
		r, s = utf8.DecodeRuneInString(str)
	}
	l.pos += s
	l.rewind.push(r)

	return r
}

// Take receives a string containing all acceptable strings and will contine
// over each consecutive character in the source until a token not in the given
// string is encountered. This should be used to quickly pull token parts.
func (l *L) Take(chars string) {
	r := l.Next()
	for strings.ContainsRune(chars, r) {
		r = l.Next()
	}
	l.Rewind() // last next wasn't a match
}

// NextToken returns the next token from the lexer and a value to denote whether
// or not the token is finished.
func (l *L) NextToken() (*Token, bool) {
	if tok, ok := <-l.tokens; ok {
		return &tok, false
	} else {
		return nil, true
	}
}

// Partial yyLexer implementation

func (l *L) Error(e string) {
	if l.ErrorHandler != nil {
		l.Err = errors.New(e)
		l.ErrorHandler(e)
	} else {
		panic(e)
	}
}
