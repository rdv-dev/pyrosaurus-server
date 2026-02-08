package ContestServer
import (
	"errors"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
	"github.com/rdv-dev/pyrosaurus-server/ContestServer/util"
	"github.com/rdv-dev/pyrosaurus-server/Database"
)

const (
	ACTIONS_PER_SECOND = 20
	NUM_SPECIES_LEN = 1
	TOTAL_DINOS_LEN = 1
	CONTEST_HEADER_RECORD_LEN = 17
	LEVEL_DATA_SIZE = 70
	SPEED_CREEP = 0
	SPEED_WALK = 1
	SPEED_RUN = 2
	CONTEST_NOT_WATCHED = 0x00
	CONTEST_TEAM1_WON = CONTEST_NOT_WATCHED | 1
	CONTEST_TEAM2_WON = CONTEST_NOT_WATCHED | 2
	CONTEST_DRAW = CONTEST_NOT_WATCHED | 3
	MAX_DELAY = 128
)

type ContestResult struct {
	Actions []byte
}

type ContestFrame struct {
	Actions []byte
	NumActions int
}

type Arena struct {
	Dinos []*util.Dino
	NumDinos int
}

type Action struct {
	code byte
	dino byte
	args []byte
}

type DecisionResult struct {
	Movement int
	Score int
	Speed byte
	DecisionId int
}

type Distance struct {
	d1 int
	d2 int
}

type DinoSense struct {
	see []byte
	hear []byte
	smell []byte
	enemy []int
	friend []int
	self int
}

type Delays struct {
	movement int
	fire int
	call int
	neck int
	tail int
}

// Vector stores position and heading using int16 integer math.
// X, Y are world coordinates; A is heading in integer degrees [0, 360].
type Vector struct {
	X int16
	Y int16
	A int16
}

type DinoMovement struct {
	count int
	movementAnimation int
	speed int16
}

func NewContestResult() *ContestResult {
	return &ContestResult {
		Actions: make([]byte, 0)}
}

func (cf * ContestFrame) Put(action *Action) {

	var encoded byte

	encoded = (action.dino * byte(12)) + action.code

	cf.Actions = append(cf.Actions, encoded)
	cf.Actions = append(cf.Actions, action.args...)
	cf.NumActions += 1

}


func (cr *ContestResult) Push(action *Action) {
	frame := ContestFrame { Actions: make([]byte, 0), NumActions: 0 }

	frame.Put(action)

	cr.Actions = append(cr.Actions, byte(frame.NumActions))
	cr.Actions = append(cr.Actions, frame.Actions...)
}

func (cr *ContestResult) EndGame() {
	cr.Actions = append(cr.Actions, byte(0))
}

func (cr *ContestResult) GenerateDelay(count int) {
	cr.Actions = append(cr.Actions, byte(-count))
}

