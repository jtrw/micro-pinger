package main

import (
	"context"
	//	"fmt"
	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	server "micro-pinger/v2/app/server"
	//"net/http"
	"os"
	"os/signal"
	//"strings"
	"syscall"
	"time"
)

type Config struct {
	Sercies map[string]interface{} `yaml:"services"`
}

type Options struct {
	Listen         string        `short:"l" long:"listen" env:"LISTEN_SERVER" default:":8080" description:"listen address"`
	Secret         string        `short:"s" long:"secret" env:"SECRET_KEY" default:"123"`
	PinSize        int           `long:"pinszie" env:"PIN_SIZE" default:"5" description:"pin size"`
	MaxExpire      time.Duration `long:"expire" env:"MAX_EXPIRE" default:"24h" description:"max lifetime"`
	MaxPinAttempts int           `long:"pinattempts" env:"PIN_ATTEMPTS" default:"3" description:"max attempts to enter pin"`
	WebRoot        string        `long:"web" env:"WEB" default:"/" description:"web ui location"`
}

var revision string

func main() {
	log.Printf("Pinger %s\n", revision)

	var opts Options
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	// Зчитування конфігураційного файлу
	config, err := readConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if x := recover(); x != nil {
			log.Printf("[WARN] run time panic:\n%v", x)
			panic(x)
		}

		// catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Printf("[WARN] interrupt signal")
		cancel()
	}()

	services := config.Sercies

	srv := server.Server{
		Listen:         opts.Listen,
		PinSize:        opts.PinSize,
		MaxExpire:      opts.MaxExpire,
		MaxPinAttempts: opts.MaxPinAttempts,
		WebRoot:        opts.WebRoot,
		Secret:         opts.Secret,
		Version:        revision,
		Services:       services,
	}
	if err := srv.Run(ctx); err != nil {
		log.Printf("[ERROR] failed, %+v", err)
	}
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
