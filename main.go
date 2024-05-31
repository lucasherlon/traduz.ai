package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"traduz.ai/geminiapi"
)

func main() {
	const command string = "!traduz"

	err := godotenv.Load(".env")
	if err != nil {
		return
	}
	token := os.Getenv("BOT_KEY")

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		args := strings.Split(m.Content, " ")
		if args[0] != command {
			return
		}

		if m.Author.ID == s.State.User.ID {
			return
		}

		input := m.Content[1:]
		response := geminiapi.GeneratePrompt(input)

		_, err := s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			log.Fatal(err)
		}

	})

	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = session.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer func(session *discordgo.Session) {
		err := session.Close()
		if err != nil {
		}
	}(session)

	fmt.Println("Traduz.ai is running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
