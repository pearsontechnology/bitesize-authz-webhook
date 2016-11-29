package main

import (
	"github.com/caarlos0/env"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
)

type envConfig struct {
	ListenAddress string `env:"LISTEN_PORT" envDefault:"8080"`
	RulesFile     string `env:"RULES_CONFIG" envDefault:"rules.hcl"`
	CertKey       string `env:"SSL_CERT" envDefault:"cert.pem"`
	PrivateKey    string `env:"SSL_KEY" envDefault:"key.pem"`
	HttpsEnabled  bool   `env:"HTTPS_ENABLED"`
}

func main() {
	cfg := envConfig{}
	env.Parse(&cfg)

	Address := ":" + cfg.ListenAddress

	logged := handlers.CombinedLoggingHandler(os.Stderr, Handlers())

	if err := LoadConfigFromFile(cfg.RulesFile); err != nil {
		log.Fatalf("Failed to read config file '%s': %s", cfg.RulesFile, err)
	}

	if cfg.HttpsEnabled {
		if err := http.ListenAndServeTLS(Address, cfg.CertKey, cfg.PrivateKey, logged); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := http.ListenAndServe(Address, logged); err != nil {
			log.Fatal(err)
		}
	}
}
