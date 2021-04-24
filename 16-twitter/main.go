package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	pathToEnv      = "./16-twitter/.env"
	consumerKey    = "CONSUMER_KEY"
	consumerSecret = "CONSUMER_SECRET"
)

func main() {
	if err := godotenv.Load(pathToEnv); err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}

	var keys struct {
		Key    string
		Secret string
	}
	keys.Key = os.Getenv(consumerKey)
	keys.Secret = os.Getenv(consumerSecret)

	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(keys.Key, keys.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var token oauth2.Token
	dec := json.NewDecoder(res.Body)
	if err = dec.Decode(&token); err != nil {
		log.Fatal(err)
	}

	var conf oauth2.Config
	tclient := conf.Client(context.Background(), &token)
	usernames, err := retweeters(tclient, "tweetId")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(usernames)
}

func retweeters(client *http.Client, tweetId string) ([]string, error) {
	url := fmt.Sprintf("https://api.twitter.com/1.1/statuses/retweets/%s.json", tweetId)

	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var retweets []struct {
		User struct {
			ScreenName string `json:"screen_name"`
		} `json:"user"`
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&retweets)
	if err != nil {
		return nil, err
	}

	usernames := make([]string, 0, len(retweets))
	for _, retweet := range retweets {
		usernames = append(usernames, retweet.User.ScreenName)
	}

	return usernames, nil
}
