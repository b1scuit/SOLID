package lexer

import (
	"strings"

	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

var LEXER_ERROR_UNEXPECTED_EOF = "Unexpected EOF"

func LexBegin(lexer *Lexer) LexFn {
	lexer.SkipWhitespace()

	for {
		l := lexer.InputToEnd()

		if strings.HasPrefix(l, lexertoken.PREFIX) {
			return LexPrefix
		}

		if strings.HasPrefix(l, lexertoken.SPARQL_PREFIX) {
			return LexSparqlPrefix
		}
		if strings.HasPrefix(l, lexertoken.BASE) {
			return LexBase
		}

		if strings.HasPrefix(l, lexertoken.SPARQL_BASE) {
			return LexSparqlBase
		}
		if strings.HasPrefix(l, lexertoken.START_IRI) {
			return LexIri
		}

		if strings.HasPrefix(l, lexertoken.END_TRIPLE) {
			return LexEndLine
		}

		if strings.HasPrefix(l, lexertoken.COMMENT) {
			return LexComment
		}

		if strings.HasPrefix(l, lexertoken.NEWLINE) {
			return LexNewLine
		}

		lexer.Ignore()

		lexer.Inc()

		if lexer.IsEOF() {
			lexer.Emit(lexertoken.TOKEN_EOF)
		}
	}
}
