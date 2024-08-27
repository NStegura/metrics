package metric

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/NStegura/metrics/internal/app/agent/models"
	"github.com/NStegura/metrics/internal/app/metricsapi"
	"github.com/NStegura/metrics/internal/business"
	"github.com/NStegura/metrics/internal/repo"
)

type testHelper struct {
	ctrl *gomock.Controller
	ts   *httptest.Server
	cli  *Client
}

func initTestHelper(t *testing.T) *testHelper {
	t.Helper()
	ctrl := gomock.NewController(t)
	ctx := context.TODO()
	l := logrus.New()
	r, err := repo.New(ctx, "", 100, "", false, l)
	require.NoError(t, err)
	businessLayer := business.New(r, l)
	server := metricsapi.New(metricsapi.NewConfig(), businessLayer, l)
	server.ConfigRouter()

	ts := httptest.NewServer(server.Router)

	metricsCli, err := New(
		ts.URL,
		"",
		"",
		l,
	)
	require.NoError(t, err)
	return &testHelper{
		ctrl: ctrl,
		ts:   ts,
		cli:  metricsCli,
	}
}

func (th *testHelper) finish() {
	th.ts.Close()
	th.ctrl.Finish()
}

func TestUpdateCounterMetric(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	err := th.cli.UpdateCounterMetric("test_counter_metric", 1)
	require.NoError(t, err)
}

func TestUpdateGaugeMetric(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	err := th.cli.UpdateGaugeMetric("test_gauge_metric", 1)
	require.NoError(t, err)
}

func TestUpdateMetrics(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()
	testDelta := int64(1)
	testValue := 1.1
	testID := "test_metric"
	testMType := "gauge"

	tests := []struct {
		metrics []Metrics
		err     error
	}{
		{
			metrics: []Metrics{
				{
					Delta: &testDelta,
					Value: &testValue,
					ID:    testID,
					MType: testMType,
				},
			},
			err: nil,
		},
	}

	for _, v := range tests {
		err := th.cli.UpdateMetrics(v.metrics)
		require.NoError(t, err)
	}
}

func TestUpdateMetric(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	tests := []struct {
		body string
		err  error
	}{
		{
			body: `{"value": 1.002, "type": "gauge", "id": "testGaugeMetric"}`,
			err:  nil,
		},
	}

	for _, v := range tests {
		err := th.cli.UpdateMetric([]byte(v.body))
		require.NoError(t, err)
	}
}

func TestCastToMetrics(t *testing.T) {
	gm := make(map[models.MetricName]*models.GaugeMetric, 1)
	cm := make(map[models.MetricName]*models.CounterMetric, 1)
	gm["test_gm_metric"] = &models.GaugeMetric{
		Name:  "test_gauge_metric",
		Type:  "gauge",
		Value: 1.01,
	}
	cm["test_counter_metric"] = &models.CounterMetric{
		Name:  "test_counter_metric",
		Type:  "counter",
		Value: 1,
	}
	expextedGaugeValue := 1.01
	expextedCounterDelta := int64(1)

	tests := []struct {
		metrics models.Metrics
		want    Metrics
	}{
		{
			metrics: models.Metrics{
				GaugeMetrics:   gm,
				CounterMetrics: nil,
			},
			want: Metrics{
				Delta: nil,
				Value: &expextedGaugeValue,
				ID:    "test_gauge_metric",
				MType: "gauge",
			},
		},
		{
			metrics: models.Metrics{
				GaugeMetrics:   nil,
				CounterMetrics: cm,
			},
			want: Metrics{
				Delta: &expextedCounterDelta,
				Value: nil,
				ID:    "test_counter_metric",
				MType: "counter",
			},
		},
	}

	for _, v := range tests {
		metrics := CastToMetrics(v.metrics)
		assert.Equal(t, v.want, metrics[0])
	}
}
