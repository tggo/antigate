package antigate

import (
	"fmt"
	"time"
	"github.com/Sirupsen/logrus"
)

const (
	createTaskBasePath = "createTask"
	getTaskbasePath = "getTaskResult"
	//taskBasePath = "test"
)

type TaskService struct {
	client *Client
}

type TaskResponse struct {
	TaskId           int64  `json:"taskId"`
	ErrorId          int64  `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`

	Status     string  `json:"status"`
	Cost       string `json:"cost"`
	FromIp     string  `json:"ip"`
	SolveCount int64   `json:"solveCount"`
	CreateTime int64   `json:"createTime"`
	EndTime    int64   `json:"endTime"`

	Solution   Solution  `json:"solution"`
}

type Solution struct {
	Text        string `json:"text"`
	GoogleResponse	string `json:"gRecaptchaResponse"`
	Url 		string `json:"url"`
}

type Task struct {
	TaskType string `json:"type,omitempty"`

	// NoCaptchaTaskProxyless
	WebsiteURL    string `json:"websiteURL,omitempty"`
	WebsiteKey    string `json:"websiteKey,omitempty"`
	WebsiteSToken string `json:"websiteSToken,omitempty"`

	// NoCaptchaTask
	ProxyType     string `json:"proxyType,omitempty"`
	ProxyAddress  string `json:"proxyAddress,omitempty"`
	ProxyPort     int64  `json:"proxyPort,omitempty"`
	ProxyLogin    string `json:"proxyLogin,omitempty"`
	ProxyPassword string `json:"proxyPassword,omitempty"`
	UserAgent     string `json:"userAgent,omitempty"`
	Cookies       string `json:"cookies,omitempty"`

	// ImageToTextTask,  other
	Body      string `json:"body,omitempty"`
	Phrase    bool   `json:"phrase,omitempty"`
	Case      bool   `json:"case,omitempty"`
	Numeric   int64  `json:"numeric,omitempty"`
	Math      bool   `json:"math,omitempty"`
	MinLength int64  `json:"minLength,omitempty"`
	MaxLength int64  `json:"maxLength,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

func (s TaskService) PutToWork(task Task) (int64, error) {
	// build url
	path := createTaskBasePath

	var taskResponse = TaskResponse{}
	req, err := s.client.NewRequest(path, 0, task)
	if err != nil {
		return 0, nil
	}
	resp, err := s.client.Do(req, &taskResponse)

	if err != nil {
		return 0, err
	}

	if taskResponse.ErrorId != 0 {
		return 0, fmt.Errorf("anti-captcha ErrorCode: %s, %s ", taskResponse.ErrorCode, taskResponse.ErrorDescription)
	}
	defer resp.Body.Close()

	return taskResponse.TaskId, nil
}

func (s TaskService) GetWork(taskId int64) (TaskResponse, error) {
	// build url
	path := getTaskbasePath

	var taskResponse = TaskResponse{}
	req, err := s.client.NewRequest(path, taskId, Task{})
	if err != nil {
		return TaskResponse{}, nil
	}
	resp, err := s.client.Do(req, &taskResponse)

	if err != nil {
		return TaskResponse{}, err
	}

	if taskResponse.ErrorId != 0 {
		return TaskResponse{}, fmt.Errorf("anti-captcha ErrorCode: %s, %s ", taskResponse.ErrorCode, taskResponse.ErrorDescription)
	}
	defer resp.Body.Close()

	return taskResponse, nil
}

func (s TaskService) GetKeyForGoogle(task Task) (string, error)  {
	responseTaskId, err := s.client.Task.PutToWork( task )
	if err != nil {
		return "", err
	}

	responseString := ""

	for sleep := 20; sleep > 0; sleep-- {
		logrus.Debugf( "wait %d sec", sleep )
		time.Sleep(time.Duration(sleep) * time.Second)
		response, err := s.client.Task.GetWork( responseTaskId)
		if err != nil {
			return "", err
		}

		if response.Status != "processing" {
			logrus.Debugf("we have response %s", response.Status)
			sleep = -10
		}

		responseString = response.Solution.GoogleResponse
	}

	return responseString, nil
}
