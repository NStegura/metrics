package metricsapi

import (
	"github.com/NStegura/metrics/internal/bll"
	"github.com/NStegura/metrics/internal/dal"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateGaugeMetricHandler(t *testing.T) {
	repo := dal.New()
	businessLayer := bll.New(repo)
	server := New(NewConfig(), businessLayer)

	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		method  string
		name    string
		request string
		want    want
	}{
		{
			method:  http.MethodPost,
			name:    "update gauge metric",
			request: "/update/gauge/SomeMetric/1.2",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			method:  http.MethodPost,
			name:    "update counter metric",
			request: "/update/counter/SomeMetric/1.2",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			method:  http.MethodPost,
			name:    "update unknown metric",
			request: "/update/dsfds/SomeMetric/1.2",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			method:  http.MethodPost,
			name:    "update",
			request: "/update/",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)

			w := httptest.NewRecorder()
			h := server.updateGaugeMetric()
			h(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
