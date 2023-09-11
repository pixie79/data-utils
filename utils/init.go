// Description: Generic utils functions
// Author: Pixie79
// ============================================================================
// package utils

package utils

import (
	// 	"fmt"
	// 	"os"

	// 	"golang.org/x/exp/slog"
	tuUtils "github.com/pixie79/tiny-utils/utils"
)

var (
	// 	Logger      *slog.Logger                   // Logger is the default logger
	// 	LogLevel    string                         // logLevel is the log level
	Environment = tuUtils.GetEnvDefault("ENVIRONMENT", "dev") // Environment is the default environment
	Prefix      = tuUtils.GetEnvDefault("PREFIX", "data")     // Prefix is the default prefix
)

// // init loads the .env file and sets the log level
// func init() {
// 	LogLevel = tuUtils.GetEnvDefault("LOG_LEVEL", "INFO")
// 	Logger = initLog()
// }

// // initLog initializes the logger
// func initLog() *slog.Logger {
// 	switch LogLevel {
// 	case "DEBUG":
// 		opts := &slog.HandlerOptions{
// 			Level: slog.LevelDebug,
// 		}
// 		logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
// 		logger.Info(fmt.Sprintf("Loglevel set to: %s", LogLevel))
// 		slog.SetDefault(logger)
// 		return logger
// 	default:
// 		opts := &slog.HandlerOptions{
// 			Level: slog.LevelInfo,
// 		}
// 		logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
// 		logger.Info(fmt.Sprintf("Loglevel set to: %s", LogLevel))
// 		slog.SetDefault(logger)
// 		return logger
// 	}
// }
