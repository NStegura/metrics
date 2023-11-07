package main

import (
	"github.com/NStegura/metrics/internal/app/agent"
	"log"
)

func main() {
	config := agent.NewConfig()
	config.ParseFlags()

	ag := agent.New(config)
	err := ag.Start()
	if err != nil {
		log.Fatal(err)
	}
}
