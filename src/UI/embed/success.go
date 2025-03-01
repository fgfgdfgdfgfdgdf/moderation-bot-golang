package embed

import (
	"bot/src/colors"

	"github.com/bwmarrin/discordgo"
)

func SuccessEmbed(s *discordgo.Session, text string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       "[✅] Успешно выполнено!",
		Description: "**" + text + "**",
		Color:       colors.Green(),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    s.State.User.Username,
			IconURL: s.State.User.AvatarURL(""),
		},
	}
	return embed
}
