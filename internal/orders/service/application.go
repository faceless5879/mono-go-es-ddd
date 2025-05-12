package service

import (
	"fmt"
	"github.com/faceless5879/mono-go-es-ddd/internal/common/metrics"
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/adapters"
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/app"
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/app/command"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
)

func getEnvOr(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
func NewApplication() app.Application {
	DB_HOST := getEnvOr("DB_HOST", "localhost")
	DB_PORT := getEnvOr("DB_PORT", "5432")
	DB_USER := getEnvOr("DB_USER", "admin")
	DB_PASS := getEnvOr("DB_PASS", "password")
	DB_NAME := getEnvOr("PG_NAME", "pg")
	SSL_MD := getEnvOr("SSL_MD", "disable")
	dataSrcStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME, SSL_MD)
	connection, err := sqlx.Connect("postgres", dataSrcStr)
	if err != nil {
		panic(err)
	}
	orderRepository := adapters.NewOrderRdbRepository(connection)
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}
	return app.Application{Commands: app.Commands{
		CreateOrderHandler: command.NewCreateOrderHandler(orderRepository, logger, metricsClient),
	}}
}
