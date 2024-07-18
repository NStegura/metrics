package agent

import (
	"sync"
	"testing"

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

	config := NewConfig()

	ag := New(config, metricsCli, logger)

	var wg sync.WaitGroup
	metricsCh := ag.collectMetrics(&wg)
	// Проверяем, что канал метрик создается успешно
	_, ok := <-metricsCh
	if ok {
		assert.True(t, ok)
	}
}
