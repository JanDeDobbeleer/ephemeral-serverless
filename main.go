package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	consumerKey        string
	consumerSecret     string
	accessToken        string
	accessTokenSecret  string
	maxTweetAge        string
	interactionTimeout string
	whitelist          []string
)

// MyResponse for AWS SAM
type MyResponse struct {
	StatusCode string `json:"StatusCode"`
	Message    string `json:"Body"`
}

func setVariables() {
	consumerKey = getenv("TWITTER_CONSUMER_KEY")
	consumerSecret = getenv("TWITTER_CONSUMER_SECRET")
	accessToken = getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")
	maxTweetAge = getenv("MAX_TWEET_AGE")
	interactionTimeout = getenv("TWEET_INTERACTION_TIMEOUT")
	whitelist = getWhitelist(os.Getenv("WHITELIST"))
}

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("Missing required environment variable " + name)
	}
	return v
}

func getWhitelist(whiteList string) []string {
	if whiteList == "" {
		return make([]string, 0)
	}
	return strings.Split(whiteList, ":")
}

func getTimeline(api *anaconda.TwitterApi) ([]anaconda.Tweet, error) {
	args := url.Values{}
	args.Add("count", "200")
	args.Add("include_rts", "true")
	timeline, err := api.GetUserTimeline(args)
	if err != nil {
		return make([]anaconda.Tweet, 0), err
	}
	return timeline, nil
}

func getRepliesForTweet(api *anaconda.TwitterApi, tweetID int64) []anaconda.Tweet {
	args := url.Values{}
	args.Add("count", "200")
	args.Add("since_id", strconv.FormatInt(tweetID, 10))
	me, err := api.GetSelf(nil)
	if err != nil {
		return make([]anaconda.Tweet, 0)
	}
	queryString := fmt.Sprintf("to:%s", me.ScreenName)
	searchResponse, err := api.GetSearch(queryString, args)
	if err != nil {
		return make([]anaconda.Tweet, 0)
	}
	replies := searchResponse.Statuses[:0]
	for _, tweet := range searchResponse.Statuses {
		if tweet.InReplyToStatusID == tweetID {
			replies = append(replies, tweet)
		}
	}
	return replies
}

func isWhitelisted(id int64) bool {
	tweetID := strconv.FormatInt(id, 10)
	for _, w := range whitelist {
		if w == tweetID {
			log.Print("TWEET IS WHITELISTED: ", tweetID)
			return true
		}
	}
	return false
}

func hasOngoingInteractions(api *anaconda.TwitterApi, tweetID int64, interactionAgeLimit time.Duration) bool {
	replies := getRepliesForTweet(api, tweetID)
	for _, reply := range replies {
		createdTime, err := reply.CreatedAtTime()
		if err != nil {
			log.Print("Could not parse time ", err)
			continue
		}
		if time.Since(createdTime) < interactionAgeLimit {
			log.Print("TWEET HAS ONGOING INTERACTIONS: ", tweetID)
			return true
		}
	}
	return false
}

func deleteFromTimeline(api *anaconda.TwitterApi, tweetAgeLimit time.Duration, interactionAgeLimit time.Duration) {
	timeline, err := getTimeline(api)
	if err != nil {
		log.Print("Could not get timeline ", err)
	}
	for _, t := range timeline {
		createdTime, err := t.CreatedAtTime()
		if err != nil {
			log.Print("Could not parse time ", err)
		} else {
			if time.Since(createdTime) > tweetAgeLimit && !isWhitelisted(t.Id) && !hasOngoingInteractions(api, t.Id, interactionAgeLimit) {
				_, err := api.DeleteTweet(t.Id, true)
				log.Print("DELETED ID ", t.Id)
				log.Print("TWEET ", createdTime, " - ", t.Text)
				if err != nil {
					log.Print("Failed to delete: ", err)
				}
			}
		}
	}
	log.Print("No more tweets to delete")
}

func ephemeral() (MyResponse, error) {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	api.SetLogger(anaconda.BasicLogger)

	tweetAgeLimit, _ := time.ParseDuration(maxTweetAge)
	interactionAgeLimit, _ := time.ParseDuration(interactionTimeout)

	deleteFromTimeline(api, tweetAgeLimit, interactionAgeLimit)

	return MyResponse{
		Message:    "No more tweets to delete",
		StatusCode: "200",
	}, nil
}

func main() {
	setVariables()
	lambda.Start(ephemeral)
}