// maybe move this under util?
func ExportContest(team1, team2 *util.ContestEntry, levelData []byte, result *ContestResult) ([]byte, error) {
	output := make([]byte, 0)

	team1ColorsNamesOffset := len(output) + CONTEST_HEADER_RECORD_LEN

	userId1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(userId1, team1.PyroUserId)

	output = append(output, userId1...)

	output = append(output, team1.TeamData[team1.ColorsNamesOffset:]...)

	team1DataOffset := len(output) + CONTEST_HEADER_RECORD_LEN

	output = append(output, byte(team1.NumSpecies))

	output = append(output, team1.TeamData[team1.SpeciesOffset:(team1.SpeciesOffset + (util.SPECIES_LEN * team1.NumSpecies))]...)

	dinosOffsetEnd := team1.DinosOffset + (team1.NumDinos * (util.TEAM_ENTRY_RECORD_LEN)) + 1
	output = append(output, team1.TeamData[team1.DinosOffset:dinosOffsetEnd]...)

	output = append(output, team1.TeamData[team1.DinoNamesOffset:team1.ColorsNamesOffset]...)


	team2ColorsNamesOffset := len(output) + CONTEST_HEADER_RECORD_LEN


	userId2 := make([]byte, 4)

	binary.LittleEndian.PutUint32(userId2, team2.PyroUserId)

	output = append(output, userId2...)

	output = append(output, team2.TeamData[team2.ColorsNamesOffset:]...)

	team2DataOffset := len(output) + CONTEST_HEADER_RECORD_LEN

	output = append(output, byte(team2.NumSpecies))

	output = append(output, team2.TeamData[team2.SpeciesOffset:(team2.SpeciesOffset + (util.SPECIES_LEN * team2.NumSpecies))]...)


	dinosOffsetEnd = team2.DinosOffset + (team2.NumDinos * (util.TEAM_ENTRY_RECORD_LEN)) + 1
	output = append(output, team2.TeamData[team2.DinosOffset:dinosOffsetEnd]...)

	output = append(output, team2.TeamData[team2.DinoNamesOffset:team2.ColorsNamesOffset]...)

	levelDataOffset := len(output) + CONTEST_HEADER_RECORD_LEN

	output = append(output, levelData...)

	contestDataOffset := len(output) + CONTEST_HEADER_RECORD_LEN

	// contestHeader := make([]byte, CONTEST_HEADER_RECORD_LEN)
	contestHeader := make([]byte, 0)
	fielduInt16 := make([]byte, 2)

	contestHeader = append(contestHeader, byte(CONTEST_DRAW))

	binary.LittleEndian.PutUint16(fielduInt16, uint16(team1ColorsNamesOffset))
	contestHeader = append(contestHeader, fielduInt16...)

	binary.LittleEndian.PutUint16(fielduInt16, uint16(team1DataOffset))
	contestHeader = append(contestHeader, fielduInt16...)

	binary.LittleEndian.PutUint16(fielduInt16, uint16(0))
	contestHeader = append(contestHeader, fielduInt16...)

	binary.LittleEndian.PutUint16(fielduInt16, uint16(team2ColorsNamesOffset))
	contestHeader = append(contestHeader, fielduInt16...)

	binary.LittleEndian.PutUint16(fielduInt16, uint16(team2DataOffset))
	contestHeader = append(contestHeader, fielduInt16...)

	binary.LittleEndian.PutUint16(fielduInt16, uint16(0))
	contestHeader = append(contestHeader, fielduInt16...)

	binary.LittleEndian.PutUint16(fielduInt16, uint16(levelDataOffset))
	contestHeader = append(contestHeader, fielduInt16...)

	binary.LittleEndian.PutUint16(fielduInt16, uint16(contestDataOffset))
	contestHeader = append(contestHeader, fielduInt16...)

	output = append(contestHeader, output...)

	output = append(output, result.Actions...)

	return output, nil
}

func SaveContest(team1, team2 *util.ContestEntry, levelData []byte, result *ContestResult, team1PlayerId uint64, team2EntryId uint64) (error) {
	outdata, err := ExportContest(team1, team2, levelData, result)

	if err != nil {
		return err
	}

	Database.InsertContest(team1PlayerId, team2EntryId, outdata)

	return nil
}

func FindOpponent(currentPlayerId uint64) (uint64, *util.ContestEntry) {
	opponentEntryId, opponentEntry, err := Database.FindOpponentEntry(currentPlayerId)

	if err != nil {
		fmt.Println("Error finding opponent", err)
		return 0, nil
	}

	opponent, err2 := util.NewContestEntry(opponentEntry)
	if err2 != nil {
		fmt.Println("Error parsing opponent contest entry", err)
		return 0, nil
	}

	return opponentEntryId, opponent

}

