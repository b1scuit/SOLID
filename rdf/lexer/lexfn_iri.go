package lexer

import (
	"strings"

	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

func LexIri(lexer *Lexer) LexFn {
	lexer.SkipWhitespace()
	// Skip the opening <
	lexer.Pos += len(lexertoken.START_IRI)
	lexer.Start = lexer.Pos

	for {
		if strings.HasPrefix(lexer.InputToEnd(), lexertoken.END_IRI) {
			lexer.Emit(lexertoken.TOKEN_IRI)
			lexer.Pos += len(lexertoken.END_IRI)
			lexer.Ignore()

			return LexBegin
		}

		lexer.Inc()

		if lexer.IsEOF() {
			return lexer.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}

func LexEndLine(lexer *Lexer) LexFn {
	lexer.SkipWhitespace()

	for {
		l := lexer.InputToEnd()
		if strings.HasPrefix(l, lexertoken.END_TRIPLE) {
			lexer.Pos += len(lexertoken.END_TRIPLE)
			lexer.Emit(lexertoken.TOKEN_END_TRIPLE)
			return LexBegin
		}

		lexer.Inc()

		if lexer.IsEOF() {
			return lexer.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}

	}
}
