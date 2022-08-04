package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// http request
type Request struct {
	Label    string
	Method   string
	Protocol string
	Host     string
	Path     string
	Body     string
	Headers  map[string]string
	Timeout  int64
	Check    string
}

// http response
type Response struct {
	Name    string
	Code    int
	Headers http.Header
	Body    string
	Time    int64
}

func (req *Request) Call() (*Response, error) {
	api := req.Protocol + "://" + req.Host + req.Path
	nr, err := http.NewRequest(req.Method, api, strings.NewReader(req.Body))
	if err != nil {
		return nil, err
	}
	if req.Headers != nil {
		for k, v := range req.Headers {
			nr.Header.Add(k, v)
		}
	}
	sTime := time.Now()
	client := &http.Client{Timeout: time.Duration(req.Timeout) * time.Second}
	resp, err := client.Do(nr)
	if err != nil {
		return nil, err
	}
	eTime := time.Now()
	defer resp.Body.Close()
	bodyByte, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(bodyByte)
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	if req.Check != "" {
		yes := strings.Contains(bodyStr, req.Check)
		if !yes {
			msg := fmt.Sprintf("(%s) not found in %s",
				req.Check, bodyStr)
			return nil, errors.New(msg)
		}
	}
	duration := eTime.Sub(sTime).Milliseconds()
	//
	if req.Label == "" {
		req.Label = fmt.Sprintf("%s-%s",
			req.Method, req.Path)
	}
	//
	res := &Response{
		Name:    req.Label,
		Code:    resp.StatusCode,
		Headers: resp.Header,
		Body:    bodyStr,
		Time:    duration}
	return res, nil
}
