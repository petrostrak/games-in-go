package main

import (
	"fmt"
	"io/ioutil"
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
	guessedLetters := initGuessedWords(targetWord)
	hangmanState := 0
	pringGameState(targetWord, guessedLetters, hangmanState)

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

func pringGameState(targetWord string, guessedLetters map[rune]bool, hangmanState int) {
	fmt.Println(getWordGuessingProgess(targetWord, guessedLetters))
	fmt.Println(getHangmanDrawing(hangmanState))
}

func getWordGuessingProgess(targetWord string, guessedLetters map[rune]bool) string {
	result := ""
	for _, ch := range targetWord {
		if ch == ' ' {
			result += " "
		} else if guessedLetters[unicode.ToLower(ch)] {
			result += fmt.Sprintf("%c", ch)
		} else {
			result += "_"
		}
	}

	return result
}

func getHangmanDrawing(hangmanState int) string {
	data, err := ioutil.ReadFile(fmt.Sprintf("states/hangman%d", hangmanState))
	if err != nil {
		panic(err)
	}

	return string(data)
}
