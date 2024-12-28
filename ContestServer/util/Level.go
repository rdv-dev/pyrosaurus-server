package util

import (
    "encoding/binary"
)

type Level struct {
    X float64
    Y float64
    MaxFoodItems int
    MaxTeamScore int
    MaxFireRange int
    MaxSightRange int
    MaxDinosOnTeam int
    NeckCost int
    SightRangeCost int
    HearingRangeCost int
    SmellCost int
    LegCostMultiplier int
    LegFootCostMult int
    HeartSizeMult int
    BaseEndurance int
    BaseNeckSizeCost int
    FireHealth1 int
    FireHealth2 int
    FireHealth3 int
    FireHealth4 int
    Unknown1 int
    Unknown2 int
    BaseCost int
    Unused1 int
    Unused2 int
    Unused3 int
    Unused4 int
    MaxTime int
    RequiredQueens int
    EnableQuads int
    EnablePacking int
    EnableSightIncrease int
    MiniMapX int
    MiniMapY int
    RawData []byte
}

func NewLevel(leveldata []byte) *Level {
    return &Level {
        X: float64(binary.LittleEndian.Uint16(leveldata[0:2]))/2,
        Y: float64(binary.LittleEndian.Uint16(leveldata[2:4]))/2,
        MaxTime: int(binary.LittleEndian.Uint16(leveldata[56:58])),
        RawData: leveldata,
    }
}
