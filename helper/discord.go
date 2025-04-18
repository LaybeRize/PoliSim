package helper

import (
	loc "PoliSim/localisation"
	"github.com/bwmarrin/discordgo"
	"log"
	"log/slog"
	"os"
	"strings"
)

var DiscordDocumentChannelID *string = nil
var DiscordPressChannelID *string = nil
var DiscordNotesChannelID *string = nil
var UrlPrefix = strings.TrimRight(os.Getenv("URL_PREFIX"), "/")
var discord *discordgo.Session = nil

func setupDiscord() {
	discordToken, hasToken := os.LookupEnv("DISCORD_TOKEN")
	if hasToken {
		var err error
		discord, err = discordgo.New("Bot " + discordToken)
		if err != nil {
			log.Fatalf("Could not connect to discord properly: %v", err)
		}
		log.Println("Discord Bot created")
		err = discord.Open()
		if err != nil {
			log.Fatalf("Could open connection to discord properly: %v", err)
		}
		log.Println("Discord Bot connected")
		_ = discord.UpdateCustomStatus(loc.DiscordStatusText)
	}
	documentChannel, hasChannel := os.LookupEnv("DOCUMENT_CHANNEL_ID")
	if hasChannel {
		DiscordDocumentChannelID = &documentChannel
	}
	pressChannel, hasChannel := os.LookupEnv("PRESS_CHANNEL_ID")
	if hasChannel {
		DiscordPressChannelID = &pressChannel
	}
	noteChannel, hasChannel := os.LookupEnv("NOTES_CHANNEL_ID")
	if hasChannel {
		DiscordNotesChannelID = &noteChannel
	}
}

func SendDiscordEmbedMessage(channelID *string, message *discordgo.MessageEmbed) {
	if discord == nil || channelID == nil || message == nil {
		return
	}
	_, err := discord.ChannelMessageSendComplex(*channelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{message},
	})
	if err != nil {
		slog.Debug(err.Error())
	}
}
