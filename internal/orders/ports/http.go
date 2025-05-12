package ports

import (
	"github.com/faceless5879/mono-go-es-ddd/internal/common/server/httperr"
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/app"
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/app/command"
	"github.com/go-chi/render"
	"net/http"
)

type HttpServer struct {
	app app.Application
}

func (h HttpServer) CreateOrder(w http.ResponseWriter, r *http.Request) {
	requestBody := &CreateOrder{}
	if err := render.Decode(r, requestBody); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}
	orderItems := make([]command.OrderItem, 0)
	for _, item := range requestBody.OrderItems {
		orderItems = append(orderItems, command.OrderItem{Quantity: item.Quantity, SkuID: item.SkuID})
	}
	if err := h.app.Commands.CreateOrderHandler.Handle(r.Context(), command.CreateOrder{
		UserUUID:        requestBody.UserUUID,
		ReceiverName:    requestBody.ReceiverName,
		DeliveryAddress: requestBody.DeliveryAddress,
		OrderItems:      orderItems,
	}); err != nil {
		httperr.RespondWithSlugError(err, w, r)
	}
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{
		app: application,
	}
}
