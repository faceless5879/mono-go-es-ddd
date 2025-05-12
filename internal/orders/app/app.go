package app

import (
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/app/command"
)

type Commands struct {
	command.CreateOrderHandler
}

type Queries struct {
}

type Application struct {
	Commands Commands
	Queries  Queries
}
