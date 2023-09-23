package lexer

import (
	"strings"

	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

func LexPrefixedName(lexer *Lexer) LexFn {
	for {
		if strings.HasPrefix(lexer.InputToEnd(), " ") {
			lexer.Emit(lexertoken.TOKEN_PREFIXED_NAME)
			return LexBegin
		}

		if strings.HasPrefix(lexer.InputToEnd(), ";") {
			lexer.Emit(lexertoken.TOKEN_PREFIXED_NAME)
			return LexObjectList
		}

		lexer.Inc()

		if lexer.IsEOF() {
			return lexer.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}
func LexObjectList(lexer *Lexer) LexFn {
	lexer.Pos += len(lexertoken.OBJECT_LIST)
	lexer.Emit(lexertoken.TOKEN_OBJECT_LIST)
	return LexBegin
}
