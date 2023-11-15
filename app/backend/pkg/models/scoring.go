package models

// scoring.go contains the scoring logic for the game of Scramble.
import (
	"errors"
	"fmt"
	"sort"
)

// scoring calculates the score of a move in a game of Scramble.
// It takes an activeGame of type Game and a newTiles of type MoveSlice as input.
// It returns an integer score and an error if the move is invalid.
func (c *LanguageClient) scoring(activeGame Game, newTiles MoveSlice) (int, error) {
	fmt.Println("scoring: ", newTiles)

	score := 0
	setOfWords := []string{}
	scoreModifier := [15]string{}
	OriginalWord := ""
	scoreAggregateModifier := 1
	OGWordScore := 0

	sort.Sort(newTiles)

	for i := 0; i < len(newTiles); i++ {

		var x = int(newTiles[i].Col)
		var y = int(newTiles[i].Row)

		OriginalWord += string(newTiles[i].Letter)
		scoreModifier[i] = checkForScoreModifier(x, y)

		// recursively get all the possible words
		leftAndRightWord := pullLeft(activeGame, x, y) + activeGame.Board[x][y] + pullRight(activeGame, x, y)
		upAndDownWord := pullUp(activeGame, x, y) + activeGame.Board[x][y] + pullDown(activeGame, x, y)

		if (!c.CheckValidWord(leftAndRightWord) && len(leftAndRightWord) > 1) || (!c.CheckValidWord(upAndDownWord) && len(upAndDownWord) > 1) {
			return 0, errors.New("this is an invalid word")
		}

		if len(leftAndRightWord) < 2 && len(upAndDownWord) < 2 && len(newTiles) == 1 {
			return 0, errors.New("words must be longer than a single letter")
		}

		// then append to the list of words that would count towards the scores
		if !checkWordExists(setOfWords, leftAndRightWord) && len(leftAndRightWord) > 1 {
			setOfWords = append(setOfWords, leftAndRightWord)
		}

		if !checkWordExists(setOfWords, upAndDownWord) && len(upAndDownWord) > 1 {
			setOfWords = append(setOfWords, upAndDownWord)
		}
	}
	// fmt.Println(setOfWords)

	for _, word := range setOfWords {
		if c.CheckValidWord(word) {
			if OriginalWord == word {
				for i := 0; i < len(word); i++ {

					switch scoreModifier[i] {
					case "dl":
						OGWordScore += 2 * (c.GetLetterScore(string(word[i])))
					case "tl":
						OGWordScore += 3 * (c.GetLetterScore(string(word[i])))
					case "dw":
						scoreAggregateModifier *= 2
						OGWordScore += (c.GetLetterScore(string(word[i])))
					case "tw":
						scoreAggregateModifier *= 3
						OGWordScore += (c.GetLetterScore(string(word[i])))
					default:
						OGWordScore += (c.GetLetterScore(string(word[i])))
					}
				}
			} else {
				for _, letter := range word {
					score += c.GetLetterScore(string(letter))
				}
			}
		}
	}

	// fmt.Println("OGWordScore: ", OGWordScore)
	score += (scoreAggregateModifier * OGWordScore)
	fmt.Println("DONE: score: ", score)
	return score, nil
}

// pullUp recursively pulls up the letters above the given position (x,y) on the game board of the given game.
// If the position is already at the top of the board or the position is empty, it returns an empty string.
// Otherwise, it returns the concatenation of the letter at the current position and the result of calling pullUp on the position above it.
func pullUp(game Game, x int, y int) string {
	// fmt.Println("pullUp: ", x, y)

	if y <= 0 || game.Board[x][y] == "" {
		return ""
	}
	// fmt.Println(pullUp(game, x, y-1) + game.Board[x][y])
	return pullUp(game, x, y-1) + game.Board[x][y-1]
}

// pullDown recursively concatenates the letters below the given position (x,y) on the game board of the given game.
// If the position is at the bottom of the board or empty, it returns an empty string.
// Otherwise, it returns the letter at the given position concatenated with the result of calling pullDown on the position below it.
func pullDown(game Game, x int, y int) string {
	// fmt.Println("pullDown: ", x, y)

	if y <= 14 || game.Board[x][y] == "" {
		return ""
	}
	// fmt.Println(game.Board[x][y] + pullDown(game, x, y+1))
	return game.Board[x][y+1] + pullDown(game, x, y+1)
}

// pullLeft recursively pulls the letters to the left of the given position (x,y) on the game board
// and returns them as a string. If the position is already at the leftmost edge of the board or
// the position is empty, it returns an empty string.
func pullLeft(game Game, x int, y int) string {
	// fmt.Println("pullLeft: ", x, y)

	if x <= 0 || game.Board[x][y] == "" {
		return ""
	}

	// fmt.Println(pullLeft(game, x-1, y) + game.Board[x][y])
	return pullLeft(game, x-1, y) + game.Board[x-1][y]
}

// pullRight recursively pulls the letters to the right of the given position (x,y) on the game board
// and returns them as a string. If the position is at the right edge of the board or the cell is empty,
// it returns an empty string.
func pullRight(game Game, x int, y int) string {
	// fmt.Println("pullRight: ", x, y)

	if x >= 14 || game.Board[x][y] == "" {
		return ""
	}
	// fmt.Println(game.Board[x][y] + pullRight(game, x+1, y))
	return game.Board[x+1][y] + pullRight(game, x+1, y)
}

// checkWordExists checks if a given word exists in a set of words.
func checkWordExists(setOfWords []string, word string) bool {
	for _, eachWord := range setOfWords {
		if eachWord == word {
			return true
		}
	}
	return false
}

// checkForScoreModifier checks if a given position (x,y) on the game board is a score modifier (double word, double letter, triple letter, triple word).
func checkForScoreModifier(x int, y int) string {
	tw := [][2]int{{0, 0}, {0, 7}, {0, 14}, {7, 0}, {7, 14}, {14, 0}, {14, 7}, {14, 14}}
	dw := [][2]int{{1, 1}, {2, 2}, {3, 3}, {4, 4}, {10, 10}, {11, 11}, {12, 12}, {13, 13}, {1, 13}, {2, 12}, {3, 11}, {4, 10}, {10, 4}, {11, 3}, {12, 2}, {13, 1}}
	tl := [][2]int{{1, 5}, {1, 9}, {5, 1}, {5, 5}, {5, 9}, {5, 13}, {9, 1}, {9, 5}, {9, 9}, {9, 13}, {13, 5}, {13, 9}}
	dl := [][2]int{{0, 3}, {0, 11}, {2, 6}, {2, 8}, {3, 0}, {3, 7}, {3, 14}, {6, 2}, {6, 6}, {6, 8}, {6, 12}, {7, 3}, {7, 11}, {8, 2}, {8, 6}, {8, 8}, {8, 12}, {11, 0}, {11, 7}, {11, 14}, {12, 6}, {12, 8}, {14, 3}, {14, 11}}

	for i := 0; i < len(tw); i++ {
		tempVal := [2]int{x, y}
		if tempVal == tw[i] {
			return "tw"
		}
	}

	for i := 0; i < len(dw); i++ {
		tempVal := [2]int{x, y}
		if tempVal == dw[i] {
			return "dw"
		}
	}

	for i := 0; i < len(tl); i++ {
		tempVal := [2]int{x, y}
		if tempVal == tl[i] {
			return "tl"
		}
	}

	for i := 0; i < len(dl); i++ {
		tempVal := [2]int{x, y}
		if tempVal == dl[i] {
			return "dl"
		}
	}

	return "na"
}
