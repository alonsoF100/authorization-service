package main

import (
	"log"

	"github.com/alonsoF100/authorization-service/internal/config"
)

func main() {
	cfg := config.Load()
	log.Println(*cfg)
}
