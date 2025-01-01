package ModemServer

import (
	"fmt"
	//"context"
	//"syscall"
	"net"
	"os"
	"time"
	// "bufio"
	// "io/ioutil"
	"encoding/hex"
	"encoding/binary"
	"encoding/base64"
	"errors"

	"github.com/rdv-dev/pyrosaurus-server/Database"
)

type PyroUser struct {
	InternalPlayerId uint64
	PyroUserId uint32
	PyroCheckId uint16
	PyroVersion byte

	Conn net.Conn
	Mode int
	Submode int
	Data []byte

	Arena int
	Rating int
	GamesAvailable int
	LastOpponent uint32

	Contest []byte
	createdUser bool

}

type ContestMessage struct {
	Data []byte
}

type ContestEntryRaw struct {
	TeamData []byte
	Messages []*ContestMessage
}

var keyArray = []byte{}

const keyArrayEncoded = "AAAhEEIgYzCEQKVQxmDncAiBKZFKoWuxjMGt0c7h7/ExEhACczJSIrVSlEL3ctZiOZMYg3uzWqO905zD//Pe42IkQzQgBAEU5mTHdKREhVRqpUu1KIUJle7lz/WsxY3VUzZyJhEWMAbXdvZmlVa0Rlu3eqcZlziH3/f+553XvMfESOVYhmineEAIYRgCKCM4zMnt2Y7pr/lIiWmZCqkrufVa1Eq3epZqcRpQCjM6Eir929zLv/ue63mbWIs7uxqrpmyHfORMxVwiLAM8YAxBHK7tj/3szc3dKq0LvWiNSZ2XfrZu1V70ThM+Mi5RHnAOn/++793f/M8bvzqvWZ94j4iRqYHKseuhDNEtwU7xb+GAEKEAwjDjIARQJUBGcGdguYOYk/uj2rM9wxzTf+Ne87ECkBLzItIyNUIUUndiVnLqtculqJWJhW71T+Us1Q3F4jTDJKAUgQRmdEdkJFQFRNun+reZh7iXX+d+9x3HPNfTJvI2kQawFldmdnYVRjRWTNltyQ75L+nImemJirmrqURYZUgGeCdowBjhCII4oyh9y1zbP+se+/mL2Ju7q5q7dUpUWjdqFnrxCtAasyqSOi79D+1s3U3Nqr2LreidyY0mfAdsZFxFTKI8gyzgHMEMH+8+/13PfN+br7q/2Y/4nxduNn5VTnReky6yPtEO8B4="

func init() {
	var err error
	keyArray, err = base64.StdEncoding.DecodeString(keyArrayEncoded)

	if err != nil {
// 		fmt.Println("Error opening keyArray.bin")
		fmt.Println("Error decoding key Array")
		os.Exit(1)
	}

	Database.InitializeDatabase()
}

// func LoadValidationKey() { // ([]byte)
// 	keyArrayFile, err := os.Open("keyArray.bin")

// 	KeyArray, err = base64.StdEncoding.DecodeString(keyArrayEncoded)

// 	if err != nil {
// 		fmt.Println("Error opening keyArray.bin")
// 		fmt.Println("Error decoding key Array")
// 		os.Exit(1)
// 	}

// 	KeyArray, err := ioutil.ReadAll(keyArrayFile)

// 	if err != nil {
// 		fmt.Println("Error reading keyArray.bin")
// 		os.Exit(1)
// 	}

// 	keyArrayFile.Close()

// 	return keyArray
// }

