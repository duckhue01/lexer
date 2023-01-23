package lexer

import (
	"reflect"
	"testing"
)

func Test_newRuneStack(t *testing.T) {
	tests := []struct {
		name string
		want runeStack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newRuneStack(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRuneStack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runeStack_push(t *testing.T) {
	type fields struct {
		start *runeNode
	}
	type args struct {
		r rune
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &runeStack{
				start: tt.fields.start,
			}
			s.push(tt.args.r)
		})
	}
}

func Test_runeStack_pop(t *testing.T) {
	type fields struct {
		start *runeNode
	}
	tests := []struct {
		name   string
		fields fields
		want   rune
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &runeStack{
				start: tt.fields.start,
			}
			if got := s.pop(); got != tt.want {
				t.Errorf("runeStack.pop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runeStack_clear(t *testing.T) {
	type fields struct {
		start *runeNode
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &runeStack{
				start: tt.fields.start,
			}
			s.clear()
		})
	}
}
