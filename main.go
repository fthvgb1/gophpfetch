package main

import "C"
import (
	"encoding/json"
	"github.com/fthvgb1/gophpfetch/fetch"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/slice"
)

type Response struct {
	Results any    `json:"results"`
	Err     string `json:"err"`
}

//export Fetch
func Fetch(s string, concurrence int, associate bool) *C.char {
	requests, err := helper.JsonDecode[[]fetch.RequestItem]([]byte(s))
	r := Response{}
	if err != nil {
		r.Err = err.Error()
		return Return(r)
	}
	rr, err := fetch.ExecuteRequests(requests, concurrence)
	if err != nil {
		r.Err = err.Error()
		return Return(r)
	}
	r.Results = rr
	if associate {
		r.Results = slice.SimpleToMap(rr, func(v fetch.ResponseItem) string {
			return v.RequestId
		})
	}
	return Return(r)
}

func Return(r Response) *C.char {
	re, _ := json.Marshal(r)
	return C.CString(string(re))
}

func main() {
}
