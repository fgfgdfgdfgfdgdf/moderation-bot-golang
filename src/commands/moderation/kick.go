package moderation

import (
	"bot/src/UI/response"
	"bot/src/tools"

	"github.com/bwmarrin/discordgo"
)

type Kick struct{}

func (h Kick) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "kick",
		Description: "Выдать блокировку пользователю",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "пользователь",
				Description: "Пользователь (Упоминание или id)",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    true,
			},
			{
				Name:        "причина",
				Description: "Причина выдачи наказания",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	}
}

func (h Kick) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !tools.GetBanAccess(i.Member, i.GuildID) {
		err := response.AccessError(s, i)
		if err != nil {
			tools.WriteLog(err, i)
		}
		return
	}

	user := i.ApplicationCommandData().Options[0].UserValue(s)

	if user.Bot {
		err := response.ErrorWithText(s, i, "Боту невозможно выдать наказание!")
		if err != nil {
			tools.WriteLog(err, i)
		}
		return
	}

	member, err := s.State.Member(i.GuildID, user.ID)
	if err != nil {
		err := response.ErrorWithText(s, i, "Пользователь не найден на сервере!")
		if err != nil {
			tools.WriteLog(err, i)
		}
		return
	}

	if tools.GetModerationAccess(member, i.GuildID) {
		err := response.ErrorWithText(s, i, "Модератору невозможно выдать наказание!")
		if err != nil {
			tools.WriteLog(err, i)
		}
		return
	}

	if !tools.CompareTopRoles(s, s.State.User.ID, member, i.GuildID) {
		err = response.ErrorWithText(s, i, "У бота недостаточно прав для кика пользователя!")
		if err != nil {
			tools.WriteLog(err, i)
		}
		return
	}

	ban_info, err := s.GuildBan(i.GuildID, member.User.ID)
	if err == nil {
		err = response.ErrorWithText(s, i, "Пользователь уже находится в бане с причиной"+ban_info.Reason)
		if err != nil {
			tools.WriteLog(err, i)
		}
		return
	}

	reason := i.ApplicationCommandData().Options[2].StringValue()

	punish := tools.GetPunish(i, "KICK", user, reason, "", 0)

	err = response.PunishResponse(s, i, punish)
	if err != nil {
		err = response.ErrorWithText(s, i, "Произошла ошибка при выдаче наказания!")
		if err != nil {
			tools.WriteLog(err, i)
		}
		return
	}

	go response.PunishDmMessage(s, punish)
	go response.PunishLog(s, punish)
	go s.GuildMemberDelete(i.GuildID, member.User.ID)
	go tools.Insert_Punish(punish.ToDict())

}
