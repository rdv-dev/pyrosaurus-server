//go:build ignore
// Disabled: needs update for int16 Vector API

package ContestServer

import (
	"testing"
	// "encoding/binary"
	"os"
	"io"
	"github.com/rdv-dev/pyrosaurus-server/ContestServer/util"
    "math"
    "math/rand"
    "time"

	"fmt"
)

func TestCoordinates(t *testing.T) {
	cases := []string {
		"Boss1",
		//"Moves",
        //"BaseTeam",
        //"Call",
    }

    var levelData []byte

	directory, err := os.Getwd()

	if err != nil {
		t.Errorf("Unable to get current directory\n")
	}

	for i:=0; i<len(cases); i++ {

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
        t.Logf("*******************************************************")
        t.Logf("Loading Level %f %f", level.X, level.Y)

        RAD := math.Pi / 180
        //DEG := 180 / math.Pi

	rounds := float64(75)

	testPoints := make([]*Vector,0)
	testTargets := make([]*util.MovePoint,0) 
	testResults := make([]*Vector,0)

	// Test 1 - Cardinal direction
        testPoints = append(testPoints, &Vector {
            X: 0,
            Y: 0,
            A: 90 * RAD,
        })

        testTargets = append(testTargets, &util.MovePoint {
            X: 0,
            Y: 100,
            GoalSize: 0,
        })

	testResults = append(testResults, &Vector {
		X: 0,
		Y: 75,
		A: 90 * RAD,
	})

	// Test 2 - Cardinal direction
        testPoints = append(testPoints, &Vector {
            X: 0,
            Y: 0,
            A: 180 * RAD,
        })

        testTargets = append(testTargets, &util.MovePoint {
            X: -75,
            Y: 0,
            GoalSize: 0,
        })

	testResults = append(testResults, &Vector {
		X: -75,
		Y: 0,
		A: 180 * RAD,
	})

	// Test 3 - Cardinal direction
        testPoints = append(testPoints, &Vector {
            X: 0,
            Y: 0,
            A: 270 * RAD,
        })

        testTargets = append(testTargets, &util.MovePoint {
            X: 0,
            Y: -100,
            GoalSize: 0,
        })

	testResults = append(testResults, &Vector {
		X: 0,
		Y: -75,
		A: 270 * RAD,
	})
	
	// Test 4 - Cardinal direction
        testPoints = append(testPoints, &Vector {
            X: 0,
            Y: 0,
            A: 0,
        })

        testTargets = append(testTargets, &util.MovePoint {
            X: 100,
            Y: 0,
            GoalSize: 0,
        })

	testResults = append(testResults, &Vector {
		X: 75,
		Y: 0,
		A: 0,
	})
	
	// Test 5 - Turn around to get to goal
        testPoints = append(testPoints, &Vector {
            X: 0,
            Y: 0,
            A: 0,
        })

        testTargets = append(testTargets, &util.MovePoint {
            X: -75,
            Y: 0,
            GoalSize: 0,
        })

	testResults = append(testResults, &Vector {
		X: -75,
		Y: 0,
		A: 0,
	})
	
	// Test 6 - Turn around to get to goal
        testPoints = append(testPoints, &Vector {
            X: 800,
            Y: 800,
            A: 0,
        })

	s1 := rand.NewSource(time.Now().Unix())
	r1 := rand.New(s1)
        testTargets = append(testTargets, &util.MovePoint {
            X: float64((r1.Intn(10) - 5) * 4000),
            Y: float64((r1.Intn(10) - 5) * 4000),
            GoalSize: 0,
        })

	testResults = append(testResults, &Vector {
		X: 0,
		Y: 0,
		A: 0,
	})
	
        for i:=0; i<len(testPoints) - 1; i++  {

		currentPoint := testPoints[i]
		targetPoint := testTargets[i]

		for r:=0; r<int(rounds); r++ {
            boundVector := CheckBounds(currentPoint, level) 
	    //t.Logf("Bounds: %f %f", boundVector.X, boundVector.Y)
            newPos, rot := CalculatePosition(currentPoint, boundVector, targetPoint)
	    rot = rot
            //corrected, boundAngleF, boundAngle := CheckBounds(currentPoint, level) 

            currentPoint.X = currentPoint.X + (newPos.X *1)
            currentPoint.Y = currentPoint.Y + (newPos.Y *1)
            currentPoint.A = newPos.A
		r=r
	    }

	    //t.Logf("Position: %f %f %f", currentPoint.X, currentPoint.Y, currentPoint.A * DEG)

	    if i < 4 && testResults[i].X != currentPoint.X && testResults[i].Y != currentPoint.Y && testResults[i].A != currentPoint.A {
		    t.Logf("Test failed on %d expect: %f %f %f got: %f %f %f", i, testResults[i].X, testResults[i].Y, testResults[i].A, currentPoint.X, currentPoint.Y, currentPoint.A)
		    t.Fail()
	    }

	    if i >= 4 {
		    dx := testResults[i].X - currentPoint.X
		    dy := testResults[i].Y - currentPoint.Y
		    if math.Sqrt(dx*dx + dy*dy) > float64(10) {
			t.Logf("Test failed on %d expect distance less than 10 got: %f", i, math.Sqrt(dx*dx + dy*dy))
	    		t.Logf("Position: %f %f %f", currentPoint.X, currentPoint.Y, currentPoint.A)
			t.Fail()
		    }
	    }
        }

		currentPoint := testPoints[5]
		targetPoint := testTargets[5]
	for i:=0; i<90000; i++ {

            boundVector := CheckBounds(currentPoint, level) 
	    //t.Logf("Bounds: %f %f", boundVector.X, boundVector.Y)
            newPos, rot := CalculatePosition(currentPoint, boundVector, targetPoint)
	    rot = rot
            //corrected, boundAngleF, boundAngle := CheckBounds(currentPoint, level) 

            currentPoint.X = currentPoint.X + (newPos.X *5)
            currentPoint.Y = currentPoint.Y + (newPos.Y *5)
            currentPoint.A = newPos.A

	    plimit := float64(-50)


	    if currentPoint.X > level.X + plimit || currentPoint.X < -level.X - plimit {
		    t.Logf("Got out of bounds on test 6!")
	    		t.Logf("Position: %f %f %f", currentPoint.X, currentPoint.Y, currentPoint.A)
			t.Logf("Rounds: %d", i)
			t.Logf("Bounding: %f %f %f", boundVector.X, boundVector.Y, boundVector.A)
			t.Fail()
			break;
		}

	    if currentPoint.Y > level.Y + plimit || currentPoint.Y < -level.Y - plimit {
		    t.Logf("Got out of bounds on test 6!")
	    		t.Logf("Position: %f %f %f", currentPoint.X, currentPoint.Y, currentPoint.A)
			t.Logf("Rounds: %d", i)
			t.Logf("Bounding: %f %f %f", boundVector.X, boundVector.Y, boundVector.A)
			t.Fail()
			break;
		}
		
	}
	
	outdata := ""
	for z:=0; z<1; z++ {
		currentPoint = testPoints[5]
		targetPoint = testTargets[5]
            boundVector := CheckBounds(currentPoint, level) 
		outdata = ""
	    //doplot := 0
		limitz := float64(100)
	for i:=0; i<90000; i++ {

	    //t.Logf("Bounds: %f %f", boundVector.X, boundVector.Y)
            boundVector = CheckBoundsV(currentPoint, level) 
            newPos, rot := CalculatePosition(currentPoint, boundVector, targetPoint)
	    rot = rot
            //corrected, boundAngleF, boundAngle := CheckBounds(currentPoint, level) 

            currentPoint.X = currentPoint.X + (newPos.X *5)
            currentPoint.Y = currentPoint.Y + (newPos.Y *5)
            currentPoint.A = newPos.A
		outdata = outdata + fmt.Sprintf("%f,%f\n", currentPoint.X, currentPoint.Y)


	//    if currentPoint.X > level.X + limitz || currentPoint.X < -level.X - limitz {
	//	    doplot = 1
	//    }
	//    if currentPoint.Y > level.Y + limitz || currentPoint.Y < -level.Y - limitz {
	//	    doplot = 1
	//    }
	//    
	//    if doplot == 1 {
	//    	t.Logf("%f %f", currentPoint.X, currentPoint.Y)
    	//	}

		
	}

	    if currentPoint.X > level.X + limitz || currentPoint.X < -level.X - limitz {
		    t.Logf("Got out of bounds on test 7!")
	    		t.Logf("Position: %f %f %f", currentPoint.X, currentPoint.Y, currentPoint.A)
			t.Logf("Rounds: %d", z)
            boundVector = CheckBounds(currentPoint, level) 
			t.Logf("Bounding: %f %f %f", boundVector.X, boundVector.Y, boundVector.A)
			t.Fail()
			break;
		}

	    if currentPoint.Y > level.Y + limitz || currentPoint.Y < -level.Y - limitz {
		    t.Logf("Got out of bounds on test 7!")
	    		t.Logf("Position: %f %f %f", currentPoint.X, currentPoint.Y, currentPoint.A)
			t.Logf("Rounds: %d", z)
            boundVector = CheckBounds(currentPoint, level) 
			t.Logf("Bounding: %f %f %f", boundVector.X, boundVector.Y, boundVector.A)
			t.Fail()
			break;
		}
	}
	
	
	contestf, err := os.Create("Plottest1.bin")
	defer contestf.Close()
	if err != nil {
		t.Log("Unable to open contest file for writing")
		fmt.Println(err)
		t.Fail()
	}

	_, err = contestf.WriteString(outdata)
	if err != nil {
		t.Log("Unable to open contest file for writing")
		fmt.Println(err)
		t.Fail()
	}
	
    }
}
