package lexer

import (
	"strings"

	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

/*
This lexer function emits a TOKEN_VALUE with the value to be assigned
to a key.
*/
func LexLiteral(lexer *Lexer) LexFn {
	lexer.Pos += len("\"")

	for {
		if strings.HasPrefix(lexer.InputToEnd(), "\"") {
			lexer.Pos += len("\"")
			lexer.Emit(lexertoken.TOKEN_LITERAL)
			return LexBegin
		}

		lexer.Inc()

		if lexer.IsEOF() {
			return lexer.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}

}
