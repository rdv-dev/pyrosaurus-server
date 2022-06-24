package util

import (
	"encoding/binary"
	"fmt"
)

type MovePoint struct {
	x int16
	y int16
	goalSize int16
}

type Moves struct {
	goalType int
	points []*MovePoint
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
						x: int16(binary.LittleEndian.Uint16(moves[i:i+2])),
						y: int16(binary.LittleEndian.Uint16(moves[i+2:i+4])),
						goalSize: int16(moves[i+4:i+5][0])}

					fmt.Printf("X:%d Y:%d S:%d\n", int16(binary.LittleEndian.Uint16(moves[i:i+2])), int16(binary.LittleEndian.Uint16(moves[i+2:i+4])),int16(moves[i+4:i+5][0]))

					i+=5
				}

				mv := &Moves {
					goalType: newGoal,
					points: m}

				retMoves = append(retMoves, mv)
			}

		} else {
			i++
		}
	}
	return retMoves
}