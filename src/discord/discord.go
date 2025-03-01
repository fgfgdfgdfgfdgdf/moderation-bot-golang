package discord

import (
	// "economy-bot/config"

	"bot/src/handler"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bot/src/parsers"

	"github.com/bwmarrin/discordgo"
)

func Init() {

	token := os.Getenv("TOKEN")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println(err)
	}

	user, err := dg.Application("@me")
	if err != nil {
		log.Printf("Ошибка при получении id бота: %s", err)
		return
	}

	registerCommandHandlers(dg, handler.CommandHandlers, user.ID)

	dg.Identify.Intents = discordgo.IntentsAll

	err = dg.Open()
	if err != nil {
		log.Printf("error opening connection: %s", err.Error())
		return
	}

	go parsers.CheckPunishments(dg)

	fmt.Println("Бот запущен.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()

}

func registerCommandHandlers(s *discordgo.Session, commandHandlers []handler.CommandHandler, userID string) {
	existingCommands, err := s.ApplicationCommands(userID, "")
	if err != nil {
		fmt.Printf("Ошибка при получении существующих команд: %v\n", err)
		return
	}

	existingCommandsMap := make(map[string]*discordgo.ApplicationCommand)
	for _, cmd := range existingCommands {
		existingCommandsMap[cmd.Name] = cmd
	}

	newCommands := make(map[string]bool)

	for _, handler := range commandHandlers {
		cmd := handler.Command()
		newCommands[cmd.Name] = true

		if existingCmd, exists := existingCommandsMap[cmd.Name]; exists {
			if commandsEqual(existingCmd, cmd) {
				continue
			}

			err := s.ApplicationCommandDelete(userID, "", existingCmd.ID)
			if err != nil {
				fmt.Printf("Ошибка при удалении старой команды '%s': %v\n", cmd.Name, err)
				continue
			}
		}

		_, err := s.ApplicationCommandCreate(userID, "", cmd)
		if err != nil {
			fmt.Printf("Ошибка при создании команды '%s': %v\n", cmd.Name, err)
			continue
		}
		fmt.Printf("Команда '%s' успешно зарегистрирована.\n", cmd.Name)
	}

	for _, cmd := range existingCommands {
		if _, exists := newCommands[cmd.Name]; !exists {
			err := s.ApplicationCommandDelete(userID, "", cmd.ID)
			if err != nil {
				fmt.Printf("Ошибка при удалении устаревшей команды '%s': %v\n", cmd.Name, err)
			} else {
				fmt.Printf("Устаревшая команда '%s' удалена.\n", cmd.Name)
			}
		}
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.GuildID == "" {
			return
		}

		if i.Type == discordgo.InteractionApplicationCommand {
			for _, handler := range commandHandlers {
				if handler.Command().Name == i.ApplicationCommandData().Name {
					handler.Handler(s, i)
					return
				}
			}
		}
	})

}

func commandsEqual(existing *discordgo.ApplicationCommand, new *discordgo.ApplicationCommand) bool {
	if existing.Name != new.Name {
		return false
	}
	if existing.Description != new.Description {
		return false
	}
	if len(existing.Options) != len(new.Options) {
		return false
	}
	for i := range existing.Options {
		if !optionsEqual(existing.Options[i], new.Options[i]) {
			return false
		}
	}
	return true
}

func optionsEqual(existing *discordgo.ApplicationCommandOption, new *discordgo.ApplicationCommandOption) bool {
	if existing.Type != new.Type {
		return false
	}
	if existing.Name != new.Name {
		return false
	}
	if existing.Description != new.Description {
		return false
	}
	if existing.Required != new.Required {
		return false
	}
	if len(existing.Choices) != len(new.Choices) {
		return false
	}
	for i := range existing.Choices {
		if !choicesEqual(existing.Choices[i], new.Choices[i]) {
			return false
		}
	}
	if len(existing.Options) != len(new.Options) {
		return false
	}
	for i := range existing.Options {
		if !optionsEqual(existing.Options[i], new.Options[i]) {
			return false
		}
	}
	return true
}

func choicesEqual(existing *discordgo.ApplicationCommandOptionChoice, new *discordgo.ApplicationCommandOptionChoice) bool {
	if existing.Name != new.Name {
		return false
	}
	if existing.Value != new.Value {
		return false
	}
	return true
}
