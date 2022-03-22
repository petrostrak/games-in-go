package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	PaddleHeight           = 4
	PaddleSymbol           = 0x2588
	BallSymbol             = 0x25CF
	InitialBallVelocityRow = 1
	InitialBallVelocityCol = 2
)

var (
	screen        tcell.Screen
	playerPaddle1 *GameObject
	playerPaddle2 *GameObject
	debugLog      string
	ball          *GameObject
	GameObjects   []*GameObject
	isGamePaused  bool
)

type GameObject struct {
	row, col, width, height int
	symbol                  rune
	velRow, velCol          int
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

	if isGamePaused {
		return
	}

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
		HandleUserInput(ReadInput(inputChan))
		UpdateState()
		DrawState()

		time.Sleep(75 * time.Millisecond)
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
		velRow: 0,
		velCol: 0,
	}

	playerPaddle2 = &GameObject{
		row:    paddleStart,
		col:    width - 1,
		width:  1,
		height: PaddleHeight,
		symbol: PaddleSymbol,
		velRow: 0,
		velCol: 0,
	}

	ball = &GameObject{
		row:    height / 2,
		col:    width / 2,
		width:  1,
		height: 1,
		symbol: BallSymbol,
		velRow: InitialBallVelocityRow,
		velCol: InitialBallVelocityCol,
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
	} else if key == "Rune[p]" {
		isGamePaused = !isGamePaused
	}
}

func UpdateState() {

	if isGamePaused {
		return
	}

	for i := range GameObjects {
		GameObjects[i].row += GameObjects[i].velRow
		GameObjects[i].col += GameObjects[i].velCol
	}

	if CollidesWithWall(ball) {
		ball.velRow = -ball.velRow
	}

	if CollidesWithPaddle(ball, playerPaddle1) || CollidesWithPaddle(ball, playerPaddle2) {
		ball.velCol = -ball.velCol
	}
}

func CollidesWithWall(obj *GameObject) bool {
	_, screenHeight := screen.Size()
	return obj.row+obj.velRow < 0 || obj.row+obj.velRow >= screenHeight
}

func CollidesWithPaddle(ball, paddle *GameObject) bool {
	var collidesOnColumn bool
	if ball.col < paddle.col {
		collidesOnColumn = ball.col+ball.velCol >= paddle.col
	} else {
		collidesOnColumn = ball.col+ball.velCol <= paddle.col
	}

	return collidesOnColumn && ball.row+ball.velRow >= paddle.row && ball.row+ball.velRow < paddle.row+paddle.height
}
