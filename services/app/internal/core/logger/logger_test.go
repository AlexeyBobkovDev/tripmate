package core_logger

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	t.Run("invalid level fails", func(t *testing.T) {
		folder := t.TempDir()

		logger, err := NewLogger(Config{Level: "not-a-level", Folder: folder})

		require.Error(t, err)
		require.Nil(t, logger)
	})

	t.Run("unwritable folder fails", func(t *testing.T) {
		folder := filepath.Join(t.TempDir(), "file-instead-of-dir")
		require.NoError(t, os.WriteFile(folder, []byte("x"), 0o644))

		logger, err := NewLogger(Config{Level: "DEBUG", Folder: folder})

		require.Error(t, err)
		require.Nil(t, logger)
	})

	t.Run("creates log file and writes to it", func(t *testing.T) {
		folder := t.TempDir()

		logger, err := NewLogger(Config{Level: "DEBUG", Folder: folder})
		require.NoError(t, err)
		require.NotNil(t, logger)

		logger.Info("hello from the test")
		require.NoError(t, logger.Close())

		entries, err := os.ReadDir(folder)
		require.NoError(t, err)
		require.Len(t, entries, 1, "exactly one log file must be created")
		require.True(t, strings.HasSuffix(entries[0].Name(), ".log"))

		data, err := os.ReadFile(filepath.Join(folder, entries[0].Name()))
		require.NoError(t, err)
		require.Contains(t, string(data), "hello from the test")
	})

	t.Run("level filters entries", func(t *testing.T) {
		folder := t.TempDir()

		logger, err := NewLogger(Config{Level: "WARN", Folder: folder})
		require.NoError(t, err)

		logger.Debug("too verbose")
		logger.Warn("kept")
		require.NoError(t, logger.Close())

		data, err := os.ReadFile(filepath.Join(folder, dirSingleFile(t, folder)))
		require.NoError(t, err)
		require.NotContains(t, string(data), "too verbose")
		require.Contains(t, string(data), "kept")
	})
}

func dirSingleFile(t *testing.T, folder string) string {
	t.Helper()

	entries, err := os.ReadDir(folder)
	require.NoError(t, err)
	require.Len(t, entries, 1)

	return entries[0].Name()
}

func TestClose(t *testing.T) {
	t.Run("nil logger is a no-op", func(t *testing.T) {
		var logger *Logger
		require.NoError(t, logger.Close())
	})

	t.Run("double close is safe", func(t *testing.T) {
		logger, err := NewLogger(Config{Level: "DEBUG", Folder: t.TempDir()})
		require.NoError(t, err)

		require.NoError(t, logger.Close())
		// second Close may report an error from re-closing the file,
		// but must not panic
		require.NotPanics(t, func() {
			_ = logger.Close()
		})
	})
}

func TestContextRoundtrip(t *testing.T) {
	t.Run("same logger comes back", func(t *testing.T) {
		logger, err := NewLogger(Config{Level: "DEBUG", Folder: t.TempDir()})
		require.NoError(t, err)
		defer func() { _ = logger.Close() }()

		ctx := ToContext(context.Background(), logger)

		require.Same(t, logger, FromContext(ctx))
	})

	t.Run("panics without a logger", func(t *testing.T) {
		require.Panics(t, func() {
			FromContext(context.Background())
		})
	})

	t.Run("panics with a wrong value type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), key, "not a logger")

		require.Panics(t, func() {
			FromContext(ctx)
		})
	})
}

func TestNewConfig(t *testing.T) {
	t.Run("reads level and folder from the environment", func(t *testing.T) {
		t.Setenv("LOGGER_LEVEL", "INFO")
		t.Setenv("LOGGER_FOLDER", t.TempDir())

		config, err := NewConfig()

		require.NoError(t, err)
		require.Equal(t, "INFO", config.Level)
		require.NotEmpty(t, config.Folder)
	})

	t.Run("level defaults to DEBUG", func(t *testing.T) {
		t.Setenv("LOGGER_FOLDER", t.TempDir())
		require.NoError(t, os.Unsetenv("LOGGER_LEVEL"))

		config, err := NewConfig()

		require.NoError(t, err)
		require.Equal(t, "DEBUG", config.Level)
	})

	t.Run("missing folder fails", func(t *testing.T) {
		require.NoError(t, os.Unsetenv("LOGGER_FOLDER"))
		require.NoError(t, os.Unsetenv("LOGGER_LEVEL"))

		_, err := NewConfig()

		require.Error(t, err)
	})
}
