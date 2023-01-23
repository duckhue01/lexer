package lexing_test

import (
	"testing"

	"github.com/duckhue01/lexing"
)

const (
	NumberToken lexing.TokenType = iota
	OpToken
	IdentToken
)

func TestXxx(t *testing.T) {
	lexer := lexing.New(`abc=xyz`, lexKey)
	lexer.Lex()
	tok, _ := lexer.NextToken()

	t.Log(tok.Val)
}

func lexKey(l *lexing.L) lexing.StateFunc {
	l.Next()
	l.Emit(NumberToken)
	return nil
}

func lexOperation(l *lexing.L) lexing.StateFunc {
	return nil
}

func lexValue(l *lexing.L) lexing.StateFunc {
	return nil
}