func RunContest(team1, team2 *util.ContestEntry, leveldata []byte, testTime int) (*ContestResult, error) {

	if team1.Team == team2.Team {
		return &ContestResult{}, errors.New("Team Pyro ID's cannot be the same")
	}

	level := util.NewLevel(leveldata)

	cr := NewContestResult()

	s1 := rand.NewSource(time.Now().Unix())
	r1 := rand.New(s1)

	arena := &Arena {
		Dinos: make([]*util.Dino, team1.NumDinos + team2.NumDinos),
		NumDinos: team1.NumDinos + team2.NumDinos,
	}

	//testTimeLimit := 60 * 5 // 5 minutes, 300 seconds, TODO based on level data
	//testTimeLimit := 30 * ACTIONS_PER_SECOND

	arenaFrames := 0
	delayCounter := 0

	if testTime > 0 {
		arenaFrames = testTime * ACTIONS_PER_SECOND
	} else {
		arenaFrames = level.MaxTime * ACTIONS_PER_SECOND
	}
	fmt.Printf("level.MaxTime: %d\n", level.MaxTime)
	// create dinos team 1
	speciesTypeOffset := ((util.TEAM_QUEEN_ARRAY_LEN + util.TEAM_SPECIES_LEG_NUM_LEN) * team1.NumDinos) + team1.DinosOffset + 1

	for i:=0; i<team1.NumDinos; i++ {
		arena.Dinos[i] = util.NewDino(team1, int(team1.TeamData[speciesTypeOffset]), i, level.X, level.Y)
		arena.Dinos[i].LegType = int(team1.TeamData[team1.DinosOffset + 1 + team1.NumDinos + i])
		speciesTypeOffset += util.TEAM_MYSTERY_DATA
	}

	// create dinos team 2
	speciesTypeOffset = ((util.TEAM_QUEEN_ARRAY_LEN + util.TEAM_SPECIES_LEG_NUM_LEN) * team2.NumDinos) + team2.DinosOffset + 1

	for i:=team1.NumDinos; i<team1.NumDinos + team2.NumDinos; i++ {
		arena.Dinos[i] = util.NewDino(team2, int(team2.TeamData[speciesTypeOffset]), (i-team1.NumDinos), level.X, level.Y)
		arena.Dinos[i].LegType = int(team2.TeamData[team2.DinosOffset + 1 + team2.NumDinos + (i-team1.NumDinos)])
		speciesTypeOffset += util.TEAM_MYSTERY_DATA
	}

	// set up distance pairs
	// distPairs := make([]*Distance, 0)
	sense := make([]*DinoSense, 0)
	// set up delays
	delay := make([]*Delays, 0)
	pos := make([]*Vector, 0)
	newPos := make([]*Vector, 0)
	target := make([]*Vector, 0)
	velocity := make([]*Vector, 0)
	move := make([]*DinoMovement, 0)

	//initFrame := ContestFrame {Actions: make([]byte, 0), NumActions: 0}

	for i:=0; i<team1.NumDinos + team2.NumDinos; i++ {
		for j:=i+1; j<team1.NumDinos + team2.NumDinos; j++ {
			// distPairs = append(distPairs, &Distance{d1: i, d2: j})
		}

		// Away team (team 2) gets positions and heading rotated 180 degrees
		// so the two teams face each other. Matches contestReadTeamDinos behavior.
		if i >= team1.NumDinos {
			xPos := arena.Dinos[i].Xpos
			yPos := arena.Dinos[i].Ypos
			util.RotateByHeading(180, &xPos, &yPos)
			pos = append(pos, &Vector {
				X: xPos,
				Y: yPos,
				A: util.NormalizeAngle(arena.Dinos[i].Angle + 180),
			})
		} else {
			pos = append(pos, &Vector {
				X: arena.Dinos[i].Xpos,
				Y: arena.Dinos[i].Ypos,
				A: arena.Dinos[i].Angle,
			})
		}

		newPos = append(newPos, &Vector {
			X: 0,
			Y: 0,
			A: 0,
		})

		target = append(target, &Vector {
			X: 0,
			Y: 0,
			A: 0,
		})

		velocity = append(velocity, &Vector {
			X: 0,
			Y: 0,
			A: 0,
		})

		sense = append(sense, &DinoSense {
			see: make([]byte, team1.NumDinos + team2.NumDinos),
			hear: make([]byte, team1.NumDinos + team2.NumDinos),
			smell: make([]byte, team1.NumDinos + team2.NumDinos),
			enemy: make([]int, 0),
			friend: make([]int, 0),
			self: 0})

		move = append(move, &DinoMovement {count: 0, movementAnimation: 0, speed: 0,})

		for j:=0; j<team1.NumDinos + team2.NumDinos; j++ {
			// sense friend, enemy etc
			if j == i {
				sense[i].self = j
			} else {
				if i < team1.NumDinos && j < team1.NumDinos {
					// team 1 friend
					sense[i].friend = append(sense[i].friend, j)
				}
				if i < team1.NumDinos && j >= team1.NumDinos {
					// team 1 enemy
					sense[i].enemy = append(sense[i].enemy, j)
				}
				if i >= team1.NumDinos && j < team1.NumDinos {
					// team 2 enemy
					sense[i].enemy = append(sense[i].enemy, j)
				}
				if i >= team1.NumDinos && j >= team1.NumDinos {
					// team 2 friend
					sense[i].friend = append(sense[i].friend, j)
				}
			}
		}

		delay = append(delay, &Delays {
			movement: 0,
			fire: 0,
			call: 0,
			neck: 0,
			tail: 0})

		// turn on dino ?
		//initFrame.Put(&Action{code: 9, dino: byte(i), args: make([]byte, 0)})
		//initFrame.Put(&Action{code: 11, dino: byte(i), args: []byte{byte(9)}})

		// initFrame.Put(&Action{code: 2, dino: byte(i), args: []byte{byte(5), byte(0x81), byte(0)}})
		// delay[i].movement = 100
		// initFrame.Put(&Action{code: 6, dino: byte(i), args: []byte{byte(0), byte(0x64 | 0x10)}})

		// if i == 0 {
		// do a fire
		// initFrame.Put(&Action{code: 11, dino: byte(i), args: []byte{8}})
		// delay the neck movement until fire has started
		// delay[i].neck = 20
		// }
	}

	//cr.Actions = append(cr.Actions, byte(initFrame.NumActions))
	//cr.Actions = append(cr.Actions, initFrame.Actions...)

	// distPairslen := len(distPairs)

	fmt.Println("Begin arena")

	//neckLocked := 0

	gameStruct := 1

	// Suppress unused variable warnings
	_ = target
	_ = velocity

	// cf := ContestFrame { Actions: make([]byte, 0), NumActions: 0 }

	for arenaFrames > 0 {

		cf := ContestFrame { Actions: make([]byte, 0), NumActions: 0 }

		for i:=0; i<arena.NumDinos; i++ {
			if delay[i].movement > 0 {
				delay[i].movement--
				// update position
				if arena.Dinos[i].DoMove != nil {
					pi := arena.Dinos[i].DoMove.ToPoint
					if pi >= len(arena.Dinos[i].DoMove.Points) {
						pi = len(arena.Dinos[i].DoMove.Points) - 1
					}
					boundVector := CheckBoundsV(pos[i], level)
					newAngle, stepRotation, canMove := CalculatePosition(pos[i], boundVector, arena.Dinos[i].DoMove.Points[pi])
					arena.Dinos[i].Rotate = stepRotation

					if canMove {
						// Compute displacement using RotateByHeading
						// Forward motion: dx=speed, dy=0 in local frame
						dx := move[i].speed
						dy := int16(0)
						util.RotateByHeading(pos[i].A, &dx, &dy)
						pos[i].X = pos[i].X + dx
						pos[i].Y = pos[i].Y + dy
						newPos[i].X = dx
						newPos[i].Y = dy
						newPos[i].A = newAngle
					} else {
						newPos[i].X = 0
						newPos[i].Y = 0
						newPos[i].A = 0
					}

					// if we finished this movement, then we update our angle when we rotate
					if delay[i].movement == 0 {
						pos[i].A = newAngle
					} else {
						// set rotation to 0 while we're moving
						arena.Dinos[i].Rotate = 0
					}

				}
			}

			if delay[i].fire > 0 {
				delay[i].fire--
				// calcuate damage
			}

			if delay[i].call > 0 {
				delay[i].call--
			}

			if delay[i].neck > 0 {
				delay[i].neck--
			}

			if delay[i].tail > 0 {
				delay[i].tail--
				// update tail angle
			}
		}

		// for i:=0; i<distPairslen; i++ {
		// distance := math.Sqrt(math.Pow(arena.Dinos[distPairs[i].d1].Xpos - arena.dinos[distPairs[i].d2].Xpos,2) + math.Pow(arena.dinos[distPairs[i].d1].Ypos - arena.dinos[distPairs[i].d2].Ypos,2))
		// sense smell
		// sense hearing
		// sense sight
		// fire range
		// friendly queen range
		// enemy queen range
		// }

		// neck/tail movement
		//for i:=0; i<arena.NumDinos; i++ {
		//	if delay[i].neck <= 0 {
		//		if neckLocked == 0 {
		//			neckAngle := byte(30)

		//			if arenaFrames % 2 != 0 {
		//				neckAngle = byte(255) - neckAngle
		//			}

		//			// Neck
		//			cf.Put(&Action{code: byte(0), dino: byte(i), args: []byte{0x11, neckAngle}})

		//			delay[i].neck = 0xF
		//		} else {
		//			cf.Put(&Action{code: byte(0), dino: byte(i), args: []byte{0x05, byte(30)}})
		//		}
		//	}

		//	if delay[i].tail <= 0 {
		//		tailAngle := byte(30)

		//		if arenaFrames % 2 != 0 {
		//			tailAngle = byte(255) - tailAngle
		//		}

		//		// Tail
		//		cf.Put(&Action{code: 1, dino: byte(i), args: []byte{0x11, tailAngle}})


		//		delay[i].tail = 0xF
		//	}
		//}

		// decisions
		for i:=0; i<arena.NumDinos; i++ {
			var rotation int
			var arg2 int

			// fighting?

			// cf.Put(&Action{code: 11, dino: byte(i), args: []byte{byte(gameStruct)}})

			// evaluate decisions
			decisions := EvaluateDecision(arena.Dinos[i])

			//fmt.Printf("Dino %d decided on movement %d\n", i, decisions[0].Movement)
			if len(decisions) > 0 {
				chosen := 0
				maxScore := decisions[0].Score

				for j:=1; j<len(decisions); j++ {
					if decisions[j].Score > maxScore {
						chosen = j
					}
				}

				if decisions[chosen].Movement == util.MOVEMENT_CALL && delay[i].call <= 0 {
					// call
					cf.Put(&Action{code: 10, dino: byte(i), args: make([]byte, 0)})

					switch arena.Dinos[i].Decisions[decisions[chosen].DecisionId].Priority {
					case 0:
						delay[i].call = 16 * ACTIONS_PER_SECOND
					case 1:
						delay[i].call = 10 * ACTIONS_PER_SECOND
					case 2:
						delay[i].call = 8 * ACTIONS_PER_SECOND
					case 3:
						delay[i].call = 4 * ACTIONS_PER_SECOND
					case 4:
						delay[i].call = 2 * ACTIONS_PER_SECOND
					}
				}

				if decisions[chosen].Movement == util.MOVEMENT_DONT_MOVE {
					// do nothing?
				}

				if decisions[chosen].Movement == util.MOVEMENT_WANDER && delay[i].movement <= 0 {

					rotation = r1.Intn(30)
					rotation = rotation - 15

					if rotation < 0 {
						arena.Dinos[i].Rotate = byte(0xFF + rotation)
					} else {
						arena.Dinos[i].Rotate = byte(rotation)
					}
				}
				// if decisions[chosen].Movement == util.MOVEMENT_MOVE_AWAY {}
				// if decisions[chosen].Movement == util.MOVEMENT_MOVE_CLOSER {}
				// if decisions[chosen].Movement == util.MOVEMENT_MOVE_NORTH {}
				// if decisions[chosen].Movement == util.MOVEMENT_MOVE_SOUTH {}

				if decisions[chosen].Movement > util.MAX_PREDEFINED_MOVEMENT && delay[i].movement <= 0 {
					// look up movement
					mvId := decisions[chosen].Movement - util.CUSTOM_MOVEMENT_START

					if arena.Dinos[i].DoMove != nil {
						// we are running a movement right now
						// check if we are at the point exept if it is a mobile point
						// if at point
						pi := arena.Dinos[i].DoMove.ToPoint
						if arena.Dinos[i].DoMove.ToPoint < len(arena.Dinos[i].DoMove.Points) && (arena.Dinos[i].DoMove.GoalType != 1 && arena.Dinos[i].DoMove.GoalType != 11) {
							// Check the dino position x y point against current movement x y point
							// Use squared distance to avoid float64 sqrt
							dx := int32(arena.Dinos[i].DoMove.Points[pi].X) - int32(pos[i].X)
							dy := int32(arena.Dinos[i].DoMove.Points[pi].Y) - int32(pos[i].Y)
							distSq := dx*dx + dy*dy
							if arena.Dinos[i].DoMove.Points[pi].GoalSize == 0 {
								if distSq < 10000 { // 100*100
									arena.Dinos[i].DoMove.ToPoint += 1
									fmt.Printf("Dino %d next point %d\n", i, pi)
								}
							} else {
								if distSq < 10000 { // 100*100
									arena.Dinos[i].DoMove.ToPoint += 1
									fmt.Printf("Dino %d next point %d\n", i, pi)
								}
							}
							if newPos[i].X == 0 && newPos[i].Y == 0 && newPos[i].A == 0 {
								// may need to rethink how this works
								decisions[chosen].Speed = util.DECISION_DONT_MOVE
								arena.Dinos[i].DoMove = nil
							}
						} else {
							// last point
							arena.Dinos[i].DoMove = nil
						}
					} else {
						// we're not running a movement
						arena.Dinos[i].DoMove = arena.Dinos[i].Moves[mvId]
						arena.Dinos[i].DoMove.ToPoint = 0
					}

				}


				if decisions[chosen].Movement >= util.MOVEMENT_WANDER && delay[i].movement <= 0 {
					// we did some kind of moving about, post processing

					
					if arena.Dinos[i].IsRunningSpeed {
						arg2 = 0x90
					} else {
						arg2 = 0x10
					}

					goType := decisions[chosen].Speed

					switch goType {
					case util.DECISION_DONT_MOVE:
						arg2 = arg2 | 0x01
						move[i].movementAnimation = 0
						move[i].speed = 0
						move[i].count = 0
						delay[i].movement = 8

					case util.DECISION_CREEP:
						arg2 = arg2 | 0x01
						switch move[i].count {
						case 0:
							move[i].movementAnimation = 0
							move[i].speed = 5
							delay[i].movement = 8
						case 1:
							move[i].movementAnimation = 1
							move[i].speed = 5
							delay[i].movement = 8
						case 2:
							move[i].movementAnimation = 2
							move[i].speed = 5
							delay[i].movement = 8
						default:
							delay[i].movement = 8
						}

						if move[i].count <= 2 {
							move[i].count++
						}
					case util.DECISION_WALK:
						arg2 = arg2 | 0x0A
						switch move[i].count {
						case 0:
							arena.Dinos[i].IsRunningSpeed = false
							move[i].movementAnimation = 0
							move[i].speed = 1
							delay[i].movement = 16
						case 1:
							move[i].movementAnimation = 1
							move[i].speed = 2
							delay[i].movement = 16
						case 2:
							move[i].movementAnimation = 2
							move[i].speed = 3
							delay[i].movement = 16
						default:
							delay[i].movement = 16
						}

						if move[i].count <= 5 {
							move[i].count++
						}
					case util.DECISION_RUN:
						arena.Dinos[i].IsRunningSpeed = true

						arg2 = arg2 | 0x04
						switch move[i].count {
						case 0:
							move[i].movementAnimation = 0
							move[i].speed = 1
							delay[i].movement = 8
						case 1:
							move[i].movementAnimation = 9
							move[i].speed = 2
							delay[i].movement = 16
						case 2:
							move[i].movementAnimation = 0xA
							move[i].speed = 3
							delay[i].movement = 16
						case 3:
							move[i].movementAnimation = 0xA
							move[i].speed = 4
							delay[i].movement = 16
						case 4:
							move[i].movementAnimation = 0xB
							move[i].speed = 5
							delay[i].movement = 16
						default:
							delay[i].movement = 16
						}

						if move[i].count <= 4 {
							move[i].count++
						}

					}

					//fmt.Printf("Dino %d rot: %d move: %x\n", i, rotation, move[i].movementAnimation)

					cf.Put(&Action{code: 2, dino: byte(i), args: []byte{byte(arena.Dinos[i].Rotate), byte(arg2), byte(move[i].movementAnimation)}})
				}

			}
		}

		if gameStruct == 15 {
			gameStruct = 1
		} else {
			gameStruct++
		}

		if cf.NumActions <= 0 {
			if delayCounter < MAX_DELAY {
				delayCounter = delayCounter + 1
			} else {
				cr.GenerateDelay(min(delayCounter, MAX_DELAY))
				delayCounter = 0
			}
		} else {
			if delayCounter > 0 {
				cr.GenerateDelay(min(delayCounter, MAX_DELAY))
				delayCounter = 0
			}

			cr.Actions = append(cr.Actions, byte(cf.NumActions))
			cr.Actions = append(cr.Actions, cf.Actions...)
		}

		arenaFrames--;
	}

	// testDieFrame := ContestFrame {Actions: make([]byte, 0), NumActions: 0}
	// testDieFrame.Put(&Action{code: 11, dino: byte(0), args: []byte{0}})

	// cr.Actions = append(cr.Actions, byte(testDieFrame.NumActions))
	// cr.Actions = append(cr.Actions, testDieFrame.Actions...)

	// endFrame := ContestFrame {Actions: make([]byte, 0), NumActions: 0}

	// for i:=0; i<team1.NumDinos + team2.NumDinos; i++ {
	//	 // set neck and tail to 0
	//	 endFrame.Put(&Action{code: byte(0), dino: byte(i), args: []byte{0x11, 0}})
	//	 endFrame.Put(&Action{code: byte(1), dino: byte(i), args: []byte{0x11, 0}})
	//	 // turn off dino ?
	//	 endFrame.Put(&Action{code: 11, dino: byte(i), args: []byte{9}})
	// }

	// cr.Actions = append(cr.Actions, byte(endFrame.NumActions))
	// cr.Actions = append(cr.Actions, endFrame.Actions...)

	// cr.Actions = append(cr.Actions, make([]byte, 80)...)

	// Flush any remaining accumulated delay before the terminator
	if delayCounter > 0 {
		cr.GenerateDelay(min(delayCounter, MAX_DELAY))
	}

	cr.EndGame()

	return cr, nil

}

