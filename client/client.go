package client

import (
	"errors"
	"net/http"
	"net/url"
	"fmt"
	"strings"
	"context"
	"io"
	"bytes"
	"encoding/json"
	"github.com/google/go-querystring/query"
	"net/http/httputil"
)

const (
	libraryVersion = "0.1"
	userAgent      = "go-anti-captcha/" + libraryVersion
	defaultBaseURL = "https://api.anti-captcha.com"
)

var (
	// ErrUnauthorized can be returned on any call on response status code 401.
	ErrUnauthorized = errors.New("asana: unauthorized")
)

var defaultOptFields = map[string][]string{
	"tags":       {"name", "color", "notes"},
	"users":      {"name", "email", "photo"},
	"projects":   {"name", "color", "archived"},
	"workspaces": {"name", "is_organization"},
	"tasks":      {"name", "assignee", "assignee_status", "completed", "parent"},
}

type (
	// Doer interface used for doing http calls.
	// Use it as point of setting Auth header or custom status code error handling.
	Doer interface {
		Do(req *http.Request) (*http.Response, error)
	}

	// DoerFunc implements Doer interface.
	// Allow to transform any appropriate function "f" to Doer instance: DoerFunc(f).
	DoerFunc func(req *http.Request) (resp *http.Response, err error)

	Client struct {
		doer      Doer
		BaseURL   *url.URL
		UserAgent string
	}

	Filter struct {
		Archived       bool     `url:"archived,omitempty"`
		Assignee       int64    `url:"assignee,omitempty"`
		Project        int64    `url:"project,omitempty"`
		Workspace      int64    `url:"workspace,omitempty"`
		CompletedSince string   `url:"completed_since,omitempty"`
		ModifiedSince  string   `url:"modified_since,omitempty"`
		OptFields      []string `url:"opt_fields,comma,omitempty"`
		OptExpand      []string `url:"opt_expand,comma,omitempty"`
	}


	request struct {
		ClientKey string `json:"clientKey"`
	}

	Response struct {
		Data   interface{} `json:"data,omitempty"`
		ErrorId int64  `json:"errorId"`
	}

	ResponseBalance struct {
		Data   interface{} `json:"data,omitempty"`
		Balacne float64  `json:"balance"`
		ErrorId int64  `json:"errorId"`
	}

	Error struct {
		Phrase  string `json:"phrase,omitempty"`
		Message string `json:"message,omitempty"`
	}

	// Errors always has at least 1 element when returned.
	Errors []Error
)

func (f DoerFunc) Do(req *http.Request) (resp *http.Response, err error) {
	return f(req)
}

func (e Error) Error() string {
	return fmt.Sprintf("%v - %v", e.Message, e.Phrase)
}

func (e Errors) Error() string {
	var sErrs []string
	for _, err := range e {
		sErrs = append(sErrs, err.Error())
	}
	return strings.Join(sErrs, ", ")
}

// NewClient created new asana client with doer.
// If doer is nil then http.DefaultClient used intead.
func NewClient(doer Doer) *Client {
	if doer == nil {
		doer = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	client := &Client{doer: doer, BaseURL: baseURL, UserAgent: userAgent}
	return client
}


func (c *Client) GetBalance(ctx context.Context) (float64, error) {
	balance := 0.0
	//err := c.Request(ctx, "getBalance", nil, users)
	err := c.request(ctx, "POST", "getBalance", nil, nil, nil, balance)
	return balance, err
}


func (c *Client) Request(ctx context.Context, path string, opt *Filter, v interface{}) error {
	return c.request(ctx, "POST", path, nil, nil, opt, v)
}

// request makes a request to Asana API, using method, at path, sending data or form with opt filter.
// Only data or form could be sent at the same time. If both provided form will be omitted.
// Also it's possible to do request with nil data and form.
// The response is populated into v, and any error is returned.
func (c *Client) request(ctx context.Context, method string, path string, data interface{}, form url.Values, opt *Filter, v interface{}) error {
	if opt == nil {
		opt = &Filter{}
	}
	if len(opt.OptFields) == 0 {
		// We should not modify opt provided to Request.
		newOpt := *opt
		opt = &newOpt
		opt.OptFields = defaultOptFields[path]
	}
	urlStr, err := addOptions(path, opt)
	if err != nil {
		return err
	}
	rel, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	u := c.BaseURL.ResolveReference(rel)
	var body io.Reader

		b, err := json.Marshal(request{ClientKey:"ccbdc36ec3274cf8fe7b49fc2d8733e4"})
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return err
	}


	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("User-Agent", c.UserAgent)

	// Save a copy of this request for debugging.
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	resp, err := c.doer.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	res := &ResponseBalance{}

	responsetDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(responsetDump))



	err = json.NewDecoder(resp.Body).Decode(res)

	fmt.Println(res.Balacne)

	return err
}

func addOptions(s string, opt interface{}) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}
	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}
	u.RawQuery = qs.Encode()
	return u.String(), nil
}

func toURLValues(m map[string]string) url.Values {
	values := make(url.Values)
	for k, v := range m {
		values[k] = []string{v}
	}
	return values
}