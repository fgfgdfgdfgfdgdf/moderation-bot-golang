package tools

import "github.com/bwmarrin/discordgo"

func CompareTopRoles(s *discordgo.Session, botUserID string, member *discordgo.Member, guildID string) bool {
	botMember, err := s.State.Member(guildID, botUserID)
	if err != nil {
		return false
	}

	guild, err := s.State.Guild(guildID)
	if err != nil {
		return false
	}

	flag1 := false
	flag2 := false

	var bot_top_role_position int
	var user_top_role_position int

	bot_top_role_id := botMember.Roles[0]
	user_top_role_id := member.Roles[0]

	for _, role := range guild.Roles {
		if !flag1 {
			if role.ID == bot_top_role_id {
				flag1 = true
				bot_top_role_position = role.Position
			}
		}
		if !flag2 {
			if role.ID == user_top_role_id {
				flag2 = true
				user_top_role_position = role.Position
			}
		}
		if flag1 && flag2 {
			break
		}
	}
	return bot_top_role_position > user_top_role_position
}
