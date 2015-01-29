package dawa

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// DefaultHost is the default host used for queries.
var DefaultHost = "http://dawa.aws.dk"

type parameter interface {
	Param() string
}

// A generic query structure
type query struct {
	host   string
	path   string
	params []parameter
}

// Add a key/value pair as additional parameter. It will be added as key=value on the URL.
// The values should not be delivered URL-encoded, that will be handled by the library.
func (q *query) Add(key, value string) {
	q.add(textQuery{Name: key, Values: []string{value}, Multi: false, Null: true})
}

func (q *query) add(p parameter) {
	q.params = append(q.params, p)
}

// WithHost allows overriding the host for this query.
//
// The default value is http://dawa.aws.dk
func (q *query) WithHost(s string) {
	q.host = s
}

// Replace the path of the query with something else.
func (q *query) OnPath(s string) {
	q.path = s
}

// Returns the URL for the generated query.
func (q query) URL() string {
	out := q.host + q.path
	if len(q.params) == 0 {
		return out
	}
	out += "?"
	for i, value := range q.params {
		out += value.Param()
		if i != len(q.params)-1 {
			out += "&"
		}
	}
	return out
}

type textQuery struct {
	Name   string
	Values []string
	Multi  bool
	Null   bool
}

func (t textQuery) Param() string {
	out := url.QueryEscape(t.Name) + "="
	if t.Null && len(t.Values) == 0 {
		return out
	}
	for i, val := range t.Values {
		out += url.QueryEscape(val)
		if !t.Multi {
			break
		}
		if i < len(t.Values)-1 {
			out += "|"
		}
	}
	return out
}

type RequestError struct {
	Type    string        `json:"type"`
	Title   string        `json:"title"`
	Details []interface{} `json:"details"`
	URL     string
}

func (r RequestError) Error() string {
	if r.Type == "" {
		return fmt.Sprintf("Error with request %s", r.URL)
	}
	return fmt.Sprintf("%s:%s. Details:%v. Request URL:%s", r.Type, r.Title, r.Details, r.URL)
}

// Perform the Request, and return the request result.
// If an error occurs during the request, or an error is reported
// this is returned.
// In some cases the error will be a RequestError type.
func (q query) Request() (io.ReadCloser, error) {
	url := q.URL()
	resp, err := http.Get(q.URL())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 400 {
		return resp.Body, nil
	}
	u, e2 := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if e2 != nil || len(u) == 0 {
		return nil, fmt.Errorf("Error with request %s", url)
	}
	rerr := RequestError{URL: url}
	e2 = json.Unmarshal(u, &rerr)
	return nil, rerr
}
