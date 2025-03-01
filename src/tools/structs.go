package tools

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Punish struct {
	GuildId string
	//Тип наказания
	Category string
	//Пользователь
	User *discordgo.User
	//Модератор, выдавший наказание
	Moder *discordgo.Member
	//Причина выдачи наказания
	Reason string
	//Срок выдачи наказания в unix
	Duration string
	//Время окончания наказания в unix
	EndTime int64
	//Время выдачи наказания в unix
	TimeNow int64
	//Дата выдачи наказания в формате dd.mm.yyyy
	Date string
}

type DbPunish struct {
	GuildId string `bson:"guild_id"`
	//Тип наказания
	Category string `bson:"category"`
	//Пользователь
	Member string `bson:"user_id"`
	//ID модератора, выдавший наказание
	Moder string `bson:"moder_id"`
	//Причина выдачи наказания
	Reason string `bson:"reason"`
	//Срок выдачи наказания в unix
	Duration string `bson:"duration"`
	//Время окончания наказания в unix
	EndTime int64 `bson:"end_time"`
	//Время выдачи наказания в unix
	TimeNow int64 `bson:"time_now"`
	//Дата выдачи наказания в формате dd.mm.yyyy
	Date string `bson:"date"`
}

func GetPunish(i *discordgo.InteractionCreate, category string, user *discordgo.User, reason string, duration string, end_time int64) *Punish {
	return &Punish{
		GuildId:  i.GuildID,
		Category: category,
		Moder:    i.Member,
		User:     user,
		Reason:   reason,
		TimeNow:  time.Now().Unix(),
		Duration: duration,
		EndTime:  end_time,
		Date:     time.Now().Format(time.DateOnly),
	}
}

func (p *Punish) ToDict() map[string]interface{} {
	return map[string]interface{}{
		"guild_id": p.GuildId,
		"category": p.Category,
		"user_id":  p.User.ID,
		"moder_id": p.Moder.User.ID,
		"reason":   p.Reason,
		"date":     p.Date,
		"end_time": p.EndTime,
		"duration": p.Duration,
		"time_now": p.TimeNow,
	}
}

type DbConfig struct {
	GuildId string `bson:"guild_id"`
	//ID канала для логирования системы модерациии
	ModerationLog string `bson:"moderation_log"`
	//ID ролей с доступом к модерации
	ModerationPermRoles []string `bson:"moder_perm_roles"`
	//ID ролей с доступом к блокировкам
	BanPermRoles []string `bson:"ban_perm_roles"`
}
