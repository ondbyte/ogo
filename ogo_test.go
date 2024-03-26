package ogo_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ondbyte/ogo"
)

func BenchmarkJsonUnamrshal(b *testing.B) {
	b.ReportAllocs()
	var err2 error
	for i := 0; i < b.N; i++ {
		type Struct struct {
			Name string
			Age  int64
		}
		m := map[string]interface{}{
			"Name": "yadu",
			"Age":  int64(64),
		}
		s := &Struct{}
		bs, err := json.Marshal(&m)
		if err != nil {
			b.Fatal(err)
		}
		err2 = json.Unmarshal(bs, s)
		if err2 != nil {
			b.Fatal(err2)
		}
	}
}
func BenchmarkJsonUnamrshalWithStrData(b *testing.B) {
	b.ReportAllocs()
	var err2 error
	for i := 0; i < b.N; i++ {
		type Struct struct {
			Name string
			Age  int64
		}
		m := `{
			"Name":"Yadu",
			"Age": 64
		}`
		s := &Struct{}
		err2 = json.Unmarshal([]byte(m), s)
		if err2 != nil {
			b.Fatal(err2)
		}
	}
}

func BenchmarkReflection(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		type Struct struct {
			Name string
			Age  int64
		}
		// a map
		m := map[string]interface{}{
			"Name": "yadu",
			"Age":  int64(64),
		}
		s := &Struct{}
		val := reflect.ValueOf(s).Elem()
		typ := reflect.TypeOf(s).Elem()
		val.Field(0).SetString(m[typ.Field(0).Name].(string))
		val.Field(1).SetInt(m[typ.Field(1).Name].(int64))
	}
}

type TestBody struct {
	Name string
}

/*
BenchmarkOgo-8   	       1	1005883529 ns/op	  417264 B/op	     684 allocs/op
BenchmarkOgo-8   	       1	1005120390 ns/op	  417208 B/op	     690 allocs/op
*/

var mux = ogo.NewOgoMux(http.DefaultServeMux)
func BenchmarkOgo(t *testing.B) {
	t.ReportAllocs()
	mux.Handle(
		"GET /users/{UserId}/orders/{OrderId}/",
		func(req *struct {
			UserId   *ogo.PathParam[int64]
			OrderId  *ogo.PathParam[int64]
			HeaderId *ogo.Header[int64]
			QueryId  *ogo.QueryParam[int64]
			FormId   *ogo.FormField[int64]
			Body     *ogo.Body[TestBody]
		}) {
			fmt.Println(req.UserId.Val, req.OrderId.Val, req.Body.Val)
		})
	go func() {
		http.ListenAndServe(":8080", http.DefaultServeMux)
	}()
	time.Sleep(time.Second)
	req, err := http.NewRequest("GET", "http://localhost:8080/users/1/orders/2/1?QueryId=3&FormId=4", strings.NewReader(`{"Name":"Yadu2"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("HeaderId", "5")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

/*
BenchmarkNormalRest-8   	       1	1001962037 ns/op	  370288 B/op	     496 allocs/op
BenchmarkNormalRest-8   	       1	1004867556 ns/op	  370696 B/op	     499 allocs/op
*/
func BenchmarkNormalRest(t *testing.B) {
	t.ReportAllocs()
	e := gin.New()
	e.Handle(
		"GET",
		"/users/{UserId}/orders/{OrderId}/",
		func(ctx *gin.Context) {
			toInt := func(v string) int64 {
				i, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					panic(err)
				}
				return i
			}
			UserId := toInt(ctx.Param("UserId"))
			OrderId := toInt(ctx.Param("OrderId"))
			HeaderId := toInt(ctx.GetHeader("GeaderId"))
			QueryId := toInt(ctx.Query("QueryId"))
			FormId := toInt(ctx.Request.FormValue("FormId"))

			rs, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				panic(err)
			}
			body := &TestBody{}
			err = json.Unmarshal(rs, body)
			if err != nil {
				panic(err)
			}
			fmt.Print(UserId, OrderId, HeaderId, QueryId, FormId, body)
		},
	)
	go func() {
		http.ListenAndServe(":8080", http.DefaultServeMux)
	}()
	time.Sleep(time.Second)
	req, err := http.NewRequest("GET", "http://localhost:8080/users/1/orders/2/1?QueryId=3&FormId=4", strings.NewReader(`{"Name":"Yadu2"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("HeaderId", "5")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
