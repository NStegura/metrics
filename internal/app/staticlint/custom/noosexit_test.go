package custom

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestNoOsExitAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), NoOsExitAnalyzer, "./...")
}