func DoChallenge(user *PyroUser) (int, error) {

	validated := 0

    // Space out the timing of the challenge key
    _, err := user.Conn.Write([]byte{0x32}) // challenge bytes
    time.Sleep(250 * time.Millisecond)

	if err != nil {
		return -1, errors.New("Error with writing to connection")
    }

	user.Conn.Write([]byte{0x3C}) // challenge bytes
    time.Sleep(250 * time.Millisecond)

	user.Conn.Write([]byte{0x46}) // challenge bytes

	fmt.Println("Challenge key sent! ")

	// data := make([]byte, 1024)
	identity := make([]byte, 0)
	//buffer := bufio.NewReader(user.Conn)
	//dtFile, err := os.Create("pyro-data.bin")


	idTotal := 0

	for idTotal < 17 {

		nread, err := user.Conn.Read(user.Data)

		idTotal += nread

		if err != nil {
			return -1, errors.New("Error reading from socket: identity")
		}

		//fmt.Printf("Read data: %d/%d\n", nread, idTotal)
		// fmt.Printf(hex.EncodeToString(user.Data[:nread]) + "\n")
		identity = append(identity, user.Data[:nread]...)
		//fmt.Printf(hex.EncodeToString(identity) + "\n")
		//dtFile.Write(data[:nread])
	}

	checkByte1 := identity[0]
	checkByte2 := identity[1]
	pyroString := string(identity[2:8])
	pyroUserID := binary.LittleEndian.Uint32(identity[8:12])
	pyroCheckId := binary.LittleEndian.Uint16(identity[12:14])
	pyroVersion := identity[14]
	pyroDatalen := binary.LittleEndian.Uint16(identity[15:17])

	fmt.Printf("Check Byte result: %d\n", (checkByte1 + checkByte2))
	fmt.Printf("Pyro String: %s\n", pyroString)
	fmt.Printf("Pyro UserID: %d\n", pyroUserID)
	fmt.Printf("PyroCheckId: %d\n", pyroCheckId)
	fmt.Printf("Version: %d\n", pyroVersion)
	fmt.Printf("Checksum: %d\n", pyroDatalen)

	user.createdUser = false
	var playerId uint64

	if ((checkByte1 + checkByte2) == 255 && 
		pyroString == "PYROB0" &&
		(pyroVersion == 2 || pyroVersion == 3)) {
            // set the Pyro Version
		    user.PyroVersion = pyroVersion

            // perform other processing on User Id
			if pyroUserID > 0 {
                playerId, err = Database.GetPlayerByID(pyroUserID)
				if err != nil {
					user.Conn.Write([]byte{0x63})
					fmt.Println("Error during DoChallenge", err)
				}
			} else {
				playerId = 0
				user.createdUser = true
                user.PyroVersion = pyroVersion
			}
			if playerId > 0 {
				user.Conn.Write([]byte{0x27}) //validated pyroid
				validated = 1
				user.InternalPlayerId = playerId
				user.PyroUserId = pyroUserID
				user.PyroCheckId = pyroCheckId
			} else {
				pyroUserID, err = Database.CreatePlayer()

				if err != nil {
					fmt.Println("Error creating player")
					user.createdUser = false
				}
			
				playerId, err = Database.GetPlayerByID(pyroUserID)
			
				if err != nil {
					user.createdUser = false
					fmt.Println("After creating user, could not find interal id", err)
				}

				if user.createdUser == true {


					user.PyroUserId = pyroUserID
					user.InternalPlayerId = playerId
					user.PyroCheckId = uint16(0xCCD)
					user.Arena = 0xA
					user.Rating = 5
					user.GamesAvailable = 255
					user.LastOpponent = uint32(0)
				
					user.Conn.Write([]byte{0x27}) //validated pyroid
					validated = 1
				} else {
					user.Conn.Write([]byte{0x63})
				}
			}
		} else {
		user.Conn.Write([]byte{0x63})
	}

	return validated, nil
}

func DoTestConnection(user *PyroUser) (int, error) {

	success := 0

	testTotal := 0
	fmt.Printf("Sending code 7...")
	user.Conn.Write([]byte{0x07})

	contents := make([]byte,256)

	for i := 0x0; i < 256; i++ {
		contents[i] = byte(i)
	}

	time.Sleep(time.Second)

	for testTotal < 256 {

		nread, err := user.Conn.Read(user.Data)

		testTotal += nread

		if err != nil {
			return success, errors.New("Error reading from socket")
		}

		//fmt.Printf("Read data: %d/%d\n", nread, testTotal)
		//fmt.Printf(hex.EncodeToString(data[:nread]) + "\n")
	}

	user.Conn.Write([]byte{0x04})

	for i := 0; i < 256; i+=16 {
		user.Conn.Write(contents[i:i+16])
		//fmt.Printf(hex.EncodeToString(contents[i:i+16]) + "\n")
		time.Sleep(50)
	}

	success = 1
	

	fmt.Println("Sent test data!")

	return success, nil
}

