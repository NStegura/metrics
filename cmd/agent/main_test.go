package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAgent(t *testing.T) {
	t.Helper()
	errs := make(chan error, 1)

	go func(errs chan<- error) {
		err := startAgent()
		errs <- err
	}(errs)

	go func(errs chan<- error) {
		time.Sleep(2 * time.Second)
		errs <- nil
	}(errs)

	assert.NoError(t, <-errs)
}