// CalculatePosition computes the new heading angle and rotation byte for a dino
// moving toward a target point. Uses integer math throughout.
// Returns: new heading angle, contest-file rotation byte, and whether the dino can move.
func CalculatePosition(sourcePoint *Vector, boundVector *Vector, targetPoint *util.MovePoint) (int16, byte, bool) {
	// If out of bounds, signal no movement
	if boundVector.X != 0 || boundVector.Y != 0 {
		return sourcePoint.A, 0, false
	}

	// Compute target angle in integer degrees using Atan2
	dx := targetPoint.X - sourcePoint.X
	dy := targetPoint.Y - sourcePoint.Y
	targetAngle := util.Atan2Degrees(dy, dx)

	currentAngle := sourcePoint.A

	// Compute shortest angle difference in range [-180, 180]
	angleDiff := int16(targetAngle - currentAngle)
	for angleDiff > 180 {
		angleDiff -= 360
	}
	for angleDiff < -180 {
		angleDiff += 360
	}

	absDiff := util.AbsInt16(angleDiff)

	// Determine rotation step size (matching original thresholds)
	var rotStep int16
	if absDiff <= 1 {
		rotStep = 0
	} else if absDiff < 5 {
		rotStep = 1
	} else if absDiff > 50 {
		rotStep = 20
	} else {
		rotStep = 5
	}
	if rotStep > absDiff {
		rotStep = absDiff
	}

	// Apply rotation toward target
	var newAngle int16
	var stepRotation byte

	if angleDiff > 0 {
		// Turn counter-clockwise
		newAngle = util.NormalizeAngle(currentAngle + rotStep)
		stepRotation = byte(255) - byte(int(rotStep)*256/360)
	} else if angleDiff < 0 {
		// Turn clockwise
		newAngle = util.NormalizeAngle(currentAngle - rotStep)
		stepRotation = byte(int(rotStep) * 256 / 360)
	} else {
		// Already aligned
		newAngle = currentAngle
		stepRotation = 0
	}

	return newAngle, stepRotation, true
}

