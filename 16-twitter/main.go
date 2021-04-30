package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"time"
	
	"github.com/alextsa22/gophercises/16-twitter/twitter"
)

const (
	consumerKey    = "CONSUMER_KEY"
	consumerSecret = "CONSUMER_SECRET"
)

func main() {
	var (
		keyFile    string
		usersFile  string
		tweetID    string
		numWinners int
	)
	flag.StringVar(&keyFile, "key", ".env", "the file where you store your consumer key and secret for the twitter api.")
	flag.StringVar(&usersFile, "users", "users.csv", "the file where users who have retweeted the tweet are stored. this will be created if it does not exist.")
	flag.StringVar(&tweetID, "tweet", "tweetId", "the id of the Tweet you wish to find retweeters of.")
	flag.IntVar(&numWinners, "winners", 0, "the number of winners to pick for the contest.")
	flag.Parse()

	key, secret, err := keys(keyFile)
	if err != nil {
		log.Fatal(err)
	}

	client, err := twitter.NewClient(key, secret)
	if err != nil {
		log.Fatal(err)
	}

	newUsernames, err := client.Retweeters(tweetID)
	if err != nil {
		log.Fatal(err)
	}

	existUsernames := existing(usersFile)
	allUsernames := merge(newUsernames, existUsernames)
	if err = writeUsers(usersFile, allUsernames); err != nil {
		log.Fatal(err)
	}

	if numWinners == 0 {
		return
	}

	existUsernames = existing(usersFile)
	winners := pickWinners(existUsernames, numWinners)
	fmt.Println("the winners are:")
	for _, username := range winners {
		fmt.Printf("\t%s\n", username)
	}
}

func keys(keyFile string) (key, secret string, err error) {
	if err := godotenv.Load(keyFile); err != nil {
		return "", "", fmt.Errorf("error loading .env file: %s", err)
	}

	var keys struct {
		Key    string
		Secret string
	}
	keys.Key = os.Getenv(consumerKey)
	keys.Secret = os.Getenv(consumerSecret)

	return keys.Key, keys.Secret, nil
}

func existing(usersFile string) []string {
	f, err := os.Open(usersFile)
	if err != nil {
		return []string{}
	}
	defer f.Close()

	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	users := make([]string, 0, len(lines))
	for _, line := range lines {
		users = append(users, line[0])
	}

	return users
}

func merge(a, b []string) []string {
	uniq := make(map[string]struct{}, 0)
	for _, user := range a {
		uniq[user] = struct{}{}
	}

	for _, user := range b {
		uniq[user] = struct{}{}
	}

	ret := make([]string, 0, len(uniq))
	for user := range uniq {
		ret = append(ret, user)
	}

	return ret
}

func writeUsers(usersFile string, users []string) error {
	f, err := os.OpenFile(usersFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	for _, user := range users {
		if err := w.Write([]string{user}); err != nil {
			return err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

func pickWinners(users []string, numWinners int) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := r.Perm(len(users))
	winners := perm[:numWinners]
	winnerUsers := make([]string, 0, numWinners)
	for _, idx := range winners {
		winnerUsers = append(winnerUsers, users[idx])
	}

	return winnerUsers
}
