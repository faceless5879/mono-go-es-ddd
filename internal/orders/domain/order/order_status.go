package order

import (
	"fmt"
	commonError "github.com/faceless5879/mono-go-es-ddd/internal/common/errors"
)

// OrderStatus use struct to ensure immutability
type OrderStatus struct {
	s string
}

func (os OrderStatus) IsEmpty() bool {
	return os == OrderStatus{}
}

func (os OrderStatus) String() string {
	return os.s
}

var (
	Pending   = OrderStatus{"PENDING"}
	Paid      = OrderStatus{"PAID"}
	Cancelled = OrderStatus{"CANCELLED"}
	Fulfilled = OrderStatus{"FULFILLED"}
	Failed    = OrderStatus{"FAILED"}
)

func NewOrderStatusFromString(orderStatus string) (OrderStatus, error) {
	switch orderStatus {
	case "PENDING":
		return Pending, nil
	case "PAID":
		return Paid, nil
	case "CANCELLED":
		return Cancelled, nil
	case "FULFILLED":
		return Fulfilled, nil
	case "FAILED":
		return Failed, nil
	}

	return OrderStatus{s: orderStatus}, commonError.NewIncorrectInputError(
		fmt.Sprintf("invalid '%s' status", orderStatus),
		"invalid-status")
}
