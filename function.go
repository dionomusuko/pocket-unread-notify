package function

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

type FetchResult struct {
	List map[int]FetchItem
}

type FetchItem struct {
	ItemId                 string                       `json:"item_id"`
	ResolvedId             string                       `json:"resolved_id"`
	GivenUrl               string                       `json:"given_url"`
	GivenTitle             string                       `json:"given_title"`
	Favorite               string                       `json:"favorite"`
	Status                 string                       `json:"status"`
	TimeAdded              string                       `json:"time_added"`
	TimeUpdated            string                       `json:"time_updated"`
	TimeRead               string                       `json:"time_read"`
	TimeFavorited          string                       `json:"time_favorited"`
	SortId                 int                          `json:"sort_id"`
	ResolvedTitle          string                       `json:"resolved_title"`
	ResolvedUrl            string                       `json:"resolved_url"`
	Excerpt                string                       `json:"excerpt"`
	IsArticle              string                       `json:"is_article"`
	IsIndex                string                       `json:"is_index"`
	HasVideo               string                       `json:"has_video"`
	HasImage               string                       `json:"has_image"`
	WordCount              string                       `json:"word_count"`
	Lang                   string                       `json:"lang"`
	TopImageUrl            string                       `json:"top_image_url"`
	ListenDurationEstimate int                          `json:"listen_duration_estimate"`
	AmpUrl                 string                       `json:"amp_url"`
	DomainMetadata         map[string]string            `json:"domain_metadata"`
	Tags                   map[string]map[string]string `json:"tags"`
}

func getPocketItem(consumerKey, accessToken string) (list map[int]FetchItem, err error) {
	params := url.Values{}
	params.Set("state", "unread")
	params.Set("count", "5")
	params.Set("consumer_key", consumerKey)
	params.Set("access_token", accessToken)
	resp, _ := http.Get("https://getpocket.com/v3/get?" + params.Encode())
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read error: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(string(body))
		return
	}
	var result FetchResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("unmarshal error: %v", err)
		return
	}
	return result.List, nil
}

func Function(w http.ResponseWriter, r *http.Request) {
	slackApiKey := os.Getenv("SLACK_API_KEY")
	consumerKey := os.Getenv("POCKET_CONSUMER_KEY")
	accessToken := os.Getenv("POCKET_ACCESS_TOKEN")
	items, err := getPocketItem(consumerKey, accessToken)
	if err != nil {
		log.Printf("get error: %v", err)
	}
	for _, item := range items {
		if item.ResolvedId == "0" {
			continue
		}
		url := item.ResolvedUrl
		url = strings.TrimRight(url, "/")
		text := fmt.Sprintf("*%v*\n %v", item.ResolvedTitle, url)
		params := &slack.WebhookMessage{Text: text}
		err = slack.PostWebhook(slackApiKey, params)
		if err != nil {
			log.Println("param error")
		}
	}
}
