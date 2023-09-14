package parser

import (
	"fmt"
	"io"
	"os"

	"github.com/b1scuit/solid/rdf/lexer"
	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

type ClientOption func(*Client)

type Client struct {
}

func New(opts ...ClientOption) (*Client, error) {
	c := &Client{}

	for _, f := range opts {
		f(c)
	}

	return c, nil
}

func MustNew(opts ...ClientOption) *Client {
	c, err := New(opts...)

	if err != nil {
		panic(err)
	}

	return c
}

func (c *Client) Do(file *os.File) error {
	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	l, _ := lexer.New(
		lexer.WithName(file.Name()),
		lexer.WithInput(string(b)),
	)

	for {
		t := l.NextToken()

		fmt.Printf("Token: %+v\n", t)

		if t.Type == lexertoken.TOKEN_EOF {
			break
		}
	}

	return nil
}
