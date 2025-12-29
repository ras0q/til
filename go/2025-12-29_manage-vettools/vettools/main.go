package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"honnef.co/go/tools/analysis/facts/nilness"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/unused"
)

func main() {
	analyzers := make([]*analysis.Analyzer, 0)
	for _, analyzer := range staticcheck.Analyzers {
		analyzers = append(analyzers, analyzer.Analyzer)
	}
	analyzers = append(analyzers, unused.Analyzer.Analyzer)
	analyzers = append(analyzers, nilness.Analysis)
	analyzers = append(analyzers, shadow.Analyzer)

	multichecker.Main(analyzers...)
}
