package main
import (
	"fmt"
	"net"
	"os"
	"io"
	"time"
	"github.com/rdv-dev/pyrosaurus-server/ContestServer"
	"github.com/rdv-dev/pyrosaurus-server/ContestServer/util"
	"github.com/rdv-dev/pyrosaurus-server/ModemServer"
	// "github.com/rdv-dev/pyrosaurus-server/Database"
	// "strconv"
)

const (
	SERVER_HOST="127.0.0.1"
	SERVER_PORT="8888"
	SERVER_TYPE="tcp"
)

func main() {

	//config := &net.ListenConfig(Control: reusePort)

	fmt.Println("Setting up server...")

	pyroJobs := make(chan *ModemServer.PyroUser)

	//server, err := config.Listen(context.Background(), SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		fmt.Println("Error setting up socket", err.Error())
		os.Exit(1)
	}
	defer server.Close()

	fmt.Println("Listening on " + SERVER_HOST+":"+SERVER_PORT)

	// KeyArray = LoadValidationKey()
	// ModemServer.LoadValidationKey()

	for {

		conn, err := server.Accept()

		if err != nil {
			fmt.Println("Error accepting modem", err.Error())
			os.Exit(1)
		}
		// defer server.Close()

		go func() {
			fmt.Println("Got connection")

			user := &ModemServer.PyroUser {
					PyroUserId: 0,
					PyroCheckId: 0,
					PyroVersion: 0,
					// active: 1,
					Conn: conn,
					Mode: 0,
					Submode: 0,
					Data: make([]byte,1024)}

			validated, err := ModemServer.DoChallenge(user)

			if err != nil {
				fmt.Println("Error during Challenge", err.Error())
				os.Exit(1)
			}
			// defer server.Close()

			if validated == 1 {

				go handleModemJobs(pyroJobs)

				idTotal := 0
				mmode := make([]byte, 0)

				for idTotal < 2 {
					nread, err := user.Conn.Read(user.Data)

					idTotal += nread

					if err != nil {
						fmt.Println("Error reading from socket: Mode", err.Error())
						break;
					}

					mmode = append(mmode, user.Data[:nread]...)
				}

				user.Mode = int(mmode[0])

				for {

					
					pyroJobs <- user

					if user.Mode == 0 {
						fmt.Println("Closing connection to client")
						break;
					}
				}
			}
		}()
	}

	
	fmt.Println("Closing...")

	// fmt.Println("Writing fundata to file...")
	// f, err := os.Create("T.TMP.1")
	// if err != nil {
	// 	fmt.Println("Error writing file")
	// 	os.Exit(1)
	// }
	// defer server.Close()

	// f.Write(fundata)

	server.Close()
	
}


func handleModemJobs(pyroJobs chan *ModemServer.PyroUser) {
	for {
		job := <- pyroJobs

		// midCheckIndex := 4
		// lastCheckIndex := 4
		// fileTotal := 0
		// fundata := make([]byte,0)
		// doModeLoop := 1

		// for doModeLoop == 1 {

		// if job.Mode == 0 {
		// 	break;
		// }

		fmt.Printf("Select mode: %d\n", job.Mode)

		switch job.Mode {
		case 0:
			// stop looping
			break;
		case 1:
			// Backup data, Get Messages
			
			success, err := ModemServer.DoSpecialModes(job)

			if err != nil {
				fmt.Println("Error during special mode", err.Error())
			}

			time.Sleep(1*time.Second)

			if success == 1 {
				job.Mode = 0
			}

			// doModeLoop = 0

		case 2, 4:
			// TODO: MAJOR ISSUE - Contest/Message files can be handled separately now
			//		as it is coded here now, sending a message will cause an issue.
			// Server gets file
			contestEntry, err := ModemServer.GetFile(job)

			if err != nil {
				fmt.Println("Error getting file", err.Error())
			}

			// fmt.Sprintf("%s", fileData[0])
            var levelData []byte

            // levelFile, err := os.Open("/home/rob/pyro-c/db/levels/LEVEL.000")
            levelFile, err := os.Open("Assets/Levels/LEVEL.000")

            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            } else {
                levelData, err = io.ReadAll(levelFile)

                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }
            }

			// if 2, process entry, run contest
			team1, err := util.NewContestEntry(contestEntry.TeamData)

			if err != nil {
				fmt.Println("Error processing player team entry", err.Error())
			}

			team2EntryId, team2 := ContestServer.FindOpponent(job.InternalPlayerId)
            fmt.Printf("Found opponent: %d\n", team2EntryId)

			if team2 != nil {
				result, err := ContestServer.RunContest(team1, team2, levelData, 0)

				fmt.Printf("Contest length: %d\n", len(result.Actions))

				if err != nil {
					fmt.Println("Failed to run contest!", err.Error())
					job.Mode = 0
				}


				if len(result.Actions) < 0 {
					fmt.Println("No result produced")
					job.Mode = 0
				} else {

                    
                    if job.PyroVersion > 2 {
                        // Modded MODEM versions 
                        outdata, err := ContestServer.ExportContest(team1, team2, levelData, result)

                        if err != nil {
                            fmt.Println(err)
                            os.Exit(1)
                        }

                        job.Contest = outdata
                    } else {
                        // Original MODEM version
                        ContestServer.SaveContest(team1, team2, levelData, result, job.InternalPlayerId, team2EntryId)
                    }
				}
			} else {

			}


			// if 4, process message

			
			// }

		case 3:
			// Server sends data

			// export contest, send file

			job.Conn.Write([]byte{0x03})

			err := ModemServer.CheckForContest(job)
			if err != nil {
				fmt.Println("Failed to check for contest", err)
                // Contest not available
			    job.Conn.Write([]byte{0x21})
			} else {

                sentFile, err := ModemServer.SendFile(job, job.Contest)

                if err != nil {
                    fmt.Println("Error sending file", err.Error())
                }

                if sentFile == 1 {
                    fmt.Println("Sent Contest file!")
                }
            }

            job.Mode = 0

		case 7:
			success, err := ModemServer.DoTestConnection(job)

			if err != nil {
				fmt.Println("Error during test: ", err.Error())
				os.Exit(1)
			}

			if success == 1 {
				job.Mode = 0

				fmt.Println("Writing updated phone number")

				sum := 0x30 + 0x30 + 0x30 + 0x30 + 0x30 + 0x31 + 0x38 + 0x38 + 0x38 + 0x38 + 0x00 + 0x00
				sum = sum & 0xFF

				job.Conn.Write([]byte{0x2,0x02,0x30,0x30,0x30,0x30,0x30,0x31,0x38,0x38,0x38,0x38,0x00,0x00, byte(sum)})
			}

			// doModeLoop = 0

		default:
			fmt.Println("Invalid mode, exiting")
			break;
		}
	// }
	}
}
