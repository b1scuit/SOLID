package lexertoken

import "fmt"

type TokenType int

const (
	TOKEN_ERROR TokenType = iota
	TOKEN_EOF

	TOKEN_PREFIX
	TOKEN_END_PREFIX
	TOKEN_PREFIX_NAME
	TOKEN_BASE
	TOKEN_OBJECT_LIST
	TOKEN_END_TRIPLE
	TOKEN_COMMENT
	TOKEN_PREFIXED_NAME

	TOKEN_IRI
	TOKEN_IRIREF
	TOKEN_BLANK_NODE
	TOKEN_LITERAL
	TOKEN_NEWLINE
)

const (
	EOF           rune = 0
	START_IRI          = "<"
	END_IRI            = ">"
	PREFIX             = "@prefix"
	SPARQL_PREFIX      = "PREFIX"
	SPARQL_BASE        = "BASE"
	BASE               = "@base"
	OBJECT_LIST        = ";"
	END_TRIPLE         = "."
	PREFIX_END         = ":"
	COMMENT            = "#"

	NEWLINE = "\n"
)

var TokenMap = map[TokenType]string{
	TOKEN_ERROR:         "Error",
	TOKEN_EOF:           "EOF",
	TOKEN_END_PREFIX:    "End Prefix Name(:)",
	TOKEN_PREFIX_NAME:   "prefixID",
	TOKEN_BASE:          "Base",
	TOKEN_OBJECT_LIST:   "Object List",
	TOKEN_END_TRIPLE:    "End Triple (.)",
	TOKEN_COMMENT:       "Comment (#)",
	TOKEN_IRI:           "IRI",
	TOKEN_IRIREF:        "IRIREF",
	TOKEN_BLANK_NODE:    "Blank Node",
	TOKEN_LITERAL:       "Literal",
	TOKEN_NEWLINE:       "New Line (\n)",
	TOKEN_PREFIXED_NAME: "Prefixed Name",
}

type Token struct {
	Type  TokenType
	Value string
}

func (t *TokenType) String() string {
	return fmt.Sprintf("%T", t)
}
