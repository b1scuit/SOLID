package lexer

import (
	"strings"

	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

func LexPrefix(lexer *Lexer) LexFn {
	// Remove the @prefix or PREFIX statements
	if strings.HasPrefix(lexer.InputToEnd(), lexertoken.PREFIX) {
		lexer.Pos += len(lexertoken.PREFIX)
	} else {
		lexer.Pos += len(lexertoken.SPARQL_PREFIX)
	}
	// reset the counter
	lexer.Ignore()
	for {
		if strings.HasPrefix(lexer.InputToEnd(), lexertoken.PREFIX_END) {
			lexer.Emit(lexertoken.TOKEN_PREFIX_NAME)
			lexer.Pos += len(lexertoken.PREFIX_END)

			return LexIri
		}

		lexer.Inc()

		if lexer.IsEOF() {
			return lexer.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}
