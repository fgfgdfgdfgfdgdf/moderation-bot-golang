package moderation

import (
	"bot/src/UI/response"
	"bot/src/tools"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Unmute struct{}

func (h Unmute) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "unmute",
		Description: "Размутить пользователя",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "пользователь",
				Description: "Пользователь (Упоминание или id)",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    true,
			},
			{
				Name:        "причина",
				Description: "Причина снятия наказания",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	}
}

func (h Unmute) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !tools.GetModerationAccess(i.Member, i.GuildID) {
		err := response.AccessError(s, i)
		if err != nil {
			log.Println(err)
		}
		return
	}

	user := i.ApplicationCommandData().Options[0].UserValue(s)

	member, err := s.State.Member(i.GuildID, user.ID)
	if err != nil {
		err = response.ErrorWithText(s, i, "Пользователь не найден!")
		if err != nil {
			log.Println(err)
		}
		return
	}

	if member.CommunicationDisabledUntil == nil || member.CommunicationDisabledUntil.Unix() <= time.Now().Unix() {
		err = response.ErrorWithText(s, i, "У пользователя не найден таймаут!")
		if err != nil {
			log.Println(err)
		}
		return
	}

	reason := i.ApplicationCommandData().Options[1].StringValue()

	punish := tools.GetPunish(i, "UNMUTE", user, reason, "", 0)

	err = response.PunishResponse(s, i, punish)
	if err != nil {
		err = response.ErrorWithText(s, i, "Произошла ошибка при выдаче наказания!")
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = response.PunishLog(s, punish)
	if err != nil {
		return
	}

	err = s.GuildMemberTimeout(i.GuildID, member.User.ID, nil)
	if err != nil {
		return
	}

	err = tools.Insert_Punish(punish.ToDict())
	if err != nil {
		return
	}

}
