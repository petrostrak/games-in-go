package main

import (
	"fmt"
	"math/rand"
	"time"
)

var dictionary = []string{
	"Zombie",
	"Gopher",
	"United States of America",
	"Indonesia",
	"Nazism",
	"Apple",
	"Programming",
}

func main() {
	rand.Seed(time.Now().UnixNano())

	targetWord := getRandomWord()
	fmt.Println(targetWord)
}

func getRandomWord() string {
	return dictionary[rand.Intn(len(dictionary))]
}
