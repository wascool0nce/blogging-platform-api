package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/wascool0nce/blogging-platform-api/pkg/httpserver"
	"github.com/wascool0nce/blogging-platform-api/pkg/logger"
	"github.com/wascool0nce/blogging-platform-api/pkg/postgres"
)

func Run() error {
	log := logger.NewLogger("debug")
	pg, err := postgres.New(buildUrlConected())
	if err != nil {
		return fmt.Errorf("dont connetct postgres: %w", err)
	}
	defer pg.Close()

	log.Debug("connection db open")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		log.Info("request")
		w.WriteHeader(http.StatusOK)
	})

	server := httpserver.New(log, ":8080", mux)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	server.Start()

	select {
	case sig := <-interrupt:
		log.Info("Получен сигнал остановки", "signal", sig.String())
	case <-server.Notify():
		log.Info("HTTP-сервер завершил работу")
	}

	if err := server.Shutdown(); err != nil {
		return fmt.Errorf("app - httpserver: %w", err)
	}

	log.Info("Приложение остановлено")
	return nil
}

func buildUrlConected() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		"user",
		"password",
		"localhost",
		"5432",
		"postgres",
	)
}
