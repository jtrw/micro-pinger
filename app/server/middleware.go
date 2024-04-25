package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	config "micro-pinger/v2/app/service"
	"net/http"

	"github.com/go-chi/render"
)

var (
	hashConfig string
)

func ReloadConfigMiddleware(cfg config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//get SHA256 of file config.yml
			hasher := sha256.New()
			//open file
			file, err := ioutil.ReadFile("config.yml")
			if err != nil {
				log.Printf("[ERROR] failed to read file, %v", err)
				render.Status(r, http.StatusInternalServerError)
			}
			hasher.Write(file)
			hash := hasher.Sum(nil)
			hashStr := hex.EncodeToString(hash)
			log.Printf("[INFO] hash: %s", hashStr)
			log.Printf("[INFO] hashConfig: %s", hashConfig)
			if hashStr == hashConfig {
				log.Printf("[INFO] config not changed")
				// next.ServeHTTP(w, r)
				// return
			}
			hashConfig = hashStr
			//reload config
			config, err := config.LoadConfig("config.yml") // Load config from file

			if err != nil {
				log.Printf("[ERROR] failed to load config, %v", err)
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, `{Message: "failed to load config"}`)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), "config", config))
			cfg = config
			log.Printf("[INFO] config reloaded")
			//
			next.ServeHTTP(w, r)
		})
	}
}
