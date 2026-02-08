//go:build ignore
// Disabled: needs update for int16 Vector API

package ContestServer
import (

	"testing"
	//"encoding/binary"
	"os"
    "math"
	//"io"
	"github.com/rdv-dev/pyrosaurus-server/ContestServer/util"

	"fmt"
)

func TestPlotMoves(t *testing.T) {
	outdata := ""
		
    testPoints := make([]*util.MovePoint, 0)

    testPoints = append(testPoints, &util.MovePoint {
    X: -20,
    Y: -20,
    GoalSize: 0,})

    testPoints = append(testPoints, &util.MovePoint {
    X: 20,
    Y: -20,
    GoalSize: 0,})

    testPoints = append(testPoints, &util.MovePoint {
    X: -20,
    Y: -20,
    GoalSize: 0,})

    testPoints = append(testPoints, &util.MovePoint {
    X: 20,
    Y: -20,
    GoalSize: 0,})

    testPoints = append(testPoints, &util.MovePoint {
    X: -20,
    Y: -20,
    GoalSize: 0,})

    testPoints = append(testPoints, &util.MovePoint {
    X: 20,
    Y: -20,
    GoalSize: 0,})

    testPoints = append(testPoints, &util.MovePoint {
    X: -20,
    Y: -20,
    GoalSize: 0,})

    level := &util.Level {X:3000, Y:3000,}


    currentPoint := &Vector { X:0, Y:30, A:90 * RAD, }

    for i:=0; i<len(testPoints); i++  {

    targetPoint := testPoints[i]
        xDist := float64(0)
        yDist := float64(0)
        dist := float64(0)

	for r:=0; r<100; r++ {
        boundVector := CheckBoundsV(currentPoint, level) 
        //t.Logf("Bounds: %f %f", boundVector.X, boundVector.Y)
        newPos, rot := CalculatePosition(currentPoint, boundVector, targetPoint)
        rot = rot
        //corrected, boundAngleF, boundAngle := CheckBounds(currentPoint, level) 

        currentPoint.X = currentPoint.X + (newPos.X *1)
        currentPoint.Y = currentPoint.Y + (newPos.Y *1)
        currentPoint.A = newPos.A
            outdata = outdata + fmt.Sprintf("%f,%f\n", currentPoint.X, currentPoint.Y)
            xDist = currentPoint.X - testPoints[i].X
            yDist = currentPoint.Y - testPoints[i].Y
            dist = math.Sqrt(xDist*xDist + yDist*yDist)

            if dist < 10 {
               break;
            }
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

func TestPlotOutOfBounds(t *testing.T) {
	outdata := ""
		
    testPoints := make([]*util.MovePoint, 0)

    testPoints = append(testPoints, &util.MovePoint {
    X: 2800,
    Y: 0,
    GoalSize: 0,})

    testPoints = append(testPoints, &util.MovePoint {
    X: -2800,
    Y: 0,
    GoalSize: 0,})

    testPoints = append(testPoints, &util.MovePoint {
    X: 2800,
    Y: 0,
    GoalSize: 0,})

    testPoints = append(testPoints, &util.MovePoint {
    X: -2800,
    Y: 0,
    GoalSize: 0,})

    testPoints = append(testPoints, &util.MovePoint {
    X: 2800,
    Y: 0,
    GoalSize: 0,})

    level := &util.Level {X:3000, Y:3000,}


    currentPoint := &Vector { X:0, Y:30, A:90 * RAD, }

    for i:=0; i<len(testPoints); i++  {

    targetPoint := testPoints[i]
        xDist := float64(0)
        yDist := float64(0)
        dist := float64(0)

	for r:=0; r<30000; r++ {
        boundVector := CheckBoundsV(currentPoint, level) 
        //t.Logf("Bounds: %f %f", boundVector.X, boundVector.Y)
        newPos, rot := CalculatePosition(currentPoint, boundVector, targetPoint)
        rot = rot
        //corrected, boundAngleF, boundAngle := CheckBounds(currentPoint, level) 

        currentPoint.X = currentPoint.X + (newPos.X *50)
        currentPoint.Y = currentPoint.Y + (newPos.Y *50)
        currentPoint.A = newPos.A
            outdata = outdata + fmt.Sprintf("%f,%f\n", currentPoint.X, currentPoint.Y)
            xDist = currentPoint.X - testPoints[i].X
            yDist = currentPoint.Y - testPoints[i].Y
            dist = math.Sqrt(xDist*xDist + yDist*yDist)

            if dist < 150 {
               break;
            }
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
