package main

import "fmt"

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

func main() {
	fmt.Println(NewGame(40, 20))
}
