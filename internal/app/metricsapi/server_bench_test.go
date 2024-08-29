package metricsapi

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NStegura/metrics/config"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/NStegura/metrics/internal/business"
	"github.com/NStegura/metrics/internal/repo"
)

type benchHelper struct {
	ts *httptest.Server
}

func initBenchHelper() *benchHelper {
	ctx := context.TODO()
	l := logrus.New()
	r, _ := repo.New(ctx, "", 100, "", false, l)
	businessLayer := business.New(r, l)
	server := New(config.NewSrvConfig(), businessLayer, l)
	server.ConfigRouter()

	ts := httptest.NewServer(server.Router)
	return &benchHelper{
		ts: ts,
	}
}

func (bh *benchHelper) Request(t *testing.B, method, path string, body io.Reader) (int, string) {
	t.Helper()
	req, err := http.NewRequest(method, bh.ts.URL+path, body)
	require.NoError(t, err)

	resp, _ := bh.ts.Client().Do(req)
	//require.NoError(t, err)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Log(err)
		}
	}()

	respBody, _ := io.ReadAll(resp.Body)
	//require.NoError(t, err)

	return resp.StatusCode, string(respBody)
}

func (bh *benchHelper) finish() {
	bh.ts.Close()
}

func BenchmarkGetAllMetricsHandler(b *testing.B) {
	bh := initBenchHelper()
	defer bh.finish()

	bh.Request(b, "POST", "/update/gauge/SomeGaugeMetric/1", nil)
	bh.Request(b, "POST", "/update/counter/SomeCounterMetric/1", nil)

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
		bh.Request(b, tests.method, tests.url, nil)
	})
}
