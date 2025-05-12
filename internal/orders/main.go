package main

import (
	"github.com/faceless5879/mono-go-es-ddd/internal/common/logs"
	"github.com/faceless5879/mono-go-es-ddd/internal/common/servers"
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/ports"
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/service"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {
	logs.Init()
	application := service.NewApplication()
	servers.RunHTTPServer(func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(
			ports.NewHttpServer(application),
			router,
		)
	})
}
