package main

import (
	"fmt"
	"math/rand"
	"time"
	"unicode"
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
	targetWord = "United States of America"
	guessedLetters := initGuessedWords(targetWord)

	pringGameState(targetWord, guessedLetters)

	guessedLetters['s'] = true
	pringGameState(targetWord, guessedLetters)

	guessedLetters['e'] = true
	pringGameState(targetWord, guessedLetters)
}

func initGuessedWords(targetWord string) map[rune]bool {
	guessedLetters := map[rune]bool{}
	guessedLetters[unicode.ToLower(rune(targetWord[0]))] = true
	guessedLetters[unicode.ToLower(rune(targetWord[len(targetWord)-1]))] = true

	return guessedLetters
}

func getRandomWord() string {
	return dictionary[rand.Intn(len(dictionary))]
}

func pringGameState(targetWord string, guessedLetters map[rune]bool) {
	for _, ch := range targetWord {
		if ch == ' ' {
			fmt.Print(" ")
		} else if guessedLetters[unicode.ToLower(ch)] {
			fmt.Printf("%c", ch)
		} else {
			fmt.Print("_")
		}
	}
	fmt.Println()
}
