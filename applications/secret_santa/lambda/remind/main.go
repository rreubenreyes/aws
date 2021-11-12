package main

import (
	"log"
	"os"
	"sync"

	"remind/internal/crawler"
	"remind/internal/participant"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	// "github.com/aws/aws-sdk-go/service/secretsmanager"
	// "github.com/bwmarrin/discordgo"
)

var (
	awsRegion            = os.Getenv("AWS_REGION")
	discordTokenSecretId = os.Getenv("DISCORD_TOKEN_SECRET_ID")
	url                  = os.Getenv("DRAW_URL")
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
	// sm := secretsmanager.New(sess)
	// res, err := sm.GetSecretValue(&secretsmanager.GetSecretValueInput{
	// 	SecretId: aws.String(discordTokenSecretId),
	// })
	// if err != nil {
	// 	log.Println("could not get discord token secret")
	// 	return err
	// }

	// log.Printf("got secret %s\n", res.String())
	// discord, err := discordgo.New("Bot " + *res.SecretString)
	// if err != nil {
	// 	log.Println("could not connect to Discord")
	// 	return err
	// }

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
				log.Printf("error getting Participant record for %s\n", m.Name)
			}

			_, err = p.PhotoPNG(s3Svc)
			if err != nil {
				log.Printf("error getting Participant photo for %s\n", p.Name)
			}
		}(m)
	}

	wg.Wait()

	return nil
}

func main() {
	lambda.Start(handler)
}
