package lexer

import (
	"strings"

	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

func LexPrefix(lexer *Lexer) LexFn {
	lexer.Start += len(lexertoken.PREFIX)
	lexer.Ignore()
	lexer.SkipWhitespace()
	for {
		if strings.HasPrefix(lexer.InputToEnd(), lexertoken.PREFIX_END) {
			lexer.Emit(lexertoken.TOKEN_PREFIX_NAME)

			return LexIri
		}

		lexer.Inc()

		if lexer.IsEOF() {
			return lexer.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}

func LexSparqlPrefix(lexer *Lexer) LexFn {
	lexer.SkipWhitespace()
	lexer.Pos += len(lexertoken.SPARQL_PREFIX)
	for {
		if strings.HasPrefix(lexer.InputToEnd(), lexertoken.PREFIX_END) {
			lexer.Pos += len(lexertoken.PREFIX_END)
			lexer.Emit(lexertoken.TOKEN_PREFIX_NAME)

			return LexIri
		}

		lexer.Inc()

		if lexer.IsEOF() {
			return lexer.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}
func LexBase(lexer *Lexer) LexFn {
	lexer.SkipWhitespace()
	lexer.Pos += len(lexertoken.BASE)
	for {
		if strings.HasPrefix(lexer.InputToEnd(), " ") {
			lexer.Pos += len(" ")
			lexer.Emit(lexertoken.TOKEN_BASE)

			return LexIri
		}

		lexer.Inc()

		if lexer.IsEOF() {
			return lexer.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}

func LexSparqlBase(lexer *Lexer) LexFn {
	lexer.SkipWhitespace()
	lexer.Pos += len(lexertoken.SPARQL_BASE)
	for {
		if strings.HasPrefix(lexer.InputToEnd(), " ") {
			lexer.Pos += len(" ")
			lexer.Emit(lexertoken.TOKEN_BASE)

			return LexIri
		}

		lexer.Inc()

		if lexer.IsEOF() {
			return lexer.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}
