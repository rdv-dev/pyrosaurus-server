package ContestServer

import (
	"testing"
	// "encoding/binary"
	"os"
	"io"
	"github.com/rdv-dev/pyrosaurus-server/ContestServer/util"

	"fmt"
)

func TestMoves(t *testing.T) {
	cases := []string {
		"Boss1",
		"Moves",
        "BaseTeam",
        "Call",
    }

    var levelData []byte

	directory, err := os.Getwd()

	if err != nil {
		t.Errorf("Unable to get current directory\n")
	}

	for i:=0; i<len(cases); i++ {
		testTeam1Path := directory + "/TestData/" + cases[i] + "/T1.TEAM"
		testTeam2Path := directory + "/TestData/" + cases[i] + "/T2.TEAM"

        t.Logf("*******************************************************")
        t.Logf("Loading Level")
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

        level := util.NewLevel(levelData)

		t.Logf("Loading %s\n", testTeam1Path)
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

        t.Logf("Dino Names: ")
        x := 0
        for d:=0; d<team1.NumDinos; {
            name := ""
            for team1.TeamData[team1.DinoNamesOffset+x] != 0 {
                name = name + string(team1.TeamData[team1.DinoNamesOffset+x])
                x++
            }
            t.Logf("%s", name)
            d++
            x++
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

        t.Logf("Dino Names: ")
        x = 0
        for d:=0; d<team2.NumDinos; {
            name := ""
            for team2.TeamData[team2.DinoNamesOffset+x] != 0 {
                name = name + string(team2.TeamData[team2.DinoNamesOffset+x])
                x++
            }
            t.Logf("%s", name)
            d++
            x++
        }
        
        arena := &Arena {
            Dinos: make([]*util.Dino, team1.NumDinos + team2.NumDinos),
            NumDinos: team1.NumDinos + team2.NumDinos,
        }

        // create dinos team 1
        speciesTypeOffset := ((util.TEAM_QUEEN_ARRAY_LEN + util.TEAM_SPECIES_LEG_NUM_LEN) * team1.NumDinos) + team1.DinosOffset + 1

        for i:=0; i<team1.NumDinos; i++ {
            arena.Dinos[i] = util.NewDino(team1, int(team1.TeamData[speciesTypeOffset]), i, level.X, level.Y)
            speciesTypeOffset += util.TEAM_MYSTERY_DATA
        }

        // create dinos team 2
        speciesTypeOffset = ((util.TEAM_QUEEN_ARRAY_LEN + util.TEAM_SPECIES_LEG_NUM_LEN) * team2.NumDinos) + team2.DinosOffset + 1

        for i:=team1.NumDinos; i<team1.NumDinos + team2.NumDinos; i++ {
            arena.Dinos[i] = util.NewDino(team2, int(team2.TeamData[speciesTypeOffset]), (i-team1.NumDinos), level.X, level.Y)
            speciesTypeOffset += util.TEAM_MYSTERY_DATA
        }

        t.Logf("Decision Movements: ")

        for i:=0; i<arena.NumDinos; i++ {
            t.Logf("Dino %d",i)
            for y:=0; y<len(arena.Dinos[i].Decisions); y++ {
                switch arena.Dinos[i].Decisions[y].Movement {
                case 0: 
                    t.Logf("Call")
                case 1: 
                    t.Logf("Don't Move")
                case 2:
                    t.Logf("Wander")
                case 3:
                    t.Logf("Move Away")
                case 4:
                    t.Logf("Move Closer")
                case 5: 
                    t.Logf("Move North")
                case 6:
                    t.Logf("Move South")
                default:
                    t.Logf("Custom Movement %d", arena.Dinos[i].Decisions[y].Movement - 7)
                }
            }
        }
    }
}
