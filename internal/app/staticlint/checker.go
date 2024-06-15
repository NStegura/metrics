package staticlint

import (
	"fmt"
	"github.com/NStegura/metrics/internal/app/staticlint/custom"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

type Checker struct {
	Analyzers []*analysis.Analyzer
	Cfg       *Config
}

func New() (*Checker, error) {
	cfg, err := NewConfig(`staticlint.config.json`)
	if err != nil {
		return nil, fmt.Errorf("failed to start analyzer config: %w", err)
	}
	logrus.Debugf("cfg: %v", cfg)

	checks := make(map[string]bool)
	for _, v := range cfg.StaticCheck {
		checks[v] = true
	}

	analyzers := make([]*analysis.Analyzer, 0, 30)
	analyzers = append(analyzers,
		custom.NoOsExitAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer)

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

func (ch *Checker) Main() {
	multichecker.Main(ch.Analyzers...)
}
