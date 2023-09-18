package lexer

import (
	"strings"

	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

func LexComment(lexer *Lexer) LexFn {
	// Remove the # at the start
	lexer.Pos += len(lexertoken.COMMENT)
	lexer.SkipWhitespace()
	lexer.Start = lexer.Pos

	for {
		if strings.HasPrefix(lexer.InputToEnd(), lexertoken.NEWLINE) {
			lexer.Emit(lexertoken.TOKEN_COMMENT)
			return LexBegin
		}

		lexer.Inc()

		if lexer.IsEOF() {
			return lexer.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}
