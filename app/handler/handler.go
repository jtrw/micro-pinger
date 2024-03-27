package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"micro-pinger/v2/app/sender"
	config "micro-pinger/v2/app/service"
	"net/http"
	"strings"
	"sync"
	"time"
)

type JSON map[string]interface{}

var (
	thresholdMutex   sync.Mutex
	FailureThreshold = make(map[string]int)
	SuccessThreshold = make(map[string]int)
)

const (
	LIMIT_MAX_FAILURE = 10000
	LIMIT_MAX_SUCCESS = 10000
)

type Handler struct {
	Services []config.Service
	Client   HTTPClient
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewHandler(services []config.Service, client HTTPClient) Handler {
	return Handler{Services: services, Client: client}
}

func (h Handler) Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	for _, service := range h.Services {
		go h.CheckService(service)
	}
	json.NewEncoder(w).Encode(JSON{"status": "ok"})
}

func (h Handler) CheckService(service config.Service) error {
	req, err := http.NewRequest(service.Method, service.URL, strings.NewReader(service.Body))

	if err != nil {
		errMsg := sender.Response{
			Text: "Error creating HTTP request",
			Code: 500,
			Err:  err,
		}
		return sendAlerts(service, errMsg)
	}
	defer req.Body.Close()
	if service.Headers != nil {
		for _, header := range service.Headers {
			req.Header.Add(header.Name, header.Value)
		}
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		errMsg := sender.Response{
			Text: "Error making HTTP request",
			Code: 500,
			Err:  err,
		}
		return sendAlerts(service, errMsg)
	}
	defer resp.Body.Close()

	if resp.StatusCode != service.Response.Status {
		errMsg := sender.Response{
			Text: "Unexpected response status",
			Code: resp.StatusCode,
			Err:  nil,
		}
		return sendAlerts(service, errMsg)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errMsg := sender.Response{
			Text: "Error reading response body",
			Code: resp.StatusCode,
			Err:  err,
		}
		return sendAlerts(service, errMsg)
	}

	if service.Response.Body != "" {
		switch {
		case service.Response.Compare == "contains":
			if !strings.Contains(string(body), service.Response.Body) {
				errMsg := sender.Response{
					Text: "Body does not contain expected string '" + service.Response.Body + "'",
					Code: resp.StatusCode,
					Err:  nil,
				}
				return sendAlerts(service, errMsg)
			}
		default:
			if string(body) != service.Response.Body {
				errMsg := sender.Response{
					Text: "Unexpected response body",
					Code: resp.StatusCode,
					Err:  nil,
				}
				return sendAlerts(service, errMsg)
			}
		}
	}

	return sendAlerts(service, sender.Response{Code: 200})
}

func sendAlerts(service config.Service, response sender.Response) error {
	thresholdMutex.Lock()
	defer thresholdMutex.Unlock()
	errs := errors.New("")
	for _, alert := range service.Alerts {
		msg := sender.Message{
			Status:      "",
			Webhook:     alert.Webhook,
			Datetime:    time.Now().Format("2006-01-02 15:04:05"),
			Url:         service.URL,
			ServiceName: service.Name,
			Response:    response,
		}

		alertName := service.Name + "_" + alert.Name
		if len(response.Text) > 0 {
			FailureThreshold[alertName]++
			if FailureThreshold[alertName] == alert.Failure {
				message := fmt.Sprintf("[%s] Service unreachable", service.Name)
				msg.Status = message
				err := sendAlert(alert, msg)
				errs = errors.Join(errs, err)
			}
		} else {
			if SuccessThreshold[alertName]+1 >= alert.Success && FailureThreshold[alertName] != 0 {
				if alert.SendOnResolve {
					resolveMessage := fmt.Sprintf("[%s] Service has recovered", service.Name)
					msg.Status = resolveMessage
					err := sendAlert(alert, msg)
					errs = errors.Join(errs, err)
				}
				FailureThreshold[alertName] = 0
				SuccessThreshold[alertName] = 0
			}
			if FailureThreshold[alertName] > 0 {
				SuccessThreshold[alertName]++
			}
		}

		if FailureThreshold[alertName] > LIMIT_MAX_FAILURE {
			FailureThreshold[alertName] = 0
		}
		if SuccessThreshold[alertName] > LIMIT_MAX_SUCCESS {
			SuccessThreshold[alertName] = 0
		}
	}

	return errs
}

func sendAlert(alert config.Alert, message sender.Message) error {
	sendService, err := sender.NewSender(alert.Type, message)
	fmt.Println("sendService", sendService)
	if err != nil {
		log.Printf("Error creating alert sender: %s", err)
		return err
	}
	err = sendService.Send()
	if err != nil {
		log.Printf("Error sending alert: %s", err)
		return err
	}
	return nil
}
