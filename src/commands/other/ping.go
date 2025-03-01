package other

import (
	"github.com/bwmarrin/discordgo"
)

type Ping struct{}

func (h Ping) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Посмотреть отклик бота",
	}
}

func (h Ping) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "pong!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
