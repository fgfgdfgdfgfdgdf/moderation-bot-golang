package tools

import (
	"slices"

	"github.com/bwmarrin/discordgo"
)

func GetRoleMembers(s *discordgo.Session, roleId string, guildId string) []*discordgo.Member {
	membersWithRole := []*discordgo.Member{}
	guild, err := s.State.Guild(guildId)
	if err != nil {
		return membersWithRole
	}

	for _, member := range guild.Members {
		if slices.Contains(member.Roles, roleId) {
			membersWithRole = append(membersWithRole, member)
		}
	}

	return membersWithRole
}
