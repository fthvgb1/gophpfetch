package main

import "C"
import (
	"encoding/json"
	"github.com/fthvgb1/gophpfetch/fetch"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/slice"
)

type Response[Data []fetch.ResponseItem | map[string]fetch.ResponseItem] struct {
	Results Data   `json:"results"`
	Err     string `json:"err"`
}

//export Fetch
func Fetch(s string, concurrence int, associate int8) *C.char {
	isAssociate := false
	if associate > 0 {
		isAssociate = true
	}
	requests, err := helper.JsonDecode[[]fetch.RequestItem]([]byte(s))
	if err != nil {
		return Return(Response[[]fetch.ResponseItem]{
			Err: err.Error(),
		})
	}
	rr, err := fetch.ExecuteRequests(requests, concurrence)
	er := ""
	if err != nil {
		er = err.Error()
	}
	if isAssociate {
		return Return(Response[map[string]fetch.ResponseItem]{
			Results: slice.SimpleToMap(rr, func(v fetch.ResponseItem) string {
				return v.RequestId
			}),
			Err: er,
		})
	}
	return Return(Response[[]fetch.ResponseItem]{
		Results: rr,
		Err:     er,
	})
}

func Return[T []fetch.ResponseItem | map[string]fetch.ResponseItem](r Response[T]) *C.char {
	re, _ := json.Marshal(r)
	return C.CString(string(re))
}

func main() {
}
