package main

import (
	"os"
	"fmt"
	"io/ioutil"
)

func main() {
	rawFile, err := os.Open("../../../Test Teams/Call/T1.TEAM")

	if err != nil {
		fmt.Println("Error opening team entry file")
		os.Exit(1)
	}

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
	outData := make([]byte, 0)
	fileNum := 1
	outFileName := ""

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
			readPos += 4

			if int(rawData[readPos + 2]) == chunkNum {
				
				chunkNum += 1
				toReadPos = (readPos + chunkSize)
				
				outData = append(outData, rawData[readPos:toReadPos]...)

				readPos = toReadPos + 2

			}

			if int(rawData[readPos + 2]) < chunkNum {
				// repeated chunk, skip it
				toReadPos = (readPos + chunkSize) + 2
				readPos = toReadPos
			}

			if int(rawData[readPos + 2]) > chunkNum {
				// invalid file?
				fmt.Println("File error, chunkNum less than file chunk number")
				os.Exit(1)
			}			

			if rawData[readPos] == 0x4 && (rawData[readPos] + rawData[readPos + 1]) == 0xFF {
				// end of file

				chunkNum = 1

				if fileNum == 1 {
					outFileName = "TEAM.BIN"
				} else {
					outFileName = fmt.Sprintf("MSG-%d.bin", fileNum-1)
				}

				fileNum += 1

				outFile, err := os.Create(outFileName)

				if err != nil {
					fmt.Println("Error opening out file")
					os.Exit(1)
				}

				outFile.Write(outData)
				outData = make([]byte, 0)

				if readPos + 3 >= fileLen {
					doRead = 0
				} else {
					readPos += 2
				}
			}
		}
	}
}