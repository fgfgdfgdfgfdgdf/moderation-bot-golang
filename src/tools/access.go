package tools

import (
	"github.com/bwmarrin/discordgo"
)

func GetModerationAccess(member *discordgo.Member, guildId string) bool {
	config, err := Get_Config(guildId)
	if err != nil {
		return false
	}

	roleSet := make(map[string]struct{}, len(config.ModerationPermRoles))
	for _, role := range config.ModerationPermRoles {
		roleSet[role] = struct{}{}
	}

	for _, userRole := range member.Roles {
		if _, exists := roleSet[userRole]; exists {
			return true
		}
	}

	return false
}

func GetBanAccess(member *discordgo.Member, guildId string) bool {
	config, err := Get_Config(guildId)
	if err != nil {
		return false
	}

	roleSet := make(map[string]struct{}, len(config.BanPermRoles))
	for _, role := range config.BanPermRoles {
		roleSet[role] = struct{}{}
	}

	for _, userRole := range member.Roles {
		if _, exists := roleSet[userRole]; exists {
			return true
		}
	}

	return false
}
