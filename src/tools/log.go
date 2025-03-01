package tools

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func InitLog() error {
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.SetOutput(f)
	return nil
}

func WriteLog(err error, i *discordgo.InteractionCreate) {
	log.Printf("ошибка в команде %s: %s", i.ApplicationCommandData().Name, err.Error())
}
