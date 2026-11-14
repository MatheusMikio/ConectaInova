package main

import "fmt"

func main() {
	err := config.Init()

	if err != nil {
		fmt.Printf("Erro ao iniciar config: %s \n", err)
		return
	}

	db := config.GetDb()

	if err := routes.Init(db); err != nil {
		fmt.Printf("Erro ao iniciar server: %s", err)
		return
	}
}
