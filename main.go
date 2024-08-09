package main

import (
	"dhanushs3366/my-portfolio/handler"
	"dhanushs3366/my-portfolio/initializers"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Cant load env vars")
	}

	DB, err := initializers.Init()
	if err != nil {
		log.Fatal(err.Error())
	}
	h := handler.Init(DB)

	h.Run(8080)
}
