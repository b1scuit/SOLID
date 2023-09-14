package lexertoken

type TokenType int

const (
	TOKEN_ERROR TokenType = iota
	TOKEN_EOF

	TOKEN_START_IRI
	TOKEN_END_IRI
	TOKEN_PREFIX
	TOKEN_END_PREFIX
	TOKEN_PREFIX_NAME
	TOKEN_BASE
	TOKEN_OBJECT_LIST
	TOKEN_END_TRIPLE
	TOKEN_COMMENT

	TOKEN_IRI
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

type Token struct {
	Type  TokenType
	Value string
}
