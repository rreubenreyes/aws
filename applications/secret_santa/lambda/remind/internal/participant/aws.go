package participant

import (
	"log"
	"os"
	"remind/internal/crawler"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	awsRegion         = os.Getenv("AWS_REGION")
	participantsTable = os.Getenv("DYNAMODB_TABLE_PARTICIPANTS")
)

// Dynamo types
type Participant struct {
	Name       string `json:"name"`
	DiscordId  string `json:"discord_id"`
	PhotoS3Key string `json:"photo_s3_key"`
}

func FromMember(svc *dynamodb.DynamoDB, m crawler.Member) (Participant, error) {
	var p Participant
	log.Printf("getting participant from Member %+v\n", m)

	log.Printf("getting record for %s\n", m.Name)
	res, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(participantsTable),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(m.Name),
			},
		},
	})
	if err != nil {
		log.Println("could not get participant")
    log.Println(err)
		return p, err
	}

	err = dynamodbattribute.UnmarshalMap(res.Item, &p)
	if err != nil {
		log.Println("could not unmarshal participant row")
    log.Println(err)
		return p, err
	}

	log.Printf("got participant %+v\n", p)

	return p, nil
}
