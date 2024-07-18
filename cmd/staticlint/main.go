package main

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"

	"github.com/NStegura/metrics/internal/app/staticlint"
)

func main() {
	checker, err := staticlint.New(
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer)
	if err != nil {
		logrus.Fatal(err)
	}
	checker.Main()
}
