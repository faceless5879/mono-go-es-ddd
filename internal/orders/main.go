package main

import (
	"github.com/faceless5879/mono-go-es-ddd/internal/common/logs"
)

func main() {
	logs.Init()

	//ctx := context.Background()

	//app, cleanup := service.NewApplication(ctx)
	//defer cleanup()
	//
	//servers.RunHTTPServer(func(router chi.Router) http.Handler {
	//	return ports.HandlerFromMux(ports.NewHttpServer(app), router)
	//})
}
