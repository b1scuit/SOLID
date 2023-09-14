package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/b1scuit/solid/rdf/parser"
)

func main() {
	slog.Info("Running RDF Parser")

	var fileName string
	flag.StringVar(&fileName, "file", "example_rdf.ttl", "Filename to open")
	flag.Parse()

	if fileName == "" {
		slog.Error("No Filename provided")
		return
	}

	l := slog.With(
		slog.String("file name", fileName),
	)

	file, err := os.Open(fileName)
	if err != nil {
		l.Error("Error opening file", slog.String("error", err.Error()))
		return
	}

	defer file.Close()

	p, _ := parser.New()
	p.Do(file)
}
