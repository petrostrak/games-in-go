package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

func emitStr(s tcell.Screen, x, y int, str string) {
	for _, c := range str {
		s.SetContent(x, y, c, nil, tcell.StyleDefault)
		x += 1
	}
}

func displayHelloWorld(screen tcell.Screen) {
	w, h := screen.Size()
	screen.Clear()
	emitStr(screen, w/2-7, h/2, "Hello, World!")
	screen.Show()
}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
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

	displayHelloWorld(screen)

	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
			displayHelloWorld(screen)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEnter {
				screen.Fini()
				os.Exit(0)
			}
		}
	}
}
