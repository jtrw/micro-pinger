package handler

import (
	"encoding/json"
	//	"github.com/go-chi/chi/v5"
	"fmt"
	"io/ioutil"
	"log"
	"micro-pinger/v2/app/sender"
	config "micro-pinger/v2/app/service"
	"net/http"
	"strings"
	"sync"
	//"time"
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
		sendAlerts(service, false)
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
		sendAlerts(service, false)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != service.Response.Status {
		log.Printf("[%s] Unexpected response status: %d", service.Name, resp.StatusCode)
		sendAlerts(service, false)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[%s] Error reading response body: %s", service.Name, err)
		sendAlerts(service, false)
		return
	}

	if service.Response.Body != "" {
		if string(body) != service.Response.Body {
			log.Printf("[%s] Unexpected response body: %s", service.Name, string(body))
			sendAlerts(service, false)
			return
		}
	}

	log.Printf("[%s] Service is reachable and responding as expected", service.Name)
	sendAlerts(service, true)
}

func sendAlerts(service config.Service, success bool) {
	thresholdMutex.Lock()
	defer thresholdMutex.Unlock()

	for _, alert := range service.Alerts {
		alertName := service.Name + "_" + alert.Name
		if success {
			if SuccessThreshold[alertName]+1 >= alert.Success && FailureThreshold[alertName] != 0 {
				if alert.SendOnResolve {
					resolveMessage := fmt.Sprintf("[%s] Service has recovered", service.Name)
					sendAlert(alert, resolveMessage)
				}
				FailureThreshold[alertName] = 0
				SuccessThreshold[alertName] = 0
			}
			if FailureThreshold[alertName] > 0 {
				SuccessThreshold[alertName]++
			}
		} else {
			FailureThreshold[alertName]++
			log.Printf("[%s] Service Failure %d", service.Name, FailureThreshold[alert.Name])
			log.Printf("[%s] Count: %d", alertName, FailureThreshold[alertName])
			if FailureThreshold[alertName] == alert.Failure {
				message := fmt.Sprintf("[%s] Service %s", service.Name, map[bool]string{true: "recovered", false: "unreachable"}[success])
				sendAlert(alert, message)
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

func sendAlert(alert config.Alert, message string) {
	switch alert.Type {
	case "telegram":
		telegramSender := sender.NewTelegram(alert.Webhook)
		telegramSender.Send(message)
		log.Printf("[%s] Sending Telegram alert to %s: %s", alert.Name, alert.To, message)
	case "slack":
		slackSender := sender.NewSlack(alert.Webhook)
		slackSender.Send(message)
		log.Printf("[%s] Sending Slack alert to %s: %s", alert.Name, alert.To, message)
	default:
		log.Printf("[%s] Unsupported alert type: %s", alert.Name, alert.Type)
	}
}
