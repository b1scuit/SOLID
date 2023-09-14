package lexer

import (
	"testing"

	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

func TestLexer(t *testing.T) {
	l := &Lexer{
		Name:   "Test",
		Input:  "@prefix hello",
		State:  LexBegin,
		Tokens: make(chan lexertoken.Token, 3),
	}

	l.Run()
}
