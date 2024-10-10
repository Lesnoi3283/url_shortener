// Package main is a package with static analysers.
package main

import (
	shorteneranalysers "github.com/Lesnoi3283/url_shortener/cmd/staticlint/analysers"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/staticcheck"
	"strings"
)

// main func runs static analysers.
// Use command line args to choose specific analyser.
func main() {

	analysers := make([]*analysis.Analyzer, 0)
	analysers = append(analysers, nilness.Analyzer, shadow.Analyzer, unmarshal.Analyzer, unusedresult.Analyzer, copylock.Analyzer, structtag.Analyzer)
	analysers = append(analysers, shorteneranalysers.OsExitAnalyzer)

	for _, analyser := range staticcheck.Analyzers {
		if strings.HasPrefix(analyser.Analyzer.Name, "ST") {
			analysers = append(analysers, analyser.Analyzer)
		} else if analyser.Analyzer.Name == "SA4006" ||
			analyser.Analyzer.Name == "SA5000" {
			analysers = append(analysers, analyser.Analyzer)
		}
	}

	multichecker.Main(
		analysers...,
	)
}
