package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

type Location struct {
	// Where to dispatch the request
	URL string
}

type HttpConfig struct {
	Address string
	Port    int
}

type Config struct {
	Locations map[string][]Location
	Http      HttpConfig
}

func dispatch(cfg Config) http.HandlerFunc {
	client := http.Client{}
	return func(w http.ResponseWriter, r *http.Request) {
		locations, ok := cfg.Locations[r.URL.Path]
		if !ok {
			log.Printf("Unconfigured URL path got called: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(r.Body)
		}

		log.Printf("Dispatching request for path: %s", r.URL.Path)
		for _, location := range locations {
			req, err := http.NewRequest(r.Method, location.URL, ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
			if err != nil {
				log.Printf("Invalid request in config for path %s: %s", r.URL.Path, err)
				continue
			}
			_, err = client.Do(req)
			if err != nil {
				log.Printf("Request failed: %s", err)
			}
		}

        w.Write([]byte("OK\n"))
	}
}

func main() {
	cfg := LoadConfig()

	http.Handle("/", dispatch(cfg))

	log.Printf("Listening on address %s:%v", cfg.Http.Address, cfg.Http.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%v", cfg.Http.Address, cfg.Http.Port), nil))
}

func LoadConfig() (config Config) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	viper.SetDefault("http.address", "0.0.0.0")
	viper.SetDefault("http.port", "8080")

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode configuration file: %s", err)
	}

	return
}
