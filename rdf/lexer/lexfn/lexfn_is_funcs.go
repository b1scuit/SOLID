package lexfn

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

var (
	isIntegerRegexp = regexp.MustCompile(`(?m)[+-]?[0-9]+`)
	isDecimalRegexp = regexp.MustCompile(`(?m)[+-]?[0-9]*\.[0-9]+`)
)

func isTurtleDoc(s string) bool {
	return isStatement(s)
}

func isStatement(s string) bool {
	return isDirective(s) || isTriples(s)
}

// Directive
// [3]	directive	::=	prefixID | base | sparqlPrefix | sparqlBase
// [4]	prefixID	::=	'@prefix' PNAME_NS IRIREF '.'
// [5]	base	::=	'@base' IRIREF '.'
// [5s]	sparqlBase	::=	"BASE" IRIREF
// [6s]	sparqlPrefix	::=	"PREFIX" PNAME_NS IRIREF
func isDirective(s string) bool {
	if len(s) < len(lexertoken.SPARQL_PREFIX) {
		return false
	}

	return isPrefixID(s) || isBase(s) || isSparqlPrefix(s) || isSparqlBase(s)
}

func isPrefixID(s string) bool {
	return strings.HasPrefix(s, lexertoken.PREFIX)
}

func isBase(s string) bool {
	return strings.HasPrefix(s, lexertoken.BASE)
}

func isSparqlBase(s string) bool {
	return strings.EqualFold(s[:len(lexertoken.SPARQL_BASE)], lexertoken.SPARQL_BASE)
}

func isSparqlPrefix(s string) bool {
	return strings.EqualFold(s[:len(lexertoken.SPARQL_PREFIX)], lexertoken.SPARQL_PREFIX)
}

// Triples
// ::=	subject predicateObjectList | blankNodePropertyList predicateObjectList?
func isTriples(s string) bool {
	return isSubject(s) || isBlanknodePropertyList(s)
}

// predicateObjectList
// ::=	verb objectList (';' (verb objectList)?)*
func isPredicateObjectList(s string) bool {
	return isVerb(s)
}

func isObjectList(s string) bool {
	return isObject(s)
}

func isVerb(s string) bool {
	return isPredicate(s) || strings.HasPrefix(s, "a")
}

// Subject
// iri | BlankNode | collection
func isSubject(s string) bool {
	return isIri(s) || isBlankNode(s)
}

func isPredicate(s string) bool {
	return isIri(s)
}

// Object
// ::=	iri | BlankNode | collection | blankNodePropertyList | literal
func isObject(s string) bool {
	return isIri(s) || isBlankNode(s) || isCollection(s) || isBlanknodePropertyList(s) || isLiteral(s)
}

func isLiteral(s string) bool {
	return isRDFLiteral(s) || isNumericLiteral(s) || isBooleanLiteral(s)
}

func isBlanknodePropertyList(s string) bool {
	return strings.HasPrefix(s, "[")
}

func isCollection(s string) bool {
	return strings.HasPrefix(s, "(")
}

func isNumericLiteral(s string) bool {
	return isInteger(s) || isDecimal(s) || isDouble(s)
}

// RDF Literal
// ::=	String (LANGTAG | '^^' iri)?
func isRDFLiteral(s string) bool {
	// TODO: Add Langtag and ^^
	return isString(s)
}

func isBooleanLiteral(s string) bool {
	return strings.HasPrefix(s, "true") ||
		strings.HasPrefix(s, "false")
}

func isString(s string) bool {
	return isStringLiteralQuote(s) ||
		isStringLiteralSingleQuote(s) ||
		isStringLiteralLongQuote(s) ||
		isStringLiteralLongSingleQuote(s)
}

func isIri(s string) bool {
	return isIriRef(s) || isPrefixedName(s)
}

func isPrefixedName(s string) bool {
	return isPNameLn(s) || isPNameNs(s)
}

// Blank Node
// ::=	BLANK_NODE_LABEL | ANON
func isBlankNode(s string) bool {
	return isBlankNodeLabel(s) || isAnon(s)
}

// PRODUCTIONS FOR TERMINALS
// ---------------------

func isIriRef(s string) bool {
	return strings.HasPrefix(s, lexertoken.START_IRI)
}

func isPNameNs(s string) bool {
	return isPNPrefix(s)
}

func isPNameLn(s string) bool {
	return isPNameNs(s) || isPNLocal(s)
}

func isBlankNodeLabel(s string) bool {
	return strings.HasPrefix(s, "_:")
}

func isLangTag(s string) bool {
	return strings.HasPrefix(s, "@")
}

func isInteger(s string) bool {
	return isIntegerRegexp.MatchString(s)
}

func isDecimal(s string) bool {
	return isDecimalRegexp.MatchString(s)
}

func isDouble(s string) bool {
	return isDecimal(s)
}

func isExponent(s string) bool {
	return false // TODO Figure this out
}

func isStringLiteralQuote(s string) bool {
	return strings.HasPrefix(s, `"`)
}

func isStringLiteralSingleQuote(s string) bool {
	return strings.HasPrefix(s, `'`)
}

func isStringLiteralLongSingleQuote(s string) bool {
	return strings.HasPrefix(s, `'''`)
}

func isStringLiteralLongQuote(s string) bool {
	return strings.HasPrefix(s, `"""`)
}

func isUChar(s string) bool {
	return strings.HasPrefix(s, `\u`) || strings.HasPrefix(s, `\U`)
}

func isEChar(s string) bool {
	return strings.HasPrefix(s, `\`)
}

func isWs(s string) bool {
	// TODO Expand
	return unicode.IsLetter(rune(s[0]))
}

func isAnon(s string) bool {
	return strings.HasPrefix(s, "[")
}

func isPnCharsBase(s string) bool {
	return unicode.IsLetter(rune(s[0]))
}

func isPnCharsU(s string) bool {
	return isPnCharsBase(s)
}

func isPnChars(s string) bool {
	return isPnCharsU(s) || strings.HasPrefix(s, "_")
}

func isPNPrefix(s string) bool {
	return isPnCharsBase(s)
}

func isPNLocal(s string) bool {
	return isPnCharsU(s) || isPnChars(s)
}

func isPLX(s string) bool {
	return isPercent(s) || isPnLocalEsc(s)
}

func isPercent(s string) bool {
	return strings.HasPrefix(s, "%")
}

func isHex(s string) bool {
	return unicode.IsNumber(rune(s[0])) || unicode.IsLetter(rune(s[0]))
}

func isPnLocalEsc(s string) bool {
	return strings.HasPrefix(s, `\`)
}

// ###############################################################################

func isComment(s string) bool {
	return strings.HasPrefix(s, "#")
}
