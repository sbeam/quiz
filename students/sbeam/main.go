package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Problem struct {
	Question string
	Answer   string
}

func main() {
	var (
		csvFileName = flag.String("csv", "../../problems.csv", "the location of the CSV file")
		// timeLimit   = flag.Int("timer", 30, "the time limit for the user to answer all problems")
		// shuffle   = flag.Bool("shuffle", false, "whether to shuffle the order of problems or not")
	)
	flag.Parse()

	problems := readProblems(*csvFileName)

	score := quiz(problems)

	fmt.Printf("You got %d out of %d correct\n", score, len(problems))
}

func quiz(problems []Problem) int {
	scanner := bufio.NewScanner(os.Stdin)
	score := 0

	for _, problem := range problems {
		fmt.Printf("%s : ", problem.Question)
		scanner.Scan()
		response := strings.TrimSpace(scanner.Text())
		if response == problem.Answer {
			fmt.Println("good job")
			score++
		}
	}

	return score
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
