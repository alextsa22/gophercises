package main

import (
	"flag"
	"fmt"
	"github.com/alextsa22/gophercises/quiz-game/quiz"
	"log"
)

var (
	filename    = flag.String("path", "problems.csv", "path to the CSV file with questions")
	timeLimit   = flag.Int("limit", 30, "time limit for the quiz")
	shuffle = flag.Bool("shuffle", false, "shuffles questions randomly")
)

func main() {
	flag.Parse()

	q, err := quiz.NewQuiz(*filename, *timeLimit, *shuffle)
	if err != nil {
		log.Fatal(err)
	}

	score, err := q.Start()
	if err != nil {
		fmt.Printf("\n%s\n", err)
	}

	fmt.Printf("your score: %d", score)
}
