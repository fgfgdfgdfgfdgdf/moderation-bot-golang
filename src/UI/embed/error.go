package embed

import (
	"bot/src/colors"

	"github.com/bwmarrin/discordgo"
)

func ErrorEmbed(s *discordgo.Session, text string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       "⚠️ | Произошла ошибка",
		Description: "**" + text + "**",
		Color:       colors.Red(),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    s.State.User.Username + " | Ошибка",
			IconURL: s.State.User.AvatarURL(""),
		},
	}
	return embed
}
