package metricsapi

import (
	"github.com/NStegura/metrics/internal/business"
	"github.com/NStegura/metrics/internal/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestUpdateGaugeMetricHandler(t *testing.T) {
	r := repo.New()
	businessLayer := business.New(r)
	server := New(NewConfig(), businessLayer)
	server.configRouter()

	ts := httptest.NewServer(server.router)
	defer ts.Close()

	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		method string
		name   string
		url    string
		want   want
	}{
		{
			method: http.MethodPost,
			name:   "update gauge metric",
			url:    "/update/gauge/SomeMetric/1.2",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			method: http.MethodPost,
			name:   "update unknown type metric",
			url:    "/update/dsfds/SomeMetric/1.2",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			method: http.MethodPost,
			name:   "update",
			url:    "/update/",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, v := range tests {
		resp, _ := testRequest(t, ts, "POST", v.url)
		defer resp.Body.Close()
		assert.Equal(t, v.want.statusCode, resp.StatusCode)
	}
}
