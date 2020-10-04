package antigate

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

const (
	libraryVersion = "0.1.0"
	defaultBaseURL = "https://api.anti-captcha.com"
	userAgent      = "vcapi/" + libraryVersion
	mediaType      = "application/json"
)

type Params map[string]string

type Config struct {
	ClientKey  string
	APIVersion string
}

type Client struct {
	// HTTP client used to communicate with the Veracross API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// Rate contains the current rate limit for the client as determined by the most recent
	// API call.
	//Rate Rate

	// Username, Password and Client
	Config *Config

	Balance BalanceService
	Task    TaskService
	logger  *zap.Logger
}

type request struct {
	ClientKey    string `json:"clientKey,omitempty"`
	TaskId       int64  `json:"taskId,omitempty"`
	SoftId       int64  `json:"softId,omitempty"`
	Task         *Task  `json:"task,omitempty"`
	LanguagePool string `json:"languagePool,omitempty"`
	CallbackUrl  string `json:"callbackUrl,omitempty"`
}

func NewClient(config *Config, logger *zap.Logger) *Client {

	// Default to API Version 2
	if config.APIVersion == "" {
		config.APIVersion = "v2"
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	//add Version and SchoolID to URL Path
	//baseURL.Path = config.SchoolID + "/" + config.APIVersion + "/"

	c := &Client{client: http.DefaultClient, BaseURL: baseURL, UserAgent: userAgent, Config: config, logger: logger}

	c.Balance = BalanceService{client: c}
	c.Task = TaskService{client: c, logger: logger}
	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash.
func (c *Client) NewRequest(urlStr string, taskId int64, task Task) (*http.Request, error) {
	method := "POST"
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var body io.Reader

	r := request{ClientKey: c.Config.ClientKey}
	if taskId != 0 {
		r.TaskId = taskId
	} else {
		r.Task = &task
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	body = bytes.NewReader(b)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", userAgent)

	return req, nil
}

// Do sends an API request and returns the API response.
func (c *Client) Do(req *http.Request, into interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(into); err != nil {
		return nil, err
	}

	return resp, nil
}
