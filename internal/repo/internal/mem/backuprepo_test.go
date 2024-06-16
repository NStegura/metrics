package mem

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestBackupRepo__Init(t *testing.T) {
	l := logrus.New()
	_, err := NewBackupRepo(0, "", l)
	assert.NoError(t, err)
}

func TestBackupRepo__CreateCounterMetric(t *testing.T) {
	l := logrus.New()
	repo, err := NewBackupRepo(0, "", l)
	assert.NoError(t, err)

	err = repo.CreateCounterMetric(context.TODO(), "test_name", "test_type", 1)
	assert.NoError(t, err)
}

func TestBackupRepo__CreateGaugeMetric(t *testing.T) {
	l := logrus.New()
	repo, err := NewBackupRepo(0, "", l)
	assert.NoError(t, err)

	err = repo.CreateGaugeMetric(context.TODO(), "test_name", "test_type", 1.01)
	assert.NoError(t, err)
}

func TestBackupRepo__UpdateCounterMetric(t *testing.T) {
	l := logrus.New()
	repo, err := NewBackupRepo(0, "", l)
	assert.NoError(t, err)

	err = repo.CreateCounterMetric(context.TODO(), "test_name", "test_type", 1)
	assert.NoError(t, err)

	err = repo.UpdateCounterMetric(context.TODO(), "test_name", 1)
	assert.NoError(t, err)
}

func TestBackupRepo__UpdateGaugeMetric(t *testing.T) {
	l := logrus.New()
	repo, err := NewBackupRepo(0, "", l)
	assert.NoError(t, err)

	err = repo.CreateGaugeMetric(context.TODO(), "test_name", "test_type", 1.01)
	assert.NoError(t, err)

	err = repo.UpdateGaugeMetric(context.TODO(), "test_name", 1.01)
	assert.NoError(t, err)
}

func TestBackupRepo__Shutdown(t *testing.T) {
	l := logrus.New()
	repo, err := NewBackupRepo(0, "", l)
	assert.NoError(t, err)

	err = repo.CreateGaugeMetric(context.TODO(), "test_name", "test_type", 1.01)
	assert.NoError(t, err)

	err = repo.UpdateGaugeMetric(context.TODO(), "test_name", 1.01)
	assert.NoError(t, err)

	repo.Shutdown(context.TODO())
}
