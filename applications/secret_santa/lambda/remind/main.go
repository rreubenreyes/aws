package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/html"
)

var (
	awsRegion            = os.Getenv("AWS_REGION")
	url                  = os.Getenv("DRAW_URL")
	participantsTable    = os.Getenv("DYNAMODB_TABLE_PARTICIPANTS")
	discordTokenSecretId = os.Getenv("DISCORD_TOKEN_SECRET_ID")
	discordChannelId     = os.Getenv("DISCORD_CHANNEL_ID")
)

// DrawNames types
type Member struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	DrawViewed  bool   `json:"drawViewed"`
	IsOrganizer bool   `json:"isOrganizer"`
}

type DrawState struct {
	Members []Member `json:"members"`
}

// Dynamo types
type Participant struct {
	Name       string `json:"name"`
	DiscordId  string `json:"discord_id"`
	PhotoS3Key string `json:"photo_s3_key"`
}

func DOMQuery(root *html.Node, predicate func(*html.Node) bool) (tags []*html.Node) {
	// bfs to get all tags matching type
	queue := []*html.Node{root}
	for len(queue) > 0 {
		cur := queue[0]
		if predicate(cur) {
			tags = append(tags, cur)
		}
		queue = queue[1:]
		for next := cur.FirstChild; next != nil; next = next.NextSibling {
			queue = append(queue, next)
		}
	}

	return
}

func LatestDrawState() (*DrawState, error) {
  log.Println("getting draw state")

	var state *DrawState
	// get the starting page
	resp, err := http.Get(url)
	if err != nil {
		log.Println("could not get draw page")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("could not read response from draw page")
		return nil, err
	}

	// parse html
	htm := string(body)
	doc, err := html.Parse(strings.NewReader(htm))
	if err != nil {
		log.Println("could not parse draw page html")
		return nil, err
	}

	// get all script tags
	isScript := func(node *html.Node) bool {
		if node.Type != html.ElementNode {
			return false
		}
		if node.Data != "script" {
			return false
		}
		return true
	}
	scripts := DOMQuery(doc, isScript)

	// find the script tag containing draw state
	var rawState string
	exprs := "modelConfig = JSON.parse" + regexp.QuoteMeta("(") + "'(.*)'" + regexp.QuoteMeta(")")
	expr := regexp.MustCompile(exprs)
	for _, script := range scripts {
		r := expr.Find([]byte(script.FirstChild.Data))
		if r != nil {
			s := strings.Replace(string(r), "modelConfig = JSON.parse('", "", 1)
			rawState = strings.Replace(s, "')", "", 1)
			rawState = strings.ReplaceAll(rawState, "\\", "")

			// extremely hacky replacing bad strings because these things are unescaped apparently
			rawState = strings.ReplaceAll(rawState, `href="#REPLACE#"`, "")
			rawState = strings.ReplaceAll(rawState, `"u003cnobru003eChange email addressesu003c/nobru003e"`, "")
			break
		}
	}

	// unmarshal draw state
	err = json.Unmarshal([]byte(rawState), &state)
	if err != nil {
		log.Println("could not parse draw state")
		return nil, err
	}

	log.Printf("got draw state %+v\n", state)

	return state, nil
}

func RemainingParticipants(members []Member) ([]Participant, error) {
  log.Println("getting remaining participants")
	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String(awsRegion)}))
	var m []Member
	for _, member := range members {
		if !member.DrawViewed {
			m = append(m, member)
		}
	}

	log.Printf("getting records for remaining draw members %+v\n", m)

	// get all participants in parallel
	ret := make(chan Participant)
	errs := make(chan error)
	var wg sync.WaitGroup
	var participants []Participant
	for _, member := range m {
		wg.Add(1)
		go func(member Member) {
			defer wg.Done()

      log.Printf("getting record for %s\n", member.Name)
			res, err := svc.GetItem(&dynamodb.GetItemInput{
				TableName: aws.String(participantsTable),
				Key: map[string]*dynamodb.AttributeValue{
					"name": {
						S: aws.String(member.Name),
					},
				},
			})
			if err != nil {
				log.Println("could not get participant")
				errs <- err
			}

			var p Participant
			err = dynamodbattribute.UnmarshalMap(res.Item, &p)
			if err != nil {
				log.Println("could not unmarshal participant row")
				errs <- err
			}

      log.Printf("got participant %+v\n", p)

			ret <- p
		}(member)
	}

  go func() {
    wg.Wait()
    close(ret)
    close(errs)
  }()

  for p := range ret {
		participants = append(participants, p)
  }
  for err := range errs {
    if err != nil {
      return nil, err
    }
  }

	log.Printf("got participants %s", participants)

	return participants, nil
}

func SendReminders(participants []Participant) error {
	svc := secretsmanager.New(session.New(&aws.Config{Region: aws.String(awsRegion)}))
	res, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(discordTokenSecretId),
	})
	if err != nil {
		log.Println("could not get secret")
		return err
	}

	log.Printf("got secret %s\n", res.String())

	sess, err := discordgo.New("Bot " + *res.SecretString)
	if err != nil {
		log.Println("could not connect to Discord")
		return err
	}

	defer sess.Close()

	var wg sync.WaitGroup
	errs := make(chan error)
	for _, p := range participants {
		wg.Add(1)
		go func(p Participant) {
			defer wg.Done()

			// TODO: embed a wide photo
			_, err := sess.ChannelMessageSend(
				discordChannelId,
				fmt.Sprintf("hey %s go sign up for secret santa", p.Name),
			)
			if err != nil {
				log.Println("could not send message")
				errs <- err
			}
		}(p)
	}

  go func() {
    wg.Wait()
    close(errs)
  }()

  for err := range errs {
    if err != nil {
      return err
    }
  }

	return nil
}

func remind() error {
	log.Println("starting")

	state, err := LatestDrawState()
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
