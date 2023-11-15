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

	adjacentToPlacedTile := adjacentToPlacedTile(activeGame, newTiles)
	if !adjacentToPlacedTile {
		return 0, errors.New("at least one new tile must be adjacent to an already placed tile")
	}

	isSequentialInX, isSequentialInY := checkSequential(newTiles, activeGame)

	// check for bisecting tile
	if !isSequentialInX && !isSequentialInY {
		var bisectingTile *Move
		for i := 0; i < len(newTiles); i++ {
			for row := 0; row < len(activeGame.Board); row++ {
				for col := 0; col < len(activeGame.Board[row]); col++ {
					if activeGame.Board[row][col] == string(newTiles[i].Letter) {
						bisectingTile = &newTiles[i]
						break
					}
				}
			}
		}
		if bisectingTile != nil {
			// check if new tiles are placed in a single row or column on either side of the bisecting tile
			isSequentialInX = true
			isSequentialInY = true
			for _, tile := range newTiles {
				if tile.Row != bisectingTile.Row {
					isSequentialInX = false
				}
				if tile.Col != bisectingTile.Col {
					isSequentialInY = false
				}
			}
			if !isSequentialInX && !isSequentialInY {
				return 0, errors.New("new tiles must be placed in a single row or column")
			}
		} else {
			return 0, errors.New("new tiles must be placed in a single row or column")
		}
	}

	// check if at least one new tile is adjacent to an already placed tile, unless the board is completely blank
	if len(activeGame.Board) > 0 {
		isAdjacent := false
		for row := 0; row < len(activeGame.Board); row++ {
			for col := 0; col < len(activeGame.Board[row]); col++ {
				if activeGame.Board[row][col] != "" {
					isNewTile := false
					for _, tile := range newTiles {
						if tile.Row == row && tile.Col == col {
							isNewTile = true
							break
						}
					}
					if !isNewTile {
						switch {
						case row > 0 && activeGame.Board[row-1][col] != "":
							isAdjacent = true
						case row < 14 && activeGame.Board[row+1][col] != "":
							isAdjacent = true
						case col > 0 && activeGame.Board[row][col-1] != "":
							isAdjacent = true
						case col < 14 && activeGame.Board[row][col+1] != "":
							isAdjacent = true
						default:
							isAdjacent = false
						}
					}
				}
				if isAdjacent {
					break
				}
			}
			if isAdjacent {
				break
			}
		}

		if !isAdjacent {
			return 0, errors.New("at least one new tile must be adjacent to an already placed tile")
		}
	}

	// check if all new tiles are in a line either horizontally or vertically
	for i := 0; i < len(newTiles); i++ {

		var x = int(newTiles[i].Col)
		var y = int(newTiles[i].Row)

		OriginalWord += string(newTiles[i].Letter)
		scoreModifier[i] = checkForScoreModifier(x, y)

		// recursively get all the possible words
		leftAndRightWord := pullLeft(activeGame, x, y) + activeGame.Board[x][y] + pullRight(activeGame, x, y)
		upAndDownWord := pullUp(activeGame, x, y) + activeGame.Board[x][y] + pullDown(activeGame, x, y)

		fmt.Println("leftAndRightWord: ", leftAndRightWord, " upAndDownWord: ", upAndDownWord)
		if (!c.CheckValidWord(leftAndRightWord) && len(leftAndRightWord) > 1) || (!c.CheckValidWord(upAndDownWord) && len(upAndDownWord) > 1) {
			return 0, errors.New("this is an invalid word: " + leftAndRightWord + " or " + upAndDownWord)
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

	fmt.Println("Valid words: ", setOfWords)

	// calculate the score
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

	// Add 50 points if all 7 tiles are used
	if newTiles.Len() == 7 {
		score += 50
	}

	// fmt.Println("OGWordScore: ", OGWordScore)
	score += (scoreAggregateModifier * OGWordScore)
	fmt.Println("DONE: score: ", score)
	return score, nil
}

func pullLeft(game Game, x int, y int) string {
	if y <= 0 || game.Board[x][y] == "" {
		return ""
	}
	return pullLeft(game, x, y-1) + game.Board[x][y-1]
}

func pullRight(game Game, x int, y int) string {
	if y >= 14 || game.Board[x][y] == "" {
		return ""
	}
	return game.Board[x][y+1] + pullRight(game, x, y+1)
}

func pullUp(game Game, x int, y int) string {
	// fmt.Println("pullLeft: ", x, y)

	if x <= 0 || game.Board[x][y] == "" {
		return ""
	}

	// fmt.Println(pullLeft(game, x-1, y) + game.Board[x][y])
	return pullUp(game, x-1, y) + game.Board[x-1][y]
}

func pullDown(game Game, x int, y int) string {
	// fmt.Println("pullRight: ", x, y)

	if x >= 14 || game.Board[x][y] == "" {
		return ""
	}
	// fmt.Println(game.Board[x][y] + pullRight(game, x+1, y))
	return game.Board[x+1][y] + pullDown(game, x+1, y)
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

// checkSequential verifies if all tiles are in a line either horizontally or vertically.
func checkSequential(tiles MoveSlice, game Game) (bool, bool) {
	isSequentialInX := true
	isSequentialInY := true

	// Check if sequential in X direction (same Row)
	for i := 1; i < len(tiles); i++ {
		if tiles[i].Row != tiles[0].Row {
			isSequentialInX = false
			break
		}
	}

	// If not sequential in X, check if sequential in Y direction (same Col)
	if !isSequentialInX {
		for i := 1; i < len(tiles); i++ {
			if tiles[i].Col != tiles[0].Col {
				isSequentialInY = false
				break
			}
		}
	} else {
		isSequentialInY = false // If they are sequential in X, they cannot be in Y
	}

	// If sequential in either direction, the adjacent elements must have consecutive indexes
	if isSequentialInX {
		sort.SliceStable(tiles, func(i, j int) bool {
			return tiles[i].Col < tiles[j].Col
		})

		for i := 1; i < len(tiles); i++ {
			if tiles[i].Col-tiles[i-1].Col != 1 {
				return false, false
			}
		}

	} else if isSequentialInY {
		sort.SliceStable(tiles, func(i, j int) bool {
			return tiles[i].Row < tiles[j].Row
		})

		for i := 1; i < len(tiles); i++ {
			if tiles[i].Row-tiles[i-1].Row != 1 {
				return false, false
			}
		}
	} else {
		return false, false
	}

	return true, false
}

// adjacentToPlacedTile checks if at least one new tile is adjacent to an already placed tile, unless the board is completely blank.
// It takes an activeGame of type Game and a newTiles of type MoveSlice as input.
// It returns a boolean value indicating whether at least one new tile is adjacent to an already placed tile.
func adjacentToPlacedTile(activeGame Game, newTiles MoveSlice) bool {
	if len(activeGame.Board) == 0 {
		// If the board is completely blank, then any new tile is adjacent to an already placed tile.
		return true
	}

	// Check if at least one new tile is adjacent to an already placed tile.
	for _, tile := range newTiles {
		row, col := tile.Row, tile.Col
		if row > 0 && activeGame.Board[row-1][col] != "" && !containsTile(newTiles, row-1, col) {
			return true
		}
		if row < 14 && activeGame.Board[row+1][col] != "" && !containsTile(newTiles, row+1, col) {
			return true
		}
		if col > 0 && activeGame.Board[row][col-1] != "" && !containsTile(newTiles, row, col-1) {
			return true
		}
		if col < 14 && activeGame.Board[row][col+1] != "" && !containsTile(newTiles, row, col+1) {
			return true
		}
	}

	return false
}

// containsTile checks if a tile with the given row and column exists in the given MoveSlice.
func containsTile(tiles MoveSlice, row, col int) bool {
	for _, tile := range tiles {
		if tile.Row == row && tile.Col == col {
			return true
		}
	}
	return false
}
