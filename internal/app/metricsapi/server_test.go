package metricsapi

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/NStegura/metrics/internal/business"
	"github.com/NStegura/metrics/internal/repo"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	t.Helper()
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Log(err)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestUpdateGaugeMetricHandler(t *testing.T) {
	ctx := context.TODO()
	l := logrus.New()
	r, err := repo.New(ctx, "", 100, "", false, l)
	require.NoError(t, err)
	businessLayer := business.New(r, l)
	server := New(NewConfig(), businessLayer, l)
	server.configRouter()

	ts := httptest.NewServer(server.router)
	defer ts.Close()

	type want struct {
		statusCode int
	}

	tests := []struct {
		method string
		name   string
		url    string
		body   string
		want   want
	}{
		{
			method: http.MethodPost,
			name:   "update gauge metric",
			url:    "/update/gauge/SomeGaugeMetric/1.2",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			method: http.MethodPost,
			name:   "update unknown type metric",
			url:    "/update/dsfds/SomeGaugeMetric/1.2",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			method: http.MethodPost,
			name:   "update",
			url:    "/update/dsfds",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, v := range tests {
		resp, _ := testRequest(t, ts, "POST", v.url, bytes.NewBufferString(v.body))
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Log(err)
			}
		}()
		assert.Equal(t, v.want.statusCode, resp.StatusCode)
	}
}

func TestUpdateCounterMetricHandler(t *testing.T) {
	ctx := context.TODO()
	l := logrus.New()
	r, err := repo.New(ctx, "", 100, "", false, l)
	require.NoError(t, err)
	businessLayer := business.New(r, l)
	server := New(NewConfig(), businessLayer, l)
	server.configRouter()

	ts := httptest.NewServer(server.router)
	defer ts.Close()

	type want struct {
		statusCode int
	}

	tests := []struct {
		method string
		name   string
		url    string
		body   string
		want   want
	}{
		{
			method: http.MethodPost,
			name:   "update counter metric",
			url:    "/update/counter/SomeCounterMetric/1",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			method: http.MethodPost,
			name:   "update unknown type metric",
			url:    "/update/dsfds/SomeCounterMetric/1",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			method: http.MethodPost,
			name:   "update",
			url:    "/update/dsfds",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, v := range tests {
		resp, _ := testRequest(t, ts, "POST", v.url, bytes.NewBufferString(v.body))
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Log(err)
			}
		}()
		assert.Equal(t, v.want.statusCode, resp.StatusCode)
	}
}

func TestUpdateMetricHandler(t *testing.T) {
	ctx := context.TODO()
	l := logrus.New()
	r, err := repo.New(ctx, "", 100, "", false, l)
	require.NoError(t, err)
	businessLayer := business.New(r, l)
	server := New(NewConfig(), businessLayer, l)
	server.configRouter()

	successBodyGauge := `{"value": 1.002, "type": "gauge", "id": "testGaugeMetric"}`
	successBodyCounter := `{"delta": 1, "type": "counter", "id": "testCounterMetric"}`
	unknownMetricTypeBody := `{"delta": 1, "type": "sdssa", "id": "testCounterMetric"}`

	ts := httptest.NewServer(server.router)
	defer ts.Close()

	type want struct {
		statusCode int
	}

	tests := []struct {
		method string
		name   string
		url    string
		body   string
		want   want
	}{
		{
			method: http.MethodPost,
			name:   "update gauge metric",
			url:    "/update/",
			body:   successBodyGauge,
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			method: http.MethodPost,
			name:   "update counter metric",
			url:    "/update/",
			body:   successBodyCounter,
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			method: http.MethodPost,
			name:   "update bad metric",
			url:    "/update/",
			body:   unknownMetricTypeBody,
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}
	for _, v := range tests {
		resp, _ := testRequest(t, ts, "POST", v.url, bytes.NewBufferString(v.body))
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Log(err)
			}
		}()
		assert.Equal(t, v.want.statusCode, resp.StatusCode)
	}
}

func TestGetCounterMetricHandler__empty(t *testing.T) {
	ctx := context.TODO()
	l := logrus.New()
	r, err := repo.New(ctx, "", 100, "", false, l)
	require.NoError(t, err)
	businessLayer := business.New(r, l)
	server := New(NewConfig(), businessLayer, l)
	server.configRouter()

	ts := httptest.NewServer(server.router)
	defer ts.Close()

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
			name:   "get counter metric empty",
			url:    "/value/counter/SomeCounterMetric",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, v := range tests {
		resp, _ := testRequest(t, ts, v.method, v.url, nil)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Log(err)
			}
		}()
		assert.Equal(t, v.want.statusCode, resp.StatusCode)
	}
}

