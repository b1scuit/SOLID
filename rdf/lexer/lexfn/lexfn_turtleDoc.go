package lexfn

import (
	"strings"

	"github.com/b1scuit/solid/rdf/lexer"
	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

var LEXER_ERROR_UNEXPECTED_EOF = "Unexpected EOF"

// https://www.w3.org/TR/2014/REC-turtle-20140225/#grammar-production-turtleDoc
// turtleDoc	::=	statement*
// A turtle doc is defined as a series of RDF turtle staements
// What this means for the lex is keep restarting the LexStatement
// until you hit EOF and finish
func LexTurtleDoc(lex *lexer.Lexer) lexer.LexFn {
	for {
		return LexStatement
	}
}

// https://www.w3.org/TR/2014/REC-turtle-20140225/#grammar-production-statement
// statement	::=	directive | triples '.'
// A Statement is either a directive (@prefix, BASE stuff like that) or a triple that ends in a .
//
// This lexer.Lexer function does nothing except direct towards a directive or triple lex func
// I've included lexing a line wide comment here, however a comment may appear on any line
// till the end of it
func LexStatement(lex *lexer.Lexer) lexer.LexFn {
	for {
		lex.SkipWhitespace()

		if lex.IsEOF() {
			// This will also kill the processing
			lex.Emit(lexertoken.TOKEN_EOF)
		}

		l := lex.InputToEnd()

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
			lex.State = LexTriples(lex)
			lex.SkipWhitespace()
			lex.Pos += len(lexertoken.END_TRIPLE)
			lex.Ignore()

			return lex.State
		}

		lex.Inc()
	}
}

// Directive
//
//	directive	::=	prefixID | base | sparqlPrefix | sparqlBase
//
// Simple string match for this one
func LexDirective(lex *lexer.Lexer) lexer.LexFn {
	l := lex.InputToEnd()

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

	return lex.Errorf("lexdirective wasn't given a tutle directive: %v", lex.InputToEnd())
}

// prefixID	::=	'@prefix' PNAME_NS IRIREF '.'
func LexPrefixId(lex *lexer.Lexer) lexer.LexFn {
	// Move the Pos counter over the length of @prefix
	// Then Ignore() sets start == pos, omitting the prefix keyword
	lex.Pos += len(lexertoken.PREFIX)
	lex.SkipWhitespace()
	lex.Ignore()

	for {
		if strings.HasPrefix(lex.InputToEnd(), lexertoken.PREFIX_END) {
			lex.Emit(lexertoken.TOKEN_PREFIX_NAME)
			lex.Pos += len(lexertoken.PREFIX_END)

			// Lex IRI and knock off the . at the end
			lex.State = LexIriRef(lex)
			lex.SkipWhitespace()
			lex.Pos += len(lexertoken.END_TRIPLE)
			lex.Ignore()

			return lex.State
		}

		lex.Inc()

		if lex.IsEOF() {
			return lex.Errorf("Unexpected EOF")
		}
	}
}

