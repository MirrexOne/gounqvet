package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	internal "github.com/MirrexOne/unqueryvet/internal/analyzer"
)

func main() {
	singlechecker.Main(internal.NewAnalyzer())
}
