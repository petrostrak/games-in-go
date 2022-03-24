package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const SnakeSymbol = '*'
const AppleSymbol = 0x25CF
const GameFrameWidth = 30
const GameFrameHeight = 15
const GameFrameSymbol = '║'

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int
	symbol                  rune
}

var (
	screen       tcell.Screen
	isGamePaused bool
	debugLog     string
	GameObjects  []*GameObject
)

func main() {
	initScreen()
	InitGameState()
	inputChan := InitUserInput()

	for {
		HandleUserInput(ReadInput(inputChan))
		UpdateState()
		DrawState()

		time.Sleep(75 * time.Millisecond)
	}

	screen.Fini()
}

func DrawState() {

	if isGamePaused {
		return
	}

	screen.Clear()
	PrintString(0, 0, debugLog)

	for _, obj := range GameObjects {
		Print(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}
	screen.Show()
}

func PrintString(row, col int, str string) {
	for _, c := range str {
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}

func Print(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func initScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func InitGameState() {

}

func UpdateState() {

	if isGamePaused {
		return
	}

	for i := range GameObjects {
		GameObjects[i].row += GameObjects[i].velRow
		GameObjects[i].col += GameObjects[i].velCol
	}

}

func InitUserInput() chan string {

	inputChan := make(chan string)

	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()

	return inputChan
}

func HandleUserInput(key string) {
	// _, height := screen.Size()
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(1)
	}
}

func ReadInput(inputChan chan string) string {
	var key string

	// With the default case the program will no
	// longer lock waiting for an input.
	select {
	case key = <-inputChan:
	default:
		key = ""
	}

	return key
}
