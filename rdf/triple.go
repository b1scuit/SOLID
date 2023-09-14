package rdf

import "net/url"

type IRI *url.URL
type Literal any
type BlankNode any // TODO figure out what this should actually be

type Objector interface {
	IRI | Literal | BlankNode
}

// SPARQL permits RDF Literals as the subject of RDF triples
// So Subject implements the Objector interface here
type Subject Objector
type Predicate string
type Object Objector

type Triple struct {
	Subject   Subject
	Predicate Predicate
	Object    Object
}
