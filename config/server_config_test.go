package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig__ok(t *testing.T) {
	cfg := NewSrvConfig()

	assert.Equal(t, cfg.Restore, false)
	assert.Equal(t, cfg.BindAddr, ":8080")
	assert.Equal(t, cfg.LogLevel, "debug")
	assert.Equal(t, cfg.FileStoragePath, "/tmp/metrics-db.json")
	assert.Equal(t, cfg.DatabaseDSN, "")
	assert.Equal(t, cfg.BodyHashKey, "")
	assert.Equal(t, cfg.StoreInterval, defaultStoreInerval)
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

	cfg := NewSrvConfig()
	err = cfg.ParseFlags()
	require.NoError(t, err)

	assert.Equal(t, cfg.Restore, true)
	assert.Equal(t, cfg.BindAddr, ":8082")
	assert.Equal(t, cfg.LogLevel, "debug")
	assert.Equal(t, cfg.FileStoragePath, "/testpath/metrics-db.json")
	assert.Equal(t, cfg.DatabaseDSN, "postgresql://localhost/mydb?user=other&password=secret")
	assert.Equal(t, cfg.BodyHashKey, "somekey")
	assert.Equal(t, time.Duration(cfg.StoreInterval), time.Second)
}
