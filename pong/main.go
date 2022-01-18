package main

import (
	"fmt"
	"math"
	"os"

	"github.com/fatih/color"
	"github.com/gdamore/tcell"
)

type E_ERROR int

const (
	SCREEN_NOT_INITIALIZED E_ERROR = iota + 1
)

var PADDLE_SYMBOL = 0x2588

func printAlongCol(s tcell.Screen, x, y, rowlen int, str int) {
	for i := 0; i < rowlen; i++ {
		s.SetContent(x, y, rune(str), nil, tcell.StyleDefault)
		y++
	}
}

func print(s tcell.Screen, x, y int, str string) {
	for _, c := range str {
		s.SetContent(x, y, c, nil, tcell.StyleDefault)
		x += 1
	}
}

/*
	handleError handles a given err and returns
*/
func handleError(fmt string, panicCode E_ERROR, err error, args ...interface{}) {
	if err != nil {
		color.Red(fmt, err.Error(), args)
		panic(int(panicCode))
	}
}

func handleWarn(fmt string, code E_ERROR, err error, args ...interface{}) bool {
	if err != nil {
		color.Yellow(fmt, err.Error(), args)
		return true
	}
	return false
}

type Position struct {
	x int
	y int
}

type Player struct {
	pos     Position
	score   int
	bgColor tcell.Color
}

type PlayerEvent interface {
	moveUp()
	moveDown()
}

type GameState struct {
	screen       tcell.Screen
	ballPosition Position
	dbg          bool
	leftPaddle   Player
	rightPaddle  Player
	paddleLen    int
	isBallMoving bool
}

func constructGameState(s tcell.Screen) GameState {
	w, h := s.Size()
	state := GameState{screen: s,
		isBallMoving: false,
		dbg:          true,
		paddleLen:    4,
		ballPosition: Position{x: w / 2, y: h / 2}}
	state.leftPaddle = Player{pos: Position{x: 0, y: (h / 2) - state.paddleLen/2}, score: 0}
	state.rightPaddle = Player{pos: Position{x: w - 1, y: (h / 2) - state.paddleLen/2}, score: 0}
	return state
}

func main() {
	// init screen and state
	screen, err := tcell.NewScreen()
	handleError("Screen could not be initialized", SCREEN_NOT_INITIALIZED, err, "")
	handleError("Screen could not be initialized", SCREEN_NOT_INITIALIZED, screen.Init(), "")
	// @TODO: Print welcome message and ask for user input to start game or leave
	// if player wants to play then call startGame else os.Exit
	// then clear screen right before we start game again
	screen.Clear()
	startGame(screen)
	// @TODO: play again?
}

func startGame(screen tcell.Screen) {
	gameState := constructGameState(screen)
	// start loop
	gameLoop(&gameState)
}

func gameLoop(state *GameState) {
	// handle player event
	var event chan string = make(chan string, 3)
	go handleEvent(state, event)
	for {
		select {
		case e := <-event:
			if state.dbg {
				print(state.screen, 1, 2, "Event handled! redraw")
			}
			if e == "redraw" {
			}
		default:
			if state.dbg {
				print(state.screen, 1, 1, "No Event sent")
			}
		}
		draw(*state)
		state.ballPosition.x++
		print(state.screen,
			state.ballPosition.x,
			state.ballPosition.y, ".")

		state.screen.Clear()
		// player movement
		// draw ball
		// update ball movement
		// handle collisions
		//   event: when hits boudary -> point
		//   event: when hits paddle -> calculate new trajectory for ball
		// handle game over
	}
}

func draw(state GameState) {
	printAlongCol(state.screen,
		state.leftPaddle.pos.x,
		state.leftPaddle.pos.y,
		state.paddleLen, PADDLE_SYMBOL)
	printAlongCol(state.screen,
		state.rightPaddle.pos.x,
		state.rightPaddle.pos.y,
		state.paddleLen, PADDLE_SYMBOL)
}

func handleEvent(state *GameState, event chan string) {
	// eventMap
	eventMap := map[string]*Player{
		"Rune[w]": &state.leftPaddle,
		"Rune[s]": &state.leftPaddle,
		"Up":      &state.rightPaddle,
		"Down":    &state.rightPaddle}
	for {
		switch ev := state.screen.PollEvent().(type) {
		case *tcell.EventResize:
			//w, h := ev.Size()
			state.screen.Sync()
			// recalculate position
		case *tcell.EventKey:
			keyPress := ev.Name()
			p, isin := eventMap[keyPress]
			if isin {
				w, h := state.screen.Size()
				if state.dbg {
					print(state.screen, 1, 0, fmt.Sprintf("KeyActivated! Key: %s, Player: %s, (H,W): (%d,%d)", ev.Name(), printPlayer(*p), h, w))
				}
				if keyPress == "Rune[w]" || keyPress == "Up" {
					p.moveUp()
				} else {
					p.moveDown()
				}
				// this just clamps the position between 0.0 and the current w/h of the screen
				//  I'm using the math Max and Min which requires floats
				//  for height, subtract 1 so that it doesn't go off the screen
				p.pos.x = int(math.Min(math.Max(float64(p.pos.x), 0.0), float64(w)))
				p.pos.y = int(math.Min(math.Max(float64(p.pos.y), 0.0), float64(h-1-state.paddleLen)))
				event <- "redraw"
			}
			if ev.Key() == tcell.KeyEscape {
				color.Blue("Exit key clicked, exiting now!")
				state.screen.Fini()
				os.Exit(0)
			}
			state.screen.Show()
		}
	}
}

func (p *Player) moveUp() {
	p.pos.y--
}

func (p *Player) moveDown() {
	p.pos.y++
}

func printPlayer(p Player) string {
	return fmt.Sprintf("Player{x:%d,y:%d,score:%d}", p.pos.x, p.pos.y, p.score)
}
