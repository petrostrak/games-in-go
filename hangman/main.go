package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

var (
	dictionary = []string{
		"Zombie",
		"Gopher",
		"United States of America",
		"Indonesia",
		"Nazism",
		"Apple",
		"Programming",
	}

	inputReader = bufio.NewReader(os.Stdin)
)

func main() {
	rand.Seed(time.Now().UnixNano())

	targetWord := getRandomWord()
	guessedLetters := initGuessedWords(targetWord)
	hangmanState := 0

	for !isWordGuessed(targetWord, guessedLetters) && !isHangmanComplete(hangmanState) {
		pringGameState(targetWord, guessedLetters, hangmanState)
		input := readInput()
		if len(input) != 1 {
			fmt.Println("Invalid input. Please use only letters")
			continue
		}

		letter := rune(input[0])
		if isCorrectGuess(targetWord, letter) {
			guessedLetters[letter] = true
		} else {
			hangmanState++
		}
	}

}

func initGuessedWords(targetWord string) map[rune]bool {
	guessedLetters := map[rune]bool{}
	guessedLetters[unicode.ToLower(rune(targetWord[0]))] = true
	guessedLetters[unicode.ToLower(rune(targetWord[len(targetWord)-1]))] = true

	return guessedLetters
}

func isWordGuessed(targetWord string, guessedLetters map[rune]bool) bool {
	for _, ch := range targetWord {
		if !guessedLetters[ch] {
			return false
		}
	}

	return true
}

func isHangmanComplete(hangmanState int) bool {
	return hangmanState >= 9
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
func readInput() string {
	fmt.Print("> ")
	input, err := inputReader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(input)
}

func isCorrectGuess(targetWord string, letter rune) bool {
	return strings.ContainsRune(targetWord, letter)
}
