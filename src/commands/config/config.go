package config

import (
	"bot/src/UI/response"
	"bot/src/tools"
	"log"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct{}

func (h *Config) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "config",
		Description: "Создать конфиг сервера",
	}
}

func (h *Config) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		err := response.AccessError(s, i)
		if err != nil {
			log.Println(err.Error())
		}
		return
	}
	_, err := tools.Get_Config(i.GuildID)
	if err == nil {
		err = response.ErrorWithText(s, i, "Конфиг сервера был создан ранее!")
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	if err != mongo.ErrNoDocuments {
		err = response.ErrorWithText(s, i, "Произошла ошибка при получении конфига сервера!")
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	data := map[string]interface{}{
		"guildId":        i.GuildID,
		"moderationLog":  "",
		"moderPermRoles": []string{},
		"banPermRoles":   []string{},
	}

	err = tools.Insert_Config(data)
	if err != nil {
		err = response.ErrorWithText(s, i, "Возникла ошибка при создании конфига!")
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	err = response.SuccessWithText(s, i, "Конфиг сервера был создан!")
	if err != nil {
		log.Println(err.Error())
	}
}
