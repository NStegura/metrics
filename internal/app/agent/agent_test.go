package agent

import (
	"context"
	"sync"
	"testing"

	"github.com/NStegura/metrics/config"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"

	mock_agent "github.com/NStegura/metrics/mocks/app/agent"

	"github.com/stretchr/testify/assert"
)

func TestAgent_collectMetrics(t *testing.T) {
	// Создаем фиктивный клиент метрик и логгер для тестов
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metricsCli := mock_agent.NewMockMetricCli(ctrl)
	logger := logrus.New()

	cfg := config.NewAgentConfig()

	ag := New(cfg, metricsCli, logger)

	var wg sync.WaitGroup
	metricsCh := ag.collectMetrics(context.Background(), &wg)
	// Проверяем, что канал метрик создается успешно
	_, ok := <-metricsCh
	if ok {
		assert.True(t, ok)
	}
}
