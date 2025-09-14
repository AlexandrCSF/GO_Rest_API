package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
	"wb_cource/internal/app/apiserver"
	"wb_cource/internal/app/config"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()
	cfg := config.NewConfig()
	_, err := toml.DecodeFile(configPath, cfg)

	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(cfg); err != nil {
		log.Fatal(err)
	}
}
