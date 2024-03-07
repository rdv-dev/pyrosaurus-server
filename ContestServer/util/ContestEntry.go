package util
import (
	"fmt"
	"encoding/binary"
	"errors"
	// "github.com/rdv-dev/pyrosaurus-server/ContestServer/Dino"
)

const (
	ENTRY_HEADER_LEN = 15
	SPECIES_LEN = 0x20
	MOVE_DATA_LEN = 0x64 * 2
	FITE_DATA1_LEN = 0xFE
	FITE_DATA2_LEN = 0x32 * 2
	FITE_DATA3_LEN = 0x32 * 2
	DECISIONS_LEN = 0x17D

	NUM_DINOS_ON_TEAM_LEN = 1
	TEAM_QUEEN_ARRAY_LEN = 1
	TEAM_SPECIES_LEG_NUM_LEN = 1
	TEAM_DINO_RESIZE = 3
	TEAM_MYSTERY_DATA = 3
	TEAM_X_POS_LEN = 2
	TEAM_Y_POS_LEN = 2
	TEAM_ROT_LEN = 2
	TEAM_SOME_DATA_CONT_LEN = 2
	DINO_INIT_DATA_LEN = 22
	TEAM_ENTRY_RECORD_LEN = TEAM_QUEEN_ARRAY_LEN + TEAM_SPECIES_LEG_NUM_LEN + TEAM_MYSTERY_DATA + TEAM_X_POS_LEN + TEAM_Y_POS_LEN + TEAM_ROT_LEN + TEAM_SOME_DATA_CONT_LEN
	TEAM_CONTEST_RECORD_LEN = NUM_DINOS_ON_TEAM_LEN + TEAM_QUEEN_ARRAY_LEN + TEAM_SPECIES_LEG_NUM_LEN + TEAM_X_POS_LEN + TEAM_Y_POS_LEN + TEAM_ROT_LEN + TEAM_SOME_DATA_CONT_LEN

	TEAM_COLORS_LEN = 12
	TEAM_FIRE_COLORS_LEN = 6
	TEAM_COLORS_RECORD_LEN = TEAM_COLORS_LEN + TEAM_FIRE_COLORS_LEN
)

type ContestEntry struct {
	PyroUserId uint32
	SpeciesOffset int
	MovesOffset int
	FitePointsOffset int
	FiteXPosOffset int
	FiteYPosOffset int
	DecisionsOffset int
	DinosOffset int
	DinoNamesOffset int
	ColorsNamesOffset int
	NumSpecies int
	NumDinos int
	TeamData []byte
	Team uint32
	ColorsNamesSize int
	DinoNamesSize int
}

func NewContestEntry(teamData []byte) (*ContestEntry, error) {

	if len(teamData) <= 0 {
		return nil, errors.New("team file empty")
	}

	pyroUserId := binary.LittleEndian.Uint32(teamData[0:4])
	numSpecies := int(teamData[14])
	dinosOffset := int(binary.LittleEndian.Uint16(teamData[8:10]))
	dinoNamesOffset := int(binary.LittleEndian.Uint16(teamData[10:12]))
	colorsNamesOffset := int(binary.LittleEndian.Uint16(teamData[12:14]))
	numDinos := int(teamData[dinosOffset])

	fmt.Printf("User ID: %d, numSpecies: %d, dinosOffset: %d, numDinos: %d\n", pyroUserId, numSpecies, dinosOffset, numDinos)

	offsets := make([]int, 6)

	offset := ENTRY_HEADER_LEN
	offsets[0] = offset // species
	
	offset += SPECIES_LEN * numSpecies
	offsets[1] = offset // moves
	
	offset += MOVE_DATA_LEN * numSpecies
	offsets[2] = offset // fitePoints

	offset += FITE_DATA1_LEN * numSpecies
	offsets[3] = offset // fiteXPos

	offset += FITE_DATA2_LEN * numSpecies
	offsets[4] = offset // fiteYPos

	offset += FITE_DATA3_LEN * numSpecies
	offsets[5] = offset // decisions

	offset += DECISIONS_LEN * numSpecies

	if offset != dinosOffset {
		errorString := fmt.Sprintf("VALIDATION FAILED - Dino Names wrong position: readPos/teamDataOffset %d/%d\n", offset, dinosOffset)
		return nil, errors.New(errorString)
		
	}

	fmt.Printf("Offsets:\nSpecies: %d, Move: %d, FITE 1,2,3: %d,%d,%d, Decisions: %d\n",offsets[0], offsets[1], offsets[2], offsets[3], offsets[4], offsets[5])

	return &ContestEntry{
		PyroUserId: pyroUserId,
		SpeciesOffset: offsets[0],
		MovesOffset: offsets[1],
		FitePointsOffset: offsets[2],
		FiteXPosOffset: offsets[3],
		FiteYPosOffset: offsets[4],
		DecisionsOffset: offsets[5],
		DinosOffset: dinosOffset,
		DinoNamesOffset: dinoNamesOffset,
		ColorsNamesOffset: colorsNamesOffset,
		NumSpecies: numSpecies,
		NumDinos: numDinos,
		TeamData: teamData,
		Team: pyroUserId,
		ColorsNamesSize: len(teamData[colorsNamesOffset:]),
		DinoNamesSize: colorsNamesOffset}, nil
}


