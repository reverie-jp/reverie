package main

import (
	"log/slog"
	"os"

	"github.com/reverie-jp/reverie/internal/application/server"
)

func main() {
	if err := server.Run(); err != nil {
		slog.Error("api server exited with error", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
