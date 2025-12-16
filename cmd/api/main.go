package main

import (
	"log"

	"github.com/alonsoF100/authorization-service/internal/config"
	"github.com/alonsoF100/authorization-service/internal/logger"
)

func main() {
	cfg := config.Load()
	log.Println(*cfg)

	logger.Setup(cfg)
}
