package util

import (
	"encoding/base64"
	"encoding/binary"
	"math"
)

// cosTable stores cos(i degrees) * 65536 for i=0..90 (16.16 fixed-point).
// Matches the original game's 91-entry sinCos lookup table.
var cosTable [91]int32

// Original binary sinCos lookup table encoded in base64
var cosTableBase64 = "AAABAPb/AADY/wAApv8AAGD/AAAH/wAAmf4AABj+AACC/QAA2fwAABz8AABM+wAAaPoAAHD5AABl+AAAR/cAABX2AADQ9AAAePMAAA7yAACQ8AAA/+4AAFztAACm6wAA3ukAAAToAAAX5gAAGeQAAAniAADn3wAAtN0AAG/bAAAa2QAAs9YAADzUAAC00QAAHM8AAHPMAAC7yQAA88YAABvEAAA1wQAAP74AADq7AAAnuAAABbUAANWxAACXrgAATKsAAPOnAACOpAAAG6EAAJydAAARmgAAeZYAANaSAAAnjwAAbYsAAKmHAADagwAAAIAAABx8AAAveAAAOXQAADlwAAAxbAAAIGgAAAdkAADmXwAAvlsAAI9XAABYUwAAHE8AANlKAACQRgAAQkIAAO89AACWOQAAOjUAANkwAAB0LAAADCgAAKEjAAAzHwAAwhoAAFAWAADcEQAAZg0AAO8IAAB4BAAAAAAAAA=="

func init() {
	data, err := base64.StdEncoding.DecodeString(cosTableBase64)
	if err != nil {
		panic("failed to decode cosTable: " + err.Error())
	}
	for i := 0; i <= 90; i++ {
		cosTable[i] = int32(binary.LittleEndian.Uint32(data[i*4 : i*4+4]))
	}
}

// NormalizeAngle wraps an angle to [0, 360] range using the original game's
// while-loop approach.
func NormalizeAngle(angle int16) int16 {
	for angle < 0 {
		angle += 360
	}
	for angle > 360 {
		angle -= 360
	}
	return angle
}

// lookupCos returns cos(angle degrees) in 16.16 fixed-point using the
// 91-entry table with 4-quadrant symmetry.
func lookupCos(angle int16) int32 {
	angle = NormalizeAngle(angle)
	if angle == 360 {
		angle = 0
	}
	switch {
	case angle <= 90:
		return cosTable[angle]
	case angle <= 180:
		return -cosTable[180-angle]
	case angle <= 270:
		return -cosTable[angle-180]
	default:
		return cosTable[360-angle]
	}
}

// SinCos returns (cos, sin) of the given angle in 16.16 fixed-point.
// Uses 4-quadrant symmetry: sin(x) = cos(90 - x).
func SinCos(angle int16) (int32, int32) {
	return lookupCos(angle), lookupCos(90 - angle)
}

// FixedMul multiplies two 16.16 fixed-point values with rounding.
// Matches the original game's shrd-based multiplication.
func FixedMul(a, b int32) int32 {
	return int32((int64(a)*int64(b) + 0x8000) >> 16)
}

// FixedToInt16 converts a 16.16 fixed-point value to int16 with rounding.
// Uses division (not shift) to match the original game's div32 truncation
// toward zero for negative values.
func FixedToInt16(val int32) int16 {
	if val >= 0 {
		return int16((val + 0x8000) / 0x10000)
	}
	return int16((val - 0x8000) / 0x10000)
}

// RotateByHeading rotates a (dx, dy) vector by heading degrees.
// Extends int16 inputs to 16.16 fixed-point, applies rotation matrix
// [cos, -sin; sin, cos], converts back via FixedToInt16.
// Each product is rounded individually before summing (matching original).
func RotateByHeading(heading int16, dx, dy *int16) {
	cosH, sinH := SinCos(heading)

	fx := int32(*dx) << 16
	fy := int32(*dy) << 16

	newX := FixedMul(fx, cosH) - FixedMul(fy, sinH)
	newY := FixedMul(fx, sinH) + FixedMul(fy, cosH)

	*dx = FixedToInt16(newX)
	*dy = FixedToInt16(newY)
}

// Atan2Degrees computes the angle in integer degrees from (0,0) to (dx, dy).
// Returns a value in [0, 360]. This is a placeholder until the exact
// original game implementation (sub_2E56C) is reverse-engineered.
func Atan2Degrees(dy, dx int16) int16 {
	if dx == 0 && dy == 0 {
		return 0
	}
	rad := math.Atan2(float64(dy), float64(dx))
	deg := int16(math.Round(rad * 180.0 / math.Pi))
	return NormalizeAngle(deg)
}

// AbsInt16 returns the absolute value of an int16.
func AbsInt16(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}
