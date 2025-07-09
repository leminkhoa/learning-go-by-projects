package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/slack-io/slacker"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger := newLogger()

	bot := slacker.NewClient(
		os.Getenv("SLACK_BOT_TOKEN"),
		os.Getenv("SLACK_APP_TOKEN"),
		slacker.WithLogger(logger),
	)

	// Register a command
	bot.AddCommand(&slacker.CommandDefinition{
		Description: "Calculate age from year of birth",
		Command:     "my yob is <year>",
		Examples:    []string{"my yob is 2020"},
		Handler: func(ctx *slacker.CommandContext) {
			year, err := strconv.Atoi(ctx.Request().Param("year"))
			if err != nil {
				ctx.Response().Reply("Invalid year")
				return
			}
			age := 2025 - year
			r := fmt.Sprintf("You are %d years old.", age)
			ctx.Response().Reply(r)

			logger.Info("Command event",
				"timestamp", ctx.Event().TimeStamp,
				"command", ctx.Event().Text,
			)

		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start listening to Slack events
	if err := bot.Listen(ctx); err != nil {
		fmt.Println("Error:", err)
	}

}

type MyLogger struct {
	logger *slog.Logger
}

func newLogger() *MyLogger {
	return &MyLogger{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (l *MyLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *MyLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *MyLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *MyLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
