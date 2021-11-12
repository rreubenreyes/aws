package crawler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
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

func domQuery(root *html.Node, predicate func(*html.Node) bool) (tags []*html.Node) {
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

func LatestDrawState(url string) (*DrawState, error) {
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
	scripts := domQuery(doc, isScript)

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
