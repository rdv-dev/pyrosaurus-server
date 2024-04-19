package util

import (
//	"fmt"
	"encoding/binary"
)

const (
	// DECISIONS_LEN = 0x17D

	DECISION_SIZE = 19

	DECISION_MOVEMENT = 0
	DECISION_TARGET = 1
	DECISION_LEGS = 2
	DECISION_THEIR_SIZE = 3
	DECISION_RANGE = 4
	DECISION_THEIR_SKIN = 5
	DECISION_MY_SKIN = 6
	DECISION_MY_CONDITION = 7
	DECISION_MY_QUEEN_ENEMY_RANGE = 8
	DECISION_MY_QUEEN_RANGE = 9
	DECISION_ENEMY_QUEEN_RANGE = 10
	DECISION_THEIR_SPEED = 11
	DECISION_THEIR_ACTION = 12
	DECISION_CALLING = 13
	DECISION_TIME = 14
	DECISION_PRIORITY = 15
	DECISION_FOOD = 16
	DECISION_GO_SPEED = 17
	DECISION_PACK = 18

    DECISION_DONT_MOVE = byte(0)
    DECISION_CREEP = byte(1)
    DECISION_WALK = byte(2)
    DECISION_RUN = byte(3)
)

type Decision struct {
	Movement int 
	Target byte
	Legs byte
	Size byte
	InRange byte
	TheirSkin byte
	MySkin byte
	MyCondition byte
	MyQueenEnemyRange byte
	MyQueenRange byte
	EnemyQueenRange byte
	TheirSpeed byte
	TheirAction byte
	Calling byte
	Time byte
	Priority byte
	Food byte
	GoSpeed byte
	Pack byte
	Score int
}

func NewDecisions(decisions []byte) []*Decision {
	retDeci := make([]*Decision, int(decisions[0]))	

	for i:=0; i<int(decisions[0]); i++ {
		//fmt.Printf("Decision Index Movement: %d\n", int(decisions[1 + (i * DECISION_SIZE) + DECISION_MOVEMENT]))

		retDeci[i] = &Decision {
			Movement: int(decisions[1 + (i * DECISION_SIZE) + DECISION_MOVEMENT]),
			Target: decisions[1 + (i * DECISION_SIZE) + DECISION_TARGET],
			Legs: decisions[1 + (i * DECISION_SIZE) + DECISION_LEGS],
			Size: decisions[1 + (i * DECISION_SIZE) + DECISION_THEIR_SIZE],
			InRange: decisions[1 + (i * DECISION_SIZE) + DECISION_RANGE],
			TheirSkin: decisions[1 + (i * DECISION_SIZE) + DECISION_THEIR_SKIN],
			MySkin: decisions[1 + (i * DECISION_SIZE) + DECISION_MY_SKIN],
			MyCondition: decisions[1 + (i * DECISION_SIZE) + DECISION_MY_CONDITION],
			MyQueenEnemyRange: decisions[1 + (i * DECISION_SIZE) + DECISION_MY_QUEEN_ENEMY_RANGE],
			MyQueenRange: decisions[1 + (i * DECISION_SIZE) + DECISION_MY_QUEEN_RANGE],
			EnemyQueenRange: decisions[1 + (i * DECISION_SIZE) + DECISION_ENEMY_QUEEN_RANGE],
			TheirSpeed: decisions[1 + (i * DECISION_SIZE) + DECISION_THEIR_SPEED],
			TheirAction: decisions[1 + (i * DECISION_SIZE) + DECISION_THEIR_ACTION],
			Calling: decisions[1 + (i * DECISION_SIZE) + DECISION_CALLING],
			Time: decisions[1 + (i * DECISION_SIZE) + DECISION_TIME],
			Priority: decisions[1 + (i * DECISION_SIZE) + DECISION_PRIORITY],
			Food: decisions[1 + (i * DECISION_SIZE) + DECISION_FOOD],
			GoSpeed: decisions[1 + (i * DECISION_SIZE) + DECISION_GO_SPEED],
			Pack: decisions[1 + (i * DECISION_SIZE) + DECISION_PACK],
			Score: 0}
	}

	return retDeci
}

