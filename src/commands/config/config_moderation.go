package config

import (
	"bot/src/UI/response"
	"bot/src/colors"
	"bot/src/tools"
	"log"
	"slices"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config_moderation struct{}

func (h Config_moderation) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "config_moderation",
		Description: "Настроить конфигурацию модерации",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "логи_модерации",
				Description: "Канал для логирования наказаний",
				Required:    false,
				ChannelTypes: []discordgo.ChannelType{
					discordgo.ChannelTypeGuildText,
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "добавить_роль_модера",
				Description: "Добавить роль с доступом к модерации",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "удалить_роль_модера",
				Description: "Удалить роль с доступом к модерации",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "добавить_роль_с_баном",
				Description: "Добавить роль с доступом к бану",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "удалить_роль_с_баном",
				Description: "Удалить роль с доступом к бану",
				Required:    false,
			},
		},
	}
}

func (h *Config_moderation) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member.Permissions&discordgo.PermissionAdministrator == 0 {
		_ = response.AccessError(s, i)
		return
	}

	config, err := tools.Get_Config(i.GuildID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = response.ErrorWithText(s, i, "Конфиг сервера не найден, введите команду /config!")
			if err != nil {
				tools.WriteLog(err, i)
			}
			return
		} else {
			err = response.ErrorWithText(s, i, "Произошла ошибка при получении конфига сервера!")
			if err != nil {
				tools.WriteLog(err, i)
			}
			return
		}

	}

	options := i.ApplicationCommandData().Options
	embedFields := []*discordgo.MessageEmbedField{}

	for _, option := range options {
		switch option.Name {
		case "логи_модерации":
			logChannel := option.ChannelValue(s)

			err := tools.Update_Config(i.GuildID, map[string]interface{}{"moderationLog": logChannel.ID})
			if err == nil {
				embedFields = append(embedFields, &discordgo.MessageEmbedField{
					Name:   "Канал для логирования наказаний:",
					Value:  logChannel.Mention(),
					Inline: false,
				})
			} else {
				embedFields = append(embedFields, &discordgo.MessageEmbedField{
					Name:   "Канал для логирования наказаний:",
					Value:  "Возникла ошибка при добавлении в бд",
					Inline: false,
				})
			}

		case "добавить_роль_модера":
			addAccessRole := option.RoleValue(s, i.GuildID)

			oldAccessRoles := config.ModerationPermRoles

			if slices.Contains(oldAccessRoles, addAccessRole.ID) {
				embedFields = append(embedFields, &discordgo.MessageEmbedField{
					Name:   "Не удалось добавить роль с доступом к модерации:",
					Value:  "Роль была добавлена ранее",
					Inline: false,
				})
			} else {

				oldAccessRoles = append(oldAccessRoles, addAccessRole.ID)
				err := tools.Update_Config(i.GuildID, map[string]interface{}{"moderPermRoles": oldAccessRoles})
				if err == nil {
					embedFields = append(embedFields, &discordgo.MessageEmbedField{
						Name:   "Добавлена роль с доступом к модерации:",
						Value:  addAccessRole.Mention(),
						Inline: false,
					})
				} else {
					embedFields = append(embedFields, &discordgo.MessageEmbedField{
						Name:   "Роль с доступом к модерации:",
						Value:  "Возникла ошибка при добавлении в бд",
						Inline: false,
					})
				}
			}

		case "удалить_роль_модера":
			dellAccessRole := option.RoleValue(s, i.GuildID)

			oldAccessRoles := config.ModerationPermRoles

			if !slices.Contains(oldAccessRoles, dellAccessRole.ID) {
				embedFields = append(embedFields, &discordgo.MessageEmbedField{
					Name:   "Не удалось удалить роль с доступом к модерации:",
					Value:  "Добавление роли не найдена в конфигурации.",
					Inline: false,
				})
			} else {

				index := slices.Index(oldAccessRoles, dellAccessRole.ID)

				oldAccessRoles = slices.Delete(oldAccessRoles, index, index)
				err := tools.Update_Config(i.GuildID, map[string]interface{}{"moderPermRoles": oldAccessRoles})
				if err == nil {
					embedFields = append(embedFields, &discordgo.MessageEmbedField{
						Name:   "Удалена роль с доступом к модерации:",
						Value:  dellAccessRole.Mention(),
						Inline: false,
					})
				} else {
					embedFields = append(embedFields, &discordgo.MessageEmbedField{
						Name:   "Удаление роли с доступом к модерации:",
						Value:  "Возникла ошибка при добавлении в бд",
						Inline: false,
					})
				}

			}

		case "добавить_роль_с_баном":
			addAccessRole := option.RoleValue(s, i.GuildID)

			oldAccessRoles := config.BanPermRoles

			if slices.Contains(oldAccessRoles, addAccessRole.ID) {
				embedFields = append(embedFields, &discordgo.MessageEmbedField{
					Name:   "Не удалось добавить роль с доступом к бану:",
					Value:  "Роль была добавлена ранее",
					Inline: false,
				})
			} else {

				oldAccessRoles = append(oldAccessRoles, addAccessRole.ID)
				err := tools.Update_Config(i.GuildID, map[string]interface{}{"banPermRoles": oldAccessRoles})
				if err == nil {
					embedFields = append(embedFields, &discordgo.MessageEmbedField{
						Name:   "Добавлена роль с доступом к бану:",
						Value:  addAccessRole.Mention(),
						Inline: false,
					})
				} else {
					embedFields = append(embedFields, &discordgo.MessageEmbedField{
						Name:   "Роль с доступом к бану:",
						Value:  "Возникла ошибка при добавлении в бд",
						Inline: false,
					})
				}
			}

		case "удалить_роль_с_баном":
			dellAccessRole := option.RoleValue(s, i.GuildID)

			oldAccessRoles := config.BanPermRoles

			if !slices.Contains(oldAccessRoles, dellAccessRole.ID) {
				embedFields = append(embedFields, &discordgo.MessageEmbedField{
					Name:   "Не удалось удалить роль с доступом к бану:",
					Value:  "Добавление роли не найдена в конфигурации.",
					Inline: false,
				})
			} else {

				index := slices.Index(oldAccessRoles, dellAccessRole.ID)
				oldAccessRoles = slices.Delete(oldAccessRoles, index, index)
				err := tools.Update_Config(i.GuildID, map[string]interface{}{"banPermRoles": oldAccessRoles})
				if err == nil {
					embedFields = append(embedFields, &discordgo.MessageEmbedField{
						Name:   "Удалена роль с доступом к бану:",
						Value:  dellAccessRole.Mention(),
						Inline: false,
					})
				} else {
					embedFields = append(embedFields, &discordgo.MessageEmbedField{
						Name:   "Удаление роли с доступом к бану:",
						Value:  "Возникла ошибка при добавлении в бд",
						Inline: false,
					})
				}

			}
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:     "Изменения:",
		Fields:    embedFields,
		Color:     colors.Blue(),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	err = response.WithEmbedsAndButtons(s, i, nil, []*discordgo.MessageEmbed{embed}, nil, true)
	if err != nil {
		log.Println(err.Error())
	}

}