func DoSpecialModes(user *PyroUser) (int, error) {
	fmt.Println("Confirming mode 1...")
	user.Conn.Write([]byte{0x01})

	// SUB-MODE 4 - send EMSG, TEMS files


	// SUB-MODE 5 - send updated user file
	if user.createdUser == true {
		// Pyro Id is 0, create new user and send user ID to them


		fmt.Println("Selecting sub-mode 5")
		user.Conn.Write([]byte{0x5, 0x5})

		fmt.Println("Ready...")
		user.Conn.Write([]byte{0x14})

		fmt.Println("Sending updated user file")

		i := 0
		sum := 0

		userData := make([]byte, 0)

		userId := make([]byte, 4)
		checkId := make([]byte, 2)
		opponent := make([]byte, 4)

		binary.LittleEndian.PutUint32(userId, user.PyroUserId)
		binary.LittleEndian.PutUint16(checkId, user.PyroCheckId)
		binary.LittleEndian.PutUint32(opponent, user.LastOpponent)

		userData = append(userData, userId...)
		userData = append(userData, checkId...)
		userData = append(userData, byte(user.Arena))
		userData = append(userData, byte(user.Rating))
		userData = append(userData, byte(user.GamesAvailable))
		userData = append(userData, opponent...)

		for i < len(userData) {
			sum = sum + int(userData[i])
			i++;
		}

		user.Conn.Write(userData)

		user.Conn.Write([]byte{ byte(sum) })
	}

	// SUB-MODE 6 - only available for version > 2
    //      check for contest, if one is available, send it

	if user.createdUser == false && len(user.Contest) > 0 && user.PyroVersion > 2 {
		fmt.Println("Selecting sub-mode 6")
		user.Conn.Write([]byte{0x6, 0x6})

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

		user.Conn.Write([]byte{0x3})

		sentFile, err := SendFile(user, user.Contest)

		if err != nil {
			return  0, errors.New("Failed to send Contest File")
		}

		return sentFile, nil
	}


	// SUB-MODE 7 - get BACKUP data

	// SUB-MODE 8 - Set modem return code to 0x14
    //      no need to implement, causes "noise" error and 
    //      does not advance game to Receive Contest state
	// SUB-MODE 9 - Send BACKUP block

	fmt.Println("Sending 0x64 for 'we're done!'...")
	user.Conn.Write([]byte{0x64, 0x64})

	// conn.SetDeadline(time.Now().Add(0))
	// _, err := conn.Read(data)

	// if err != nil {
	// 	fmt.Println("Error reading socket")
	// 	os.Exit(1)
	// }
	// defer server.Close()

	return 1, nil
}

func CheckForContest(user *PyroUser) (error) {

	contest_data, err := Database.GetContestByPlayerId(user.InternalPlayerId)
	if err != nil {
		return err
	}

	user.Contest = contest_data
	return nil
}

