package lexfn

import (
	"strings"

	"github.com/b1scuit/solid/rdf/lexer"
	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

func LexComment(lexer *lexer.Lexer) lexer.LexFn {
	// Remove the # at the start
	lexer.Pos += len(lexertoken.COMMENT)
	lexer.SkipWhitespace()
	lexer.Start = lexer.Pos

	for {
		if strings.HasPrefix(lexer.InputToEnd(), lexertoken.NEWLINE) {
			lexer.Emit(lexertoken.TOKEN_COMMENT)
			return LexStatement
		}

		lexer.Inc()

		if lexer.IsEOF() {
			return lexer.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}
