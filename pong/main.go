package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gdamore/tcell"
)

type E_ERROR int

const (
	SCREEN_NOT_INITIALIZED E_ERROR = iota + 1
)

const PADDLE_SYMBOL = 0x2588
const BALL_SYMBOL = 0x25CF

func printChar(s tcell.Screen, x, y, h, w int, str rune) {
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			s.SetContent(x, y, rune(str), nil, tcell.StyleDefault)
			y++
		}
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
	x    int
	y    int
	char rune
}

type VecMath interface {
	add(pos1, pos2 *Position)
}

type Ball struct {
	pos      Position
	velocity Position
	radius   int
}

type Player struct {
	pos     Position
	bgColor tcell.Color
	score   int
}

type PlayerEvent interface {
	moveUp()
	moveDown()
}

type GameEvent struct {
	p     Player
	event string
	key   string
}

type GameState struct {
	screen       tcell.Screen
	ball         Ball
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
		paddleLen:    4}

	state.ball = Ball{pos: Position{x: w / 2, y: h / 2, char: BALL_SYMBOL}, velocity: Position{x: 1, y: 2}, radius: 1}
	state.leftPaddle = Player{pos: Position{x: 0, y: (h / 2) - state.paddleLen/2, char: PADDLE_SYMBOL}, score: 0}
	state.rightPaddle = Player{pos: Position{x: w - 1, y: (h / 2) - state.paddleLen/2, char: PADDLE_SYMBOL}, score: 0}
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
	event := make(chan GameEvent)
	go handleEvent(state, event)
	for {
		state.screen.Clear()
		select {
		case e := <-event:
			if state.dbg {
				h, w := state.screen.Size()
				print(state.screen, 1, 0, fmt.Sprintf("KeyActivated! EventType: %s, Key: %s, Player: %s, (H,W): (%d,%d)", e.event, e.key, printPlayer(e.p), h, w))
			}
		default:
			time.Sleep(time.Millisecond * 75)
		}
		update(state)
		draw(*state)
		state.screen.Show()
		// player movement
		// draw ball
		// update ball movement
		// handle collisions
		//   event: when hits boudary -> point
		//   event: when hits paddle -> calculate new trajectory for ball
		// handle game over
	}
}

func update(state *GameState) {
	state.ball.pos.add(state.ball.velocity)
}

func draw(state GameState) {
	printChar(state.screen,
		state.leftPaddle.pos.x,
		state.leftPaddle.pos.y,
		state.paddleLen, 1,
		state.leftPaddle.pos.char)
	printChar(state.screen,
		state.rightPaddle.pos.x,
		state.rightPaddle.pos.y,
		state.paddleLen, 1,
		state.rightPaddle.pos.char)
	printChar(state.screen,
		state.ball.pos.x,
		state.ball.pos.y,
		state.ball.radius,
		state.ball.radius,
		state.ball.pos.char)
}

func handleEvent(state *GameState, event chan GameEvent) {
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
				event <- GameEvent{event: "redraw", p: *p, key: keyPress}
			}
			if ev.Key() == tcell.KeyEscape {
				color.Blue("Exit key clicked, exiting now!")
				state.screen.Fini()
				os.Exit(0)
			}
		}
	}
}

func (p *Player) moveUp() {
	p.pos.y--
}

func (p *Player) moveDown() {
	p.pos.y++
}

func (pos *Position) add(other Position) {
	pos.x += other.x
	pos.y += other.y
}

func printPlayer(p Player) string {
	return fmt.Sprintf("Player{x:%d,y:%d,score:%d}", p.pos.x, p.pos.y, p.score)
}
