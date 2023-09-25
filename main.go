package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/b1scuit/solid/rdf/lexer"
	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
	"github.com/b1scuit/solid/rdf/lexer/lexfn"
	"github.com/b1scuit/solid/rdf/parser"
	"github.com/olekukonko/tablewriter"
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

	lexer, _ := lexer.New(
		lexer.WihInitalState(lexfn.LexTurtleDoc),
	)

	p, _ := parser.New(
		parser.WithLexeror(lexer),
	)
	if err := p.Do(file); err != nil {
		l.Error("Error parsing file", slog.Any("error", err))
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Prefix name", "IRI"})

	for k, v := range p.GetPrefixMap() {
		table.Append([]string{k, v.Value})
	}

	table.Render()

	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Type", "Value"})
	for _, v := range p.GetLexemes() {
		table.Append([]string{lexertoken.TokenMap[v.Type], v.Value})
	}

	table.Render()
}
