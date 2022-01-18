package state

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type GameState struct {
	IncGuesses int // number of incorrect guesses, this corresponds to the hangman files in data and after 9 its gamer over
	CorGuesses int
	Option     rune     // whether player wants to play or quit
	Word       string   // word that the player has to guess
	Guesses    []string // stores every guess and will cross reference Word for what was correct
	QuotedBy   string   // The person that the quote came from if the Word is a quote
}

func GetCurrentState(incGuesses int) string {
	hangmanBytes, err := os.ReadFile("./state/data/hangman" + fmt.Sprint(incGuesses) + ".txt")
	if err != nil {
		color.Red("Error retrieving hangman picture. ", err.Error())
	}
	if len(hangmanBytes) > 0 {
		return string(hangmanBytes)
	} else {
		return ""
	}
}
