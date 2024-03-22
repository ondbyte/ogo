package ogo_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

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

func TestOgo(t *testing.T) {
	mux := ogo.NewOgoMux(http.DefaultServeMux)
	type Body struct {
		Name string
	}
	mux.Handle(
		"GET /users/{UserId}/orders/{OrderId}/{IsOkay}",
		func(req *struct {
			UserId  ogo.Int64PathParam
			OrderId ogo.Float64PathParam
			IsOkay  ogo.BoolPathParam
			Body    *Body
		}) {
			fmt.Println(req)
		})
	go func() {
		http.ListenAndServe(":8080", http.DefaultServeMux)
	}()
	time.Sleep(time.Second)
	res, err := http.Get("http://localhost:8080/users/1/orders/2.1/1")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func TestSomeRandomShit(t *testing.T) {
	var a interface{} = struct {
		a string
	}{
		a: "a",
	}
	var b interface{} = Do
	b.(func(any))(a)
}

func Do(a struct {
	a string
}) {
	fmt.Println(a)
}
