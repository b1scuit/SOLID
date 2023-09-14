package rdf

type Prefix Objector

// An RDF Graph is a collection of RDF Triples
type Graph struct {
	Prefixes map[string]Prefix
	Triples  map[string]Triple
}
