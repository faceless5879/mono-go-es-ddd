package test_test

import (
	"github.com/faceless5879/mono-go-es-ddd/internal/orders/domain/order"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrder_NewOrderStatusFromString(t *testing.T) {
	t.Parallel()
	orderStatus := "PENDING"

	os, err := order.NewOrderStatusFromString(orderStatus)
	require.NoError(t, err)

	assert.Equal(t, os.String(), orderStatus)
}

func TestOrder_NewOrderStatusFromString_invalid(t *testing.T) {
	t.Parallel()
	orderStatus := "invalid status"

	_, err := order.NewOrderStatusFromString(orderStatus)
	assert.Error(t, err)
}
