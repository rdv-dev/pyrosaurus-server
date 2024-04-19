package ContestServer

import (
	"testing"
	"encoding/binary"
	"os"
	"io"
	"github.com/rdv-dev/pyrosaurus-server/ContestServer/util"
	"github.com/rdv-dev/pyrosaurus-server/ContestServer"
	// "https://github.com/stretchr/testify/assert"

	"fmt"
)

func TestRunContest(t *testing.T) {
	cases := []string {
		"Call",
		"Boss1",
        "BaseTeam",
		"Moves",
    }

	//caseType := []int {
	//	1, 0, 1}

	var levelData []byte

	directory, err := os.Getwd()

	if err != nil {
		t.Errorf("Unable to get current directory\n")
	}

	for i:=0; i<len(cases); i++ {
		testTeam1Path := directory + "/TestData/" + cases[i] + "/T1.TEAM"
		testTeam2Path := directory + "/TestData/" + cases[i] + "/T2.TEAM"

		fmt.Printf("Loading %s\n", testTeam1Path)
        team1DataFile, err := os.Open(testTeam1Path)
        if err != nil {
            t.Fail()
            t.Logf("Unable to open team file %s\n", testTeam1Path)
            fmt.Println(err)
            os.Exit(1)
        }

        team1Data, err := io.ReadAll(team1DataFile)

        if err != nil {
            t.Fail()
            t.Logf("Unable to read team file %s\n", testTeam1Path)
            fmt.Println(err)
            os.Exit(1)
        } 

		team1, err := util.NewContestEntry(team1Data)

		if err != nil {
			t.Fail()
			t.Logf("unable to parse team file %s\n", testTeam1Path)
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Loading %s\n", testTeam2Path)
        team2DataFile, err := os.Open(testTeam2Path)
        if err != nil {
            t.Fail()
            t.Logf("Unable to open team file %s\n", testTeam1Path)
            fmt.Println(err)
            os.Exit(1)
        }

        team2Data, err := io.ReadAll(team2DataFile)

        if err != nil {
            t.Fail()
            t.Logf("Unable to read team file %s\n", testTeam1Path)
            fmt.Println(err)
            os.Exit(1)
        } 

		team2, err := util.NewContestEntry(team2Data)

		if err != nil {
			t.Fail()
			t.Logf("unable to parse team file %s\n", testTeam2Path)
			fmt.Println(err)
			os.Exit(1)
		}

        fmt.Println("Loading Level")
		levelFile, err := os.Open(directory + "/TestData/" + cases[i] + "/LEVEL.000")

		if err != nil {
			t.Fail()
			t.Logf("Unable to open Level File\n")
			fmt.Println(err)
			os.Exit(1)
		} else {
			levelData, err = io.ReadAll(levelFile)

			if err != nil {
				t.Fail()
				t.Logf("Failed to load Level File\n")
				fmt.Println(err)
				os.Exit(1)
			}
		}

        fmt.Println("Running Contest")

		result, err := ContestServer.RunContest(team1, team2, levelData, 0)

		if err != nil {
            fmt.Println(err)
			t.Fail()
			t.Logf("Failed to run contest!\n")
			os.Exit(1)
		}


		if len(result.Actions) < 0 {
			fmt.Println("No result produced")
			t.Fail()
			t.Logf("No result produced\n")
			os.Exit(1)
		}

		//outFile, err := os.Create(directory + "/TestData/" + cases[i] + "/CONT.000")

		//if err != nil {
		//	t.Fail()
		//	t.Logf("Unable to open contest file for writing\n")
		//	os.Exit(1)
		//}


        fmt.Println("Exporting Contest")

		outdata, err := ContestServer.ExportContest(team1, team2, levelData, result)
        

		if err != nil {
			t.Fail()
			t.Logf("Failed to export contest\n")
			fmt.Println(err)
			os.Exit(1)
		}

		//outFile.Write(outdata)

		// outFile.Close()
        contestDataOffset := int(binary.LittleEndian.Uint16(outdata[0xF:0xF+2]))

        pos := contestDataOffset

        t.Logf("Contest Data Offset: %d, size: %d", pos, len(outdata))

        for pos < len(outdata) {
            fsize := int(outdata[pos])
            t.Logf("Frame size: %d\n", fsize)
            pos++

            frameCount := 0

            for frameCount < fsize && pos < len(outdata) {

                encodedByte := int(outdata[pos])
                dinoIndex := encodedByte / 12
                action := encodedByte - (dinoIndex * 12)

                frameCount++
                pos++

                switch action {
                case 0:
                    //t.Logf("Move Neck: ")
                    pos += 2
                case 1:
                    //t.Logf("Move Tail: ")
                    pos += 2
                case 2:
                    t.Logf("Move Dino(%d): Rot: %d Arg1: %d Arg2: %d", dinoIndex, int(outdata[pos]), int(outdata[pos+1]), int(outdata[pos+2]))
                    pos += 3
                case 3:
                    t.Logf("Set Breath Rate: ")
                    pos += 2
                case 4:
                    t.Logf("Step Left/Right: ")
                    pos++
                case 5:
                    t.Logf("Step Forward/Back: ")
                    pos++
                case 6:
                    t.Logf("Die")
                    pos++
                case 7:
                    t.Logf("Jump Left/Right")
                    pos++
                case 8:
                    t.Logf("Jump Forward/Back")
                    pos++
                case 9:
                    t.Logf("Locks neck movement?")
                case 10:
                    t.Logf("Call")
                case 11:
                    //t.Logf("Special Actions...")
                    pos++
                default:
                    t.Logf("Unknown action encountered!")
                    t.Fail()
                }
            } // frame details
        } // for each frame
        t.Log("Writing to file...")
        contestf, err := os.Create(cases[i] + ".bin")
        if err != nil {
            t.Log("Unable to open contest file for writing")
            fmt.Println(err)
            t.Fail()
        }

        _, err = contestf.Write(outdata)
        if err != nil {
            t.Log("Unable to open contest file for writing")
            fmt.Println(err)
            t.Fail()
        }

        defer contestf.Close()

	} // for each test
}
