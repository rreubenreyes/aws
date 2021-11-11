package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

  "github.com/aws/aws-lambda-go/lambda"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb/expression"
  "golang.org/x/net/html"
)

var url = os.Getenv("DRAW_URL")

type Member struct {
  Name        string `json:"name"`
  Email       string `json:"email"`
  DrawViewed  bool   `json:"drawViewed"`
  IsOrganizer bool   `json:"isOrganizer"`
}

type DrawState struct {
  Members []Member `json:"members"`
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

func LatestDrawState() (state *DrawState) {
  // get the starting page
  resp, err := http.Get(url)
  if err != nil {
    panic("invalid url")
  }
  defer resp.Body.Close()

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    panic("invalid response")
  }

  // parse html
  htm := string(body)
  doc, err := html.Parse(strings.NewReader(htm))
  if err != nil {
    panic("could not parse html")
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
    content := script.FirstChild.Data
    r := expr.Find([]byte(content))
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
    panic(err)
  }

  return
}

func RemainingParticipants(members []Member) {
  svc := dynamodb.New(session.New())
  for _, member := range members {
    go func() {
      svc.GetItem(&dynamodb.GetItemInput{
        Key: map[string]*dynamodb.AttributeValue{
          "name": {
            S: aws.String(member.Name),
          },
        },
      })
    }()
  }
}

func remind(ctx context.Context) (string, error) {
  state := LatestDrawState()
  for _, member := range state.Members {
    if !member.DrawViewed {
      fmt.Printf("will remind %s\n", member.Name)
    }
  }
  return "ok", nil
}

func main() {
  lambda.Start(remind)
}
