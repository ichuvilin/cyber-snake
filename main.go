package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type Point struct {
	x, y int
}

type Game struct {
	snake         []Point
	food          Point
	malware       []Point
	dir           Point
	score         int
	level         int
	gameOver      bool
	width, height int
	quit          chan struct{}
}

func NewGame(width, height int) *Game {
	return &Game{
		snake: []Point{{x: width / 2, y: height / 2}},
		dir: Point{
			x: 1,
			y: 0,
		},
		width:    width,
		height:   height,
		level:    1,
		score:    0,
		gameOver: false,
		quit:     make(chan struct{}),
	}
}

func (g *Game) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for x := 0; x < g.width; x++ {
		termbox.SetCell(x, 0, '─', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, g.height-1, '─', termbox.ColorWhite, termbox.ColorDefault)
	}
	for y := 0; y < g.height; y++ {
		termbox.SetCell(0, y, '│', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(g.width-1, y, '│', termbox.ColorWhite, termbox.ColorDefault)
	}
	termbox.SetCell(0, 0, '┌', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(g.width-1, 0, '┐', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(0, g.height-1, '└', termbox.ColorWhite, termbox.ColorDefault)
	termbox.SetCell(g.width-1, g.height-1, '┘', termbox.ColorWhite, termbox.ColorDefault)
	if len(g.snake) > 0 {
		head := g.snake[0]
		termbox.SetCell(head.x, head.y, '█', termbox.ColorGreen, termbox.ColorDefault)
	}
	scoreText := fmt.Sprintf("Score: %d Level: %d", g.score, g.level)
	for i, char := range scoreText {
		termbox.SetCell(i+2, 1, char, termbox.ColorYellow, termbox.ColorDefault)
	}
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	defer termbox.Close()

	g := NewGame(40, 20)
	g.draw()
}
