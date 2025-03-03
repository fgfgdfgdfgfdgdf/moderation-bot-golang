package moderation

import (
	"bot/src/UI/response"
	"bot/src/tools"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Mute struct{}

func (h Mute) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "mute",
		Description: "Замутить пользователя",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "пользователь",
				Description: "Пользователь (Упоминание или id)",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionUser,
			},
			{
				Name:        "длительность",
				Description: "Срок мута (m/h/d/w)",
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

func (h Mute) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	if tools.GetModerationAccess(member, i.GuildID) {
		err := response.ErrorWithText(s, i, "Модератору невозможно выдать наказание!")
		if err != nil {
			log.Println(err)
		}
		return
	}

	if !tools.CompareTopRoles(s, s.State.User.ID, member, i.GuildID) {
		err = response.ErrorWithText(s, i, "У бота недостаточно прав для выдачи мута пользователю!")
		if err != nil {
			tools.WriteLog(err, i)
		}
		return
	}

	duration := i.ApplicationCommandData().Options[1].StringValue()

	if !tools.Check_duration(duration) {
		err = response.ErrorWithText(s, i, "Неверный ввод времени наказания!")
		if err != nil {
			log.Println(err)
		}
		return
	}

	if member.CommunicationDisabledUntil == nil || member.CommunicationDisabledUntil.Unix() >= time.Now().Unix() {
		err = response.ErrorWithText(s, i, "Пользователю уже выдан таймаут!")
		if err != nil {
			log.Println(err)
		}
		return
	}

	reason := i.ApplicationCommandData().Options[2].StringValue()

	end_time := tools.StringToSeconds(duration) + time.Now().Unix()

	punish := tools.GetPunish(i, "MUTE", user, reason, duration, end_time)

	punish_duration := time.Unix(punish.EndTime, 0)

	go response.PunishResponse(s, i, punish)
	go response.PunishDmMessage(s, punish)
	go response.PunishLog(s, punish)
	go s.GuildMemberTimeout(i.GuildID, member.User.ID, &punish_duration)
	go tools.Insert_Punish(punish.ToDict())
}
