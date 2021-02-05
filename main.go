package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'questions,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds.")
	flag.Parse()

	//Open and read the file
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
		os.Exit(1)
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()

	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	//split questions and answers
	problems := parseLines(lines)

	// set timer
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	//counter for correct responses
	correct := 0

	for i, problem := range problems {
		fmt.Printf("Question %d: %s = ", i+1, problem.question)
		answerChannel := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			//send the answer to the answer channel
			answerChannel <- answer
		}()

		select {
		//block and stop the program, and wait until there is a message from the channel
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
			return
		case answer := <-answerChannel:
			if answer == problem.answer {
				correct++
			}

		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

//format questions and answers and return slice
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for index, line := range lines {
		ret[index] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return ret
}

type problem struct {
	question string
	answer   string
}

//exit the program
func exit(message string) {
	fmt.Println(message)
	os.Exit(1)
}
