package parser

import (
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/b1scuit/solid/rdf/lexer"
	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

type ClientOption func(*Client)

type Client struct {
	l       *lexer.Lexer
	lexemes []lexertoken.Token

	prefixMap map[string]lexertoken.Token
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

func (c *Client) Do(file io.Reader) error {
	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	c.l, _ = lexer.New(
		lexer.WithInput(string(b)),
	)

	if err := c.CollectTokens(); err != nil {
		slog.Error("Error in collecting tokens", slog.Any("error", err))
		//return err
	}

	c.ParsePrefixes()

	//c.SwapPrefixForIRI()

	return nil
}

func (c *Client) GetPrefixMap() map[string]lexertoken.Token {
	return c.prefixMap
}

func (c *Client) GetLexemes() []lexertoken.Token {
	return c.lexemes
}

func (c *Client) CollectTokens() error {
	go c.l.Run()

	for t := range c.l.NextToken() {
		if t.Type == lexertoken.TOKEN_ERROR {
			return fmt.Errorf("lexer error: %v", t.Value)
		}

		if t.Type == lexertoken.TOKEN_EOF {
			break
		}

		c.lexemes = append(c.lexemes, t)
	}

	return nil
}

func (c *Client) ParsePrefixes() {
	if c.prefixMap == nil {
		c.prefixMap = make(map[string]lexertoken.Token)
	}

	// Loop through the tokens
	// if we find a prefix, we know 100% the next token is the IRI for the prefix
	// add that in and skip processing it
	for i := 0; i < len(c.lexemes); i++ {
		if c.lexemes[i].Type == lexertoken.TOKEN_PREFIX_NAME {
			c.prefixMap[c.lexemes[i].Value] = c.lexemes[i+1]
			i = i + 1
		}
	}
}

func (c *Client) SwapPrefixForIRI() {
	for i := 0; i < len(c.lexemes); i++ {

		if c.lexemes[i].Type == lexertoken.TOKEN_PREFIX_NAME {
			// Does this match anything in the prefix table
			prefix := strings.Split(c.lexemes[i].Value, ":")

			if iri, ok := c.prefixMap[prefix[0]]; ok {
				c.lexemes[i] = lexertoken.Token{
					Type:  lexertoken.TOKEN_IRI,
					Value: strings.Replace(c.lexemes[i].Value, prefix[0]+":", iri.Value, -1),
				}
			}
		}
	}
}

func (c *Client) FlattenPrefixObjectLists() {}
