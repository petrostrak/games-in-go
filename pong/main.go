package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

const (
	PaddleHeight = 4
	PaddleSymbol = 0x2588
)

var (
	screen  tcell.Screen
	player1 *Paddle
	player2 *Paddle
)

type Paddle struct {
	row, col, width, height int
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

func DrawState() {
	screen.Clear()
	Print(player1.row, player1.col, player1.width, player1.height, PaddleSymbol)
	Print(player2.row, player2.col, player2.width, player2.height, PaddleSymbol)
	screen.Show()
}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	initScreen()
	InitGameState()

	DrawState()

	for {
		DrawState()

		switch ev := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
			DrawState()
		case *tcell.EventKey:
			if ev.Rune() == 'q' {
				screen.Fini()
				os.Exit(0)
			} else if ev.Rune() == 'w' {
				player1.row--
			} else if ev.Rune() == 's' {
				player1.row++
			} else if ev.Key() == tcell.KeyUp {
				player2.row--
			} else if ev.Key() == tcell.KeyDown {
				player2.row++
			}
		}
	}
}

func InitGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - PaddleHeight/2

	player1 = &Paddle{
		row:    paddleStart,
		col:    0,
		width:  1,
		height: PaddleHeight,
	}

	player2 = &Paddle{
		row:    paddleStart,
		col:    width - 1,
		width:  1,
		height: PaddleHeight,
	}
}
