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

func initScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
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

	return screen
}

func Print(s tcell.Screen, row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			s.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func displayPaddles(screen tcell.Screen) {
	screen.Clear()
	width, height := screen.Size()
	paddleStart := height/2 - PaddleHeight/2
	Print(screen, paddleStart, 0, 1, PaddleHeight, PaddleSymbol)
	Print(screen, paddleStart, width-1, 1, PaddleHeight, PaddleSymbol)
	screen.Show()
}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	screen := initScreen()
	displayPaddles(screen)

	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
			displayPaddles(screen)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEnter {
				screen.Fini()
				os.Exit(0)
			}
		}
	}
}
