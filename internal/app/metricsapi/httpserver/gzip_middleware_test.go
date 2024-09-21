package httpserver

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readGZIP(s string) string {
	r := strings.NewReader(s)
	reader, err := gzip.NewReader(r)
	if err != nil {
		return ""
	}
	out, err := io.ReadAll(reader)
	if err != nil {
		return ""
	}
	return string(out)
}

func TestGZIPMiddleware(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	th.Request(t, "POST", "/update/gauge/SomeGaugeMetric/1", nil, nil)

	type want struct {
		statusCode int
	}

	tests := []struct {
		method string
		name   string
		url    string
		want   want
	}{
		{
			method: http.MethodGet,
			name:   "get gauge metric",
			url:    "/value/gauge/SomeGaugeMetric",
			want: want{
				statusCode: http.StatusOK,
			},
		},
	}

	for _, v := range tests {
		statusCode, resp := th.Request(t, v.method, v.url, nil, map[string]string{"Accept-Encoding": "gzip"})
		assert.Equal(t, v.want.statusCode, statusCode)
		assert.Equal(t, "1", readGZIP(resp))
	}
}
