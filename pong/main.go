package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	PaddleHeight = 4
	PaddleSymbol = 0x2588
	BallSymbol   = 0x25CF
)

var (
	screen        tcell.Screen
	playerPaddle1 *GameObject
	playerPaddle2 *GameObject
	debugLog      string
	ball          *GameObject
	GameObjects   []*GameObject
)

type GameObject struct {
	row, col, width, height int
	symbol                  rune
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

func Print(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func PrintString(row, col int, str string) {
	for _, c := range str {
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}

func DrawState() {
	screen.Clear()
	PrintString(0, 0, debugLog)

	for _, obj := range GameObjects {
		Print(obj.row, obj.col, obj.width, obj.height, PaddleSymbol)
	}
	screen.Show()
}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	initScreen()
	InitGameState()
	inputChan := InitUserInput()
	DrawState()

	for {
		DrawState()
		time.Sleep(50 * time.Millisecond)

		key := ReadInput(inputChan)
		HandleUserInput(key)
	}
}

func InitGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - PaddleHeight/2

	playerPaddle1 = &GameObject{
		row:    paddleStart,
		col:    0,
		width:  1,
		height: PaddleHeight,
		symbol: PaddleSymbol,
	}

	playerPaddle2 = &GameObject{
		row:    paddleStart,
		col:    width - 1,
		width:  1,
		height: PaddleHeight,
		symbol: PaddleSymbol,
	}

	ball = &GameObject{
		row:    height / 2,
		col:    width / 2,
		width:  1,
		height: 1,
		symbol: BallSymbol,
	}

	GameObjects = []*GameObject{
		playerPaddle1, playerPaddle2, ball,
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

func checkTopBoundry(height int, player *GameObject) bool {
	return player.row > 0
}
func checkBottomBoundry(height int, player *GameObject) bool {
	return player.row+player.height < height
}

func HandleUserInput(key string) {
	_, height := screen.Size()
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(1)
	} else if key == "Rune[w]" && checkTopBoundry(height, playerPaddle1) {
		playerPaddle1.row--
	} else if key == "Rune[s]" && checkBottomBoundry(height, playerPaddle1) {
		playerPaddle1.row++
	} else if key == "Up" && checkTopBoundry(height, playerPaddle2) {
		playerPaddle2.row--
	} else if key == "Down" && checkBottomBoundry(height, playerPaddle2) {
		playerPaddle2.row++
	}
}