type Dino struct {
	Team uint32
	species []byte
    HasLegs bool
	Moves []*Moves
	DoMove *Moves
	fitePoints []byte
	fiteXPos []byte
	fiteYPos []byte
	Decisions []*Decision
	dino []byte
	name []byte
	Xpos float64
	Ypos float64
	Angle float64
	neckAngle int
	attacking []byte
	attackedBy []byte
	
}

func NewDino(inTeam *ContestEntry, species, dino int, xBound, yBound float64) *Dino {
	//fmt.Printf("Moves Offset: %d Species: %d\n", inTeam.MovesOffset, species)
	movesStart := inTeam.MovesOffset + (MOVE_DATA_LEN * species)
	movesEnd := inTeam.MovesOffset + (MOVE_DATA_LEN * species) + MOVE_DATA_LEN

	//fmt.Printf("Moves Start offset: %d Moves End Offset: %d Diff: %d\n", movesStart, movesEnd, movesEnd-movesStart)

	decisionStart := inTeam.DecisionsOffset + (DECISIONS_LEN * species)
	decisionEnd := inTeam.DecisionsOffset + (DECISIONS_LEN * species) + DECISIONS_LEN

	//fmt.Printf("Decision Start offset: %d Decision End Offset: %d Diff: %d\n", decisionStart, decisionEnd, decisionEnd-decisionStart)

	dinoXPosIndex := NUM_DINOS_ON_TEAM_LEN + (inTeam.NumDinos * (TEAM_QUEEN_ARRAY_LEN + TEAM_SPECIES_LEG_NUM_LEN + TEAM_DINO_RESIZE)) + (dino * TEAM_X_POS_LEN) + decisionEnd
	dinoYPosIndex := NUM_DINOS_ON_TEAM_LEN + (inTeam.NumDinos * (TEAM_QUEEN_ARRAY_LEN + TEAM_SPECIES_LEG_NUM_LEN + TEAM_DINO_RESIZE + TEAM_X_POS_LEN)) + (dino * TEAM_Y_POS_LEN) + decisionEnd
	dinoRotnIndex := NUM_DINOS_ON_TEAM_LEN + (inTeam.NumDinos * (TEAM_QUEEN_ARRAY_LEN + TEAM_SPECIES_LEG_NUM_LEN + TEAM_DINO_RESIZE + TEAM_X_POS_LEN + TEAM_Y_POS_LEN)) + (dino * TEAM_ROT_LEN) + decisionEnd

	//fmt.Printf("Dino Offsets XPos: %d YPos: %d Rotn: %d\n", dinoXPosIndex, dinoYPosIndex, dinoRotnIndex)
	//fmt.Printf("Dino Values  XPos: %f YPos: %f Rotn: %f\n",float64(int16(binary.LittleEndian.Uint16(inTeam.TeamData[dinoXPosIndex:dinoXPosIndex+2]))),float64(int16(binary.LittleEndian.Uint16(inTeam.TeamData[dinoYPosIndex:dinoYPosIndex+2]))),float64(int16(binary.LittleEndian.Uint16(inTeam.TeamData[dinoRotnIndex:dinoRotnIndex+2]))%360))

	return &Dino {
		Team: inTeam.Team,
		species: make([]byte, SPECIES_LEN),
        HasLegs: true,
		Moves: NewMoves(inTeam.TeamData[movesStart:movesEnd], xBound, yBound),
		fitePoints: make([]byte, FITE_DATA1_LEN),
		fiteXPos: make([]byte, FITE_DATA2_LEN),
		fiteYPos: make([]byte, FITE_DATA3_LEN),
		Decisions: NewDecisions(inTeam.TeamData[decisionStart:decisionEnd]),
		dino: make([]byte, TEAM_ENTRY_RECORD_LEN),
		name: make([]byte, 50),
		Xpos: float64(int16(binary.LittleEndian.Uint16(inTeam.TeamData[dinoXPosIndex:dinoXPosIndex+2]))),
		Ypos: float64(int16(binary.LittleEndian.Uint16(inTeam.TeamData[dinoYPosIndex:dinoYPosIndex+2]))),
		Angle: float64(int16(binary.LittleEndian.Uint16(inTeam.TeamData[dinoRotnIndex:dinoRotnIndex+2]))%360),
		neckAngle: 0,
		attacking: make([]byte, 0),
		attackedBy: make([]byte, 0)}
}
