package parsers

import (
	"bot/src/tools"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func CheckPunishments(s *discordgo.Session) {
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

				user, err := s.User(punish.Member)

				_, exists := PUNISH_TYPES[punish.Category]
				if exists {

					switch punish.Category {
					case "BAN":
						if err != nil {
							continue
						}
						err = s.GuildBanDelete(guild.ID, user.ID)
						if err != nil {
							continue
						}
					}

					filter := bson.M{
						"category": punish.Category,
						"guild_id": punish.GuildId,
						"user_id":  punish.Member,
						"end_time": bson.M{
							"$ne": nil,
						},
					}

					data := bson.M{
						"end_time": nil,
					}

					err = tools.Update_Punish(filter, data)
					if err != nil {

					}

					if logChannelId != "" {

						embed := &discordgo.MessageEmbed{
							Title:     fmt.Sprintf("[UN%s] %s", punish.Category, user.Username),
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
									Value:  user.Mention(),
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
