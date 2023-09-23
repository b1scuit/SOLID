package lexer

import "github.com/b1scuit/solid/rdf/lexer/lexertoken"

func LexNewLine(lexer *Lexer) LexFn {
	lexer.Pos += len(lexertoken.NEWLINE)

	if lexer.IsEOF() {
		lexer.Emit(lexertoken.TOKEN_EOF)
		return nil
	}
	return LexBegin
}
