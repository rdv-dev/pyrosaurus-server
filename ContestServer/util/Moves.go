package util

import (
	"encoding/binary"
	"fmt"
)

const (
	MOVEMENT_CALL = 0
	MOVEMENT_DONT_MOVE = 1
	MOVEMENT_WANDER = 2
	MOVEMENT_MOVE_AWAY = 3
	MOVEMENT_MOVE_CLOSER = 4
	MOVEMENT_MOVE_NORTH = 5
	MOVEMENT_MOVE_SOUTH = 6
	
	MAX_PREDEFINED_MOVEMENT = 6
	CUSTOM_MOVEMENT_START = 7
)

type MovePoint struct {
	X int16
	Y int16
	GoalSize int16
}

type Moves struct {
	GoalType int
	Points []*MovePoint
	ToPoint int
}

func NewMoves(moves []byte) []*Moves {
	retMoves := make([]*Moves, 0)
	i := 0
	for i<len(moves) {
		if moves[i] != 0 {
			numPoints := int(moves[i])
			i++

			if numPoints <= 10 {

				fmt.Printf("Loading moves...\nNum Points: %d\n", numPoints)

				m := make([]*MovePoint, numPoints)

				newGoal := int(moves[i])
				i++

				for p:=0; p<numPoints; p++ {
					m[p] = &MovePoint {
						X: int16(binary.LittleEndian.Uint16(moves[i:i+2])),
						Y: int16(binary.LittleEndian.Uint16(moves[i+2:i+4])),
						GoalSize: int16(moves[i+4:i+5][0])}

					fmt.Printf("X:%d Y:%d S:%d\n", int16(binary.LittleEndian.Uint16(moves[i:i+2])), int16(binary.LittleEndian.Uint16(moves[i+2:i+4])),int16(moves[i+4:i+5][0]))

					i+=5
				}

				mv := &Moves {
					GoalType: newGoal,
					Points: m}

				retMoves = append(retMoves, mv)
			}

		} else {
			i++
		}
	}
	return retMoves
}