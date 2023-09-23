package lexer

import (
	"strings"

	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

var LEXER_ERROR_UNEXPECTED_EOF = "Unexpected EOF"

// https://www.w3.org/TR/2014/REC-turtle-20140225/#grammar-production-turtleDoc
// turtleDoc	::=	statement*
// A turtle doc is defined as a series of RDF turtle staements
// What this means for the lexer is keep restarting the LexStatement
// until you hit EOF and finish
func LexTurtleDoc(lexer *Lexer) LexFn {
	for {
		return LexStatement
	}
}

// https://www.w3.org/TR/2014/REC-turtle-20140225/#grammar-production-statement
// statement	::=	directive | triples '.'
// A Statement is either a directive (@prefix, BASE stuff like that) or a triple that ends in a .
//
// This Lexer function does nothing except direct towards a directive or triple lex func
// I've included lexing a line wide comment here, however a comment may appear on any line
// till the end of it
func LexStatement(lexer *Lexer) LexFn {
	for {
		lexer.SkipWhitespace()

		if lexer.IsEOF() {
			// This will also kill the processing
			lexer.Emit(lexertoken.TOKEN_EOF)
		}

		l := lexer.InputToEnd()

		if isComment(l) {
			return LexComment
		}

		if isDirective(l) {
			return LexDirective
		}

		// If the start of the statement is an IRIREF
		// high chance this is a triple ahead
		// It is also the responsibility of the statement
		// to lex the "." at the end of a triple
		if isTriples(l) {
			lexer.State = LexTriples(lexer)
			lexer.SkipWhitespace()
			lexer.Pos += len(lexertoken.END_TRIPLE)
			lexer.Ignore()

			return lexer.State
		}

		lexer.Inc()
	}
}

// Directive
//
//	directive	::=	prefixID | base | sparqlPrefix | sparqlBase
//
// Simple string match for this one
func LexDirective(lexer *Lexer) LexFn {
	l := lexer.InputToEnd()

	if isPrefixID(l) {
		return LexPrefixId
	}

	if isBase(l) {
		return LexBase
	}

	if isSparqlPrefix(l) {
		return LexSparqlPrefix
	}

	if isSparqlBase(l) {
		return LexSparqlBase
	}

	return lexer.Errorf("lexdirective wasn't given a tutle directive: %v", lexer.InputToEnd())
}

// prefixID	::=	'@prefix' PNAME_NS IRIREF '.'
func LexPrefixId(lexer *Lexer) LexFn {
	// Move the Pos counter over the length of @prefix
	// Then Ignore() sets start == pos, omitting the prefix keyword
	lexer.Pos += len(lexertoken.PREFIX)
	lexer.SkipWhitespace()
	lexer.Ignore()

	for {
		if strings.HasPrefix(lexer.InputToEnd(), lexertoken.PREFIX_END) {
			lexer.Emit(lexertoken.TOKEN_PREFIX_NAME)
			lexer.Pos += len(lexertoken.PREFIX_END)

			// Lex IRI and knock off the . at the end
			lexer.State = LexIriRef(lexer)
			return lexer.State
		}

		lexer.Inc()

		if lexer.IsEOF() {
			return lexer.Errorf("Unexpected EOF")
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

			// Lex IRI and knock off the . at the end
			lexer.State = LexIriRef(lexer)
			lexer.SkipWhitespace()
			lexer.Pos += len(lexertoken.END_TRIPLE)
			lexer.SkipWhitespace()
			lexer.Ignore()

			if isComment(lexer.InputToEnd()) {
				return LexComment
			}

			return lexer.State
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

			// Lex IRI and knock off the . at the end
			lexer.State = LexIriRef(lexer)
			lexer.SkipWhitespace()
			lexer.Pos += len(lexertoken.NEWLINE)
			lexer.Ignore()

			if isComment(lexer.InputToEnd()) {
				return LexComment
			}

			return lexer.State
		}

		lexer.Inc()

		if lexer.IsEOF() {
			return lexer.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}

func LexSparqlPrefix(lexer *Lexer) LexFn {
	// Move the Pos counter over the length of @prefix
	// Then Ignore() sets start == pos, omitting the prefix keyword
	lexer.Pos += len(lexertoken.SPARQL_PREFIX)
	lexer.Ignore()

	for {
		if strings.HasPrefix(lexer.InputToEnd(), lexertoken.PREFIX_END) {
			lexer.Emit(lexertoken.TOKEN_PREFIX_NAME)
			lexer.Pos += len(lexertoken.PREFIX_END)

			lexer.State = LexIriRef(lexer)
			lexer.SkipWhitespace()
			lexer.Pos += len(lexertoken.END_TRIPLE)
			lexer.SkipWhitespace()
			lexer.Ignore()

			if isComment(lexer.InputToEnd()) {
				return LexComment
			}

			return lexer.State
		}

		lexer.Inc()

		if lexer.IsEOF() {
			return lexer.Errorf("Unexpected EOF")
		}
	}
}

func LexTriples(lexer *Lexer) LexFn {
	l := lexer.InputToEnd()

	if isSubject(l) {
		return lexer.Errorf("No Subject lex func")
	}

	if isBlanknodePropertyList(l) {
		return lexer.Errorf("No Blank node property list lex fuc")
	}

	return lexer.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
}

func LexIriRef(lexer *Lexer) LexFn {
	lexer.SkipWhitespace()
	lexer.Pos += len(lexertoken.START_IRI)
	lexer.Ignore()

	for {
		if strings.HasPrefix(lexer.InputToEnd(), lexertoken.END_IRI) {
			lexer.Emit(lexertoken.TOKEN_IRIREF)
			lexer.Pos += len(lexertoken.END_IRI)
			lexer.Ignore()

			return LexStatement
		}

		lexer.Inc()

		if lexer.IsEOF() {
			return lexer.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}
