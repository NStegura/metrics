package main

import (
	"github.com/NStegura/metrics/internal/app/agent"
	"log"
)

func main() {
	ag := agent.New()
	err := ag.Start()
	if err != nil {
		log.Fatal(err)
	}
}
