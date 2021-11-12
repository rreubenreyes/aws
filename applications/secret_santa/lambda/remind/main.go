package main

import (
	"log"
	"os"

	"remind/internal/crawler"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	url                  = os.Getenv("DRAW_URL")
	discordTokenSecretId = os.Getenv("DISCORD_TOKEN_SECRET_ID")
)

func remind() error {
	log.Println("starting")

	state, err := crawler.LatestDrawState(url)
	if err != nil {
		log.Println("could not get latest draw state")
		return err
	}

	// TODO: get wide photos
	_, err = RemainingParticipants(state.Members)
	if err != nil {
		log.Println("could not get remaining participants")
		return err
	}

	// SendReminders(participants)

	return nil
}

func main() {
	lambda.Start(remind)
}
