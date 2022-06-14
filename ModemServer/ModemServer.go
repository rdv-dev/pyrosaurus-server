package main
import (
	"fmt"
	//"context"
	//"syscall"
	"net"
	"os"
	"time"
	"bufio"
	"io/ioutil"
	"encoding/hex"
	"encoding/binary"
)

const (
	SERVER_HOST="127.0.0.1"
	SERVER_PORT="8888"
	SERVER_TYPE="tcp"
)

func main() {

	//config := &net.ListenConfig(Control: reusePort)

	fmt.Println("Setting up server...")
	//server, err := config.Listen(context.Background(), SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		fmt.Println("Error setting up socket", err.Error())
		os.Exit(1)
	}
	defer server.Close()

	fmt.Println("Listening on " + SERVER_HOST+":"+SERVER_PORT)

	keyArrayFile, err := os.Open("keyArray.bin")

	if err != nil {
		fmt.Println("Error opening keyArray.bin")
		os.Exit(1)
	}

	// keyArrayBytes, err := ioutil.ReadAll(keyArrayFile)
	keyArray, err := ioutil.ReadAll(keyArrayFile)

	if err != nil {
		fmt.Println("Error reading keyArray.bin")
		os.Exit(1)
	}

	keyArrayFile.Close()

	//for {
	conn, err := server.Accept()
	if err != nil {
		fmt.Println("Error accepting modem", err.Error())
		os.Exit(1)
	}
	defer server.Close()

	fmt.Println("Got connection")
	//go processConnection(conn)

	//x, err := conn.Write([]byte{0x32, 0x3C, 0x46, 0x27})
	x, err := conn.Write([]byte{0x32, 0x3C, 0x46}) // challenge bytes

	if err != nil {
		fmt.Println("Error with writing to connection");
		os.Exit(1)
	}
	defer server.Close()

	fmt.Println("Challenge key sent! " + string(x))

	contents := make([]byte,256)

	for i := 0x0; i < 256; i++ {
		contents[i] = byte(i)
	}

	time.Sleep(time.Second)

	data := make([]byte, 1024)
	identity := make([]byte, 0)
	buffer := bufio.NewReader(conn)
	//dtFile, err := os.Create("pyro-data.bin")

	if err != nil {
		fmt.Println("Error opening file")
		os.Exit(1)
	}
	defer server.Close()

	idTotal := 0

	//for idTotal < 275 {
	for idTotal < 17 {

		nread, err := buffer.Read(data)

		idTotal += nread

		if err != nil {
			fmt.Println("Error reading from socket: identity")
			os.Exit(1)
		}
		defer server.Close()

		//fmt.Printf("Read data: %d/%d\n", nread, idTotal)
		//fmt.Printf(hex.EncodeToString(data[:nread]) + "\n")
		identity = append(identity, data[:nread]...)
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
	fmt.Printf("Data len: %d\n", pyroDatalen)

	if ((checkByte1 + checkByte2) == 255 && 
		pyroString == "PYROB0" &&
		pyroUserID > 0 &&
		pyroCheckId > 0 &&
		pyroVersion == 2) {
			conn.Write([]byte{0x27}) //validated pyroid
		} else {
		conn.Write([]byte{0x63})
	}

	idTotal = 0
	mmode := make([]byte, 0)

	for idTotal < 2 {
		nread, err := buffer.Read(data)

		idTotal += nread

		if err != nil {
			fmt.Println("Error reading from socket: Mode")
			os.Exit(1)
		}
		defer server.Close()

		mmode = append(mmode, data[:nread]...)
	}

	midCheckIndex := 4
	lastCheckIndex := 4
	fileTotal := 0
	fundata := make([]byte,0)
	doModeLoop := 1

	for doModeLoop == 1 {

		fmt.Printf("Select mode: %d\n", mmode[0])

		switch mmode[0] {
		case 1:
			// Backup data, Get Messages
			fmt.Println("Confirming mode 1...")
			conn.Write([]byte{0x01})

			fmt.Println("ToDo: Implement more here...")
			fmt.Println("Sending 0x64 for 'we're done!'...")

			conn.Write([]byte{0x64, 0x64, 0x64, 0x64})

			// conn.SetDeadline(time.Now().Add(0))
			// _, err := conn.Read(data)

			// if err != nil {
			// 	fmt.Println("Error reading socket")
			// 	os.Exit(1)
			// }
			// defer server.Close()

			time.Sleep(10*time.Second)

			doModeLoop = 0

		case 2, 4:
			fmt.Printf("Confirming mode %d...\n", int(mmode[0]))
			conn.Write(mmode[0:1])

			fmt.Println("Sending ready code (1)...")
			conn.Write([]byte{0x01})

			if int(mmode[0]) == 2 {
				fmt.Println("No contest available, sending 1...")
				conn.Write([]byte{0x1})
			}

			errorCount := 0
			found14 := 0

			for found14 == 0 {

				// fmt.Println("Reading for 0x14...")

				nread, err := conn.Read(data)

				if err != nil {
					fmt.Println("Error reading mode")
					os.Exit(1)
				}
				defer server.Close()

				for i:=0; i<nread; i++ {
					if data[i] == 0x14 {
						fmt.Println("Got 0x14, sending 0x14 and 0x47/0xB8 response...")
						conn.Write([]byte{0x14, 0x47, 0xB8})
						found14 += 1
						break;
					}
				}
			}

			// fundata = append(fundata, data[1:nread]...)

			time.Sleep(time.Second)

			timeout := 10*time.Millisecond
			doGetFileLoop := 1
			modeSelectCount := 0

			for errorCount < 3 && doGetFileLoop == 1 {

				// conn.SetDeadline(time.Now().Add(timeout))
				nread, err := conn.Read(data)

				fileTotal += nread

				if err != nil {
					errorCount += 1
					timeout += 4*time.Second
				} else { 
					errorCount -= 1
					timeout = 10*time.Millisecond; 
				}
				//defer server.Close()

				fundata = append(fundata,data[:nread]...)

				//fmt.Printf("Read data: %d/%d/%d/%d\n", nread, midCheckIndex,lastCheckIndex, len(fundata))
				fmt.Printf(hex.EncodeToString(data[:nread]) + "\n")

				for midCheckIndex <= len(fundata) - 4 {
					if int(fundata[midCheckIndex]) == 1 || int(fundata[midCheckIndex]) == 2 || int(fundata[midCheckIndex]) == 3 {
						if (int(fundata[midCheckIndex]) + int(fundata[midCheckIndex+1]) == 255) || (int(fundata[midCheckIndex+2]) + int(fundata[midCheckIndex+3]) == 255) {
							fmt.Println("Sending Check In 0x06F9")
							conn.Write([]byte{0x06, 0xF9})
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
							conn.Write([]byte{0x06, 0xF9})
							lastCheckIndex += 2
							// timeout = 2*time.Second
							selectMode := -1

							for doGetFileLoop == 1 && errorCount < 4 {

								// conn.SetDeadline(time.Now().Add(timeout))

								nread, err := conn.Read(data)

								if err != nil {
									errorCount += 1
									timeout += 4*time.Second
									fmt.Println("Nothing to read...")
								} else {

									fmt.Printf(hex.EncodeToString(data[:nread]) + "\n")
									lastCheckIndex = 0

									for lastCheckIndex < len(data[:nread]) {
										if int(data[lastCheckIndex]) == 1 || int(data[lastCheckIndex]) == 4 {
											modeSelectCount += 1
											selectMode = int(data[lastCheckIndex])
										} else {
											modeSelectCount = 0
										}
										lastCheckIndex += 1
									}

									if modeSelectCount >= 2 {
										doGetFileLoop = 0
										mmode[0] = byte(selectMode)
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

		case 3:
			fmt.Println("Sending code 3...")
			conn.Write([]byte{0x03})

			fmt.Println("Normal Contest available (0x14)...")
			conn.Write([]byte{0x14})

			// fmt.Println("Sending server ready (1)...")
			// conn.Write([]byte{0x1})

			nread, err := conn.Read(data)

			if err != nil {
				fmt.Println("Error reading mode")
				os.Exit(1)
			}
			defer server.Close()

			if data[0] == 0x14 {
				fmt.Println("Got 0x14")
			} else {
				fmt.Printf("Got this number: %d", int(data[0]))
			}

			fundata = append(fundata, data[:nread]...)

			// contFile, err := os.Open("CONT.000")
			contFile, err := os.Open("CONT.TEST")

			if err != nil {
				fmt.Println("Error opening contest file")
				os.Exit(1)
			}
			defer server.Close()

			contents, err := ioutil.ReadAll(contFile)

			if err != nil {
				fmt.Println("Error reading contest file")
				os.Exit(1)
			}
			defer server.Close()

			//time.Sleep(time.Second)

			//conn.Write(contents)

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
					conn.Write([]byte{0x1, 0xFE})
					chunkSize = 0x80
				} else {
					conn.Write([]byte{0x2, 0xFD})
					chunkSize = 0x400
				}

				conn.Write([]byte{byte(j+1), (0xFF - byte(j+1))})

				for byteCount = 0; byteCount < chunkSize; byteCount++ {

					if i<len(contents) { 
						conn.Write([]byte{contents[i]})
						bx = int(contents[i])   // mov bl, [si-1]
					} else {
						conn.Write([]byte{0})
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

				fmt.Printf("Check Hash: %x\n", int(binary.BigEndian.Uint16(dx)))

				conn.Write([]byte{dx[1], dx[0]})
			}

			conn.Write([]byte{0x4, 0xFB})


			// timeout := 100*time.Second
			doGetStatusLoop := 1
			errorCount := 0
			modeSelectCount := 0
			selectMode := 0

			for doGetStatusLoop == 1 && errorCount < 4 {

				// conn.SetDeadline(time.Now().Add(timeout))

				nread, err := conn.Read(data)

				if err != nil {
					errorCount += 1
					// timeout += 4*time.Second
					fmt.Println("Nothing to read...")
				} else {

					fmt.Printf(hex.EncodeToString(data[:nread]) + "\n")
					lastCheckIndex := 0

					for lastCheckIndex < len(data[:nread]) {
						if int(data[lastCheckIndex]) == 1 || int(data[lastCheckIndex]) == 4 {
							modeSelectCount += 1
							selectMode = int(data[lastCheckIndex])
						} else {
							modeSelectCount = 0
						}
						lastCheckIndex += 1
					}

					if modeSelectCount >= 2 {
						doGetStatusLoop = 0
						mmode[0] = byte(selectMode)
						break
					}
				}
			}


		case 7:
			testTotal := 0
			fmt.Printf("Sending code 7...")
			conn.Write([]byte{0x07})

			for testTotal < 256 {

				nread, err := buffer.Read(data)

				testTotal += nread

				if err != nil {
					fmt.Println("Error reading from socket")
					os.Exit(1)
				}
				defer server.Close()

				//fmt.Printf("Read data: %d/%d\n", nread, testTotal)
				//fmt.Printf(hex.EncodeToString(data[:nread]) + "\n")
			}

			conn.Write([]byte{0x04})

			for i := 0; i < 256; i+=16 {
				conn.Write(contents[i:i+16])
				//fmt.Printf(hex.EncodeToString(contents[i:i+16]) + "\n")
				time.Sleep(50)
			}
			

			fmt.Println("Sent test data!")
			doModeLoop = 0
		}
	}
	fmt.Println("Closing...")

	fmt.Println("Writing fundata to file...")
	f, err := os.Create("T.TMP.1")
	if err != nil {
		fmt.Println("Error writing file")
		os.Exit(1)
	}
	defer server.Close()

	f.Write(fundata)

	server.Close()
	
}

