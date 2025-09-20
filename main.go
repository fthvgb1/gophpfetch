package main

import "C"
import (
	"encoding/json"
	"errors"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/httptool"
	"github.com/fthvgb1/wp-go/stream"
	"io"
	"net/http"
	"path"
	"strings"
	"time"
)

type ResponseItem struct {
	RequestId      string            `json:"requestId"`
	HttpStatusCode int               `json:"httpStatusCode"`
	Header         map[string]string `json:"header"`
	Res            string            `json:"res"`
	Err            string            `json:"err"`
}

type Response struct {
	Res []ResponseItem `json:"res"`
	Err string         `json:"err"`
}

type Config struct {
	Requests    []Request `json:"requests,omitempty"`
	Concurrence int       `json:"concurrence,omitempty"`
	Associate   bool      `json:"associate"`
}

type Request struct {
	Id                string            `json:"id,omitempty"`
	Url               string            `json:"url,omitempty"`
	Method            string            `json:"method,omitempty"`
	Query             map[string]any    `json:"query,omitempty"`
	Header            map[string]string `json:"header,omitempty"`
	Body              map[string]any    `json:"body,omitempty"`
	MaxRedirectNum    int               `json:"maxRedirectNum,omitempty"`
	Timeout           int               `json:"timeout,omitempty"`
	SaveFilename      string            `json:"saveFilename,omitempty"`
	GetResponseHeader bool              `json:"getResponseHeader,omitempty"`
}

func (r *Response) Return() *C.char {
	re, _ := json.Marshal(r)
	return C.CString(string(re))
}

var rr = Response{}

//export Fetch
func Fetch(s string) *C.char {
	conf, err := helper.JsonDecode[Config]([]byte(s))
	if err != nil {
		rr.Err = err.Error()
		return rr.Return()
	}
	if len(conf.Requests) < 1 {
		rr.Err = "no valid request"
		return rr.Return()
	}
	requestStream := stream.NewStream(conf.Requests)
	concurrence := conf.Concurrence
	if concurrence < 1 {
		concurrence = len(conf.Requests)
	}
	resStream := stream.ParallelFilterAndMap(requestStream, request, concurrence)
	rr.Res = resStream.Result()
	return rr.Return()
}

func request(reqConf Request) (res ResponseItem, ok bool) {
	ok = true
	res.RequestId = helper.Defaults(reqConf.Id, reqConf.Url)
	cli, req, err := httptool.BuildClient(reqConf.Url, helper.Defaults(reqConf.Method, "get"), reqConf.Query)
	if err != nil {
		res.Err = err.Error()
		return
	}
	if reqConf.Timeout > 0 {
		cli.Timeout = time.Duration(reqConf.Timeout) * time.Millisecond
	}
	if reqConf.MaxRedirectNum > 0 {
		cli.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= reqConf.MaxRedirectNum {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		}
	}
	re, err := cli.Do(req)
	if err != nil {
		res.Err = err.Error()
		return
	}
	defer re.Body.Close()
	if reqConf.SaveFilename != "" {
		err = helper.IsDirExistAndMkdir(path.Dir(reqConf.SaveFilename), 0666)
		if err != nil {
			res.Err = err.Error()
		}
		return
	}
	bytes, err := io.ReadAll(re.Body)
	if err != nil {
		res.Err = err.Error()
		return
	}
	res.HttpStatusCode = re.StatusCode
	res.Res = string(bytes)
	if reqConf.GetResponseHeader {
		m := make(map[string]string)
		for k, v := range re.Header {
			m[k] = strings.Join(v, "; ")
		}
		res.Header = m
	}
	return
}

func main() {
}
