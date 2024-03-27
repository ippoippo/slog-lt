package main

import (
	"log"
	"log/slog"
	"os"
)

func main() {
	// Default loggers
	log.Print("1a: I am a basic out of the box `log`")                                  // Pre Go1.21 slog formatting
	slog.Info("1b: I am a basic out of the box `slog`", slog.String("arg", "argValue")) // Only difference to `log` is additional INFO, and the appended text

	// Create a text handler
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}
	slogger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	log.Print("2a: After creating a TextHandler slogger `log`")                                                 // Default logger: Same as 1a
	slog.Info("2b: After creating a TextHandler slogger `slog`", slog.String("arg", "argValue"))                // Default slogger: Same as 2a
	slogger.Info("2c: After creating a TextHandler slogger `slogger instance`", slog.String("arg", "argValue")) // TextHandler output

	// Change the default logger
	slog.SetDefault(slogger)
	log.Print("3a: After changing the Default `log`")                                  // Changes to TextHandler format
	slog.Info("3b: After changing the Default `slog`", slog.String("arg", "argValue")) // "Same" as 3a

	slogger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(slogger)
	log.Print("4a: After changing the Default to JSON Handler `log`")                                  // Changes to JSON format
	slog.Info("4b: After changing the Default to JSON Handler `slog`", slog.String("arg", "argValue")) // Changes to JSON format
}
