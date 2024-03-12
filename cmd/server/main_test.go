package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	addr = "http://127.0.0.1:8080"
)

func TestInitServer(t *testing.T) {
	t.Helper()
	go func() {
		main()
	}()
	time.Sleep(time.Second)
	gaugeValue := 123.1231
	url := fmt.Sprintf(
		"%s/update/gauge/test/%v",
		addr, gaugeValue)

	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = url

	_, err := req.Send()
	require.NoError(t, err)

	url = fmt.Sprintf(
		"%s/value/gauge/test",
		"http://127.0.0.1:8080")

	req = resty.New().R()
	req.Method = http.MethodGet
	req.URL = url

	res, err := req.Send()
	require.NoError(t, err)
	assert.Equal(t, string(res.Body()), "123.1231")
}