func TestGetCounterMetricHandler__ok(t *testing.T) {
	ctx := context.TODO()
	l := logrus.New()
	r, err := repo.New(ctx, "", 100, "", false, l)
	require.NoError(t, err)
	businessLayer := business.New(r, l)
	server := New(NewConfig(), businessLayer, l)
	server.configRouter()

	ts := httptest.NewServer(server.router)
	defer ts.Close()

	testRequest(t, ts, "POST", "/update/counter/SomeCounterMetric/1", nil)

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
			name:   "get counter metric empty",
			url:    "/value/counter/SomeCounterMetric",
			want: want{
				statusCode: http.StatusOK,
			},
		},
	}
	for _, v := range tests {
		resp, _ := testRequest(t, ts, v.method, v.url, nil)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Log(err)
			}
		}()
		assert.Equal(t, v.want.statusCode, resp.StatusCode)
	}
}

func TestGetGaugeMetricHandler__empty(t *testing.T) {
	ctx := context.TODO()
	l := logrus.New()
	r, err := repo.New(ctx, "", 100, "", false, l)
	require.NoError(t, err)
	businessLayer := business.New(r, l)
	server := New(NewConfig(), businessLayer, l)
	server.configRouter()

	ts := httptest.NewServer(server.router)
	defer ts.Close()

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
			name:   "get gauge metric empty",
			url:    "/value/gauge/SomeGaugeMetric",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, v := range tests {
		resp, _ := testRequest(t, ts, v.method, v.url, nil)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Log(err)
			}
		}()
		assert.Equal(t, v.want.statusCode, resp.StatusCode)
	}
}

func TestGetGaugeMetricHandler__ok(t *testing.T) {
	ctx := context.TODO()
	l := logrus.New()
	r, err := repo.New(ctx, "", 100, "", false, l)
	require.NoError(t, err)
	businessLayer := business.New(r, l)
	server := New(NewConfig(), businessLayer, l)
	server.configRouter()

	ts := httptest.NewServer(server.router)
	defer ts.Close()

	testRequest(t, ts, "POST", "/update/gauge/SomeGaugeMetric/1", nil)

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
		resp, _ := testRequest(t, ts, v.method, v.url, nil)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Log(err)
			}
		}()
		assert.Equal(t, v.want.statusCode, resp.StatusCode)
	}
}

func TestGetAllMetricsHandler__ok(t *testing.T) {
	ctx := context.TODO()
	l := logrus.New()
	r, err := repo.New(ctx, "", 100, "", false, l)
	require.NoError(t, err)
	businessLayer := business.New(r, l)
	server := New(NewConfig(), businessLayer, l)
	server.configRouter()

	ts := httptest.NewServer(server.router)
	defer ts.Close()

	testRequest(t, ts, "POST", "/update/gauge/SomeGaugeMetric/1", nil)
	testRequest(t, ts, "POST", "/update/counter/SomeCounterMetric/1", nil)

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
			name:   "get metrics",
			url:    "/",
			want: want{
				statusCode: http.StatusOK,
			},
		},
	}
	for _, v := range tests {
		resp, _ := testRequest(t, ts, v.method, v.url, nil)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Log(err)
			}
		}()
		assert.Equal(t, v.want.statusCode, resp.StatusCode)
	}
}

func TestGetAllMetricsHandler__empty(t *testing.T) {
	ctx := context.TODO()
	l := logrus.New()
	r, err := repo.New(ctx, "", 100, "", false, l)
	require.NoError(t, err)
	businessLayer := business.New(r, l)
	server := New(NewConfig(), businessLayer, l)
	server.configRouter()

	ts := httptest.NewServer(server.router)
	defer ts.Close()

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
			name:   "get metrics not found",
			url:    "/",
			want: want{
				statusCode: http.StatusOK,
			},
		},
	}
	for _, v := range tests {
		resp, _ := testRequest(t, ts, v.method, v.url, nil)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Log(err)
			}
		}()
		assert.Equal(t, v.want.statusCode, resp.StatusCode)
	}
}
