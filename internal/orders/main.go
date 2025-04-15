package main

import (
	"github.com/faceless5879/mono-go-es-ddd/internal/common/logs"
	_ "github.com/lib/pq"
)

func main() {
	logs.Init()
	//app, cleanup := service.NewApplication(ctx)
	//defer cleanup()
	//
	//servers.RunHTTPServer(func(router chi.Router) http.Handler {
	//	return ports.HandlerFromMux(ports.NewHttpServer(app), router)
	//})
}
