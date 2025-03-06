package order

import (
	commonError "github.com/faceless5879/mono-go-es-ddd/internal/common/errors"
)

type OrderItem struct {
	skuID    string
	quantity int
}

func (oi OrderItem) SkuID() string {
	return oi.skuID
}

func (oi OrderItem) Quantity() int {
	return oi.quantity
}

func NewOrderItem(skuID string, quantity int) (OrderItem, error) {
	if skuID == "" {
		return OrderItem{}, commonError.NewIncorrectInputError(
			"empty sku ID",
			"empty-sku-id")
	}

	if quantity <= 0 || quantity > 99 {
		return OrderItem{}, commonError.NewIncorrectInputError(
			"invalid item quantity",
			"invalid-item-quantity")
	}

	return OrderItem{skuID: skuID, quantity: quantity}, nil
}
