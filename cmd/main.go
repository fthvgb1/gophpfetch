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
	var input, single bool
	var c int
	flag.StringVar(&file, "f", "", "json file path")
	flag.StringVar(&jsonContent, "j", "", "json content")
	flag.BoolVar(&input, "i", false, "use stdin input json content")
	flag.Int("c", 1, "concurrence number")
	flag.BoolVar(&single, "s", false, "single request")
	flag.Parse()
	if file == "" && jsonContent == "" && !input && len(os.Args) < 2 {
		fmt.Println("invalid parameter")
		return
	}
	var requests any
	if single {
		requests = any(fetch.RequestItem{})
	} else {
		requests = any([]fetch.RequestItem{})
	}
	if file != "" {
		bytes, err := os.ReadFile(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = json.Unmarshal(bytes, &requests)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if jsonContent != "" {
		err := json.Unmarshal([]byte(jsonContent), &requests)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if input {
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = json.Unmarshal(bytes, &requests)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		file = os.Args[1]
		bytes, err := os.ReadFile(file)
		if err != nil {
			jsonContent = file
			err = json.Unmarshal([]byte(jsonContent), &requests)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			err = json.Unmarshal(bytes, &requests)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
	if single {
		r, _ := fetch.Request(requests.(fetch.RequestItem))
		fmt.Print(r)
		return
	}
	r, err := fetch.ExecuteRequests(requests.([]fetch.RequestItem), c)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(r)
}
