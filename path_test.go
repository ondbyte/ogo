package ogo_test

import (
	"testing"

	"github.com/ondbyte/ogo"
)

func TestPathParams(t *testing.T) {
	path := `/users/{UserId}/orders/{OrderId}/deliveries/{DeliveryId}`
	m := ogo.PathParams(path)
	if m["UserId"] != 1 || m["OrderId"] != 1 || m["DeliveryId"] != 1 {
		t.Fail()
	}
}
