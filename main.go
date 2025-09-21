package main

import "C"
import (
	"encoding/json"
	"github.com/fthvgb1/gophpfetch/fetch"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/slice"
)

//export Fetch
func Fetch(s string, concurrence int, associate int8) *C.char {
	isAssociate := false
	if associate > 0 {
		isAssociate = true
	}
	requests, err := helper.JsonDecode[[]fetch.RequestItem]([]byte(s))
	if err != nil {
		return Return(fetch.Response[[]fetch.ResponseItem]{
			Err: err.Error(),
		})
	}
	rr := fetch.ExecuteRequests(requests, concurrence)
	if isAssociate {
		return Return(fetch.Response[map[string]fetch.ResponseItem]{
			Results: slice.SimpleToMap(rr.Results, func(v fetch.ResponseItem) string {
				return v.RequestId
			}),
			Err: rr.Err,
		})
	}
	return Return(rr)
}

func Return[T []fetch.ResponseItem | map[string]fetch.ResponseItem](r fetch.Response[T]) *C.char {
	re, _ := json.Marshal(r)
	return C.CString(string(re))
}

func main() {
}
