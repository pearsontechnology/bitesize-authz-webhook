package main

// import "log"
import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
)

func main() {
	listenAddress := ":8080"
	optPort := os.Getenv("LISTEN_PORT")
	if optPort != "" {
		listenAddress = ":" + os.Getenv("LISTEN_PORT")
	}

	HTTPS_ENABLED := false
	optHttpsEnabled := os.Getenv("HTTPS_ENABLED")
	parsedOutput, _ := strconv.ParseBool(optHttpsEnabled)
	if parsedOutput == true {
		HTTPS_ENABLED = true
	}

	configFile := os.Getenv("RULES_CONFIG")

	if configFile == "" {
		configFile = "rules.hcl"
	}

	logged := handlers.CombinedLoggingHandler(os.Stderr, Handlers())

	if err := LoadConfigFromFile(configFile); err != nil {
		log.Fatalf("Failed to read config file '%s': %s", configFile, err)
	}

	if HTTPS_ENABLED == true {
		if err := http.ListenAndServeTLS(listenAddress, "cert.pem", "key.pem", logged); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := http.ListenAndServe(listenAddress, logged); err != nil {
			log.Fatal(err)
		}
	}
}
