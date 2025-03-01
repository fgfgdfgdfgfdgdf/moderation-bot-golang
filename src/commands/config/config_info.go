package config

import (
	"bot/src/UI/response"
	"bot/src/colors"
	"bot/src/tools"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Config_info struct{}

func (h Config_info) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "config_info",
		Description: "Информация о конфигурации сервера",
	}
}

func (h Config_info) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		_ = response.AccessError(s, i)
		return
	}

	config, err := tools.Get_Config(i.GuildID)
	if err != nil {
		_ = response.ErrorWithText(s, i, "Произошла ошибка при получении конфигурации!")
		return
	}

	moderPermRoles := ""
	banPermRoles := ""

	for _, moderPermRole := range config.ModerationPermRoles {
		moderPermRoles += "<@&" + moderPermRole + ">"
	}
	for _, banPermRole := range config.BanPermRoles {
		banPermRoles += "<@&" + banPermRole + ">"
	}

	embedDescription := fmt.Sprintf("**Модерация:**\nКанал для логирования:  <#%s>\nРоли модерации: %s\nРоли с доступом к бану: %s", config.ModerationLog, moderPermRoles, banPermRoles)

	embed := &discordgo.MessageEmbed{
		Title:       "Информация о конфигурации сервера",
		Description: embedDescription,
		Color:       colors.Blue(),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    s.State.User.Username,
			IconURL: s.State.User.AvatarURL(""),
		},
	}

	_ = response.WithEmbedsAndButtons(s, i, nil, []*discordgo.MessageEmbed{embed}, nil, true)

}