func GetFile(user *PyroUser) (*ContestEntryRaw, error) {

	fundata := make([]byte, 0)

	fmt.Printf("Confirming mode %d...\n", user.Mode)
	user.Conn.Write([]byte{byte(user.Mode)})

	fmt.Println("Sending ready code (1)...")
	user.Conn.Write([]byte{0x01})

	// if user.Mode == 2 {
	// 	fmt.Println("No contest available, sending 1...")
	// 	user.Conn.Write([]byte{0x1})
	// }

	errorCount := 0
	found14 := 0

	for found14 == 0 {

		// fmt.Println("Reading for 0x14...")

		nread, err := user.Conn.Read(user.Data)

		if err != nil {
			return nil, errors.New("Error reading mode")
			// os.Exit(1)
		}

		for i:=0; i<nread; i++ {
			if user.Data[i] == 0x14 {
				fmt.Println("Got 0x14, sending 0x14 and 0x47/0xB8 response...")
				user.Conn.Write([]byte{0x14, 0x47, 0xB8})
				found14 += 1
				break;
			}
		}
	}

	// fundata = append(fundata, data[1:nread]...)

	time.Sleep(time.Second)

	midCheckIndex := 4
	lastCheckIndex := 4

	fileTotal := 0
	timeout := 10*time.Millisecond
	doGetFileLoop := 1
	modeSelectCount := 0

	for errorCount < 3 && doGetFileLoop == 1 {

		// conn.SetDeadline(time.Now().Add(timeout))
		nread, err := user.Conn.Read(user.Data)

		fileTotal += nread

		if err != nil {
			errorCount += 1
			timeout += 4*time.Second
		} else { 
			errorCount -= 1
			timeout = 10*time.Millisecond; 
		}
		//defer server.Close()

		fundata = append(fundata,user.Data[:nread]...)

		//fmt.Printf("Read data: %d/%d/%d/%d\n", nread, midCheckIndex,lastCheckIndex, len(fundata))
		fmt.Printf(hex.EncodeToString(user.Data[:nread]) + "\n")

		for midCheckIndex <= len(fundata) - 4 {
			if int(fundata[midCheckIndex]) == 1 || int(fundata[midCheckIndex]) == 2 || int(fundata[midCheckIndex]) == 3 {
				if (int(fundata[midCheckIndex]) + int(fundata[midCheckIndex+1]) == 255) || (int(fundata[midCheckIndex+2]) + int(fundata[midCheckIndex+3]) == 255) {
					fmt.Println("Sending Check In 0x06F9")
					user.Conn.Write([]byte{0x06, 0xF9})
					midCheckIndex += 4
		// Set the lastCheckIndex to mid because we have already dealt with this data
		// the next loop doesn't need to consider this
					lastCheckIndex = midCheckIndex
					break;
				}
			}
			midCheckIndex += 1
		}

		for lastCheckIndex <= len(fundata) - 2 {
			if int(fundata[lastCheckIndex]) == 4 {
				if (int(fundata[lastCheckIndex]) + int(fundata[lastCheckIndex+1])) == 255 {
					fmt.Println("Final Chunk! Sending Check In 0x06F9")
					user.Conn.Write([]byte{0x06, 0xF9})
					lastCheckIndex += 2
					// timeout = 2*time.Second
					selectMode := -1

					for doGetFileLoop == 1 && errorCount < 4 {

						// conn.SetDeadline(time.Now().Add(timeout))

						nread, err := user.Conn.Read(user.Data)

						if err != nil {
							errorCount += 1
							timeout += 4*time.Second
							fmt.Println("Nothing to read...")
						} else {

							fmt.Printf(hex.EncodeToString(user.Data[:nread]) + "\n")
							lastCheckIndex = 0

							for lastCheckIndex < len(user.Data[:nread]) {
								if int(user.Data[lastCheckIndex]) == 1 || int(user.Data[lastCheckIndex]) == 4 {
									modeSelectCount += 1
									selectMode = int(user.Data[lastCheckIndex])
								} else {
									modeSelectCount = 0
								}
								lastCheckIndex += 1
							}

							if modeSelectCount >= 2 {
								doGetFileLoop = 0
								user.Mode = selectMode
								break
							}
						}
					}
				} else {
					// If we found a 4 but did not find 0xFB after it, move on
					if len(fundata) >= lastCheckIndex + 2 {
						lastCheckIndex += 1
						break
					}
				}
			} else {
				lastCheckIndex += 1
			}
		}
	}

	rawEntry := parseFile(fundata)
	
	binary.LittleEndian.PutUint32(rawEntry.TeamData[0:4], uint32(user.PyroUserId))
	}

	err := Database.CreateContestEntry(rawEntry.TeamData, user.InternalPlayerId)

	if err != nil {
		fmt.Println("Error creating contest entry, InternalPlayerId:",user.InternalPlayerId,":",err)
	}

	return rawEntry, nil
}

