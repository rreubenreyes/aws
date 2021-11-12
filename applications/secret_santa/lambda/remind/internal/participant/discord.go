package participant

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var discordChannelId = os.Getenv("DISCORD_CHANNEL_ID")

func (p *Participant) SendReminder(sess *discordgo.Session) error {
	_, err := sess.ChannelMessageSend(
		discordChannelId,
		fmt.Sprintf("hey %s go sign up for secret santa", p.Name),
	)
	if err != nil {
		log.Println("could not send message")
		return err
	}

	return nil
}