// CheckBoundsV checks if a dino is near the arena boundary.
// Returns a Vector with non-zero X/Y components indicating which bounds are violated.
func CheckBoundsV(sourcePoint *Vector, level *util.Level) *Vector {
	limit := int16(367)

	var newX, newY int16

	if sourcePoint.X > level.X - limit {
		newX = -1
	}

	if sourcePoint.Y > level.Y - limit {
		newY = -1
	}

	if sourcePoint.X < -level.X + limit {
		newX = 1
	}

	if sourcePoint.Y < -level.Y + limit {
		newY = 1
	}

	return &Vector {
		X: newX,
		Y: newY,
		A: 0,
	}
}

// DoMoveStateMachine mirrors the client's doMove state machine.
// Given the current movementAnimation state and movement parameters,
// returns the new movementAnimation state that the client will transition to.
// This is the interface spec between server simulation and client animation.
func DoMoveStateMachine(state, mode, dx, cx, legType int, isRunningSpeed bool) int {
	// Snake (legType 2) can never use running speed table
	if legType == util.LEG_TYPE_NONE {
		isRunningSpeed = false
	}

	if isRunningSpeed {
		return doMoveRunningTable(state, mode, dx, cx, legType)
	}
	return doMoveNormalTable(state, mode, dx, cx, legType)
}

