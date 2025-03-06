package order

import (
	commonError "github.com/faceless5879/mono-go-es-ddd/internal/common/errors"
)

type DeliveryAddress struct {
	receiverName string
	address      string
}

func NewDeliveryAddress(receiverName string, address string) (DeliveryAddress, error) {
	if receiverName == "" {
		return DeliveryAddress{}, commonError.NewIncorrectInputError(
			"empty receiver name",
			"empty-receiver-name")
	}

	if address == "" {
		return DeliveryAddress{}, commonError.NewIncorrectInputError(
			"empty address", "empty-address")
	}

	return DeliveryAddress{receiverName: receiverName, address: address}, nil
}

func (da DeliveryAddress) ReceiverName() string {
	return da.receiverName
}

func (da DeliveryAddress) Address() string {
	return da.address
}
