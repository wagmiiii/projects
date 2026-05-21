package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func main() {
	w := strconv.Itoa(rand.Intn(10))
	x := strconv.Itoa(rand.Intn(10))
	y := strconv.Itoa(rand.Intn(10))
	z := strconv.Itoa(rand.Intn(10))
	secretNum := w + x + y + z
	// fmt.Println(secretNum)
	fmt.Println("Welcome to Crack The Code")
	fmt.Println("Try to determine the 4-digit code")
	reader := bufio.NewReader(os.Stdin)
	count := 0
	for {
		fmt.Print("Enter a try: ")
		
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		input = strings.TrimSpace(input)
		if len(input) != 4 {
			fmt.Println("enter a 4 digit number only")
			continue
		}
		if input == secretNum {
			fmt.Println("congrats, you are correct")
			count++
			fmt.Println("no of tries:", count)
			break
		}
		correct := 0
		for i := 0; i < 4; i++ {
			if input[i] == secretNum[i] {
				fmt.Println(string(input[i]), "is correct")
				correct++
			} else {
				fmt.Println(string(input[i]), "is not correct")
			}
		}
		count++
		if correct == 4 {
			fmt.Println("congrats, you are correct..")
			fmt.Println("no of tries:", count)
			break
		}

	}

}
