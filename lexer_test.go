package lexer_test

import (
	"testing"

	"github.com/duckhue01/lexer"
)

const (
	NumberToken lexer.TokenType = iota
	OpToken
	IdentToken
)

func TestXxx(t *testing.T) {
	l := lexer.New(`abc=xyz`, lexKey)
	l.Lex()
	tok, _ := l.NextToken()

	t.Log(tok.Val)
}

func lexKey(l *lexer.L) lexer.StateFunc {
	l.Next()
	l.Emit(NumberToken)
	return nil
}

func lexOperation(l *lexer.L) lexer.StateFunc {
	return nil
}

func lexValue(l *lexer.L) lexer.StateFunc {
	return nil
}
