package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"hangman.com/state"
)

func main() {
	s := setupGame()
	for int('Q') != int(s.Option) && int('q') != int(s.Option) {
		printStartMessage(s.Word)
		reader := bufio.NewReader(os.Stdin)
		tempOpt, sizeT, err := reader.ReadRune()
		if err != nil {
			color.Red("Error found taking in player guess. Inputs received: %d, Error: %s\n", sizeT, err.Error())
		}
		s.Option = tempOpt

		if int('p') == int(s.Option) || int('P') == int(s.Option) {
			fmt.Println("Game Starting!")
			playGame(&s)
			s = setupGame()
		}
	}
	fmt.Println("Thank you for playing")
}

func playGame(s *state.GameState) {
	for s.IncGuesses <= 9 {
		mp, _ := printHangman(s)
		accepted, isGuessInWord := receiveGuess(s)
		if !accepted {
			fmt.Println("You have already made that guess, please enter a new one")
			time.Sleep(time.Second)
			continue
		}
		if mp != nil && isGuessInWord {
			verifyHangman(s, mp)
			s.CorGuesses++
			if s.CorGuesses == len(s.Word) {
				fmt.Println("Congratualations! You've guessed the correct word! ", s.Word)
				time.Sleep(time.Second * 2)
				s.IncGuesses = 100
			}
		}
	}

	if s.IncGuesses == 10 {
		fmt.Println("Sorry, you've run out of chances! Game Over!")
		time.Sleep(time.Second * 2)
	}
}

func setupGame() state.GameState {
	temp := state.GameState{
		IncGuesses: 0,
		Option:     'p',
		Word:       getWord(),
		Guesses:    make([]string, 1)}
	// if a quote, we'll split by the designated separator and then fill QuotedBy field of GameState
	if strings.Contains(temp.Word, "-") {
		tempWord := strings.Split(temp.Word, "-")
		temp.QuotedBy = tempWord[1]
		temp.Word = tempWord[0]
	}
	return temp
}

func getWord() string {
	wordBytes, err := os.ReadFile("./words.txt")
	if err != nil {
		color.Red("Error found during opening a file: %s", err.Error())
	}
	words := strings.Split(string(wordBytes), "\r\n")
	rand.Seed(time.Now().UnixNano())

	return words[rand.Intn(len(words)-1)]
}

func printStartMessage(word string) {
	fmt.Println("Hey, welcome to the Hangman game! You'll be prompted to enter guesses one at a time to guess a word or phrase.")
	if strings.Contains(word, "-") {
		fmt.Println("The phrase is a famous quote. You'll be told who said it to help guess it")
	}
	fmt.Printf("Here are your options\n1.(P)lay\n2.(Q)uit\n\nPlease select one by writing in the first letter of the option you'd like to choose.\n\n")
}

func verifyHangman(s *state.GameState, guessMap map[string]bool) {
	s.CorGuesses = 0
	for _, char := range s.Word {
		didGuess := guessMap[string(char)]
		if char == ' ' || char == '\'' || char == '.' || didGuess {
			s.CorGuesses++
		}
	}
}

func printHangman(s *state.GameState) (map[string]bool, error) {
	s.CorGuesses = 0
	guessMap := make(map[string]bool, 1)
	for _, elem := range s.Guesses {
		guessMap[elem] = true
	}

	fmt.Println(state.GetCurrentState(s.IncGuesses))

	for _, char := range s.Word {
		didGuess := guessMap[string(char)]
		if char == ' ' || char == '\'' || char == '.' || didGuess {
			fmt.Printf("%c ", char)
		} else {
			fmt.Printf("_ ")
		}
	}
	fmt.Printf("\n")
	// print quote
	if s.QuotedBy != "" {
		fmt.Printf("- %s\n", s.QuotedBy)
	}
	fmt.Println("Guesses: ", s.Guesses)
	return guessMap, nil
}

func receiveGuess(s *state.GameState) (bool, bool) {
	isGuessInWord := false
	fmt.Printf("Please enter your guess> ")
	reader := bufio.NewReader(os.Stdin)

	guess, sizeT, err := reader.ReadRune()
	fmt.Printf("\n")
	if err != nil {
		color.Red("Error while receiving guess, # of inputs: %d, Error: %s\n\n", sizeT, err.Error())
		panic(1)
	}

	// convert any UpperCase guesses into their lower case counterpoints

	if int(guess) >= 65 && int(guess) <= 91 {
		guess = rune(int(guess) + 32)
	}

	if alreadyGuessed(s.Guesses, string(guess)) {
		return false, isGuessInWord
	}

	s.Guesses = append(s.Guesses, string(guess))
	for _, char := range s.Word {
		if int(char) == int(guess) {
			isGuessInWord = true
		}
	}

	if !isGuessInWord {
		s.IncGuesses++
	}
	return true, isGuessInWord
}

func alreadyGuessed(sl []string, g string) bool {
	for _, e := range sl {
		if g == e {
			return true
		}
	}
	return false
}
