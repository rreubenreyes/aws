package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"remind/internal/crawler"
	"remind/internal/participant"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/bwmarrin/discordgo"
)

var (
	awsRegion              = os.Getenv("AWS_REGION")
	discordChannelId       = os.Getenv("DISCORD_CHANNEL_ID")
	discordReminderMessage = os.Getenv("DISCORD_REMINDER_MESSAGE")
	discordTokenSecretId   = os.Getenv("DISCORD_TOKEN_SECRET_ID")
	url                    = os.Getenv("DRAW_URL")
	registerUrl            = os.Getenv("REGISTER_URL")
)

func handler() error {
	log.Println("starting")

	state, err := crawler.LatestDrawState(url)
	if err != nil {
		log.Println("could not get latest draw state")
		return err
	}

	// start discord client
	sess := session.New(&aws.Config{Region: aws.String(awsRegion)})
	smSvc := secretsmanager.New(sess)
	res, err := smSvc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(discordTokenSecretId),
	})
	if err != nil {
		log.Println("could not get discord token secret")
		return err
	}

	log.Printf("got secret %s\n", res.String())
	discord, err := discordgo.New("Bot " + *res.SecretString)
	if err != nil {
		log.Println("could not connect to Discord")
		return err
	}

	// start dynamodb client
	ddbSvc := dynamodb.New(sess)
	s3Svc := s3.New(sess)
	var wg sync.WaitGroup
	for _, m := range state.Members {
		if m.DrawViewed {
			continue
		}

		wg.Add(1)
		go func(m crawler.Member) {
			defer wg.Done()

			p, err := participant.FromMember(ddbSvc, m)
			if err != nil {
				log.Println(err)
				log.Printf("error getting Participant record for %s\n", m.Name)
			}

			photo, err := p.PhotoPNG(s3Svc)
			if err != nil {
				log.Println(err)
				log.Printf("error getting Participant photo for %s\n", p.Name)
			}

			widened, err := participant.WidenPhotoPNG(photo, p.PhotoDX, p.PhotoDY)
			if err != nil {
				log.Println(err)
				log.Printf("error widening photo for %s\n", p.Name)
			}

			f := &discordgo.File{
				Name:        p.Name + ".png",
				ContentType: "image/png",
				Reader:      bytes.NewReader(widened),
			}
			_, err = discord.ChannelMessageSendComplex(discordChannelId, &discordgo.MessageSend{
				Content: fmt.Sprintf(discordReminderMessage, p.DiscordId, registerUrl),
				Files:   []*discordgo.File{f},
			})
			if err != nil {
				log.Println("could not send message for", p.Name)
				log.Println(err)
			}

			err = p.Update(ddbSvc, &participant.ParticipantUpdate{
				DiscordId:  p.DiscordId,
				PhotoS3Key: p.PhotoS3Key,
				PhotoDX:    p.PhotoDX * 1.2,
				PhotoDY:    p.PhotoDY * 0.8,
			})
			if err != nil {
				log.Println("could not update photo scale for participant", p.Name)
				log.Println(err)
			}
		}(m)
	}

	wg.Wait()

	return nil
}

func main() {
	lambda.Start(handler)
}
