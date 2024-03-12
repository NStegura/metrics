package metricsapi

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig__ok(t *testing.T) {
	config := NewConfig()

	assert.Equal(t, config.Restore, false)
	assert.Equal(t, config.BindAddr, ":8080")
	assert.Equal(t, config.LogLevel, "debug")
	assert.Equal(t, config.FileStoragePath, "/tmp/metrics-db.json")
	assert.Equal(t, config.DatabaseDSN, "")
	assert.Equal(t, config.RequestKey, "")
	assert.Equal(t, config.StoreInterval, defaultStoreInerval)
}

func TestConfig__ParseEnvs(t *testing.T) {
	err := os.Setenv("ADDRESS", ":8082")
	require.NoError(t, err)
	err = os.Setenv("LOG_LEVEL", "debug")
	require.NoError(t, err)
	err = os.Setenv("FILE_STORAGE_PATH", "/testpath/metrics-db.json")
	require.NoError(t, err)
	err = os.Setenv("DATABASE_DSN", "postgresql://localhost/mydb?user=other&password=secret")
	require.NoError(t, err)
	err = os.Setenv("STORE_INTERVAL", "1")
	require.NoError(t, err)
	err = os.Setenv("RESTORE", "true")
	require.NoError(t, err)
	err = os.Setenv("KEY", "somekey")
	require.NoError(t, err)

	config := NewConfig()
	err = config.ParseFlags()
	require.NoError(t, err)

	assert.Equal(t, config.Restore, true)
	assert.Equal(t, config.BindAddr, ":8082")
	assert.Equal(t, config.LogLevel, "debug")
	assert.Equal(t, config.FileStoragePath, "/testpath/metrics-db.json")
	assert.Equal(t, config.DatabaseDSN, "postgresql://localhost/mydb?user=other&password=secret")
	assert.Equal(t, config.RequestKey, "somekey")
	assert.Equal(t, config.StoreInterval, time.Second)
}
