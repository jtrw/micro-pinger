package server

import (
	"context"
	"crypto/sha256"
	"log"
	config "micro-pinger/v2/app/service"
	"net/http"

	"github.com/go-chi/render"
)

var (
	hashConfig []byte
)

func ReloadConfigMiddleware(cfg config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//get SHA256 of file config.yml
			hasher := sha256.New()
			hasher.Write([]byte("config.yml"))
			hash := hasher.Sum(nil)
			// if string(hash) == string(hashConfig) {
			// 	log.Printf("[INFO] config not changed")
			// 	next.ServeHTTP(w, r)
			// 	return
			// }
			hashConfig = hash
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
