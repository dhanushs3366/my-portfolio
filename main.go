package main

import (
	"dhanushs3366/my-portfolio/handler"
	"dhanushs3366/my-portfolio/services/db"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Cant load env vars")
	}

	store, err := db.Init()
	if err != nil {
		log.Fatal(err.Error())
	}
	h := handler.Init(store)

	h.Run(8080)
}
