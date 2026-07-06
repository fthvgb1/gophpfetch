package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fthvgb1/gophpfetch/fetch"
	"io"
	"os"
)

func main() {
	var file, jsonContent string
	var input, single, rr bool
	var c int
	flag.StringVar(&file, "f", "", "json file path")
	flag.StringVar(&jsonContent, "j", "", "json content")
	flag.BoolVar(&input, "i", false, "use stdin input json content")
	flag.Int("c", 1, "concurrence number")
	flag.BoolVar(&single, "s", false, "single request")
	flag.BoolVar(&rr, "r", false, "directly output result not to json")
	flag.Parse()
	if file == "" && jsonContent == "" && !input && len(os.Args) < 2 {
		fmt.Println("invalid parameter")
		return
	}
	var requests any
	bytesFn := func(bytes []byte) (err error) {
		if single {
			var a fetch.RequestItem
			err = json.Unmarshal(bytes, &a)
			requests = a
			return
		}
		var a []fetch.RequestItem
		err = json.Unmarshal(bytes, &a)
		requests = a
		return
	}
	fileFn := func(path string) error {
		bytes, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		return bytesFn(bytes)
	}
	var err error
	if file != "" {
		err = fileFn(file)
	} else if jsonContent != "" {
		err = bytesFn([]byte(jsonContent))
	} else if input {
		bytes, er := io.ReadAll(os.Stdin)
		if er != nil {
			err = er
			return
		}
		err = bytesFn(bytes)
	} else {
		file = os.Args[len(os.Args)-1]
		_ = fileFn(file)
		err = bytesFn([]byte(file))
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	if single {
		r, _ := fetch.Request(requests.(fetch.RequestItem))
		if rr {
			fmt.Print(r.Result)
			return
		}
		bytes, err := json.Marshal(r)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(string(bytes))
		return
	}
	r, err := fetch.ExecuteRequests(requests.([]fetch.RequestItem), c)
	if err != nil {
		fmt.Println(err)
		return
	}
	if rr {
		for _, item := range r {
			fmt.Println(item.Result)
		}
		return
	}
	bytes, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(string(bytes))
}
