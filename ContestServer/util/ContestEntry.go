package util
import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/binary"
	"errors"
	// "github.com/algae-disco/pyrosaurus-server/ContestServer/Dino"
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

type contestMessage struct {
	Data []byte
}

type contestEntryRaw struct {
	TeamData []byte
	Messages []*contestMessage
}

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

func NewContestEntry(filePath string, isRawFile int) (*ContestEntry, error) {
	rawFile, err := os.Open(filePath)

	if err != nil {
		return nil, errors.New("unable to open team file")
	}

	var teamData []byte

	if isRawFile == 1 {
		teamData = parseFile(rawFile).TeamData
	} else { 
		teamData, err = ioutil.ReadAll(rawFile)
	}

	if err != nil {
		return nil, errors.New("unable to read team file")
	}

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


func parseFile(rawFile *os.File) *contestEntryRaw {

	rawData, err := ioutil.ReadAll(rawFile)

	if err != nil {
		fmt.Println("Error reading team entry file")
		os.Exit(1)
	}

	if len(rawData) <= 0 {
		fmt.Println("File length 0, skipping")
		os.Exit(1)
	}

	readPos := 0
	toReadPos := 0
	fileLen := len(rawData)
	doRead := 1
	chunkSize := 0
	chunkNum := 1
	fileChunkNum := 0
	outData := make([]byte, 0)
	fileNum := 1
	// outFileName := ""

	entryData := &contestEntryRaw{TeamData: make([]byte, 0), Messages: make([]*contestMessage, 0)}

	for doRead == 1 {
		if rawData[readPos] == 0x2 && ((rawData[readPos] + rawData[readPos+1]) == 0xFF) && ((rawData[readPos+2] + rawData[readPos+3]) == 0xFF) {
			chunkSize = 0x400
		} else {
			if rawData[readPos] == 0x1 && ((rawData[readPos] + rawData[readPos+1]) == 0xFF) && ((rawData[readPos+2] + rawData[readPos+3]) == 0xFF) {
				chunkSize = 0x80
			} else {
				doRead = 0
			}
		}

		if doRead == 0 {
			fmt.Println("file not able to be parsed or already parsed")
			os.Exit(1)
		}

		if doRead == 1 {

			fileChunkNum = int(rawData[readPos + 2])

			if fileChunkNum == chunkNum {

				readPos += 4
				
				chunkNum += 1
				toReadPos = (readPos + chunkSize)
				
				outData = append(outData, rawData[readPos:toReadPos]...)

				readPos = toReadPos + 2

			} else {
				if fileChunkNum < chunkNum {
					// repeated chunk, skip it
					toReadPos = (readPos + chunkSize) + 2
					readPos = toReadPos
					fmt.Println("Skipping chunk...")
				}

				if fileChunkNum > chunkNum {
					// invalid file?
					fmt.Printf("File error, chunkNum less than file chunk number (%d, %d) readPos: %d\n", fileChunkNum, chunkNum, readPos)
					os.Exit(1)
				}	
			}		

			if rawData[readPos] == 0x4 && (rawData[readPos] + rawData[readPos + 1]) == 0xFF {
				// end of file

				chunkNum = 1

				if fileNum == 1 {
					// outFileName = "TEAM.BIN"
					entryData.TeamData = append(entryData.TeamData, outData...)
				} else {
					// outFileName = fmt.Sprintf("MSG-%d.bin", fileNum-1)
					entryData.Messages = append(entryData.Messages, &contestMessage{Data: outData})
				}

				fileNum += 1

				// outFile, err := os.Create(outFileName)

				// if err != nil {
				// 	fmt.Println("Error opening out file")
				// 	os.Exit(1)
				// }

				// outFile.Write(outData)
				// outData = make([]byte, 0)

				if readPos + 3 >= fileLen {
					doRead = 0
				} else {
					readPos += 2
				}
			}
		}
	}

	return entryData
}