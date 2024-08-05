package main

import (
	"dhanushs3366/my-portfolio/api"
	"dhanushs3366/my-portfolio/services"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Cant load env vars")
	}

	err = services.Init()
	if err != nil {
		log.Fatal(err.Error())
	}
	h := api.Init()

	h.Run(8080)
}
