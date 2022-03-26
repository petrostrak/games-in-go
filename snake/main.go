package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	SnakeSymbol                = '█'
	AppleSymbol                = 0x25CF
	GameFrameWidth             = 30
	GameFrameHeight            = 15
	GameFrameSymbolHorizontal  = '═'
	GameFrameSymbolVertical    = '║'
	GameFrameSymbolTopLeft     = '╔'
	GameFrameSymbolTopRight    = '╗'
	GameFrameSymbolBottomLeft  = '╚'
	GameFrameSymbolBottomRight = '╝'
)

type Point struct {
	row, col int
}

type Snake struct {
	parts          []*Point
	velRow, velCol int
	symbol         rune
}

type Apple struct {
	point  *Point
	symbol rune
}

var (
	screen       tcell.Screen
	isGamePaused bool
	isGameOver   bool
	debugLog     string
	snake        *Snake
	apple        *Apple
	score        int
)

func main() {
	rand.Seed(time.Now().UnixNano())

	initScreen()
	InitGameState()
	inputChan := InitUserInput()

	for !isGameOver {
		HandleUserInput(ReadInput(inputChan))
		UpdateState()
		DrawState()

		time.Sleep(75 * time.Millisecond)
	}

	sWidth, sHeight := screen.Size()
	PrintStringCentered(sHeight/2, sWidth/2, "Game Over!")
	PrintStringCentered(sHeight/2, sWidth/2, fmt.Sprintf("Your score is %d!", score))
	screen.Show()
	time.Sleep(2 * time.Second)
	screen.Fini()
}

func DrawState() {

	if isGamePaused {
		return
	}

	screen.Clear()
	PrintString(0, 0, debugLog)
	PrintGameFrame()

	PrintSnake()
	PrintApple()

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
	snake = &Snake{
		parts: []*Point{
			{row: 9, col: 3}, // head
			{row: 8, col: 3},
			{row: 7, col: 3},
			{row: 6, col: 3},
			{row: 5, col: 3}, // tail
		},
		velRow: -1,
		velCol: 0,
		symbol: SnakeSymbol,
	}

	apple = &Apple{
		point: &Point{
			row: 10,
			col: 10,
		},
		symbol: AppleSymbol,
	}
}

func UpdateState() {

	if isGamePaused {
		return
	}

	// Update Snake + Apple
	UpdateSnake()
	UpdateApple()

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
	debugLog = key
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(1)
	} else if key == "Rune[p]" {
		isGamePaused = true
	} else if key == "Rune[w]" && snake.velRow != 1 {
		snake.velRow = -1
		snake.velCol = 0
	} else if key == "Rune[a]" && snake.velCol != 1 {
		snake.velRow = 0
		snake.velCol = -1
	} else if key == "Rune[s]" && snake.velRow != -1 {
		snake.velRow = 1
		snake.velCol = 0
	} else if key == "Rune[d]" && snake.velCol != -1 {
		snake.velRow = 0
		snake.velCol = 1
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

func PrintStringCentered(row, col int, s string) {
	col = col - len(s)/2
	PrintString(row, col, s)
}

func PrintFilledRect(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func PrintUnfilledRect(row, col, width, height int) {
	for c := 0; c < width; c++ {
		if c == 0 {
			screen.SetContent(col+c, row, GameFrameSymbolTopLeft, nil, tcell.StyleDefault)
			continue
		}

		if c == width-1 {
			screen.SetContent(col+c, row, GameFrameSymbolTopRight, nil, tcell.StyleDefault)
			continue
		}

		screen.SetContent(col+c, row, GameFrameSymbolHorizontal, nil, tcell.StyleDefault)
	}

	for r := 1; r < height-1; r++ {
		screen.SetContent(col, row+r, GameFrameSymbolVertical, nil, tcell.StyleDefault)
		screen.SetContent(col+width-1, row+r, GameFrameSymbolVertical, nil, tcell.StyleDefault)
	}

	for c := 0; c < width; c++ {
		if c == 0 {
			screen.SetContent(col+c, row+height-1, GameFrameSymbolBottomLeft, nil, tcell.StyleDefault)
			continue
		}

		if c == width-1 {
			screen.SetContent(col+c, row+height-1, GameFrameSymbolBottomRight, nil, tcell.StyleDefault)
			continue
		}

		screen.SetContent(col+c, row+height-1, GameFrameSymbolHorizontal, nil, tcell.StyleDefault)
	}
}

func PrintGameFrame() {
	gameFrameTopLeftRow, gameFrameTopLeftCol := GetGameFrameTopLeft()
	row, col := gameFrameTopLeftRow-1, gameFrameTopLeftCol-1
	width, height := GameFrameWidth+2, GameFrameHeight+2

	PrintUnfilledRect(row, col, width, height)
	// PrintUnfilledRect(row+1, col+1, GameFrameWidth, GameFrameHeight, '*')
}

func PrintFilledRectInGameFrame(row, col, width, height int, ch rune) {
	r, c := GetGameFrameTopLeft()
	PrintFilledRect(row+r, col+c, width, height, ch)
}

func PrintSnake() {
	for _, p := range snake.parts {
		PrintFilledRectInGameFrame(p.row, p.col, 1, 1, snake.symbol)
	}
}

func PrintApple() {
	PrintFilledRectInGameFrame(apple.point.row, apple.point.col, 1, 1, apple.symbol)
}

func GetGameFrameTopLeft() (int, int) {
	// Get game frames top left point (row, col)
	sWidth, sHeight := screen.Size()
	return sHeight/2 - GameFrameHeight/2, sWidth/2 - GameFrameWidth/2
}

func UpdateSnake() {
	// add a new element
	head := GetSnakeHead()
	snake.parts = append(snake.parts, &Point{
		row: head.row + snake.velRow,
		col: head.col + snake.velCol,
	})

	// if the head of the snake is equal to the apple
	// it means that the snake has eaten the apple. In
	// that case we skip the deletion of the last element
	// of the snake, making it longer.
	if !AppleIsInsideSnake() {
		snake.parts = snake.parts[1:]
	} else {
		score++
	}

	if IsOutSideGameFrame() || IsEatingItself() {
		isGameOver = true
	}

}

func UpdateApple() {
	if AppleIsInsideSnake() {
		apple.point.row, apple.point.col = rand.Intn(GameFrameHeight), rand.Intn(GameFrameWidth)
	}
}

func AppleIsInsideSnake() bool {
	for _, p := range snake.parts {
		if p.row == apple.point.row && p.col == apple.point.col {
			return true
		}
	}

	return false
}

func IsOutSideGameFrame() bool {
	head := GetSnakeHead()
	return head.row < 0 ||
		head.row >= GameFrameHeight ||
		head.col < 0 ||
		head.col >= GameFrameWidth
}

func IsEatingItself() bool {
	head := GetSnakeHead()

	// range through all snake parts except the last one which
	// is the head.
	for _, p := range snake.parts[:GetSnakeHeadIndex()] {
		if p.row == head.row && p.col == head.col {
			return true
		}
	}

	return false
}

func GetSnakeHeadIndex() int {
	return len(snake.parts) - 1
}

func GetSnakeHead() *Point {
	return snake.parts[GetSnakeHeadIndex()]
}
