package tools

func GetModerationLogChannel(guildId string) string {
	config, err := Get_Config(guildId)
	if err != nil {
		return ""
	}

	return config.ModerationLog
}
