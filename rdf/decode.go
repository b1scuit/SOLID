package rdf

import (
	"log/slog"
	"reflect"
	"regexp"
	"strings"
)

var (
	removeWhitespace = regexp.MustCompile(`(?m)\n|\t`)
	prefixRegex      = regexp.MustCompile(`(?m)^(.*):.*<(.*)>$`)
)

func Unmarshal(b []byte, v any) error {
	// Is v a pointer reference
	if err := isPointer(v); err != nil {
		return err
	}

	l := &lukeFilewalker{
		data:  b,
		graph: v.(*Graph),
	}

	return l.Parse()
}

type lukeFilewalker struct {
	data  []byte
	lines []string
	graph *Graph
}

func (l *lukeFilewalker) Parse() error {
	// A RDF file uses a "." as it's end of line
	l.lines = strings.Split(string(l.data), ".\n")

	for _, line := range l.lines {
		line = removeWhitespace.ReplaceAllString(line, "")

		// Is this a prefix line
		if strings.HasPrefix(line, "@prefix") {
			l.ParsePrefix(line)
		} else {
			l.ParseCollection(line)
		}

		slog.Info(line)
	}

	return nil
}

// A RDF collection is a series of triples that share a subject
func (l *lukeFilewalker) ParseCollection(line string) {
	line = strings.TrimSpace(line)
	// Is the start of the collection/triple a IRI or a prefix, we can tell by looking for an opening <

	if !strings.HasPrefix(line, "<") {
		// It's using a prefix or literal
	} else {
		// It's a IRI
	}

}

func (l *lukeFilewalker) ParsePrefix(line string) {
	line = strings.TrimPrefix(line, "@prefix")

	line = strings.TrimSpace(line)

	prefix := strings.Split(prefixRegex.ReplaceAllString(line, "$1,$2"), ",")

	if l.graph.Prefixes == nil {
		l.graph.Prefixes = make(map[string]Prefix)
	}

	l.graph.Prefixes[prefix[0]] = Prefix(prefix[1])
}

type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "rdf: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Pointer {
		return "rdf: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "rdf: Unmarshal(nil " + e.Type.String() + ")"
}

func isPointer(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	return nil
}
