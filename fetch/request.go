package fetch

import (
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/httptool"
	"github.com/fthvgb1/wp-go/stream"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"
)

type ResponseItem struct {
	RequestId      string              `json:"requestId,omitempty"`
	HttpStatusCode int                 `json:"httpStatusCode,omitempty"`
	Header         []map[string]string `json:"header,omitempty"`
	Result         string              `json:"result,omitempty"`
	Err            string              `json:"err,omitempty"`
}

type RequestItem struct {
	Id                string            `json:"id"`
	Url               string            `json:"url"`
	Method            string            `json:"method"`
	Query             map[string]any    `json:"query"`
	Header            map[string]string `json:"header"`
	Body              map[string]any    `json:"body"`
	NoReturn          bool              `json:"noReturn"`
	Host              string            `json:"host"`
	Jar               bool              `json:"jar"`
	MaxRedirectNum    int               `json:"maxRedirectNum"`
	Timeout           int               `json:"timeout"`
	SaveFile          FileSave          `json:"saveFile"`
	GetResponseHeader bool              `json:"getResponseHeader"`
}

type FileSave struct {
	Path    string `json:"path"`
	Mode    string `json:"mode"`
	DirMode string `json:"dirMode"`
}

func getMode(mode string) (os.FileMode, error) {
	p, err := strconv.ParseUint(mode, 8, 32)
	if err != nil {
		return 0, err
	}
	return os.FileMode(p), nil
}

func saveFile(request RequestItem, body io.ReadCloser) error {
	fileMode, err := getMode(helper.Defaults(request.SaveFile.Mode, os.Getenv("uploadFileMod"), "0666"))
	if err != nil {
		return err
	}
	dirMode, err := getMode(helper.Defaults(request.SaveFile.DirMode, os.Getenv("uploadDirMod"), "0755"))
	if err != nil {
		return err
	}
	err = helper.IsDirExistAndMkdir(path.Dir(request.SaveFile.Path), dirMode)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(request.SaveFile.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, body)
	if err != nil {
		return err
	}
	return nil
}

func setCheckRedirect(jar bool, num int, cli *http.Client) {
	if jar {
		j, _ := cookiejar.New(nil)
		cli.Jar = j
	}
	cli.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= num {
			return fmt.Errorf("stopped after %d redirects", num)
		}

		if !jar {
			return nil
		}
		for k, v := range req.Response.Header {
			if k != "Set-Cookie" {
				req.Header[k] = v
			}
		}
		re := via[len(via)-1]
		if req.Response.Header.Get("Set-Cookie") == "" {
			return nil
		}
		icookies := make(map[string][]*http.Cookie)
		for _, cookie := range req.Response.Header["Set-Cookie"] {
			item := strings.Split(cookie, "=")
			cc := http.Cookie{Name: item[0], Value: item[1]}
			icookies[item[0]] = []*http.Cookie{&cc}
		}

		for _, c := range re.Cookies() {
			if _, ok := icookies[c.Name]; !ok {
				icookies[c.Name] = []*http.Cookie{c}
			}
		}
		var ss []string
		for _, cs := range icookies {
			for _, c := range cs {
				ss = append(ss, c.Name+"="+c.Value)
			}
		}
		slices.Sort(ss) // Ensure deterministic headers
		req.Header.Set("Cookie", strings.Join(ss, "; "))
		return nil
	}
}

func Request(request RequestItem) (res ResponseItem, ok bool) {
	ok = true
	res.RequestId = helper.Defaults(request.Id, request.Url)
	cli, req, err := httptool.BuildClient(request.Url, helper.Defaults(request.Method, "get"), request.Query)
	if err != nil {
		res.Err = err.Error()
		return
	}
	if request.Host != "" {
		req.Host = request.Host
	}
	if request.Jar || request.MaxRedirectNum > 0 {
		setCheckRedirect(request.Jar, helper.Defaults(request.MaxRedirectNum, 10), cli)
	}
	if request.Timeout > 0 {
		cli.Timeout = time.Duration(request.Timeout) * time.Millisecond
	}

	req.Header = http.Header{}
	if len(request.Header) > 0 {
		for k, v := range request.Header {
			req.Header.Set(k, v)
		}
	}
	var fns []func()

	if len(request.Body) > 0 {
		if err = SetBody(request, req, &fns); err != nil {
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
	res.HttpStatusCode = re.StatusCode
	if request.GetResponseHeader {
		resp := re
		for {
			m := make(map[string]string)
			for k, v := range resp.Header {
				m[k] = strings.Join(v, "; ")
			}
			res.Header = append(res.Header, m)
			if resp.Request.Response == nil {
				break
			}
			resp = re.Request.Response
		}
	}
	if request.SaveFile.Path != "" {
		if err = saveFile(request, re.Body); err != nil {
			res.Err = err.Error()
		}
		return
	}
	if request.NoReturn {
		return
	}
	bytes, err := io.ReadAll(re.Body)
	if err != nil {
		res.Err = err.Error()
		return
	}

	res.Result = string(bytes)
	return
}

var typeInt = map[string]int{
	"x-www-form-urlencoded": 1,
	"form-data":             2,
	"json":                  3,
	"plain":                 4,
	"binary":                5,
}

var contentMap = map[int]string{
	3: "application/json",
	4: "text/plain",
}

func SetBody(r RequestItem, req *http.Request, fns *[]func()) (err error) {
	t, ok := typeInt[r.Header["Content-Type"]]
	if !ok && strings.ToLower(r.Method) == "post" {
		t = 1
	}
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
	case 3, 4:
		if d, ok := r.Body["__Data"]; ok && d != "" {
			b := strings.NewReader(d.(string))
			req.Body = io.NopCloser(b)
			req.ContentLength = int64(b.Len())
			req.Header.Set("Content-Type", contentMap[t])
			return nil
		}
	}
	err = httptool.SetBody(req, t, r.Body)
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
