package lexer

import (
	"bytes"
	"testing"

	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
	"github.com/olekukonko/tablewriter"
)

func TestLexer(t *testing.T) {
	l := &Lexer{
		Input:  "@prefix hello",
		State:  LexTurtleDoc,
		Tokens: make(chan lexertoken.Token, 3),
	}

	b := bytes.Buffer{}
	table := tablewriter.NewWriter(&b)
	table.SetHeader([]string{"Token", "Value"})
	go l.Run()

	for v := range l.Tokens {
		table.Append([]string{lexertoken.TokenMap[v.Type], v.Value})

		if v.Type == lexertoken.TOKEN_ERROR {
			t.Error(v.Value)
			break
		}

		if v.Type == lexertoken.TOKEN_EOF {
			break
		}
	}

	table.Render()

	t.Logf("\n%v\n", b.String())
}
