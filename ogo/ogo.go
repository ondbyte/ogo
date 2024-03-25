package ogo

import (
	"fmt"
	"net/http"
	"strconv"
)

type S struct {
}
type OCtx struct {
	ogo  bool
	req  *http.Request
	errs []string
}

var PATH_VALUE_ERR = `path value '$PV' is required and it should be '$PVT' type`

func (octx *OCtx) PathValue(
	name string,
	v any,
) {
	d := octx.req.PathValue(name)
	switch v.(type) {
	case *int64:
		vv, err := strconv.ParseInt(d, 10, 64)
		if err != nil {
			fmt.Println(err)
			octx.errs = append(octx.errs, PATH_VALUE_ERR)
			return
		}
		v = &vv
		fmt.Println(v)
	}
	panic(v)
}

type Ctx struct {
}
type OH[T any] func(c *OCtx) T
type H[T any] func(c *Ctx, data T)

func Handle[T any](
	s *S,
	method string,
	oh OH[T],
	h H[T],
) {

}
