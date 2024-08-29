package staticlint

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

const (
	defaultAnalyzersCnt = 10
)

type Checker struct {
	Cfg       *Config
	Analyzers []*analysis.Analyzer
}

// New конфигурирование анализаторов.
// По дефолту использутются staticcheck.Analyzers, simple.Analyzers и custom.NoOsExitAnalyzer.
func New(additionalAnalyzers ...*analysis.Analyzer) (*Checker, error) {
	cfg, err := NewConfig(`config.staticlint.json`)
	if err != nil {
		return nil, fmt.Errorf("failed to start analyzer config: %w", err)
	}

	checks := make(map[string]bool)
	for _, v := range cfg.StaticCheck {
		checks[v] = true
	}

	analyzers := make([]*analysis.Analyzer, 0, defaultAnalyzersCnt)

	analyzers = append(analyzers, additionalAnalyzers...)

	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] || cfg.StaticAll {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	for _, v := range simple.Analyzers {
		if checks[v.Analyzer.Name] || cfg.SimpleAll {
			analyzers = append(analyzers, v.Analyzer)
		}
	}
	logrus.Debugf("count analyzers: %v", len(analyzers))

	return &Checker{
		Analyzers: analyzers,
		Cfg:       cfg,
	}, nil
}

// Main запуск анализаторов.
func (ch *Checker) Main() {
	multichecker.Main(ch.Analyzers...)
}
