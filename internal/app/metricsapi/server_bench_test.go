package metricsapi

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/NStegura/metrics/internal/business"
	"github.com/NStegura/metrics/internal/repo"
)

func testRequestBench(t *testing.B, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
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

func BenchmarkGetAllMetricsHandler(b *testing.B) {
	ctx := context.TODO()
	l := logrus.New()
	r, _ := repo.New(ctx, "", 100, "", false, l)
	businessLayer := business.New(r, l)
	server := New(NewConfig(), businessLayer, l)
	server.configRouter()

	ts := httptest.NewServer(server.router)
	defer ts.Close()

	testRequestBench(b, ts, "POST", "/update/gauge/SomeGaugeMetric/1", nil)
	testRequestBench(b, ts, "POST", "/update/counter/SomeCounterMetric/1", nil)

	type want struct {
		statusCode int
	}

	tests := struct {
		method string
		name   string
		url    string
		want   want
	}{
		method: http.MethodGet,
		name:   "get metrics",
		url:    "/",
		want:   want{statusCode: http.StatusOK},
	}

	b.ResetTimer()
	b.Run("", func(b *testing.B) {
		resp, _ := testRequestBench(b, ts, tests.method, tests.url, nil)
		defer func() {
			if err := resp.Body.Close(); err != nil {
				b.Log(err)
			}
		}()
	})
}
