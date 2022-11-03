package main
import (
	"fmt"
	"net"
	"os"
	"time"
	"github.com/algae-disco/pyrosaurus-server/ContestServer"
	"github.com/algae-disco/pyrosaurus-server/ModemServer"
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

			time.Sleep(10*time.Second)

			if success == 1 {
				// job.Mode = 0
			}

			// doModeLoop = 0

		case 2, 4:
			// Server gets file
			fileData, err := ModemServer.GetFile(job)

			if err != nil {
				fmt.Println("Error getting file", err.Error())
			}

			fmt.Sprintf("%s", fileData[0])

			// if 2, process entry, run contest
			// if 4, process message

		case 3:
			// Server sends data

			// export contest, send file

			success, err := ModemServer.SendFile(job, 0) // contest

			if err != nil {
				fmt.Println("Error sending file", err.Error())
			}

			if success == 1 {
				success = 1
			}


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