func LexBase(lex *lexer.Lexer) lexer.LexFn {
	lex.SkipWhitespace()
	lex.Pos += len(lexertoken.BASE)
	for {
		if strings.HasPrefix(lex.InputToEnd(), " ") {
			lex.Pos += len(" ")
			lex.Emit(lexertoken.TOKEN_BASE)

			// Lex IRI and knock off the . at the end
			lex.State = LexIriRef(lex)
			lex.SkipWhitespace()
			lex.Pos += len(lexertoken.END_TRIPLE)
			lex.SkipWhitespace()
			lex.Ignore()

			if isComment(lex.InputToEnd()) {
				return LexComment
			}

			return lex.State
		}

		lex.Inc()

		if lex.IsEOF() {
			return lex.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}

func LexSparqlBase(lex *lexer.Lexer) lexer.LexFn {
	lex.SkipWhitespace()
	lex.Pos += len(lexertoken.SPARQL_BASE)
	for {
		if strings.HasPrefix(lex.InputToEnd(), " ") {
			lex.Pos += len(" ")
			lex.Emit(lexertoken.TOKEN_BASE)

			// Lex IRI and knock off the . at the end
			lex.State = LexIriRef(lex)
			lex.SkipWhitespace()
			lex.Pos += len(lexertoken.NEWLINE)
			lex.Ignore()

			if isComment(lex.InputToEnd()) {
				return LexComment
			}

			return lex.State
		}

		lex.Inc()

		if lex.IsEOF() {
			return lex.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}

func LexSparqlPrefix(lex *lexer.Lexer) lexer.LexFn {
	// Move the Pos counter over the length of @prefix
	// Then Ignore() sets start == pos, omitting the prefix keyword
	lex.Pos += len(lexertoken.SPARQL_PREFIX)
	lex.Ignore()

	for {
		if strings.HasPrefix(lex.InputToEnd(), lexertoken.PREFIX_END) {
			lex.Emit(lexertoken.TOKEN_PREFIX_NAME)
			lex.Pos += len(lexertoken.PREFIX_END)

			lex.State = LexIriRef(lex)
			lex.SkipWhitespace()
			lex.Pos += len(lexertoken.END_TRIPLE)
			lex.SkipWhitespace()
			lex.Ignore()

			if isComment(lex.InputToEnd()) {
				return LexComment
			}

			return lex.State
		}

		lex.Inc()

		if lex.IsEOF() {
			return lex.Errorf("Unexpected EOF")
		}
	}
}

// Triples
// ::=	subject predicateObjectList | blankNodePropertyList predicateObjectList?
func LexTriples(lex *lexer.Lexer) lexer.LexFn {
	l := lex.InputToEnd()

	if isSubject(l) {
		// On a triple the subject is always followed by a predicateObjectList
		lex.State = LexSubject(lex)
		lex.SkipWhitespace()
		lex.State = LexPredicateObjectList(lex)
	} else if isBlanknodePropertyList(l) {
		lex.State = LexBlankNodePropertyList(lex)
		lex.SkipWhitespace()

		if isPredicateObjectList(lex.InputToEnd()) {
			lex.State = LexPredicateObjectList(lex)
		}
	} else {
		lex.State = lex.Errorf("input to LexTriples was not a triple")
	}

	return lex.State
}

// predicateObjectList
// ::=	verb objectList (';' (verb objectList)?)*
func LexPredicateObjectList(lex *lexer.Lexer) lexer.LexFn {

	for {
		lex.SkipWhitespace()
		lex.Ignore()

		if isVerb(lex.CurrentInput()) {
			lex.State = LexVerb(lex)
		}

		if isObjectList(lex.InputToEnd()) {
			lex.State = LexObjectList(lex)
		}

		if strings.HasPrefix(lex.InputToEnd(), lexertoken.OBJECT_LIST) {
			lex.Pos += len(lexertoken.OBJECT_LIST)
			lex.Ignore()
		}

		if strings.HasPrefix(lex.InputToEnd(), lexertoken.END_TRIPLE) {
			return lex.State
		}

		lex.Inc()
	}
}

func LexObjectList(lex *lexer.Lexer) lexer.LexFn {
	for {
		lex.SkipWhitespace()

		if strings.HasPrefix(lex.InputToEnd(), lexertoken.END_TRIPLE) || strings.HasPrefix(lexertoken.OBJECT_LIST, lex.InputToEnd()) {
			return lex.State
		}

		if isObject(lex.InputToEnd()) {
			lex.State = LexObject(lex)
		}

		if strings.HasPrefix(lex.InputToEnd(), lexertoken.OBJECT) {
			lex.Pos += len(lexertoken.OBJECT)
			lex.Ignore()
		}

		lex.Inc()
	}
}

func LexVerb(lex *lexer.Lexer) lexer.LexFn {
	lex.SkipWhitespace()

	if isPredicate(lex.InputToEnd()) {
		lex.State = LexPredicate(lex)
		return lex.State
	}

	if strings.HasPrefix(lex.InputToEnd(), "a") {
		lex.Pos += len("a")
		lex.Emit(lexertoken.TOKEN_PREDICATE)
	}

	return lex.Errorf("value passed to LexVerb unknown")
}

// Subject
// ::=	iri | BlankNode | collection
func LexSubject(lex *lexer.Lexer) lexer.LexFn {
	l := lex.InputToEnd()

	if isIri(l) {
		lex.State = LexIri(lex)
	} else if isBlankNode(l) {
		lex.State = LexBlankNode(lex)
	} else if isCollection(l) {
		lex.State = LexCollection(lex)
	} else {
		lex.State = lex.Errorf("input to LexSubject was not a RDF subject")
	}

	return lex.State
}

func LexPredicate(lex *lexer.Lexer) lexer.LexFn {
	return LexIri
}

func LexObject(lex *lexer.Lexer) lexer.LexFn {
	l := lex.InputToEnd()

	if isIri(l) {
		lex.State = LexIri(lex)
	} else if isBlankNode(l) {
		lex.State = LexBlankNode(lex)
	} else if isCollection(l) {
		lex.State = LexCollection(lex)
	} else if isBlanknodePropertyList(l) {
		lex.State = LexBlankNodePropertyList(lex)
	} else if isLiteral(l) {
		lex.State = LexLiteral(lex)
	} else {
		lex.State = lex.Errorf("invalid input passed to LexObject")
	}

	return lex.State
}

func LexLiteral(lex *lexer.Lexer) lexer.LexFn {
	return lex.Errorf("literal unimplemented")
}

func LexBlankNodePropertyList(lex *lexer.Lexer) lexer.LexFn {
	return lex.Errorf("blankNodePropertyList unimplemented")
}

func LexCollection(lex *lexer.Lexer) lexer.LexFn {
	return lex.Errorf("collection unimplemented")
}

func LexNumericLiteral(lex *lexer.Lexer) lexer.LexFn {
	return lex.Errorf("NumericLiteral unimplemented")
}

func LexRDFLiteral(lex *lexer.Lexer) lexer.LexFn {
	return lex.Errorf("RDFLiteral unimplemented")
}

func LexBooleanLiteral(lex *lexer.Lexer) lexer.LexFn {
	return lex.Errorf("BooleanLiteral unimplemented")
}

func LexString(lex *lexer.Lexer) lexer.LexFn {
	return lex.Errorf("String unimplemented")
}

func LexIri(lex *lexer.Lexer) lexer.LexFn {
	if isIriRef(lex.InputToEnd()) {
		lex.State = LexIriRef(lex)
	} else if isPrefixedName(lex.InputToEnd()) {
		lex.State = LexPrefixedName(lex)
	} else {
		lex.State = lex.Errorf("input to LexIRI was not an RDF IRI")
	}

	return lex.State
}

func LexPrefixedName(lex *lexer.Lexer) lexer.LexFn {
	return lex.Errorf("PrefixedName unimplemented")
}

func LexBlankNode(lex *lexer.Lexer) lexer.LexFn {
	return lex.Errorf("BlankNode unimplemented")
}

func LexIriRef(lex *lexer.Lexer) lexer.LexFn {
	lex.SkipWhitespace()
	lex.Pos += len(lexertoken.START_IRI)
	lex.Ignore()

	for {
		if strings.HasPrefix(lex.InputToEnd(), lexertoken.END_IRI) {
			lex.Emit(lexertoken.TOKEN_IRIREF)
			lex.Pos += len(lexertoken.END_IRI)
			lex.Ignore()

			return LexStatement
		}

		lex.Inc()

		if lex.IsEOF() {
			return lex.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}
