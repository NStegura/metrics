package main

import (
	"github.com/NStegura/metrics/internal/app/staticlint"
	"github.com/sirupsen/logrus"
)

func main() {
	checker, err := staticlint.New()
	if err != nil {
		logrus.Fatal(err)
	}
	checker.Main()
}
