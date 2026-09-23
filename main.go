package main

import (
	"fmt"
	"math/rand"

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
	g := &Game{
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

	g.placeFood()
	g.placeMalware()

	return g
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

	termbox.SetCell(
		g.food.x,
		g.food.y,
		'●',
		termbox.ColorGreen,
		termbox.ColorDefault,
	)

	for _, malware := range g.malware {
		termbox.SetCell(
			malware.x,
			malware.y,
			'✗',
			termbox.ColorRed,
			termbox.ColorDefault,
		)
	}

	for i, segment := range g.snake {
		char := '○'
		if i == 0 {
			char = g.dir.ToRune()
		}
		termbox.SetCell(segment.x, segment.y, char, termbox.ColorGreen, termbox.ColorDefault)
	}

	scoreText := fmt.Sprintf("Score: %d Level: %d", g.score, g.level)
	for i, char := range scoreText {
		termbox.SetCell(i+2, 1, char, termbox.ColorYellow, termbox.ColorDefault)
	}
	termbox.Flush()
}

func (p Point) ToRune() rune {
	switch {
	case p.x == 0 && p.y == -1:
		return '▲'
	case p.x == 0 && p.y == 1:
		return '▼'
	case p.x == -1 && p.y == 0:
		return '◀'
	case p.x == 1 && p.y == 0:
		return '▶'
	default:
		return '●'
	}
}

func (g *Game) handleInput(ev termbox.Event) {
	if ev.Type != termbox.EventKey {
		return
	}
	switch ev.Key {
	case termbox.KeyEsc:
		close(g.quit)
		return
	}
	var newDir Point
	switch ev.Key {
	case termbox.KeyArrowUp:
		newDir = Point{x: 0, y: -1}
	case termbox.KeyArrowDown:
		newDir = Point{x: 0, y: 1}
	case termbox.KeyArrowLeft:
		newDir = Point{x: -1, y: 0}
	case termbox.KeyArrowRight:
		newDir = Point{x: 1, y: 0}
	default:
		switch ev.Ch {
		case 'w', 'W':
			newDir = Point{x: 0, y: -1}
		case 's', 'S':
			newDir = Point{x: 0, y: 1}
		case 'a', 'A':
			newDir = Point{x: -1, y: 0}
		case 'd', 'D':
			newDir = Point{x: 1, y: 0}
		case 'q', 'Q':
			close(g.quit)
			return
		default:
			return
		}
	}
	if newDir.x == -g.dir.x && newDir.y == -g.dir.y {
		return
	}
	g.dir = newDir
}

func (g *Game) isOnSnake(p Point) bool {
	for _, segment := range g.snake {
		if segment.x == p.x && segment.y == p.y {
			return true
		}
	}

	return false
}

func (g *Game) isOnMalware(p Point) bool {
	for _, malware := range g.malware {
		if malware.x == p.x && malware.y == p.y {
			return true
		}
	}

	return false
}

func (g *Game) isOutOfBounds(p Point) bool {
	return p.x < 1 ||
		p.x > g.width-2 ||
		p.y < 1 ||
		p.y > g.height-2
}

func (g *Game) placeFood() {
	for {
		p := Point{
			x: rand.Intn(g.width-2) + 1,
			y: rand.Intn(g.height-2) + 1,
		}

		if !g.isOnSnake(p) && !g.isOnMalware(p) {
			g.food = p
			return
		}
	}
}

func (g *Game) placeMalware() {
	for {
		p := Point{x: rand.Intn(g.width-2) + 1, y: rand.Intn(g.height-2) + 1}
		if !g.isOnSnake(p) && !g.isOnMalware(p) && (p.x != g.food.x || p.y != g.food.y) {
			g.malware = append(g.malware, p)
			return
		}
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer termbox.Close()
	g := NewGame(40, 20)
	events := make(chan termbox.Event)
	go func() {
		for {
			ev := termbox.PollEvent()
			events <- ev
		}
	}()
	g.draw()
	for {
		select {
		case ev := <-events:
			g.handleInput(ev)
			g.draw()
		case <-g.quit:
			return
		}
	}
}
