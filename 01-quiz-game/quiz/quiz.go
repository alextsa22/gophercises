package quiz

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

type question struct {
	question string
	answer   string
}

type Quiz struct {
	questions []question
	timeLimit time.Duration
}

func NewQuiz(filename string, timeLimit int, shuffle bool) (*Quiz, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var qs []question
	reader := csv.NewReader(f)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		qs = append(qs, question{record[0], cleanLine(record[1])})
	}

	if shuffle {
		shuffleQuestions(qs)
	}

	return &Quiz{qs, time.Duration(timeLimit)}, nil
}

// Knuth Fisher-Yates shuffle
func shuffleQuestions(qs []question) {
	rand.Seed(time.Now().Unix())

	l := len(qs)
	for i := 0; i < l; i++ {
		r := i + rand.Intn(l-i)
		qs[r], qs[i] = qs[i], qs[r]
	}
}

type answer struct {
	input string
	error error
}

func (q *Quiz) Start() (int, error) {
	score := 0
	timer := time.NewTimer(q.timeLimit * time.Second)
	done := make(chan answer)

	go readInput(done)

	for _, q := range q.questions {
		res, err := askQuestion(q, *timer, done)
		if err != nil {
			return score, err
		}
		score += res
	}

	return score, nil
}

func readInput(done chan<- answer) {
	for {
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		done <- answer{cleanLine(input), err}
	}
}

func askQuestion(q question, timer time.Timer, done <-chan answer) (int, error) {
	fmt.Printf("%s: ", q.question)

	for {
		select {
		case <-timer.C:
			return 0, errors.New("time out")
		case ans := <-done:
			if ans.error != nil {
				return 0, ans.error
			}
			if ans.input != q.answer {
				return 0, nil
			}
			return 1, nil
		}
	}
}

func cleanLine(in string) string {
	return strings.ToLower(strings.TrimSpace(strings.Trim(in, "\n")))
}