func SendFile(user *PyroUser, contents []byte) (int, error) {

	fundata := make([]byte, 0)

	//fmt.Println("Sending code 3...")
	//user.Conn.Write([]byte{0x03})

	//fmt.Println("Sending code 14...")
	//user.Conn.Write([]byte{0x14})

	//nread, err := user.Conn.Read(user.Data)

	//if err != nil {
	//	return 0, errors.New("Error reading mode")
	//}

    found14 := 0
    
    user.Conn.Write([]byte{0x14})


    //if user.Data[0] == 0x03 {
    //    fmt.Println("Mode 3 confirmed")
    //} else {
    //    fmt.Printf("Got this number: %d\n", int(user.Data[0]))
    //}

    fmt.Println("Normal Contest available (0x14)...")

    fmt.Println("Sending server ready (1)...")
    user.Conn.Write([]byte{0x1})

    nread, err := user.Conn.Read(user.Data)

    if err != nil {
        return 0, errors.New("Error reading mode")
    }

    if user.Data[0] == 0x14 {
        fmt.Println("Got 0x14")
        found14 += 1
    } else {
        fmt.Printf("Got this number: %d\n", int(user.Data[0]))
    }

	fundata = append(fundata, user.Data[:nread]...)

	bx := 0
	dx := make([]byte, 2)
	shortLastChunk := 0
	numChunks := 0

	if len(contents) % 0x400 == 0 {
		shortLastChunk = 0
	} else {
		shortLastChunk = 1
	}

	numChunks = int((len(contents)/0x400))
	extraChunks := 0
	chunkSize := 0x400
	byteCount := 0
	i := 0

	byteLen := numChunks * 0x400
	for byteLen < len(contents) {
		byteLen += 80
		extraChunks++
	}

	extraChunks -= 1


	for j:=0; j<numChunks+extraChunks; j++ {
		dx[0] = 0
		dx[1] = 0

		if shortLastChunk == 1 && j > numChunks - 1 {
			user.Conn.Write([]byte{0x1, 0xFE})
			chunkSize = 0x80
		} else {
			user.Conn.Write([]byte{0x2, 0xFD})
			chunkSize = 0x400
		}

		user.Conn.Write([]byte{byte(j+1), (0xFF - byte(j+1))})

		for byteCount = 0; byteCount < chunkSize; byteCount++ {

			if i<len(contents) { 
				user.Conn.Write([]byte{contents[i]})
				bx = int(contents[i])   // mov bl, [si-1]
			} else {
				user.Conn.Write([]byte{0})
				bx = 0
			}
		
			bx = bx ^ int(dx[1])  // xor bl, dh

			bx = bx << 1 	// shl bx, 1
			

			dx[1] = dx[0]
			dx[0] = 0

			dx[0] = dx[0] ^ keyArray[bx:bx+2][0]
			dx[1] = dx[1] ^ keyArray[bx:bx+2][1]

			i++

		}

		// fmt.Printf("Check Hash: %x\n", int(binary.BigEndian.Uint16(dx)))

		user.Conn.Write([]byte{dx[1], dx[0]})
	}

	user.Conn.Write([]byte{0x4, 0xFB})

	return 1, nil

}


func parseFile(rawData []byte) *ContestEntryRaw {

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

	entryData := &ContestEntryRaw{TeamData: make([]byte, 0), Messages: make([]*ContestMessage, 0)}

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
			fmt.Printf("file not able to be parsed or already parsed at %d\n", readPos)
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
					fmt.Printf("Skipping chunk %d...\n", chunkNum)
					fmt.Printf("Before: readPos: %d toReadPos: %d\n", readPos, toReadPos)
					toReadPos = (readPos + chunkSize) + 2 + 4
					readPos = toReadPos
					fmt.Printf("After: readPos: %d toReadPos: %d\n", readPos, toReadPos)
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
					entryData.Messages = append(entryData.Messages, &ContestMessage{Data: outData})
				}

				fileNum += 1

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
