package quoridor

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func BuildQuoridorBoardFromString(board string) (*Game, error) {
	return BuildQuoridorBoard(bufio.NewReader(strings.NewReader(board)))
}

func BuildQuoridorBoard(reader *bufio.Reader) (*Game, error) {
	game := &Game{
		Board:   make(Board),
		Players: make(map[PlayerPosition]*Player),
	}
	var row = 0
	for {
		line, _, err := reader.ReadLine()
		if line == nil || err != nil {
			break
		}
		if len(line) != BoardSize {
			return nil, errors.New(fmt.Sprint("wrong number of columns in row ", row))
		}
		for col, char := range line {
			if char == '-' || char == '|' {
				position := Position{
					Y: row,
					X: col,
				}
				if char == '-' && (row%2 == 0) {
					return nil, errors.New(fmt.Sprintf("horizontal barrier at %v is invalid", position))
				}
				if char == '|' && (col%2 == 0) {
					return nil, errors.New(fmt.Sprintf("vertical barrier at %v is invalid", position))
				}
				game.Board[position] = Piece{Position: position}
			} else if unicode.IsDigit(rune(char)) {
				playerPos := PlayerPosition(char - '0')
				piecePos := Position{col, row}
				if !(col%2 == 0 && row%2 == 0) {
					return nil, errors.New(fmt.Sprintf("invalid pawn position for %d at %v", playerPos, piecePos))
				}
				player := &Player{
					Barriers:   5,
					Pawn:       Piece{Position: piecePos, Owner: playerPos},
					PlayerName: fmt.Sprint(" Player", playerPos),
				}
				game.Players[playerPos] = player
				game.Board[piecePos] = player.Pawn
			} else if char != '.' {
				panic(fmt.Sprintf("Unexpected character %b", char))
			}
		}
		row++
	}
	if row != BoardSize {
		return nil, errors.New("wrong number of rows")
	}
	return game, nil
}
