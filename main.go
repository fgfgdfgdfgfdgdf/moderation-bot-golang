package main

import (
	"bot/src/discord"
	"bot/src/tools"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Print("Не найден файл .env")
	}

	err := tools.InitLog()
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = tools.InitDB()
	if err != nil {
		log.Println(err.Error())
		return
	}

	discord.Init()

}
