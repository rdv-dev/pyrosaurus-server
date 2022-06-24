package ContestServer
import (
	"errors"
	"encoding/binary"
	// "fmt"
	// "math"
	"github.com/algae-disco/pyrosaurus-server/ContestServer/util"
	// "github.com/algae-disco/pyrosaurus-server/ContestServer"
	// "github.com/algae-disco/pyrosaurus-server/ContestServer/ContestEntry"
	// "github.com/algae-disco/pyrosaurus-server/ContestServer/Dino"
)

const (
	ACTIONS_PER_SECOND = 20
	NUM_SPECIES_LEN = 1
	TOTAL_DINOS_LEN = 1
	CONTEST_HEADER_RECORD_LEN = 16
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
}

type Delays struct {
	movement int
	fire int
	call int
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

func ExportContest(team1, team2 *util.ContestEntry, result *ContestResult) ([]byte, error) {
	output := make([]byte, 0)


	levelData := []byte {
	0xB8, 0x0B, 0x10, 0x27, 0x03, 0x00, 0x88, 0x13, 0x64, 0x00, 0x2C, 0x01,
	0x0A, 0x00, 0x0E, 0x00, 0x21, 0x00, 0xD5, 0x00, 0x15, 0x02, 0x1B, 0x00,
	0x64, 0x00, 0x0D, 0x00, 0x0A, 0x00, 0x0A, 0x00, 0x11, 0x00, 0x25, 0x00,
	0x0E, 0x00, 0x15, 0x00, 0x1E, 0x00, 0x14, 0x00, 0x00, 0x00, 0x48, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x68, 0x01, 0x01, 0x00,
	0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x50, 0x00};

	// output = append(output, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})

	team1ColorsNamesOffset := len(output) + 17

	userId1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(userId1, team1.PyroUserId)

	// if sizeput < 4 {
	// 	return make([]byte, 0), errors.New("Failed to write PyroUserId for Team1")
	// }

	output = append(output, userId1...)

	output = append(output, team1.TeamData[team1.ColorsNamesOffset:]...)

	team1DataOffset := len(output) + 17

	// sizeput = binary.PutUvarint(output, uint64(team1.NumSpecies))
	output = append(output, byte(team1.NumSpecies))

	output = append(output, team1.TeamData[team1.SpeciesOffset:(team1.SpeciesOffset + (util.SPECIES_LEN * team1.NumSpecies))]...)

	dinosOffsetEnd := team1.DinosOffset + (team1.NumDinos * (util.TEAM_ENTRY_RECORD_LEN)) + 1
	output = append(output, team1.TeamData[team1.DinosOffset:dinosOffsetEnd]...)

	output = append(output, team1.TeamData[team1.DinoNamesOffset:team1.ColorsNamesOffset]...)


	team2ColorsNamesOffset := len(output) + 17


	userId2 := make([]byte, 4)
	// sizeput = binary.PutUvarint(output, uint64(team2.PyroUserId))
	binary.LittleEndian.PutUint32(userId2, team2.PyroUserId)

	// if sizeput < 4 {
	// 	return make([]byte, 0), errors.New("Failed to write PyroUserId for Team1")
	// }

	output = append(output, userId2...)

	output = append(output, team2.TeamData[team2.ColorsNamesOffset:]...)

	team2DataOffset := len(output) + 17

	// sizeput = binary.PutUvarint(output, uint64(team2.NumSpecies))
	output = append(output, byte(team2.NumSpecies))

	output = append(output, team2.TeamData[team2.SpeciesOffset:(team2.SpeciesOffset + (util.SPECIES_LEN * team2.NumSpecies))]...)


	dinosOffsetEnd = team2.DinosOffset + (team2.NumDinos * (util.TEAM_ENTRY_RECORD_LEN)) + 1
	output = append(output, team2.TeamData[team2.DinosOffset:dinosOffsetEnd]...)

	output = append(output, team2.TeamData[team2.DinoNamesOffset:team2.ColorsNamesOffset]...)

	levelDataOffset := len(output) + 17

	output = append(output, levelData...)

	contestDataOffset := len(output) + 17

	contestHeader := make([]byte, 17)
	fielduInt16 := make([]byte, 2)

	// contestHeader = append(contestHeader, byte(0))
	contestHeader[0] = byte(0)
	binary.LittleEndian.PutUint16(fielduInt16, uint16(team1ColorsNamesOffset))
	contestHeader[1] = fielduInt16[0]
	contestHeader[2] = fielduInt16[1]
	binary.LittleEndian.PutUint16(fielduInt16, uint16(team1DataOffset))
	contestHeader[3] = fielduInt16[0]
	contestHeader[4] = fielduInt16[1]
	binary.LittleEndian.PutUint16(fielduInt16, uint16(0))
	contestHeader[5] = fielduInt16[0]
	contestHeader[6] = fielduInt16[1]
	binary.LittleEndian.PutUint16(fielduInt16, uint16(team2ColorsNamesOffset))
	contestHeader[7] = fielduInt16[0]
	contestHeader[8] = fielduInt16[1]
	binary.LittleEndian.PutUint16(fielduInt16, uint16(team2DataOffset))
	contestHeader[9] = fielduInt16[0]
	contestHeader[10] = fielduInt16[1]
	binary.LittleEndian.PutUint16(fielduInt16, uint16(0))
	contestHeader[11] = fielduInt16[0]
	contestHeader[12] = fielduInt16[1]
	binary.LittleEndian.PutUint16(fielduInt16, uint16(levelDataOffset))
	contestHeader[13] = fielduInt16[0]
	contestHeader[14] = fielduInt16[1]
	binary.LittleEndian.PutUint16(fielduInt16, uint16(contestDataOffset))
	contestHeader[15] = fielduInt16[0]
	contestHeader[16] = fielduInt16[1]

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

	initFrame := ContestFrame {Actions: make([]byte, 0), NumActions: 0}

	for i:=0; i<team1.NumDinos + team2.NumDinos; i++ {
		for j:=i+1; j<team1.NumDinos + team2.NumDinos; j++ {
			// distPairs = append(distPairs, &Distance{d1: i, d2: j})
		}

		sense = append(sense, &DinoSense {
			see: make([]byte, team1.NumDinos + team2.NumDinos),
			hear: make([]byte, team1.NumDinos + team2.NumDinos),
			smell: make([]byte, team1.NumDinos + team2.NumDinos)})

		delay = append(delay, &Delays {
			movement: 0,
			fire: 0,
			call: 0})

		// turn on dino ?
		initFrame.Put(&Action{code: 9, dino: byte(i), args: make([]byte, 0)})
	}

	cr.Actions = append(cr.Actions, byte(initFrame.NumActions))
	cr.Actions = append(cr.Actions, initFrame.Actions...)

	// distPairslen := len(distPairs)

	for arenaFrames > 0 {

		cf := ContestFrame { Actions: make([]byte, 0), NumActions: 0 }

		for i:=0; i<arena.numDinos; i++ {
			if delay[i].movement > 0 {
				delay[i].movement--
			}

			if delay[i].fire > 0 {
				delay[i].fire--
			}

			if delay[i].call > 0 {
				delay[i].call--
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

				if decisions[chosen].Movement == 0 && delay[i].call <= 0 {
					// call
					cf.Put(&Action{code: 10, dino: byte(i), args: make([]byte, 0)})

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

				// update position
			}
		}

		if cf.NumActions <= 0 {
			cf.Put(&Action{code: 11, dino: byte(0), args: []byte{10}})
		}

		cr.Actions = append(cr.Actions, byte(cf.NumActions))
		cr.Actions = append(cr.Actions, cf.Actions...)

		arenaFrames--;
	}

	testDieFrame := ContestFrame {Actions: make([]byte, 0), NumActions: 0}
	testDieFrame.Put(&Action{code: 11, dino: byte(0), args: []byte{0}})

	cr.Actions = append(cr.Actions, byte(testDieFrame.NumActions))
	cr.Actions = append(cr.Actions, testDieFrame.Actions...)

	endFrame := ContestFrame {Actions: make([]byte, 0), NumActions: 0}

	for i:=0; i<team1.NumDinos + team2.NumDinos; i++ {
		// turn off dino ?
		endFrame.Put(&Action{code: 11, dino: byte(i), args: []byte{9}})
	}



	cr.Actions = append(cr.Actions, byte(endFrame.NumActions))
	cr.Actions = append(cr.Actions, endFrame.Actions...)

	// cr.Actions = append(cr.Actions, make([]byte, 80)...)


	if len(cr.Actions) < 0xF000 {
		cr.Actions = append(cr.Actions, make([]byte, 0xF000 - (len(cr.Actions)%0xF000))...)
	} else {
		cr.Actions = append(cr.Actions, make([]byte, len(cr.Actions)%0xF000)...)
	}

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