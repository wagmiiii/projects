package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func main() {
	randomNumber := rand.Intn(101)

	fmt.Println("welcome to higher or lower (0 - 100)")

	reader := bufio.NewReader(os.Stdin)

	success := false

	guesses := 0

	for guesses < 5 {
		fmt.Print("enter a number: ")

		input, err := reader.ReadString('\n')
													
		if err != nil {

			log.Fatal(err)
		}

		input = strings.TrimSpace(input)

		guess, err := strconv.Atoi(input)

		if err != nil {
			fmt.Println("please enter only numbers")

			continue
		}
		if guess < randomNumber {
			fmt.Println("higher")

			guesses++

			fmt.Println("you have", 5-guesses, "guesses left")

		} else if guess > randomNumber {

			fmt.Println("lower")

			guesses++

			fmt.Println("you have", 5-guesses, "guesses left")

		} else {

			fmt.Println("good job, you guessed it")

			success = true

			break
		}
	}

	if !success {
		fmt.Println("sorry, you are out of guesses. the answer is: ", randomNumber)
	}
}