// doMoveNormalTable implements the normal-speed doMove switch (19 states).
func doMoveNormalTable(state, mode, dx, cx, legType int) int {
	switch state {
	case 0:
		// Stop override
		if mode == 0 {
			return 18
		}
		// Creep with small dx override
		if mode == 1 && dx < 3 {
			return 3
		}
		if dx > 0 {
			return 4
		}
		if dx == 0 {
			return 3
		}
		return state
	case 1, 7:
		if dx > 0 {
			return 2
		}
		return state
	case 2:
		if dx == 0 {
			return 1
		}
		return state
	case 3:
		if dx > 0 {
			return 4
		}
		return state
	case 4:
		if dx == 0 {
			return 3
		}
		return state
	case 5:
		if dx == 0 {
			return 6
		}
		return state
	case 6:
		if dx > 0 {
			return 5
		}
		return state
	case 8, 16, 17, 18:
		if dx <= 0 {
			return 1
		}
		return 2
	case 9:
		if dx != 0 {
			return 4
		}
		return 3
	case 10, 11, 14, 15:
		return 15
	case 12, 13:
		return 13
	default:
		return state
	}
}

// doMoveRunningTable implements the running-speed doMove switch (19 states).
func doMoveRunningTable(state, mode, dx, cx, legType int) int {
	switch state {
	case 0:
		if mode == 0 {
			return state
		}
		if legType == util.LEG_TYPE_TWO_SPRAWL {
			return 11
		}
		return 9
	case 1, 2, 7, 16, 17, 18:
		return 8
	case 3, 4, 5, 6:
		return 9
	case 8, 9, 10, 11:
		return state
	case 12, 13:
		return 10
	case 14, 15:
		return 11
	default:
		return state
	}
}

