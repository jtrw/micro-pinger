package service

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Service []Service `yaml:"services"`
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
	Webhook       string `yaml:"webhook"`
	To            string `yaml:"to"`
	Failure       int    `yaml:"failure"`
	Success       int    `yaml:"success"`
	SendOnResolve bool   `yaml:"send-on-resolve"`
}

func LoadConfig(filename string) (Config, error) {
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
