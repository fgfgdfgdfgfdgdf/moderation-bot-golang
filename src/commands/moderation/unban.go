package moderation

import (
	"bot/src/UI/response"
	"bot/src/tools"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Unban struct{}

func (h Unban) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "unban",
		Description: "Разбанить пользователя",
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

func (h Unban) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !tools.GetModerationAccess(i.Member, i.GuildID) {
		err := response.AccessError(s, i)
		if err != nil {
			log.Println(err)
		}
		return
	}

	user := i.ApplicationCommandData().Options[0].UserValue(s)

	user, err := s.User(user.ID)
	if err != nil {
		err = response.ErrorWithText(s, i, "Пользователь не найден!")
		if err != nil {
			log.Println(err)
		}
		return
	}

	_, err = s.GuildBan(i.GuildID, user.ID)
	if err != nil {
		err = response.ErrorWithText(s, i, "Пользователь не находится в бане!")
		if err != nil {
			tools.WriteLog(err, i)
		}
		return
	}

	reason := i.ApplicationCommandData().Options[1].StringValue()

	punish := tools.GetPunish(i, "UNBAN", user, reason, "", 0)

	err = response.PunishResponse(s, i, punish)
	if err != nil {
		err = response.ErrorWithText(s, i, "Произошла ошибка при выдаче наказания!")
		if err != nil {
			log.Println(err)
		}
		return
	}

	go response.PunishLog(s, punish)
	go s.GuildBanDelete(i.GuildID, user.ID)
	go tools.Insert_Punish(punish.ToDict())

}
