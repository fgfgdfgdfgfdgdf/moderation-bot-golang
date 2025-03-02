package response

import (
	"bot/src/UI/embed"
	"bot/src/colors"
	"bot/src/tools"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	TYPES_WITH_DURATION = map[string]struct{}{
		"BAN":  {},
		"MUTE": {},
	}

	ResponseData = map[string]interface{}{
		"KICK":   []string{"Кик пользователя", "был кикнут модератором"},
		"BAN":    []string{"Блокировка пользователя", "был заблокирован модератором"},
		"UNBAN":  []string{"Разблокировка пользователя", "был разблокирован модератором"},
		"MUTE":   []string{"Блокировка чата", "была выдана блокировка чата модератором"},
                "UNMUTE": []string{"Снятие блокировки чата", "была снята блокировка чата модератором"},
	}

	DM_DATA = map[string]interface{}{
		"KICK": "Вы были кикнуты с сервера",
		"BAN":  "Вы были забанены на сервере",
		"MUTE": "Вы были замучены на сервере",
	}
)

func SuccessWithText(s *discordgo.Session, i *discordgo.InteractionCreate, text string) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				embed.SuccessEmbed(s, text),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("ошибка при отправке ответа %w", err)
	}
	return nil
}

func ErrorWithText(s *discordgo.Session, i *discordgo.InteractionCreate, text string) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				embed.ErrorEmbed(s, text),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("ошибка при отправке ответа %w", err)
	}
	return nil
}

func AccessError(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				embed.ErrorEmbed(s, "У вас недостаточно прав для использования данной команды!"),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("ошибка при отправке ответа %w", err)
	}
	return nil
}

func WithEmbedsAndButtons(s *discordgo.Session, i *discordgo.InteractionCreate, text *string, embeds []*discordgo.MessageEmbed, buttons []*discordgo.Button, ephemeral bool) error {
	var components []discordgo.MessageComponent
	if len(buttons) > 0 {
		actionRow := discordgo.ActionsRow{
			Components: make([]discordgo.MessageComponent, len(buttons)),
		}
		for index, button := range buttons {
			actionRow.Components[index] = *button
		}
		components = append(components, actionRow)
	}

	var content string
	if text != nil {
		content = *text
	}

	var flag discordgo.MessageFlags

	if ephemeral {
		flag = discordgo.MessageFlagsEphemeral
	} else {
		flag = 0
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:      flag,
			Content:    content,
			Embeds:     embeds,
			Components: components,
		},
	})
	if err != nil {
		return fmt.Errorf("ошибка при отправке ответа: %w", err)
	}

	return nil
}

func PunishResponse(s *discordgo.Session, i *discordgo.InteractionCreate, punish *tools.Punish) error {

	embedFields := []*discordgo.MessageEmbedField{
		{
			Name:   "Причина",
			Value:  punish.Reason,
			Inline: true,
		},
	}

	if _, exists := TYPES_WITH_DURATION[punish.Category]; exists {

		fields := []*discordgo.MessageEmbedField{
			{
				Name:   "Срок окончания",
				Value:  fmt.Sprintf("<t:%d>", punish.EndTime),
				Inline: false,
			},
		}

		embedFields = append(embedFields, fields...)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       ResponseData[punish.Category].([]string)[0],
					Description: punish.User.Mention() + " " + ResponseData[punish.Category].([]string)[1] + " " + punish.Moder.Mention(),
					Fields:      embedFields,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    s.State.User.Username,
						IconURL: s.State.User.AvatarURL(""),
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("ошибка при отправке ответа: %w", err)
	}
	return nil
}

func PunishDmMessage(s *discordgo.Session, punish *tools.Punish) error {
	channel, err := s.UserChannelCreate(punish.User.ID)
	if err != nil {
		return fmt.Errorf("ошибка при создании dm канала: %w", err)
	}

	guild_name := ""

	guild, err := s.State.Guild(punish.GuildId)
	if err == nil {
		guild_name = guild.Name
	}

	embedDescription := fmt.Sprintf("Модератор: `%s [%s]`\nПричина: %s ", punish.Moder.DisplayName(), punish.Moder.User.ID, punish.Reason)
	_, exists := DM_DATA[punish.Category]
	if exists {
		embedDescription += fmt.Sprintf("\nДата окончания наказания:<t:%d>", punish.EndTime)
	}

	embed := &discordgo.MessageEmbed{
		Title:       DM_DATA[punish.Category].(string) + " " + guild_name,
		Description: embedDescription,
		Color:       colors.Red(),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    s.State.User.Username,
			IconURL: s.State.User.AvatarURL(""),
		},
	}

	_, err = s.ChannelMessageSendEmbed(channel.ID, embed)
	if err != nil {
		return fmt.Errorf("ошибка при отправке сообщения: %w", err)
	}
	return nil
}

func PunishLog(s *discordgo.Session, punish *tools.Punish) error {
	logChannelId := tools.GetModerationLogChannel(punish.GuildId)
	if logChannelId != "" {

		embedFields := []*discordgo.MessageEmbedField{
			{
				Name:   "Модератор",
				Value:  punish.Moder.Mention(),
				Inline: true,
			},
			{
				Name:   "Пользователь",
				Value:  punish.User.Mention(),
				Inline: true,
			},
			{
				Name:   "Причина",
				Value:  punish.Reason,
				Inline: true,
			},
		}

		_, exists := TYPES_WITH_DURATION[punish.Category]
		if exists {

			fields := []*discordgo.MessageEmbedField{
				{
					Name:   "Срок наказания",
					Value:  punish.Duration,
					Inline: true,
				},
				{
					Name:   "Дата окончания наказания",
					Value:  fmt.Sprintf("<t:%d>", punish.EndTime),
					Inline: true,
				}}

			embedFields = append(embedFields, fields...)
		}

		embed := &discordgo.MessageEmbed{
			Title:     "[" + punish.Category + "]" + punish.User.Username,
			Timestamp: time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text:    s.State.User.Username,
				IconURL: s.State.User.AvatarURL(""),
			},
			Fields: embedFields,
		}
		_, err := s.ChannelMessageSendEmbed(logChannelId, embed)
		if err != nil {
			return fmt.Errorf("ошибка при отправке сообщения: %w", err)
		}

	}
	return nil
}
