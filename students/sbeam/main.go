package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Problem struct {
	Question string
	Answer   string
}

type Quiz struct {
	Problems []Problem
	Score    int
}

func main() {
	var (
		csvFileName = flag.String("csv", "../../problems.csv", "the location of the CSV file")
		timeLimit   = flag.Int("timer", 30, "the time limit for the user to answer all problems")
		shuffle     = flag.Bool("shuffle", false, "whether to shuffle the order of problems or not")
	)
	flag.Parse()

	problems := readProblems(*csvFileName)

	if *shuffle {
		shuffleProblems(&problems)
	}

	quiz := Quiz{
		Problems: problems,
		Score:    0,
	}

	fmt.Printf("Answer %d problems in %d seconds! Press ENTER to Begin", len(problems), *timeLimit)
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	go func() {
		<-time.After(time.Duration(*timeLimit) * time.Second)
		fmt.Println("*** TIME IS UP! ***")
		final(quiz.Score, len(quiz.Problems))
		os.Exit(0)
	}()

	run(&quiz)

	final(quiz.Score, len(quiz.Problems))
}

func final(score int, possible int) {
	fmt.Printf("You got %d out of %d correct. %s.\n",
		score,
		possible,
		grade(score, possible),
	)
}

func grade(score int, possible int) string {
	ratio := float64(score) / float64(possible)
	if ratio < 0.9 {
		return "PATHETIC"
	} else if int(ratio) == 1 {
		return "PERFECT"
	} else {
		return "Not bad"
	}
}

func shuffleProblems(problems *[]Problem) {
	rand.Seed(time.Now().Unix())
	for i := len(*problems) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		(*problems)[i], (*problems)[j] = (*problems)[j], (*problems)[i]
	}
}

func run(quiz *Quiz) {
	scanner := bufio.NewScanner(os.Stdin)

	for _, problem := range quiz.Problems {
		fmt.Printf("%s : ", problem.Question)
		scanner.Scan()
		response := strings.TrimSpace(scanner.Text())
		if response == problem.Answer {
			fmt.Println("good job")
			quiz.Score++
		} else {
			fmt.Println("BZZZZT")
		}
	}
}

func readProblems(path string) []Problem {
	var problems []Problem

	csvfile, err := os.Open(path)

	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	reader := csv.NewReader(csvfile)

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		problems = append(problems, Problem{
			Question: record[0],
			Answer:   record[1],
		})
	}
	return problems

}
