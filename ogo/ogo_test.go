package ogo_test

import (
	"testing"

	"github.com/ondbyte/ogo/ogo"
)

func TestOgo(t *testing.T) {
	s := &ogo.S{}
	type ReqData struct {
		OrderId int64
	}
	ogo.Handle[*ReqData](
		s,
		"GET",
		func(c *ogo.OCtx) *ReqData {
			data := &ReqData{}
			c.PathValue("orderId", data.OrderId)
			return data
		},
		func(c *ogo.Ctx, data *ReqData) {

		},
	)
}
