package lexer

import (
	"strings"

	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

var LEXER_ERROR_UNEXPECTED_EOF = "Unexpected EOF"

func LexBegin(lexer *Lexer) LexFn {
	for {
		lexer.SkipWhitespace()
		l := lexer.InputToEnd()

		// @prefix or PREFIX
		if strings.HasPrefix(l, lexertoken.PREFIX) || strings.HasPrefix(l, lexertoken.SPARQL_PREFIX) {
			return LexPrefix
		}

		// @base or BASE
		if strings.HasPrefix(l, lexertoken.BASE) {
			return LexBase
		}

		if strings.HasPrefix(l, lexertoken.SPARQL_BASE) {
			return LexSparqlBase
		}

		// If the line is a straight <iri>
		if strings.HasPrefix(l, lexertoken.START_IRI) {
			return LexIri
		}

		// If we've arrived at a comment
		if strings.HasPrefix(l, lexertoken.COMMENT) {
			return LexComment
		}

		// If we're at the end of a line + some basic ignores
		if strings.HasPrefix(l, lexertoken.END_TRIPLE) {
			return LexEndLine
		}

		if strings.HasPrefix(l, lexertoken.NEWLINE) {
			return LexNewLine
		}

		if strings.HasPrefix(l, lexertoken.OBJECT_LIST) {
			return LexObjectList
		}

		if strings.HasPrefix(l, "\"") {
			return LexLiteral
		}

		if lexer.IsCharacter() {
			return LexPrefixedName
		}

		lexer.Inc()

		if lexer.IsEOF() {
			lexer.Emit(lexertoken.TOKEN_EOF)
		}
	}
}
