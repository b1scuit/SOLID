package lexer

import "github.com/b1scuit/solid/rdf/lexer/lexertoken"

func LexNewLine(lexer *Lexer) LexFn {
	lexer.Start += len(lexertoken.NEWLINE)
	lexer.Ignore()
	return LexBegin
}
