# Pyrosaurus Function Map

## PYRO Executable Notable Functions
* callModemMain - Entry point for all modem functionality
* RunContestFileHandle - Entry point for running a selected Contest
* calculateSpecPoints - function for calculating species points
* contestReadAllocate - the program reads a set number amount of the contest based on the memory available, this is occasionally called if more of the contest needs to be read
* **contestReadFrame** - this function reads each frame, decodes it, and calls the specified function on the selected dino
* pyroMain - the main function for the game
* j_mainGameLoop - the main game loop
* loadLevelFile - loads the level data
* speciesDefault - this determines the default values when a species is first created
* writeContestEntryFile - this creates the Contest Entry file which is sent on via the MODEM
* dinoMoveNeck - corresponds to Contest Action 0
* dinoMoveTail - corresponds to Contest Action 1
* doMove - corresponds to Contest Action 2
* setBreathRate - corresponds to Contest Action 3
* doStepSide - corresponds to Contest Action 4
* doStepForBak - corresponds to Contest Action 5
* dinoAction_6 - corresponds to Contest Action 6
* doJumpSide - corresponds to Contest Action 7
* doJumpForBak - corresponds to Contest Action 8
* doCall - corresponds to Contest Action 10
* doEatFood - corresponds to Contest Special Action 7
* doFire - correspnds to Contest Special Action 8

## MODEM Executable Notable Functions
* modemDial - dials the IP address
* modemChallenge - the challenge procedure for authenticating users
* modemGetSend - determines whether modem is sending or getting data
* getUsrDBackup - poorly named function which handles various functions based on the servers input
* modemTest - the test procedure for the modem
* doCheckSum - creates the checksum when sending data
* calcCheckSum - calculates the checksum from data
