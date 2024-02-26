package handler

import (
	"encoding/json"
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
}

func NewHandler(services []config.Service) Handler {
	return Handler{Services: services}
}

func (h Handler) Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	for _, service := range h.Services {
		go checkService(service)
	}
	json.NewEncoder(w).Encode(JSON{"status": "ok"})
}

func checkService(service config.Service) {
	log.Printf("[%s] Checking service...", service.Name)

	client := &http.Client{}
	req, err := http.NewRequest(service.Method, service.URL, strings.NewReader(service.Body))
	defer req.Body.Close()
	if err != nil {
		log.Printf("[%s] Error creating HTTP request: %s", service.Name, err)
		errMsg := sender.Response{
			Text: "Error creating HTTP request",
			Code: 500,
			Err:  err,
		}
		sendAlerts(service, errMsg)
		return
	}

	if service.Headers != nil {
		for _, header := range service.Headers {
			req.Header.Add(header.Name, header.Value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[%s] Error making HTTP request", service.Name)
		errMsg := sender.Response{
			Text: "Error making HTTP request",
			Code: 500,
			Err:  err,
		}
		sendAlerts(service, errMsg)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != service.Response.Status {
		log.Printf("[%s] Unexpected response status: %d", service.Name, resp.StatusCode)
		errMsg := sender.Response{
			Text: "Unexpected response status",
			Code: resp.StatusCode,
			Err:  nil,
		}
		sendAlerts(service, errMsg)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[%s] Error reading response body: %s", service.Name, err)
		errMsg := sender.Response{
			Text: "Error reading response body",
			Code: resp.StatusCode,
			Err:  err,
		}
		sendAlerts(service, errMsg)
		return
	}

	if service.Response.Body != "" {
		if string(body) != service.Response.Body {
			log.Printf("[%s] Unexpected response body: %s", service.Name, string(body))
			errMsg := sender.Response{
				Text: "Unexpected response body",
				Code: resp.StatusCode,
				Err:  nil,
			}
			sendAlerts(service, errMsg)
			return
		}
	}

	log.Printf("[%s] Service is reachable and responding as expected", service.Name)
	sendAlerts(service, sender.Response{Code: 200})
}

func sendAlerts(service config.Service, response sender.Response) {
	thresholdMutex.Lock()
	defer thresholdMutex.Unlock()

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
				sendAlert(alert, msg)
			}
		} else {
			if SuccessThreshold[alertName]+1 >= alert.Success && FailureThreshold[alertName] != 0 {
				if alert.SendOnResolve {
					resolveMessage := fmt.Sprintf("[%s] Service has recovered", service.Name)
					msg.Status = resolveMessage
					sendAlert(alert, msg)
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
}

func sendAlert(alert config.Alert, message sender.Message) {
	sendService := sender.NewSender(alert.Type, message)
	err := sendService.Send()
	if err != nil {
		log.Printf("Error sending alert: %s", err)
	}
}
