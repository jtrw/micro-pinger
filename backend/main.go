package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	Services []Service `yaml:"services"`
}

type Service struct {
	Name     string   `yaml:"name"`
	URL      string   `yaml:"url"`
	Method   string   `yaml:"method"`
	Type     string   `yaml:"type"`
	Body     string   `yaml:"body"`
	Interval string   `yaml:"interval"`
	Headers  []Header `yaml:"headers"`
	Response Response `yaml:"response"`
	Alerts   []Alert  `yaml:"alerts"`
}

type Header struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Response struct {
	Status int    `yaml:"status"`
	Body   string `yaml:"body"`
}

type Alert struct {
	Name          string `yaml:"name"`
	Type          string `yaml:"type"`
	To            string `yaml:"to"`
	Failure       int    `yaml:"failure"`
	Success       int    `yaml:"success"`
	SendOnResolve bool   `yaml:"send-on-resolve"`
}

func main() {
	// Зчитування конфігураційного файлу
	config, err := readConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	// Запуск сервісу для кожного сервісу з конфігурації
	for _, service := range config.Services {
		go monitorService(service)
	}

	// Некінцевий цикл для утримання головного потоку активним
	select {}
}

func readConfig(filename string) (Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func monitorService(service Service) {
	interval, err := time.ParseDuration(service.Interval)
	if err != nil {
		log.Printf("[%s] Error parsing interval: %s", service.Name, err)
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkService(service)
		}
	}
}

func checkService(service Service) {
	log.Printf("[%s] Checking service...", service.Name)

	// Виконати HTTP-запит до сервісу
	client := &http.Client{}
	req, err := http.NewRequest(service.Method, service.URL, strings.NewReader(service.Body))
	if err != nil {
		log.Printf("[%s] Error creating HTTP request: %s", service.Name, err)
		return
	}

	// Додати заголовки до запиту
	for _, header := range service.Headers {
		req.Header.Add(header.Name, header.Value)
	}

	// Виконати HTTP-запит
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[%s] Error making HTTP request: %s", service.Name, err)
		sendAlerts(service, false)
		return
	}
	defer resp.Body.Close()

	// Перевірити статус відповіді
	if resp.StatusCode != service.Response.Status {
		log.Printf("[%s] Unexpected response status: %d", service.Name, resp.StatusCode)
		sendAlerts(service, false)
		return
	}

	// Перевірити тіло відповіді
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[%s] Error reading response body: %s", service.Name, err)
		sendAlerts(service, false)
		return
	}

	// if string(body) != service.Response.Body {
	// 	log.Printf("[%s] Unexpected response body: %s", service.Name, string(body))
	// 	sendAlerts(service, false)
	// 	return
	// }

	log.Printf("[%s] Service is reachable and responding as expected", service.Name)
	sendAlerts(service, true)
}

func sendAlerts(service Service, success bool) {
	for _, alert := range service.Alerts {
		if success {
			if alert.Success > 0 {
				alert.Success--
				continue
			}
		} else {
			if alert.Failure > 0 {
				alert.Failure--
				continue
			}
		}

		message := fmt.Sprintf("[%s] Service %s", service.Name, map[bool]string{true: "recovered", false: "unreachable"}[success])
		sendAlert(alert, message)

		if success && alert.SendOnResolve {
			// Відправити додатковий алерт про відновлення
			resolveMessage := fmt.Sprintf("[%s] Service has recovered", service.Name)
			sendAlert(alert, resolveMessage)
		}
	}
}

func sendAlert(alert Alert, message string) {
	switch alert.Type {
	case "email":
		// Реалізуйте відправку електронної пошти
		log.Printf("[%s] Sending email alert to %s: %s", alert.Name, alert.To, message)
	case "slack":
		// Реалізуйте відправку повідомлення в Slack
		log.Printf("[%s] Sending Slack alert to %s: %s", alert.Name, alert.To, message)
	default:
		log.Printf("[%s] Unsupported alert type: %s", alert.Name, alert.Type)
	}
}
