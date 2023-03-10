package lexer_test

import (
	"testing"

	"github.com/duckhue01/lexer"
	"github.com/google/go-cmp/cmp"
)

const (
	Key lexer.TokenType = iota
	Op
	Value
	Divider
)

func TestNew(t *testing.T) {
	type args struct {
		src       string
		initState lexer.StateFunc
	}
	tests := []struct {
		name string
		args args
		want *lexer.L
	}{
		{
			name: "int a lexer with initState function",
			args: args{
				src: "key=value",
				initState: func(*lexer.L) lexer.StateFunc {
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lexer.New(tt.args.src, tt.args.initState, nil); got == nil {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateFunc(t *testing.T) {
	tests := []struct {
		name  string
		lexer *lexer.L
		want  *lexer.Token
	}{
		{
			name: "state func lex one character with Next function",
			lexer: lexer.New("key=value", func(l *lexer.L) lexer.StateFunc {
				l.Next()
				l.Emit(Key)
				return nil
			}, nil),
			want: &lexer.Token{
				Typ: Key,
				Val: "k",
			},
		},
		{
			name: "state func lex the key tokens with Take function",
			lexer: lexer.New("key=value", func(l *lexer.L) lexer.StateFunc {
				l.Take("key")
				l.Emit(Key)
				return nil
			}, nil),
			want: &lexer.Token{
				Typ: Key,
				Val: "key",
			},
		},
		{
			name: "state func lex Emit the first char after Peek Peek function",
			lexer: lexer.New("key=value", func(l *lexer.L) lexer.StateFunc {
				t.Log("l.Peek():", string(l.Peek()))
				l.Next()
				l.Emit(Key)
				return nil
			}, nil),
			want: &lexer.Token{
				Typ: Key,
				Val: "k",
			},
		},
		{
			name: "state func ignore a character then lex the second char",
			lexer: lexer.New("key=value", func(l *lexer.L) lexer.StateFunc {
				l.Next()
				l.Ignore()
				l.Next()
				l.Emit(Key)
				return nil
			}, nil),
			want: &lexer.Token{
				Typ: Key,
				Val: "e",
			},
		},
		{
			name: "state func take two first char",
			lexer: lexer.New("key=value", func(l *lexer.L) lexer.StateFunc {
				l.Next()
				l.Emit(Key)
				l.Rewind()
				l.Next()
				l.Emit(Key)
				return nil
			}, nil),
			want: &lexer.Token{
				Typ: Key,
				Val: "k",
			},
		},
		{
			name: "state func take three first char and skip one",
			lexer: lexer.New("key=value", func(l *lexer.L) lexer.StateFunc {
				l.Next()
				l.Next()
				l.Next()
				l.Skip()
				l.Emit(Key)
				return nil
			}, nil),
			want: &lexer.Token{
				Typ: Key,
				Val: "ke",
			},
		},
		{
			name: "state func through error",
			lexer: lexer.New("key=value", func(l *lexer.L) lexer.StateFunc {
				l.ErrorHandler = func(e string) {

				}
				l.Error("asdasd")

				return nil
			}, nil),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.lexer.Lex()
			got, _ := tt.lexer.NextToken()

			if !cmp.Equal(got, tt.want) {
				t.Errorf("(tt.lexer.Lex() = %v, want %v", got, tt.want)
			}
		})
	}
}
