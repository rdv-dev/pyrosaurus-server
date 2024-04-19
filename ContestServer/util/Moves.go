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
	X float64
	Y float64
	GoalSize float64
}

type Moves struct {
	GoalType int
	Points []*MovePoint
	ToPoint int
}

func (p * Moves) PrintPoints() {
    for i:=0; i<len(p.Points); i++ {
        fmt.Printf("X: %f Y: %f GoalSize: %f\n", p.Points[i].X, p.Points[i].Y, p.Points[i].GoalSize)
    }
}

func NewMoves(moves []byte, levelX, levelY float64) []*Moves {
	retMoves := make([]*Moves, 0)
	i := 0
	for i<len(moves) {
		if moves[i] != 0 {
			numPoints := int(moves[i])
			i++

			if numPoints <= 10 {

				//fmt.Printf("Loading moves...\nNum Points: %d\n", numPoints)

				m := make([]*MovePoint, numPoints)

				newGoal := int(moves[i])
				i++

				for p:=0; p<numPoints; p++ {
                    var x, y, gs float64
                    x = float64(int16(binary.LittleEndian.Uint16(moves[i:i+2])))
                    y = float64(int16(binary.LittleEndian.Uint16(moves[i+2:i+4])))
                    gs = float64(int16(moves[i+4:i+5][0]))

                    if x < 0 {
                        if x < (float64(levelX) * -1) {
                            x = float64(levelX) + 100 + gs
                        }
                    } else {
                        if x > float64(levelX) {
                            x = float64(levelX) - 100 - gs
                        }
                    }
                    if y < 0 {
                        if y < (float64(levelY) * -1) {
                            y = float64(levelY) + 100 + gs
                        }
                    } else {
                        if y > float64(levelY) {
                            y = float64(levelY) - 100 - gs
                        }
                    }

					m[p] = &MovePoint {
						X: x,
						Y: y,
						GoalSize: float64(int16(moves[i+4:i+5][0])),
                    }

					//fmt.Printf("X:%f Y:%f S:%f\n", x, y, gs)

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
