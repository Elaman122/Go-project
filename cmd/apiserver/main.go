package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/BurntSushi/toml"

	"github.com/Elaman122/Go-project/internal/app/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
