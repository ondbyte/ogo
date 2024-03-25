package ogo_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
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

type TestBody struct {
	Name string
}

var mux = ogo.NewOgoMux(http.DefaultServeMux)

var maxReq = 10000

func doreq() {
	var wg sync.WaitGroup
	for i := range maxReq {
		wg.Add(1)
		req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/users/%v/orders/2?QueryId=3&FormId=4", i), strings.NewReader(`{"Name":"Yadu2"}`))
		if err != nil {
			panic(err)
		}
		req.Header.Set("HeaderId", "5")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
		wg.Done()
	}
	wg.Wait()
}

/*
1	6054114890 ns/op	158732632 B/op	 2133714 allocs/op
*/

func BenchmarkOgo(t *testing.B) {
	var wg sync.WaitGroup
	wg.Add(maxReq)
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
			fmt.Println(req)
			wg.Done()
		})
	go func() {
		http.ListenAndServe(":8080", http.DefaultServeMux)
	}()
	time.Sleep(time.Second)
	go func() {
		doreq()
	}()
	wg.Wait()
}

/*
1	3699736187 ns/op	85722296 B/op	 1141446 allocs/op
*/
func BenchmarkNormalRest(t *testing.B) {
	t.ReportAllocs()
	var wg sync.WaitGroup
	wg.Add(maxReq)
	http.DefaultServeMux.HandleFunc(
		"GET /users/{UserId}/orders/{OrderId}",
		func(w http.ResponseWriter, r *http.Request) {
			defer wg.Done()
			toInt := func(v string) int64 {
				i, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					panic(err)
				}
				return i
			}
			UserId := toInt(r.PathValue("UserId"))
			OrderId := toInt(r.PathValue("OrderId"))
			HeaderId := toInt(r.Header.Get("HeaderId"))
			QueryId := toInt(r.URL.Query().Get("QueryId"))
			FormId := toInt(r.FormValue("FormId"))

			rs, err := io.ReadAll(r.Body)
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
	go func() {
		doreq()
	}()
	wg.Wait()
}
