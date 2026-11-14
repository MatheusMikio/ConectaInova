package main

import (
	"log"

	"conectainova/config"
	"conectainova/internal/routes"
)

func main() {
	err := config.Init()

	if err != nil {
		log.Fatalf("Erro ao iniciar config: %s", err)
	}

	db := config.GetDb()

	if err := routes.Init(db); err != nil {
		log.Fatalf("Erro ao iniciar server: %s", err)
	}
}