func EvaluateDecision(d *util.Dino) []*DecisionResult {
	var execute bool
	var score int

	result := make([]*DecisionResult, 0)
	execute = true

	for i:=0; i<len(d.Decisions); i++ {

		score = 0

		if d.Decisions[i].Target == byte(0x04) { execute = true }
		if d.Decisions[i].Legs == byte(0x30) && execute == true { execute = true }
		if d.Decisions[i].Size == byte(0x05) && execute == true { execute = true }
		if d.Decisions[i].InRange == byte(0x40) && execute == true { execute = true }
		if d.Decisions[i].TheirSkin == byte(0x05) && execute == true { execute = true }
		if d.Decisions[i].MySkin == byte(0x50) && execute == true { execute = true }
		if d.Decisions[i].MyCondition == byte(0x40) && execute == true { execute = true }
		if d.Decisions[i].MyQueenEnemyRange == byte(0x30) && execute == true { execute = true }
		if d.Decisions[i].MyQueenRange == byte(0x30) && execute == true { execute = true }
		if d.Decisions[i].EnemyQueenRange == byte(0x40) && execute == true { execute = true }
		if d.Decisions[i].TheirSpeed == byte(0x03) && execute == true { execute = true }
		if d.Decisions[i].TheirAction == byte(0x50) && execute == true { execute = true }
		if d.Decisions[i].Calling == byte(0x02) && execute == true { execute = true }
		if d.Decisions[i].Time == byte(0x30) && execute == true { execute = true }

		if execute == true {
			result = append(result, &DecisionResult{
				Movement: d.Decisions[i].Movement,
				Score: score,
				Speed: d.Decisions[i].GoSpeed,
				DecisionId: i,
			})
		}
	}

	return result
}

func RandomMovement(r *rand.Rand, dinoNum int) (*Action) {
	var movement int
	var dx int
	// moveIndex := r.Intn(4)

	movement = 1

	// if moveIndex == 0 { movement = 10 }

	// if moveIndex == 1 { movement = 0 }

	// if moveIndex == 2 { movement = 4 }

	// if moveIndex == 3 { movement = 1 }

	rotation := r.Intn(80)

	rotation = rotation - 40

	dx_prep := r.Intn(10)

	if dx_prep <= 7 {
		dx = 8
	} else {
		dx = 8
	}

	var arg1 byte
	var arg2 byte

	if rotation < 0 {
		arg1 = byte(0xFF + rotation)
	} else {
		arg1 = byte(rotation)
	}

	arg2 = byte(dx << 7)
	arg2 = 0x80 | byte(movement)

	fmt.Printf("Dino %d rot: %d move: %x\n", dinoNum, rotation, arg2)

	return &Action{code: 2, dino: byte(dinoNum), args: []byte{arg1, arg2, byte(0x0)}}

}
