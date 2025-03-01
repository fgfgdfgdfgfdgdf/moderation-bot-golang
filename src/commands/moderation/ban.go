package moderation

import (
	"bot/src/UI/response"
	"bot/src/tools"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Ban struct{}

func (h Ban) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ban",
		Description: "Выдать блокировку пользователю",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "пользователь",
				Description: "Пользователь (Упоминание или id)",
				Type:        discordgo.ApplicationCommandOptionUser,
				Required:    true,
			},
			{
				Name:        "длительность",
				Description: "Длительность блокировки (m/h/d/w)",
				Type:        discordgo.ApplicationCommandOptionString,
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

func (h Ban) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	if err == nil {

		if tools.GetModerationAccess(member, i.GuildID) {
			err := response.ErrorWithText(s, i, "Модератору невозможно выдать наказание!")
			if err != nil {
				tools.WriteLog(err, i)
			}
			return
		}

		if !tools.CompareTopRoles(s, s.State.User.ID, member, i.GuildID) {
			err = response.ErrorWithText(s, i, "У бота недостаточно прав для блокировки пользователя!")
			if err != nil {
				tools.WriteLog(err, i)
			}
			return
		}
	}

	duration := i.ApplicationCommandData().Options[1].StringValue()

	if !tools.Check_duration(duration) {
		err = response.ErrorWithText(s, i, "Неверный ввод времени наказания!")
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

	end_time := tools.StringToSeconds(duration) + time.Now().Unix()

	punish := tools.GetPunish(i, "BAN", user, reason, duration, end_time)

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
	go s.GuildBanCreateWithReason(i.GuildID, member.User.ID, reason, 0)
	go tools.Insert_Punish(punish.ToDict())

}
