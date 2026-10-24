package testutils

import (
	"testing"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
)

func SetupLogger(t *testing.T, level, folder string) *core_logger.Logger {
	t.Helper()

	logger, err := core_logger.NewLogger(core_logger.Config{
		Level:  level,
		Folder: folder,
	})
	if err != nil {
		t.Fatal(err)
	}

	return logger
}
