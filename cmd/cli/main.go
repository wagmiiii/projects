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
	fmt.Println(randomNumber)

	fmt.Println("welcome to higher or lower (0 - 100)")

	reader := bufio.NewReader(os.Stdin)

	success := false
	count := 0



	for {
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
			count++


		} else if guess > randomNumber {

			fmt.Println("lower")
			count++

		} else {
			count++
			fmt.Println("good job, you guessed it with",count,"tries")

			success = true

			break
		}
	}

	if !success {
		fmt.Println("sorry, you are out of guesses. the answer is: ", randomNumber)
	}
}
