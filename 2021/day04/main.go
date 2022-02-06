package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func parseDrawnNumbers(s string) []uint {
	drawnNumbersAsStrings := strings.Split(s, ",")
	drawnNumbers := make([]uint, len(drawnNumbersAsStrings))
	for i, s := range drawnNumbersAsStrings {
		number, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			panic("help")
		}
		drawnNumbers[i] = uint(number)
	}
	return drawnNumbers
}

type boardNumber struct {
	// value of the board number
	value uint
	// the row the number is on
	row uint
	// the column the number is on
	col uint
	// if the value has been marked
	marked bool
}

// collection of numbers on the bingoBoard
type bingoBoard [25]boardNumber

// Convert a sequence of 5 lines with whitespace-separated numbers into a board
// struct, recording the position of each number as we go.
func parseBoard(s string) bingoBoard {
	var b bingoBoard
	rowsAsStrings := strings.Split(s, "\n")
	numberPos := 0
	for i, rowAsString := range rowsAsStrings {
		rowNumbersAsStrings := strings.Fields(rowAsString)
		for j, rowNumberAsString := range rowNumbersAsStrings {
			number, err := strconv.ParseUint(rowNumberAsString, 10, 8)
			if err != nil {
				panic(err)
			}
			b[numberPos] = boardNumber{
				value: uint(number),
				row:   uint(i),
				col:   uint(j),
			}
			numberPos++
		}
	}
	return b
}

func bingoDraw(board bingoBoard, drawnNumbers []uint) (uint, uint) {
	rowHits := [5]uint{}
	colHits := [5]uint{}

	bingoDraw := uint(len(drawnNumbers))

drawloop:
	for draw, drawnNumber := range drawnNumbers {
		// check number on the board
		for i, number := range board {
			if number.value == drawnNumber {
				// cross number
				board[i].marked = true
				rowHits[number.row] += 1
				colHits[number.col] += 1
				// bingo check
				if rowHits[number.row] == 5 || colHits[number.col] == 5 {
					bingoDraw = uint(draw)
					break drawloop
				}
			}
		}
	}
	// compute score
	score := uint(0)
	for _, number := range board {
		if !number.marked {
			score += number.value
		}
	}
	score *= drawnNumbers[bingoDraw]
	return bingoDraw, score
}

func fastestBoardScore() uint {
	data, err := ioutil.ReadFile("data.txt")
	if err != nil {
		panic(err)
	}

	s := strings.Split(strings.TrimSpace(string(data)), "\n\n")

	drawnNumbers := parseDrawnNumbers(s[0])
	fmt.Println(drawnNumbers)

	fastestDraw := uint(len(drawnNumbers))
	fastestScore := uint(0)

	boardsAsText := s[1:]
	for i, boardAsText := range boardsAsText {
		board := parseBoard(boardAsText)
		boardDraw, boardScore := bingoDraw(board, drawnNumbers)
		if boardDraw < fastestDraw {
			fmt.Println("board", i, "draw", boardDraw, "score", boardScore)
			fastestDraw = boardDraw
			fastestScore = boardScore
		}
	}

	return fastestScore
}

func slowestBoardScore() uint {
	data, err := ioutil.ReadFile("data.txt")
	if err != nil {
		panic(err)
	}

	s := strings.Split(strings.TrimSpace(string(data)), "\n\n")

	drawnNumbers := parseDrawnNumbers(s[0])
	fmt.Println(drawnNumbers)

	slowestDraw := uint(0)
	slowestScore := uint(0)

	boardsAsText := s[1:]
	for i, boardAsText := range boardsAsText {
		board := parseBoard(boardAsText)
		boardDraw, boardScore := bingoDraw(board, drawnNumbers)
		if boardDraw > slowestDraw {
			fmt.Println("board", i, "draw", boardDraw, "score", boardScore)
			slowestDraw = boardDraw
			slowestScore = boardScore
		}
	}

	return slowestScore
}

func main() {
	winningScore := fastestBoardScore()
	fmt.Println("Score on fastest board: ", winningScore)

	losingScore := slowestBoardScore()
	fmt.Println("Score on slowest board: ", losingScore)
}
