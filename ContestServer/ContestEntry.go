package ContestEntry
import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/binary"
)

const (
	ENTRY_HEADER_LEN = 15
	MAX_DINOS_ON_TEAM = 10
	SPECIES_LEN = 0x20
	MOVE_DATA_LEN = 0x64 * 2
	FITE_DATA1_LEN = 0xFE
	FITE_DATA2_LEN = 0x32 * 2
	FITE_DATA3_LEN = 0x32 * 2
	DECISIONS_LEN = 0x17D

	NUM_DINOS_ON_TEAM_LEN = 1
	TEAM_QUEEN_ARRAY_LEN = 1
	TEAM_SPECIES_LEG_NUM_LEN = 1
	TEAM_MYSTERY_DATA = 3
	TEAM_X_POS_LEN = 2
	TEAM_Y_POS_LEN = 2
	TEAM_ROT_LEN = 2
	TEAM_SOME_DATA_CONT_LEN = 2
	DINO_INIT_DATA_LEN = 22
	TEAM_ENTRY_RECORD_LEN = TEAM_QUEEN_ARRAY_LEN + TEAM_SPECIES_LEG_NUM_LEN + TEAM_MYSTERY_DATA + TEAM_X_POS_LEN + TEAM_Y_POS_LEN + TEAM_ROT_LEN + TEAM_SOME_DATA_CONT_LEN + DINO_INIT_DATA_LEN
	TEAM_CONTEST_RECORD_LEN = TEAM_QUEEN_ARRAY_LEN + TEAM_SPECIES_LEG_NUM_LEN + TEAM_X_POS_LEN + TEAM_Y_POS_LEN + TEAM_ROT_LEN + TEAM_SOME_DATA_CONT_LEN

	TEAM_COLORS_LEN = 12
	TEAM_FIRE_COLORS_LEN = 6
	TEAM_COLORS_RECORD_LEN = TEAM_COLORS_LEN + TEAM_FIRE_COLORS_LEN
)

type TeamData struct {
	dinoData []byte
	dinoNames []byte
	playerData []byte
}

func main() {
	rawFile, err := os.Open("TEAM.BIN")

	if err != nil {
		fmt.Println("unable to open team file")
		os.Exit(1)
	}

	teamData, err := ioutil.ReadAll(rawFile)

	if err != nil {
		fmt.Println("unable to read team file")
		os.Exit(1)
	}

	if len(teamData) <= 0 {
		fmt.Println("team file empty")
		os.Exit(1)
	}

	readPos := 0

	header := teamData[readPos:ENTRY_HEADER_LEN]
	readPos += ENTRY_HEADER_LEN

	pyroUserID := int(binary.LittleEndian.Uint32(header[0:4]))
	levelCheckSum := int(binary.LittleEndian.Uint32(header[4:8]))
	teamDataOffset := int(binary.LittleEndian.Uint16(header[8:10]))
	dinoNamesOffset := int(binary.LittleEndian.Uint16(header[10:12]))
	teamColorsNames := int(binary.LittleEndian.Uint16(header[12:14]))
	numSpecies := int(header[14])

	fmt.Printf("Pyro User ID: %d\n", pyroUserID)
	fmt.Printf("Level Checksum: %d\n", levelCheckSum)
	fmt.Printf("Team Data Offset: %d\n", teamDataOffset)
	fmt.Printf("Dino Names Offset: %d\n", dinoNamesOffset)
	fmt.Printf("Team Colors and Names Offset: %d\n", teamColorsNames)
	fmt.Printf("Species count: %d\n", numSpecies)

	speciesData := make([][]byte, 0)
	moveData := make([][]byte, 0)
	fitePointData := make([][]byte, 0)
	fiteXPos := make([][]byte, 0)
	fiteYPos := make([][]byte, 0)
	decisionData := make([][]byte, 0)

	for i:=0; i<numSpecies; i++ {
		sData := make([]byte, 0)
		sData = append(sData, teamData[readPos:(readPos + SPECIES_LEN)]...)
		speciesData = append(speciesData,sData)
		readPos += SPECIES_LEN
	}

	for i:=0; i<numSpecies; i++ {
		sData := make([]byte, 0)
		sData = append(sData, teamData[readPos:(readPos + MOVE_DATA_LEN)]...)
		moveData = append(moveData,sData)
		readPos += MOVE_DATA_LEN
	}

	for i:=0; i<numSpecies; i++ {
		sData := make([]byte, 0)
		sData = append(sData, teamData[readPos:(readPos + FITE_DATA1_LEN)]...)
		fitePointData = append(fitePointData,sData)
		readPos += FITE_DATA1_LEN
	}

	for i:=0; i<numSpecies; i++ {
		sData := make([]byte, 0)
		sData = append(sData, teamData[readPos:(readPos + FITE_DATA2_LEN)]...)
		fiteXPos = append(fiteXPos,sData)
		readPos += FITE_DATA2_LEN
	}

	for i:=0; i<numSpecies; i++ {
		sData := make([]byte, 0)
		sData = append(sData, teamData[readPos:(readPos + FITE_DATA3_LEN)]...)
		fiteYPos = append(fiteYPos,sData)
		readPos += FITE_DATA3_LEN
	}

	for i:=0; i<numSpecies; i++ {
		sData := make([]byte, 0)
		sData = append(sData, teamData[readPos:(readPos + DECISIONS_LEN)]...)
		decisionData = append(decisionData,sData)
		readPos += DECISIONS_LEN
	}

	if readPos != teamDataOffset {
		fmt.Println("CONTEST ENTRY VALIDATION FAILED - Dino Names wrong position")
		fmt.Printf("readPos/teamDataOffset %d/%d\n", readPos, teamDataOffset)
		os.Exit(1)
	}



	// dinoData := teamData[teamDataOffset:dinoNamesOffset]
	// dinoNames := teamData[dinoNamesOffset:teamColorsNames]
	// playerData := teamData[teamColorsNames:]


}