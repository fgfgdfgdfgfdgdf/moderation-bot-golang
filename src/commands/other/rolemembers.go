package other

import (
	"bot/src/UI/response"
	"bot/src/tools"

	"github.com/bwmarrin/discordgo"
)

type Rolemembers struct{}

func (h *Rolemembers) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "rolemembers",
		Description: "Получить информацию об участниках роли",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "роль",
				Description: "Роль (Упоминание или id)",
				Required:    true,
			},
		},
	}
}

func (h *Rolemembers) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		_ = response.AccessError(s, i)
		return
	}

	var role discordgo.Role
	for _, option := range i.ApplicationCommandData().Options {
		if option.Name == "роль" {
			role = *option.RoleValue(s, i.GuildID)
			break
		}
	}

	if role.ID == "" {
		_ = response.ErrorWithText(s, i, "Указанная роль не найдена.")
		return
	}

	membersWithRole := tools.GetRoleMembers(s, role.ID, i.GuildID)

	embed := &discordgo.MessageEmbed{
		Color:       int(role.Color),
		Description: "**Участники с ролью** " + role.Mention() + "\n\n",
		Footer: &discordgo.MessageEmbedFooter{
			Text:    s.State.User.Username,
			IconURL: s.State.User.AvatarURL(""),
		},
	}

	if len(membersWithRole) != 0 {
		for _, member := range membersWithRole {
			embed.Description += member.Mention() + " "
			if len(embed.Description) > 3950 {
				embed.Description += "\nИ еще несколько..."
				break
			}
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
