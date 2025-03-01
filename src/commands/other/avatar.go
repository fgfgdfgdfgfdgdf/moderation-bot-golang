package other

import (
	"github.com/bwmarrin/discordgo"
)

type Avatar struct{}

func (h *Avatar) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "avatar",
		Description: "Получить аватар пользователя",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "пользователь",
				Description: "Пользователь (Упоминание или id)",
				Required:    false,
			},
		},
	}
}

func (h *Avatar) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var avatarUrl string
	var userName string
	var user *discordgo.User

	for _, option := range i.ApplicationCommandData().Options {
		if option.Name == "пользователь" {
			user = option.UserValue(s)
			break
		}
	}

	if user == nil {
		avatarUrl = i.Member.AvatarURL("512")
		userName = i.Member.Nick
	} else {
		avatarUrl = user.AvatarURL("512")
		userName = user.GlobalName
	}

	embed := &discordgo.MessageEmbed{
		Title: "Аватар пользователя " + userName,
		Image: &discordgo.MessageEmbedImage{
			URL: avatarUrl,
		},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				embed,
			},
		},
	})

	if err != nil {
		return
	}

}
