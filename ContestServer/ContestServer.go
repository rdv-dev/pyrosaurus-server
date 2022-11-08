package ContestServer
import (
	"errors"
	"encoding/binary"
	// "fmt"
	// "math"
	"github.com/algae-disco/pyrosaurus-server/ContestServer/util"
)

const (
	ACTIONS_PER_SECOND = 20
	NUM_SPECIES_LEN = 1
	TOTAL_DINOS_LEN = 1
	CONTEST_HEADER_RECORD_LEN = 17
	LEVEL_DATA_SIZE = 70
)

type ContestResult struct {
	Actions []byte
}

type ContestFrame struct {
	Actions []byte
	NumActions int
}

type Arena struct {
	dinos []*util.Dino
	numDinos int
}

type Action struct {
	code byte
	dino byte
	args []byte
}

type DecisionResult struct {
	Movement int
	Score int
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

type Vector struct {
	x float64
	y float64
	a float64
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

	contestHeader = append(contestHeader, byte(0))

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

func RunContest(team1, team2 *util.ContestEntry) (*ContestResult, error) {

	if team1.Team == team2.Team {
		return &ContestResult{}, errors.New("Team Pyro ID's cannot be the same")
	}

	cr := NewContestResult()

	arena := &Arena {
		dinos: make([]*util.Dino, team1.NumDinos + team2.NumDinos),
		numDinos: team1.NumDinos + team2.NumDinos}

	testTimeLimit := 60 * 5 // 5 minutes, 300 seconds, TODO based on level data

	// arenaFrames := testTimeLimit * ACTIONS_PER_SECOND
	arenaFrames := 50 + (0* testTimeLimit)

	// create dinos team 1
	speciesTypeOffset := ((util.TEAM_QUEEN_ARRAY_LEN + util.TEAM_SPECIES_LEG_NUM_LEN) * team1.NumDinos) + team1.DinosOffset + 1

	for i:=0; i<team1.NumDinos; i++ {
		arena.dinos[i] = util.NewDino(team1, int(team1.TeamData[speciesTypeOffset]), i)
		speciesTypeOffset += util.TEAM_MYSTERY_DATA
	}

	// create dinos team 2
	speciesTypeOffset = ((util.TEAM_QUEEN_ARRAY_LEN + util.TEAM_SPECIES_LEG_NUM_LEN) * team2.NumDinos) + team2.DinosOffset + 1

	for i:=team1.NumDinos; i<team1.NumDinos + team2.NumDinos; i++ {
		arena.dinos[i] = util.NewDino(team2, int(team2.TeamData[speciesTypeOffset]), (i-team1.NumDinos))
		speciesTypeOffset += util.TEAM_MYSTERY_DATA
	}

	// set up distance pairs
	// distPairs := make([]*Distance, 0)
	sense := make([]*DinoSense, 0)
	// set up delays
	delay := make([]*Delays, 0)
	pos := make([]*Vector, 0)
	velocity := make([]*Vector, 0)

	initFrame := ContestFrame {Actions: make([]byte, 0), NumActions: 0}

	for i:=0; i<team1.NumDinos + team2.NumDinos; i++ {
		for j:=i+1; j<team1.NumDinos + team2.NumDinos; j++ {
			// distPairs = append(distPairs, &Distance{d1: i, d2: j})
		}

		pos = append(pos, &Vector {
			x: arena.dinos[i].Xpos,
			y: arena.dinos[i].Ypos,
			a: arena.dinos[i].Angle})

		velocity = append(velocity, &Vector {x: 0, y: 0, a: 0})

		sense = append(sense, &DinoSense {
			see: make([]byte, team1.NumDinos + team2.NumDinos),
			hear: make([]byte, team1.NumDinos + team2.NumDinos),
			smell: make([]byte, team1.NumDinos + team2.NumDinos),
			enemy: make([]int, 0),
			friend: make([]int, 0),
			self: 0})

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
		initFrame.Put(&Action{code: 9, dino: byte(i), args: make([]byte, 0)})

		if i == 0 {
			// do a fire
			initFrame.Put(&Action{code: 11, dino: byte(i), args: []byte{8}})
			// delay the neck movement until fire has started
			delay[i].neck = 20
		}
	}

	cr.Actions = append(cr.Actions, byte(initFrame.NumActions))
	cr.Actions = append(cr.Actions, initFrame.Actions...)

	// distPairslen := len(distPairs)

	for arenaFrames > 0 {

		cf := ContestFrame { Actions: make([]byte, 0), NumActions: 0 }

		for i:=0; i<arena.numDinos; i++ {
			if delay[i].movement > 0 {
				delay[i].movement--
				// update position
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
				// update neck angle
			}

			if delay[i].tail > 0 {
				delay[i].tail--
				// update tail angle
			}
		}

		// for i:=0; i<distPairslen; i++ {
			// distance := math.Sqrt(math.Pow(arena.dinos[distPairs[i].d1].Xpos - arena.dinos[distPairs[i].d2].Xpos,2) + math.Pow(arena.dinos[distPairs[i].d1].Ypos - arena.dinos[distPairs[i].d2].Ypos,2))
			// sense smell
			// sense hearing
			// sense sight
			// fire range
			// friendly queen range
			// enemy queen range
		// }

		// neck/tail movement
		for i:=0; i<arena.numDinos; i++ {
			if delay[i].neck <= 0 {
				if i == 0 {
					shakeAngle := byte(10)

					if arenaFrames % 2 != 0 {
						shakeAngle = byte(255) - shakeAngle
					} 

					cf.Put(&Action{code: byte(0), dino: byte(i), args: []byte{0x4, shakeAngle}})

					delay[i].neck = 0x3

				} else {
					neckAngle := byte(30)
					
					if arenaFrames % 2 != 0 {
						neckAngle = byte(255) - neckAngle
					} 

					// Neck
					cf.Put(&Action{code: byte(0), dino: byte(i), args: []byte{0x11, neckAngle}})

					delay[i].neck = 0xF
				}
			}
			

			if delay[i].tail <= 0 {
				tailAngle := byte(30)

				if arenaFrames % 2 != 0 {
					tailAngle = byte(255) - tailAngle
				}

				// Tail
				cf.Put(&Action{code: 1, dino: byte(i), args: []byte{0x11, tailAngle}})

				
				delay[i].tail = 0xF
			}
		}

		// decisions
		for i:=0; i<arena.numDinos; i++ {

			// fighting?

			// evaluate decisions
			decisions := EvaluateDecision(arena.dinos[i])
			

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
					// cf.Put(&Action{code: 10, dino: byte(i), args: make([]byte, 0)})

					switch arena.dinos[i].Decisions[decisions[chosen].DecisionId].Priority {
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

				// if decisions[chosen].Movement == util.MOVEMENT_WANDER {}
				// if decisions[chosen].Movement == util.MOVEMENT_MOVE_AWAY {}
				// if decisions[chosen].Movement == util.MOVEMENT_MOVE_CLOSER {}
				// if decisions[chosen].Movement == util.MOVEMENT_MOVE_NORTH {}
				// if decisions[chosen].Movement == util.MOVEMENT_MOVE_SOUTH {}

				if decisions[chosen].Movement > util.MAX_PREDEFINED_MOVEMENT && delay[i].movement <= 0 {
					// look up movement
					mvId := decisions[chosen].Movement - util.CUSTOM_MOVEMENT_START

					if arena.dinos[i].DoMove != nil {
						// we are running a movement right now
						// check if we are at the point
						// if at point
						if arena.dinos[i].DoMove.ToPoint < len(arena.dinos[i].DoMove.Points) {
							arena.dinos[i].DoMove.ToPoint += 1
						} else {
							// last point
							arena.dinos[i].DoMove = nil
						}
					} else { 
						// we're not running a movement
						arena.dinos[i].DoMove = arena.dinos[i].Moves[mvId]
						arena.dinos[i].DoMove.ToPoint = 0
					}

					if arena.dinos[i].DoMove != nil {
						// calculate angle to point
						// x1 = x + cos(ang) * distance;
						// y1 = y + sin(ang) * distance;
						// or
						// x1 = x + sin(ang) * distance;
						// y1 = y + cos(ang) * distance;
						// execute movement
						// setup delay
					}
				}

				
				if decisions[chosen].Movement >= util.MOVEMENT_WANDER {
					// we did some kind of moving about

				}

			}
		}

		if cf.NumActions <= 0 {
			cf.Put(&Action{code: 11, dino: byte(0), args: []byte{10}})
		}

		cr.Actions = append(cr.Actions, byte(cf.NumActions))
		cr.Actions = append(cr.Actions, cf.Actions...)

		arenaFrames--;
	}

	// testDieFrame := ContestFrame {Actions: make([]byte, 0), NumActions: 0}
	// testDieFrame.Put(&Action{code: 11, dino: byte(0), args: []byte{0}})

	// cr.Actions = append(cr.Actions, byte(testDieFrame.NumActions))
	// cr.Actions = append(cr.Actions, testDieFrame.Actions...)

	endFrame := ContestFrame {Actions: make([]byte, 0), NumActions: 0}

	for i:=0; i<team1.NumDinos + team2.NumDinos; i++ {
		// set neck and tail to 0
		endFrame.Put(&Action{code: byte(0), dino: byte(i), args: []byte{0x11, 0}})
		endFrame.Put(&Action{code: byte(1), dino: byte(i), args: []byte{0x11, 0}})
		// turn off dino ?
		endFrame.Put(&Action{code: 11, dino: byte(i), args: []byte{9}})
	}

	cr.Actions = append(cr.Actions, byte(endFrame.NumActions))
	cr.Actions = append(cr.Actions, endFrame.Actions...)

	// cr.Actions = append(cr.Actions, make([]byte, 80)...)

	return cr, nil

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
			result = append(result, &DecisionResult{Movement: d.Decisions[i].Movement, Score: score, DecisionId: i})
		}
	}

	return result
}