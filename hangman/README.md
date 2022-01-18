# Hangman game

## Intro

The idea is simple. Create a command line game that lets you play hangman.

1. Choose a random word
2. Show player how many letters and current status
3. Receive letter input
   1. if letter is in the word, then prompt until the word is shown
   2. if letter is not in word, increment incorrect counter and picture until picture is complete


## Motivation

Learning Go, Go modules, working with files, taking user input

## resources

under state/data you'll find a file for each state of the game. hangman0.txt is the beginning of the game with no failures and hangman9 is the end of the game and loss state.

words.txt is our dictionary of words. One thing we could easily do to improve the performance would be to save the dictionary in an array rather then open the file every new game. 

words_perf.txt was a test to see how much of an impact a larger file would have. There are 200k+ lines this file, but on my machine it still runs very fast. Maybe on a machine with less resources it would run much slower? I could test this on my chrome book (@TODO?)