package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
)

var (
	Logger      *slog.Logger
	logLevel    string
	Err         error
	Hostname, _ = os.Hostname()
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		slog.Debug("No .env file found")
	}
	logLevel = GetEnv("LOG_LEVEL", "INFO")
	Logger = initLog()
}

func initLog() *slog.Logger {
	switch logLevel {
	case "DEBUG":
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
		logger.Info(fmt.Sprintf("Loglevel set to: %s", logLevel))
		slog.SetDefault(logger)
		return logger
	default:
		opts := &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
		logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
		logger.Info(fmt.Sprintf("Loglevel set to: %s", logLevel))
		slog.SetDefault(logger)
		return logger
	}
}
