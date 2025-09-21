package fetch

import "C"
import (
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
	RequestId      string            `json:"requestId,omitempty"`
	HttpStatusCode int               `json:"httpStatusCode,omitempty"`
	Header         map[string]string `json:"header,omitempty"`
	Result         string            `json:"result,omitempty"`
	Err            string            `json:"err,omitempty"`
}

type RequestItem struct {
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

func Request(request RequestItem) (res ResponseItem, ok bool) {
	ok = true
	res.RequestId = helper.Defaults(request.Id, request.Url)
	cli, req, err := httptool.BuildClient(request.Url, helper.Defaults(request.Method, "get"), request.Query)
	if err != nil {
		res.Err = err.Error()
		return
	}
	if request.Timeout > 0 {
		cli.Timeout = time.Duration(request.Timeout) * time.Millisecond
	}
	if request.MaxRedirectNum > 0 {
		cli.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= request.MaxRedirectNum {
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
	if request.SaveFilename != "" {
		err = helper.IsDirExistAndMkdir(path.Dir(request.SaveFilename), 0666)
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
	res.Result = string(bytes)
	if request.GetResponseHeader {
		m := make(map[string]string)
		for k, v := range re.Header {
			m[k] = strings.Join(v, "; ")
		}
		res.Header = m
	}
	return
}

func ExecuteRequests(requests []RequestItem, concurrence int) (rr []ResponseItem, err error) {
	if len(requests) < 1 {
		err = errors.New("no valid request")
		return
	}
	requestStream := stream.NewStream(requests)
	if concurrence < 1 {
		concurrence = len(requests)
	}
	resStream := stream.ParallelFilterAndMap(requestStream, Request, concurrence)
	rr = resStream.Result()
	return
}
