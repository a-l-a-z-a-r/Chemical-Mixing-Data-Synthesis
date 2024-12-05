package main

import (
	"bufio"
	"fmt" // main package used to do basic processs
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var status string
	fmt.Print("Enter grade: ")
	reader := bufio.NewReader(os.Stdin)   // operating system input
	input, err := reader.ReadString('\n') // Read input until the newline characte
	input = strings.TrimSpace(input)
	grade, err := strconv.ParseFloat(input, 64)
	if err != nil {
		log.Fatal(err)
	}
	if grade >= 60 {
		status = "Congratulations   we are happy to tell you that you are passing"
	} else {
		status = "I am sorry to tell you that you are failing"
	}
	fmt.Println(status)
}

func main() {
	var status string
	reader := bufio.NewReader(os.Stdin)
	seconds := time.Now().Unix()
	rand.Seed(seconds)
	target := rand.Intn(100) + 1
	fmt.Println("guess my random number: ")
	for guesses := 0; guesses < 10; guesses++ {

		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input)
		guess, err := strconv.Atoi(input)
		if err != nil {
			log.Fatal(err)
		}
		if guess < target {
			status = "guess to low"

		}
		if guess > target {
			status = "guess was to high"
		}

		if guess == target {
			status = "Congratulations you have won"
			fmt.Println(status)
			break
		}
		fmt.Println(status)
	}
	fmt.Println(target)
}
