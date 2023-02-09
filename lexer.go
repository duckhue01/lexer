package lexer

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
	buffer       []rune
	stop         bool
	charC        map[rune]int
}

type TokenType int

type Token struct {
	Typ TokenType
	Val string
}

type StateFunc func(*L) StateFunc

const (
	EOFRune rune = -1
)

const (
	DefaultBufferCap = 10
)

// New creates new lexer instance
func New(src string, initState StateFunc, funcHandler func(e string)) *L {
	return &L{
		source:       src,
		start:        0,
		pos:          0,
		tokens:       make(chan Token),
		initState:    initState,
		rewind:       newRuneStack(),
		buffer:       make([]rune, 0, DefaultBufferCap),
		ErrorHandler: funcHandler,
		charC:        make(map[rune]int),
	}
}

// Lex starts the lexer machine
func (l *L) Lex() {
	go l.run()
}

func (l *L) run() {
	state := l.initState
	for state != nil && !l.stop {
		state = state(l)
	}

	close(l.tokens)
}

// Next pulls the next rune from the Lexer and returns it, moving the pos forward in the source.
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
	l.buffer = append(l.buffer, r)

	return r
}

// Take receives a string containing all acceptable strings and will continue over each consecutive character in the source until a token not in the given string is encountered. This should be used to quickly pull token parts.
func (l *L) Take(chars string) {
	r := l.Next()
	for strings.ContainsRune(chars, r) {
		r = l.Next()
	}
	l.Rewind() // last next wasn't a match
}

// Emit will receive a token type and push a new token with the current analyzed value into the tokens channel.
func (l *L) Emit(t TokenType) {
	tok := Token{
		Typ: t,
		Val: l.Current(),
	}
	l.tokens <- tok
	l.start = l.pos
	l.rewind.clear()
	l.buffer = make([]rune, 0, DefaultBufferCap)

}

// Peek performs a Next operation immediately followed by a Rewind returning the peeked rune.
func (l *L) Peek() rune {
	r := l.Next()
	l.Rewind()

	return r
}

// Ignore clears the rewind stack and then sets the current start pos to the current pos in the source which effectively ignores the section of the source being analyzed.
func (l *L) Ignore() {
	l.rewind.clear()
	l.buffer = make([]rune, 0, DefaultBufferCap)
	l.start = l.pos
}

// Rewind will take the last rune read (if any) and rewind back. Rewinds can occur more than once per call to Next but you can never rewind past the last point a token was emitted.
func (l *L) Rewind() {
	r := l.rewind.pop()
	if len(l.buffer) > 0 {
		l.buffer = l.buffer[:len(l.buffer)-1]
	}

	if r > EOFRune {
		size := utf8.RuneLen(r)
		l.pos -= size
		if l.pos < l.start {
			l.pos = l.start
		}
	}
}

// NextToken returns the next token from the lexer and a value to denote whether or not the token is finished.
func (l *L) NextToken() (*Token, bool) {
	if tok, ok := <-l.tokens; ok {
		return &tok, false
	} else {
		return nil, true
	}
}

// Error create new error instance and assign it to Err property.
func (l *L) Error(e string) {
	if l.ErrorHandler != nil {
		l.Err = errors.New(e)
		l.ErrorHandler(e)
		l.Stop()
	} else {
		panic(e)
	}
}

// Current returns the value being analyzed at this moment.
func (l *L) Current() string {
	return string(l.buffer)
}

// Skip skips the char at this time and move the pos to the next index.
func (l *L) Skip() {
	if len(l.buffer) > 0 {
		l.buffer = l.buffer[:len(l.buffer)-1]
	}
}

func (l *L) Inc(c rune) {
	l.charC[c]++

}

func (l *L) Dec(c rune) {
	if l.charC[c] > 0 {
		l.charC[c]--
	}
}

func (l *L) Count(c rune) int {
	return l.charC[c]
}

// Stop stop the lexer engine
func (l *L) Stop() {
	l.stop = true
}
