package ContestServer

import (
	"testing"
	"os"
	"io/ioutil"
	"github.com/algae-disco/pyrosaurus-server/ContestServer/util"
	"github.com/algae-disco/pyrosaurus-server/ContestServer"
	// "https://github.com/stretchr/testify/assert"

	"fmt"
)

func TestRunContest(t *testing.T) {
	cases := []string {
		"Call",
		"BaseTeam",
		"Moves"}

	caseType := []int {
		1, 0, 1}

	var levelData []byte

	directory, err := os.Getwd()

	if err != nil {
		t.Errorf("Unable to get current directory\n")
	}

	for i:=0; i<len(cases); i++ {
		testTeam1Path := directory + "/TestData/" + cases[i] + "/T1.TEAM"
		testTeam2Path := directory + "/TestData/" + cases[i] + "/T2.TEAM"

		fmt.Printf("Loading %s\n", testTeam1Path)

		team1, err := util.NewContestEntry(testTeam1Path, caseType[i])

		if err != nil {
			t.Fail()
			t.Logf("unable to load team file %s\n", testTeam1Path)
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Loading %s\n", testTeam2Path)

		team2, err := util.NewContestEntry(testTeam2Path, caseType[i])

		if err != nil {
			t.Fail()
			t.Logf("unable to load team file %s\n", testTeam2Path)
			fmt.Println(err)
			os.Exit(1)
		}

		result, err := ContestServer.RunContest(team1, team2)

		if err != nil {
			t.Fail()
			t.Logf("Failed to run contest!\n")
			os.Exit(1)
		}


		if len(result.Actions) < 0 {
			t.Fail()
			t.Logf("No result produced\n")
			os.Exit(1)
		}

		outFile, err := os.Create(directory + "/TestData/" + cases[i] + "/CONT.000")

		if err != nil {
			t.Fail()
			t.Logf("Unable to open contest file for writing\n")
			os.Exit(1)
		}

		levelFile, err := os.Open(directory + "/TestData/" + cases[i] + "/LEVEL.000")

		if err != nil {
			t.Fail()
			t.Logf("Unable to open Level File\n")
			fmt.Println(err)
			os.Exit(1)
		} else {
			levelData, err = ioutil.ReadAll(levelFile)

			if err != nil {
				t.Fail()
				t.Logf("Failed to load Level File\n")
				fmt.Println(err)
				os.Exit(1)
			}
		}

		outdata, err := ContestServer.ExportContest(team1, team2, levelData, result)

		if err != nil {
			t.Fail()
			t.Logf("Failed to export contest\n")
			fmt.Println(err)
			os.Exit(1)
		}

		outFile.Write(outdata)

		// outFile.Close()
	}

}