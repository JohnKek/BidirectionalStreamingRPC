package main

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}
func main() {
	logger.Info("Starting server")
	err := StartGrpcServer()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
