package main

import (
	"context"
	"dift_backend_driver/driver-rewards-service/internal/servicecore"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"dift_backend_driver/driver-rewards-service/config"
	natsadmin "dift_backend_driver/driver-rewards-service/internal/adapter/nats_admin"
	"dift_backend_driver/driver-rewards-service/internal/handler"
	"dift_backend_driver/driver-rewards-service/internal/route"
	"dift_backend_driver/driver-rewards-service/internal/service"
	"github.com/nats-io/nats.go"
)

func main() {
	_ = servicecore.NewEngineUnifiedBundle(servicecore.LoadEngineUnifiedConfigFromEnv("driver-rewards-service"))
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		panic(err)
	}
	logger := slog.Default()
	svc := service.NewRewardsService()
	h := handler.NewRewardsHandler(svc)
	mux := http.NewServeMux()
	route.Register(mux, h)
	if natsURL := os.Getenv("NATS_URL"); natsURL != "" {
		if nc, err := nats.Connect(natsURL); err == nil {
			consumer := natsadmin.NewAdminConsumer(svc)
			subject := os.Getenv("ADMIN_NATS_SUBJECT")
			if subject == "" {
				subject = "admin.control.driver-rewards-service.command"
			}
			_ = consumer.Subscribe(context.Background(), nc, subject)
		}
	}
	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	logger.Info("driver-rewards-service started", "port", cfg.Server.HTTPPort)
	if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
