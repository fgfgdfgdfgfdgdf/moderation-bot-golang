package moderation

import (
	"bot/src/UI/response"
	"bot/src/tools"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Clear struct{}

func (h Clear) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "clear",
		Description: "Очистить сообщения в канале",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "количество",
				Description: "Количество сообщений для удаления (не более 20)",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionInteger,
			},
		},
	}
}

func (h Clear) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !tools.GetModerationAccess(i.Member, i.GuildID) {
		err := response.AccessError(s, i)
		if err != nil {
			log.Println(err)
		}
		return
	}

	amount := i.ApplicationCommandData().Options[0].IntValue()

	if amount > 20 || amount < 1 {
		err := response.ErrorWithText(s, i, "Введено неверное количество удаляемых сообщений!")
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	messages, err := s.ChannelMessages(i.ChannelID, int(amount), "", "", "")
	if err != nil {
		err := response.ErrorWithText(s, i, "Произошла ошибка при получении сообщений!")
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	minimum_time := int((time.Now().Unix()-14*24*60*60)*1000-1420070400000) << 22

	messagesToBulkDelete := []string{}

	messagesToSingleDelete := []string{}

	flag := false

	for _, message := range messages {
		if !flag {
			message_id_int, err := strconv.Atoi(message.ID)
			if err != nil {
				continue
			}
			if message_id_int < minimum_time {
				flag = true
				messagesToSingleDelete = append(messagesToSingleDelete, message.ID)
			} else {
				messagesToBulkDelete = append(messagesToBulkDelete, message.ID)
			}
		} else {
			messagesToSingleDelete = append(messagesToSingleDelete, message.ID)
		}
	}

	err = s.ChannelMessagesBulkDelete(i.ChannelID, messagesToBulkDelete)
	if err != nil {
		err := response.ErrorWithText(s, i, "Произошла ошибка при удалении сообщений!")
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	for _, messageId := range messagesToSingleDelete {
		err = s.ChannelMessageDelete(i.ChannelID, messageId)
		if err != nil {
			continue
		}
	}

	err = response.SuccessWithText(s, i, fmt.Sprintf("Удалено %d сообщений!", amount))
	if err != nil {
		err := response.ErrorWithText(s, i, "Произошла ошибка при отправке ответа об успешном выполении команды!")
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

}
