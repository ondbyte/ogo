package ogo_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ondbyte/ogo"
)

type ReqData struct {
	Token int64
}

type RespData struct {
	Error   string
	Success string
}

/*
BenchmarkOgo2-8   	       1	3177675126 ns/op	286402104 B/op	 1562880 allocs/op
BenchmarkOgo2-8   	       1	2359794270 ns/op	285044272 B/op	 1588718 allocs/op
BenchmarkOgo2-8   	       1	3126377003 ns/op	286909912 B/op	 1573521 allocs/op
BenchmarkOgo2-8   	       1	3046848244 ns/op	286352840 B/op	 1574939 allocs/op

*/

func TestB(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(maxReq)
	allocs := 0
	timeTaken := 0

	{
		http.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {

			// Measure the starting statistics
			var memstats runtime.MemStats
			runtime.ReadMemStats(&memstats)
			mallocs := 0 - memstats.Mallocs
			then := time.Now()
			data := &ReqData{}
			ts := r.URL.Query().Get("token")
			t, err := strconv.ParseInt(ts, 10, 64)
			if err != nil {
				panic(err)
			}
			data.Token = t
			bs, err := json.Marshal(&RespData{Success: "done dana done"})
			if err != nil {
				panic(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(bs)
			wg.Done()

			// Read the final statistics
			runtime.ReadMemStats(&memstats)
			mallocs += memstats.Mallocs
			allocs += int(mallocs)
			timeTaken += int(time.Since(then).Microseconds())
		})
	}
	go func() {
		http.ListenAndServe(":8080", http.DefaultServeMux)
	}()
	time.Sleep(time.Second)
	go func() {
		doreq()
	}()
	wg.Wait()
	fmt.Println("avg allocation per req handle", allocs/maxReq, "avg time taken per req handle", timeTaken/maxReq)
}

/*
 */
func TestA(t *testing.T) {
	o := ogo.NewServer(nil)
	var wg sync.WaitGroup
	wg.Add(maxReq)
	allocs := 0
	timeTaken := 0
	{
		ogo.SetupHandler[ReqData, RespData](
			o, "GET", "/users/me",
			func(ctx *ogo.RequestValidator[ReqData, RespData], data *ReqData) {
				// Measure the starting statistics
				var memstats runtime.MemStats
				runtime.ReadMemStats(&memstats)
				mallocs := 0 - memstats.Mallocs
				then := time.Now()
				ctx.Description(
					func() string {
						return "### description of the endpoint"
					},
				)
				ctx.QueryParam(
					"token",
					&data.Token,
					func(param *ogo.Param) {
						param.Description("description of the param 'token'")
						param.Required(http.StatusNotFound, "the query param 'token' is required")
					},
				)

				// Read the final statistics
				runtime.ReadMemStats(&memstats)
				mallocs += memstats.Mallocs
				allocs += int(mallocs)
				timeTaken += int(time.Since(then).Microseconds())
			},
			func(validatedStatus int, validatedErr string) (resp *ogo.Response[RespData]) {
				return nil
			},

			func(reqData *ReqData) (resp *ogo.Response[RespData]) {
				return nil
			},
		)
	}
	go func() {
		o.Run(8080, func(info *ogo.ServerInfo) {})
	}()
	time.Sleep(time.Second)
	go func() {
		doreq()
	}()
	wg.Wait()
	fmt.Println("avg allocation per req handle", allocs/maxReq, "avg time taken per req handle", timeTaken/maxReq)
}

var maxReq = 10000

func doreq() {
	var wg sync.WaitGroup
	for i := 0; i < maxReq; i++ {
		wg.Add(1)
		req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/users/me?token=%v", i), strings.NewReader(`{"Name":"Yadu2"}`))
		if err != nil {
			panic(err)
		}
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		//fmt.Println(res)
		wg.Done()
	}
	wg.Wait()
}
