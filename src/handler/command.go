package handler

import (
	"bot/src/commands/config"
	"bot/src/commands/moderation"
	"bot/src/commands/other"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler interface {
	Command() *discordgo.ApplicationCommand
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var CommandHandlers = []CommandHandler{
	&moderation.Mute{},
	&moderation.Unmute{},
	&moderation.Ban{},
	&moderation.Clear{},
	&moderation.Unban{},
	&moderation.Kick{},
	&other.Avatar{},
	&other.Ping{},
	&other.Rolemembers{},
	&config.Config_info{},
	&config.Config_moderation{},
	&config.Config{},
}
