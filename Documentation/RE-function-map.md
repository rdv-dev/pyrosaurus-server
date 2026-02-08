# Pyrosaurus Function Map

## PYRO Executable Notable Functions
* callModemMain - Entry point for all modem functionality
* RunContestFileHandle - Entry point for running a selected Contest
* calculateSpecPoints - function for calculating species points
* contestReadAllocate - the program reads a set number amount of the contest based on the memory available, this is occasionally called if more of the contest needs to be read
* **contestReadFrame** - this function reads each frame, decodes it, and calls the specified animation function on the selected dino
* pyroMain - the main function for the game
* j_mainGameLoop - the main game loop
* loadLevelFile - loads the level data
* speciesDefault - this determines the default values when a species is first created
* writeContestEntryFile - this creates the Contest Entry file which is sent on via the MODEM
* dinoMoveNeck - corresponds to Contest Animation 0
* dinoMoveTail - corresponds to Contest Animation 1
* doMove - corresponds to Contest Animation 2
* setBreathRate - corresponds to Contest Animation 3
* doStepSide - corresponds to Contest Animation 4
* doStepForBak - corresponds to Contest Animation 5
* dinoAction_6 - corresponds to Contest Animation 6
* doJumpSide - corresponds to Contest Animation 7
* doJumpForBak - corresponds to Contest Animation 8
* doCall - corresponds to Contest Animation 10
* doEatFood - corresponds to Contest Special Animation 7
* doFire - correspnds to Contest Special Animation 8

## MODEM Executable Notable Functions
* modemDial - dials the IP address
* modemChallenge - the challenge procedure for authenticating users
* modemGetSend - determines whether modem is sending or getting data
* getUsrDBackup - poorly named function which handles various functions based on the servers input
* modemTest - the test procedure for the modem
* doCheckSum - creates the checksum when sending data
* calcCheckSum - calculates the checksum from data
