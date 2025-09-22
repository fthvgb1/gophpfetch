package fetch

import "C"
import (
	"errors"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/httptool"
	"github.com/fthvgb1/wp-go/stream"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
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
	if len(request.Header) > 0 {
		req.Header = http.Header{}
		for k, v := range request.Header {
			req.Header.Set(k, v)
		}
	}
	var fns []func()

	if len(request.Body) > 0 {
		err = SetBody(request, req, &fns)
		if err != nil {
			res.Err = err.Error()
			return
		}
	}
	defer func() {
		if len(fns) > 0 {
			for _, fn := range fns {
				fn()
			}
		}
	}()
	re, err := cli.Do(req)
	if err != nil {
		res.Err = err.Error()
		return
	}
	defer re.Body.Close()
	bytes, err := io.ReadAll(re.Body)
	if err != nil {
		res.Err = err.Error()
		return
	}
	if request.SaveFilename != "" {
		file := strings.Split(request.SaveFilename, "|")
		name := file[0]
		perm := helper.Defaults(os.Getenv("uploadFileMod"), "0666")
		if len(file) > 1 {
			perm = file[1]
		}
		p, err := strconv.ParseUint(perm, 8, 32)
		if err != nil {
			res.Err = err.Error()
			return
		}
		mod := os.FileMode(p)

		dirMode := helper.Defaults(os.Getenv("uploadDirMod"), "0755")
		d, err := strconv.ParseUint(dirMode, 8, 32)
		if err != nil {
			res.Err = err.Error()
			return
		}
		err = helper.IsDirExistAndMkdir(path.Dir(name), os.FileMode(d))
		if err != nil {
			res.Err = err.Error()
			return
		}

		err = os.WriteFile(name, bytes, mod)
		if err != nil {
			res.Err = err.Error()
		}
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

var typeInt = map[string]int{
	"x-www-form-urlencoded": 1,
	"form-data":             2,
	"json":                  3,
	"binary":                4,
}

func SetBody(r RequestItem, req *http.Request, fns *[]func()) (err error) {
	if t, ok := typeInt[r.Header["Content-Type"]]; ok {
		switch t {
		case 2:
			if files, ok := r.Body["__uploadFiles"].(map[string]any); ok && files != nil {
				delete(r.Body, "__uploadFiles")
				for filename, field := range files {
					fd, err := os.Open(filename)
					if err != nil {
						return err
					}
					r.Body[field.(string)] = fd
					*fns = append(*fns, func() {
						fd.Close()
					})
				}
			}
		case 3:
			if d, ok := r.Body["jsonData"]; ok && d != "" {
				b := strings.NewReader(d.(string))
				req.Body = io.NopCloser(b)
				req.ContentLength = int64(b.Len())
				req.Header.Set("Content-Type", "application/json")
				return nil
			}
		}
		err = httptool.SetBody(req, t, r.Body)
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
