package task

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

var jobs chan job

type job struct {
	Name   string
	URL    string
	Method string
	Body   string
	Time   time.Time
}

func (j job) Consume() {
	opt := &CheckOpt{j.Method, j.Name, nil}
	ok := HealthCheck(j.URL, opt)
	if ok {
		return
	}
	p := plans[j.Name]
	callbacks <- callback{j.Name, p.Callback, p.Mail, time.Now()}
	return
}

// Checker checker is used for resp.
type Checker interface {
	Ok(*http.Response) bool
}

func status(code int) func(*http.Response) bool {
	return func(r *http.Response) bool {
		return r.StatusCode == code
	}
}

func body(b string) func(*http.Response) bool {
	return func(r *http.Response) bool {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return string(body) == b
		}
		return false
	}
}

// DefaultChecker default checker (head checker)
type DefaultChecker struct{}

// Ok is the response ok
func (d DefaultChecker) Ok(r *http.Response) bool {
	return status(200)(r)
}

var defaultChecker DefaultChecker

// CheckOpt options for checker
type CheckOpt struct {
	Method  string
	Body    string
	Checker Checker
}

// NewCheckOpt return a pointer to CheckOpt
func NewCheckOpt(method string, body string, ck Checker) *CheckOpt {
	return &CheckOpt{method, body, ck}
}

// HeadCheck simple checker
func HeadCheck(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	return status(200)(resp)
}

// HealthCheck checkers
func HealthCheck(url string, o *CheckOpt) bool {
	client := &http.Client{}
	method := o.Method
	req, _ := http.NewRequest(method, url, bytes.NewBuffer([]byte(o.Body)))
	ok := make(chan bool, 1)
	go func() {
		resp, err := client.Do(req)
		if err != nil {
			ok <- false
			return
		}
		ck := o.Checker
		if ck == nil {
			ck = &defaultChecker
		}
		ok <- ck.Ok(resp)
	}()
	select {
	case r := <-ok:
		return r
	case <-time.After(3 * time.Second): // timeout 3s
		return false
	}
}

// Consume consume jobs
func Consume() {
	for v := range jobs {
		v.Consume()
	}
}

func init() {
	go Consume()
}
