package main

import (
	"fmt"
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
    Port int
}

type Config struct {
    Locations map[string]Location
    Http HttpConfig
}

func dispatch(cfg Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

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
