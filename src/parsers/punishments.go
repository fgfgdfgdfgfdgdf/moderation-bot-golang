package parsers

import (
	"bot/src/tools"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func CheckPunishments(s *discordgo.Session) {
	var punishRole string

	PUNISH_TYPES := map[string]struct{}{
		"BAN": {},
	}

	for {
		for _, guild := range s.State.Guilds {
			punishments, err := tools.Get_Punishments(bson.M{"guild_id": guild.ID, "end_time": map[string]interface{}{"$ne": 0}})

			if err != nil {
				continue
			}

			logChannelId := tools.GetModerationLogChannel(guild.ID)

			for _, punish := range punishments {

				if punish.EndTime > time.Now().Unix() {
					continue
				}

				member, err := s.State.Member(punish.GuildId, punish.Member)
				if err != nil {
					continue
				}

				_, exists := PUNISH_TYPES[punish.Category]
				if exists {

					switch punish.Category {
					case "BAN":

					}

					err = s.GuildMemberRoleRemove(punish.GuildId, punish.Member, punishRole)
					if err != nil {

					}

					filter := bson.M{
						"category": punish.Category,
						"guildId":  punish.GuildId,
						"memberId": punish.Member,
						"endTime": bson.M{
							"$ne": nil,
						},
					}

					data := bson.M{
						"endTime": nil,
					}

					err = tools.Update_Punish(filter, data)
					if err != nil {

					}

					if logChannelId != "" {

						embed := &discordgo.MessageEmbed{
							Title:     fmt.Sprintf("[UN%s] %s", punish.Category, member.DisplayName()),
							Timestamp: time.Now().Format(time.RFC3339),
							Footer: &discordgo.MessageEmbedFooter{
								Text:    s.State.User.Username,
								IconURL: s.State.User.AvatarURL(""),
							},
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:   "Модератор",
									Value:  "<@" + punish.Moder + ">",
									Inline: true,
								},
								{
									Name:   "Пользователь",
									Value:  member.Mention(),
									Inline: true,
								},
								{
									Name:   "Причина",
									Value:  "Срок наказания вышел",
									Inline: true,
								},
							},
						}

						_, err = s.ChannelMessageSendEmbed(logChannelId, embed)
						if err != nil {
							continue
						}

					}
				}

			}
		}
	}
}
