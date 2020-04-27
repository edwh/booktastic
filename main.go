package main

import "fmt"

// How good our fuzzy match needs to be.
const confidence = 75

func main() {
	_, err := HelloWorld("Hello World!")

	if err != nil {
		panic(err)
	}
}

func HelloWorld(message string) (string, error) {
	_, err := fmt.Println(message)
	return message, err
}
