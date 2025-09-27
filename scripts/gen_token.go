package main

import (
	"fmt"
	"log"

	"github.com/nabind47/sensor_api/internal/config"
	"github.com/nabind47/sensor_api/internal/util"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	token := util.GenerateHash(cfg.Auth.ClientID, cfg.Auth.ClientSecret)
	fmt.Println(token)
}
