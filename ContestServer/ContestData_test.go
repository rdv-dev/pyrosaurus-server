package ContestServer

import (
	"encoding/binary"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// animationArgCount returns the number of argument bytes for a given animation code.
// Animation 11 (special) has variable args handled separately.
var animationArgCount = [12]int{
	2, // 0: Set Neck Angle
	2, // 1: Set Tail Angle
	3, // 2: Move Dino
	2, // 3: Set Breath Rate
	1, // 4: Step Left/Right
	1, // 5: Step Forward/Back
	1, // 6: Dino Die
	1, // 7: Jump Left/Right
	1, // 8: Jump Forward/Back
	0, // 9: Lock neck movement
	0, // 10: Call
	1, // 11: Special Animations (base: 1 for the sub-action byte)
}

// specialAnimationExtraArgs returns additional argument bytes for special animation sub-actions.
func specialAnimationExtraArgs(subAction byte) int {
	if subAction == 7 { // Eat Food takes 1 extra argument
		return 1
	}
	return 0
}

// ContestDataValidation holds validation results for a contest data stream.
type ContestDataValidation struct {
	TotalFrames      int
	TotalAnimations  int
	TotalDelayTicks  int
	HasTerminator    bool
	MaxDinoIndex     int
	Errors           []string
}

// validateContestData walks the contest data byte stream and checks structural integrity.
// It simulates how contestReadFrame in the original game parses the data.
func validateContestData(data []byte) *ContestDataValidation {
	v := &ContestDataValidation{
		MaxDinoIndex: -1,
	}

	pos := 0

	for pos < len(data) {
		frameHeader := data[pos]
		pos++

		// Signed interpretation of the frame header byte
		signedHeader := int8(frameHeader)

		if signedHeader == 0 {
			// Terminator: end of contest
			v.HasTerminator = true
			break
		}

		if signedHeader < 0 {
			// Delay frame: magnitude is number of ticks to skip
			// Use -int() to avoid int8 overflow when signedHeader is -128
			delayTicks := -int(signedHeader)
			v.TotalDelayTicks += delayTicks
			v.TotalFrames++
			continue
		}

		// Positive: number of animation entries in this frame
		numAnimations := int(signedHeader)
		v.TotalFrames++

		if numAnimations > 127 {
			v.Errors = append(v.Errors, fmt.Sprintf("frame %d: numAnimations %d exceeds max 127", v.TotalFrames, numAnimations))
			break
		}

		for i := 0; i < numAnimations; i++ {
			if pos >= len(data) {
				v.Errors = append(v.Errors, fmt.Sprintf("frame %d, animation %d: unexpected end of data while reading encoded byte", v.TotalFrames, i))
				return v
			}

			encodedByte := data[pos]
			pos++

			dinoIndex := int(encodedByte) / 12
			animCode := int(encodedByte) - (dinoIndex * 12)

			if dinoIndex > 19 {
				v.Errors = append(v.Errors, fmt.Sprintf("frame %d, animation %d: dino index %d exceeds max 19 (encoded byte 0x%02X)", v.TotalFrames, i, dinoIndex, encodedByte))
			}

			if animCode > 11 {
				v.Errors = append(v.Errors, fmt.Sprintf("frame %d, animation %d: animation code %d exceeds max 11 (encoded byte 0x%02X)", v.TotalFrames, i, animCode, encodedByte))
				break
			}

			if dinoIndex > v.MaxDinoIndex {
				v.MaxDinoIndex = dinoIndex
			}

			// Determine how many argument bytes to consume
			argCount := animationArgCount[animCode]

			// Special animation (code 11) has variable args
			if animCode == 11 {
				if pos >= len(data) {
					v.Errors = append(v.Errors, fmt.Sprintf("frame %d, animation %d: unexpected end of data reading special animation sub-action", v.TotalFrames, i))
					return v
				}
				subAction := data[pos]
				pos++ // consumed the sub-action byte (already counted in argCount=1)

				extra := specialAnimationExtraArgs(subAction)
				if pos+extra > len(data) {
					v.Errors = append(v.Errors, fmt.Sprintf("frame %d, animation %d: unexpected end of data reading special animation %d extra args", v.TotalFrames, i, subAction))
					return v
				}
				pos += extra
			} else {
				if pos+argCount > len(data) {
					v.Errors = append(v.Errors, fmt.Sprintf("frame %d, animation %d: unexpected end of data reading %d args for animation %d", v.TotalFrames, i, argCount, animCode))
					return v
				}
				pos += argCount
			}

			v.TotalAnimations++
		}
	}

	if !v.HasTerminator && len(v.Errors) == 0 {
		v.Errors = append(v.Errors, "contest data has no 0x00 terminator byte")
	}

	return v
}

// TestContestFileIntegrity finds all CONT.* files under TestData/ recursively
// and validates the contest data structure within each file.
func TestContestFileIntegrity(t *testing.T) {
	directory, err := os.Getwd()
	if err != nil {
		t.Fatalf("Unable to get current directory: %v", err)
	}

	testDataDir := filepath.Join(directory, "TestData")

	var contFiles []string

	err = filepath.WalkDir(testDataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasPrefix(d.Name(), "CONT.") {
			contFiles = append(contFiles, path)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Error walking TestData directory: %v", err)
	}

	if len(contFiles) == 0 {
		t.Skip("No CONT.* files found under TestData/")
	}

	for _, contFile := range contFiles {
		relPath, _ := filepath.Rel(testDataDir, contFile)
		t.Run(relPath, func(t *testing.T) {
			data, err := os.ReadFile(contFile)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			// Validate header
			if len(data) < CONTEST_HEADER_RECORD_LEN {
				t.Fatalf("File too small for contest header: %d bytes", len(data))
			}

			statusByte := data[0]
			if statusByte == 0 {
				t.Errorf("Contest status byte is 0x00 (game requires non-zero to end properly)")
			}
			// Lower 2 bits encode result: 1=team1 won, 2=team2 won, 3=draw
			// Bit 7 (0x80) is the "not watched" flag
			resultBits := statusByte & 0x03
			watched := statusByte & 0x80
			t.Logf("Status breakdown: Watched=0x%02X, result=%d", watched, resultBits)

			// Read all offsets from header
			team1ColorsOffset := int(binary.LittleEndian.Uint16(data[1:3]))
			team1DataOffset := int(binary.LittleEndian.Uint16(data[3:5]))
			// bytes 5-6 are zeros
			team2ColorsOffset := int(binary.LittleEndian.Uint16(data[7:9]))
			team2DataOffset := int(binary.LittleEndian.Uint16(data[9:0xB]))
			// bytes 0xB-0xC are zeros
			levelDataOffset := int(binary.LittleEndian.Uint16(data[0xD:0xF]))
			contestDataOffset := int(binary.LittleEndian.Uint16(data[0xF:0xF+2]))

			fileSize := len(data)

			t.Logf("File size: %d bytes", fileSize)
			t.Logf("Status: 0x%02X", statusByte)
			t.Logf("Offsets - Team1Colors: %d, Team1Data: %d, Team2Colors: %d, Team2Data: %d, Level: %d, ContestData: %d",
				team1ColorsOffset, team1DataOffset, team2ColorsOffset, team2DataOffset, levelDataOffset, contestDataOffset)

			// Validate offsets are within file bounds
			offsets := map[string]int{
				"team1Colors": team1ColorsOffset,
				"team1Data":   team1DataOffset,
				"team2Colors": team2ColorsOffset,
				"team2Data":   team2DataOffset,
				"levelData":   levelDataOffset,
				"contestData": contestDataOffset,
			}

			for name, offset := range offsets {
				if offset >= fileSize {
					t.Errorf("Offset %s (%d) exceeds file size (%d)", name, offset, fileSize)
				}
				if offset < CONTEST_HEADER_RECORD_LEN {
					t.Errorf("Offset %s (%d) points inside the header (< %d)", name, offset, CONTEST_HEADER_RECORD_LEN)
				}
			}

			// Validate contest data stream
			if contestDataOffset >= fileSize {
				t.Fatalf("Contest data offset %d exceeds file size %d", contestDataOffset, fileSize)
			}

			contestData := data[contestDataOffset:]
			v := validateContestData(contestData)

			t.Logf("Frames: %d, Animations: %d, Delay ticks: %d, Max dino index: %d, Has terminator: %v",
				v.TotalFrames, v.TotalAnimations, v.TotalDelayTicks, v.MaxDinoIndex, v.HasTerminator)

			for _, e := range v.Errors {
				t.Errorf("Contest data error: %s", e)
			}

			if !v.HasTerminator {
				t.Errorf("Contest data missing 0x00 terminator byte")
			}

			if v.TotalFrames == 0 && v.HasTerminator {
				t.Logf("Warning: contest data contains only a terminator (0 frames)")
			}
		})
	}
}

// TestContestDataStructure validates specific properties of the contest data encoding
// using hand-crafted byte sequences.
func TestContestDataStructure(t *testing.T) {

	t.Run("terminator only", func(t *testing.T) {
		data := []byte{0x00}
		v := validateContestData(data)
		if !v.HasTerminator {
			t.Error("Expected terminator")
		}
		if v.TotalFrames != 0 {
			t.Errorf("Expected 0 frames, got %d", v.TotalFrames)
		}
	})

	t.Run("single frame with call animation", func(t *testing.T) {
		// Frame: 1 animation, encoded byte for dino 0 animation 10 (Call, 0 args)
		// dino=0, anim=10 → encoded = 0*12+10 = 10
		data := []byte{
			0x01,       // 1 animation in this frame
			0x0A,       // encoded: dino 0, animation 10 (Call)
			0x00,       // terminator
		}
		v := validateContestData(data)
		if !v.HasTerminator {
			t.Error("Expected terminator")
		}
		if v.TotalFrames != 1 {
			t.Errorf("Expected 1 frame, got %d", v.TotalFrames)
		}
		if v.TotalAnimations != 1 {
			t.Errorf("Expected 1 animation, got %d", v.TotalAnimations)
		}
		if len(v.Errors) > 0 {
			t.Errorf("Unexpected errors: %v", v.Errors)
		}
	})

	t.Run("frame with move dino (3 args)", func(t *testing.T) {
		// dino=1, anim=2 (Move) → encoded = 1*12+2 = 14
		data := []byte{
			0x01,       // 1 animation
			0x0E,       // encoded: dino 1, animation 2 (Move Dino)
			0x05,       // arg1: heading
			0x91,       // arg2: flags
			0x00,       // arg3: speed
			0x00,       // terminator
		}
		v := validateContestData(data)
		if !v.HasTerminator {
			t.Error("Expected terminator")
		}
		if v.TotalAnimations != 1 {
			t.Errorf("Expected 1 animation, got %d", v.TotalAnimations)
		}
		if v.MaxDinoIndex != 1 {
			t.Errorf("Expected max dino index 1, got %d", v.MaxDinoIndex)
		}
		if len(v.Errors) > 0 {
			t.Errorf("Unexpected errors: %v", v.Errors)
		}
	})

	t.Run("delay frames", func(t *testing.T) {
		// 0xFF = -1 signed → delay 1 tick
		// 0xFB = -5 signed → delay 5 ticks
		// 0x80 = -128 signed → delay 128 ticks
		data := []byte{
			0xFF,       // delay 1 tick
			0xFB,       // delay 5 ticks
			0x80,       // delay 128 ticks
			0x00,       // terminator
		}
		v := validateContestData(data)
		if !v.HasTerminator {
			t.Error("Expected terminator")
		}
		if v.TotalDelayTicks != 134 { // 1 + 5 + 128
			t.Errorf("Expected 134 delay ticks, got %d", v.TotalDelayTicks)
		}
		if v.TotalFrames != 3 {
			t.Errorf("Expected 3 delay frames, got %d", v.TotalFrames)
		}
		if len(v.Errors) > 0 {
			t.Errorf("Unexpected errors: %v", v.Errors)
		}
	})

	t.Run("special animation eat food", func(t *testing.T) {
		// dino=2, anim=11 (Special) → encoded = 2*12+11 = 35
		data := []byte{
			0x01,       // 1 animation
			0x23,       // encoded: dino 2, animation 11 (Special)
			0x07,       // sub-action 7 (Eat Food)
			0x40,       // extra arg for eat food
			0x00,       // terminator
		}
		v := validateContestData(data)
		if !v.HasTerminator {
			t.Error("Expected terminator")
		}
		if v.TotalAnimations != 1 {
			t.Errorf("Expected 1 animation, got %d", v.TotalAnimations)
		}
		if len(v.Errors) > 0 {
			t.Errorf("Unexpected errors: %v", v.Errors)
		}
	})

	t.Run("special animation fire (no extra args)", func(t *testing.T) {
		// dino=0, anim=11 → encoded = 11
		data := []byte{
			0x01,       // 1 animation
			0x0B,       // encoded: dino 0, animation 11 (Special)
			0x08,       // sub-action 8 (Fire)
			0x00,       // terminator
		}
		v := validateContestData(data)
		if !v.HasTerminator {
			t.Error("Expected terminator")
		}
		if v.TotalAnimations != 1 {
			t.Errorf("Expected 1 animation, got %d", v.TotalAnimations)
		}
		if len(v.Errors) > 0 {
			t.Errorf("Unexpected errors: %v", v.Errors)
		}
	})

	t.Run("missing terminator", func(t *testing.T) {
		// dino=0, anim=10 (Call) → encoded = 10
		data := []byte{
			0x01,       // 1 animation
			0x0A,       // Call
			// no terminator!
		}
		v := validateContestData(data)
		if v.HasTerminator {
			t.Error("Should not have terminator")
		}
		if len(v.Errors) == 0 {
			t.Error("Expected error about missing terminator")
		}
	})

	t.Run("truncated animation args", func(t *testing.T) {
		// dino=0, anim=2 (Move, needs 3 args) but only 1 arg provided
		data := []byte{
			0x01,       // 1 animation
			0x02,       // Move Dino (needs 3 args)
			0x05,       // only 1 arg
			// truncated
		}
		v := validateContestData(data)
		if len(v.Errors) == 0 {
			t.Error("Expected error about truncated data")
		}
	})

	t.Run("multiple animations per frame", func(t *testing.T) {
		// 3 animations in one frame:
		// dino 0 neck (2 args) + dino 0 tail (2 args) + dino 0 call (0 args)
		data := []byte{
			0x03,       // 3 animations
			0x00,       // dino 0, anim 0 (Neck)
			0x11, 0x1E, // neck args
			0x01,       // dino 0, anim 1 (Tail)
			0x11, 0x1E, // tail args
			0x0A,       // dino 0, anim 10 (Call)
			0x00,       // terminator
		}
		v := validateContestData(data)
		if !v.HasTerminator {
			t.Error("Expected terminator")
		}
		if v.TotalAnimations != 3 {
			t.Errorf("Expected 3 animations, got %d", v.TotalAnimations)
		}
		if v.TotalFrames != 1 {
			t.Errorf("Expected 1 frame, got %d", v.TotalFrames)
		}
		if len(v.Errors) > 0 {
			t.Errorf("Unexpected errors: %v", v.Errors)
		}
	})

	t.Run("encoding roundtrip", func(t *testing.T) {
		// Verify encoding: encoded = dino * 12 + animation
		// Then decode: dino = encoded / 12, animation = encoded - dino * 12
		for dino := 0; dino <= 19; dino++ {
			for anim := 0; anim <= 11; anim++ {
				encoded := byte(dino*12 + anim)
				decodedDino := int(encoded) / 12
				decodedAnim := int(encoded) - (decodedDino * 12)

				if decodedDino != dino || decodedAnim != anim {
					t.Errorf("Encoding roundtrip failed: dino=%d anim=%d → encoded=%d → dino=%d anim=%d",
						dino, anim, encoded, decodedDino, decodedAnim)
				}
			}
		}
	})

	t.Run("max encoded byte boundary", func(t *testing.T) {
		// dino=19, anim=11 → 19*12+11 = 239 (0xEF), fits in a byte
		// dino=20, anim=0 → 20*12 = 240 (0xF0), still fits but invalid dino
		encoded := byte(19*12 + 11)
		if encoded != 0xEF {
			t.Errorf("Max valid encoded byte should be 0xEF, got 0x%02X", encoded)
		}
	})
}

// AnimationEntry holds the fully decoded details of a single animation command,
// matching how contestReadFrame parses and dispatches each animation.
type AnimationEntry struct {
	Frame       int    // which frame this animation belongs to
	DinoIndex   int    // decoded dino index (0-19)
	AnimCode    int    // decoded animation code (0-11)
	EncodedByte byte   // raw encoded byte
	RawArgs     []byte // raw argument bytes as they appear in the stream

	// Move Dino (animation 2) decoded fields — mirrors contestReadFrame dispatch to doMove
	MoveHeading  int8 // ax: signed heading change (arg1, sign-extended via cbw)
	MoveBX       int  // bx/di: arg2 & 0x0F — movement mode (1=creep, 4=run, 0xA=walk per doc)
	MoveDX       int  // dx: (arg2 & 0x70) >> 4 — stop(0) / moving(1+)
	IsRunningSpeed  int  // stack arg_0: arg2 >> 7 — whether dino is at running speed during either running or walking
	MoveMoveCode byte // ds:65b8[dinoIndex]: arg3 — stored to movementAnimation array
}

// FrameEntry represents a single frame in the contest data stream.
type FrameEntry struct {
	FrameNum   int
	IsDelay    bool
	DelayTicks int
	Animations []AnimationEntry
}

var animationNames = [12]string{
	"Set Neck Angle",
	"Set Tail Angle",
	"Move Dino",
	"Set Breath Rate",
	"Step Left/Right",
	"Step Forward/Back",
	"Dino Die",
	"Jump Left/Right",
	"Jump Forward/Back",
	"Lock Neck",
	"Call",
	"Special",
}

// parseContestAnimations parses contest data and returns detailed animation entries.
// This mirrors the parsing logic of contestReadFrame in the original game binary.
func parseContestAnimations(data []byte) ([]FrameEntry, []string) {
	var frames []FrameEntry
	var errors []string
	pos := 0
	frameNum := 0

	for pos < len(data) {
		frameHeader := data[pos]
		pos++
		signedHeader := int8(frameHeader)

		if signedHeader == 0 {
			// Terminator
			break
		}

		frameNum++

		if signedHeader < 0 {
			delayTicks := -int(signedHeader)
			frames = append(frames, FrameEntry{
				FrameNum:   frameNum,
				IsDelay:    true,
				DelayTicks: delayTicks,
			})
			continue
		}

		numAnimations := int(signedHeader)
		frame := FrameEntry{
			FrameNum: frameNum,
			IsDelay:  false,
		}

		for i := 0; i < numAnimations; i++ {
			if pos >= len(data) {
				errors = append(errors, fmt.Sprintf("frame %d: truncated at encoded byte", frameNum))
				return frames, errors
			}

			encodedByte := data[pos]
			pos++

			dinoIndex := int(encodedByte) / 12
			animCode := int(encodedByte) - (dinoIndex * 12)

			entry := AnimationEntry{
				Frame:       frameNum,
				DinoIndex:   dinoIndex,
				AnimCode:    animCode,
				EncodedByte: encodedByte,
			}

			// Read argument bytes based on animation code, matching contestReadFrame logic
			argCount := animationArgCount[animCode]

			if animCode == 11 {
				// Special animation: read sub-action byte
				if pos >= len(data) {
					errors = append(errors, fmt.Sprintf("frame %d: truncated at special sub-action", frameNum))
					return frames, errors
				}
				subAction := data[pos]
				entry.RawArgs = append(entry.RawArgs, subAction)
				pos++
				extra := specialAnimationExtraArgs(subAction)
				for j := 0; j < extra; j++ {
					if pos >= len(data) {
						errors = append(errors, fmt.Sprintf("frame %d: truncated at special extra arg", frameNum))
						return frames, errors
					}
					entry.RawArgs = append(entry.RawArgs, data[pos])
					pos++
				}
			} else {
				if pos+argCount > len(data) {
					errors = append(errors, fmt.Sprintf("frame %d: truncated reading %d args", frameNum, argCount))
					return frames, errors
				}
				entry.RawArgs = make([]byte, argCount)
				copy(entry.RawArgs, data[pos:pos+argCount])
				pos += argCount
			}

			// For Move Dino (animation 2), decode arguments exactly as contestReadFrame does
			if animCode == 2 && len(entry.RawArgs) == 3 {
				arg1 := entry.RawArgs[0] // varNextByte — heading change
				arg2 := entry.RawArgs[1] // varNextByte2 — flags
				arg3 := entry.RawArgs[2] // movementAnimation[dinoIndex]

				// contestReadFrame decoding (see contestReadFrame.asm lines 342-373):
				//   bx = arg2 & 0x0F
				//   stack (arg_0) = arg2 >> 7
				//   ax = cbw(arg1)  (sign-extend arg1 to 16-bit)
				//   dx = (arg2 & 0x70) >> 4
				entry.MoveHeading = int8(arg1)
				entry.MoveBX = int(arg2 & 0x0F)
				entry.MoveDX = int((arg2 & 0x70) >> 4)
				entry.IsRunningSpeed = int(arg2 >> 7)
				entry.MoveMoveCode = arg3
			}

			frame.Animations = append(frame.Animations, entry)
		}

		frames = append(frames, frame)
	}

	return frames, errors
}

// TestCONT018MoveDetails parses CONT.018 and logs detailed analysis of every animation,
// with special focus on Move Dino commands to verify argument encoding matches what
// contestReadFrame passes to doMove.
func TestCONT018MoveDetails(t *testing.T) {
	contFile := filepath.Join("TestData", "Contests", "CONT.018")
	data, err := os.ReadFile(contFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", contFile, err)
	}

	if len(data) < CONTEST_HEADER_RECORD_LEN {
		t.Fatalf("File too small: %d bytes", len(data))
	}

	contestDataOffset := int(binary.LittleEndian.Uint16(data[0xF:0x11]))
	t.Logf("Contest data starts at file offset 0x%X (%d)", contestDataOffset, contestDataOffset)

	if contestDataOffset >= len(data) {
		t.Fatalf("Contest data offset 0x%X exceeds file size %d", contestDataOffset, len(data))
	}

	contestData := data[contestDataOffset:]
	frames, parseErrors := parseContestAnimations(contestData)

	for _, e := range parseErrors {
		t.Errorf("Parse error: %s", e)
	}

	// Summary stats
	totalAnimFrames := 0
	totalDelayFrames := 0
	totalDelayTicks := 0
	moveCount := 0

	for _, f := range frames {
		if f.IsDelay {
			totalDelayFrames++
			totalDelayTicks += f.DelayTicks
		} else {
			totalAnimFrames++
		}
	}

	t.Logf("=== Contest Data Summary ===")
	t.Logf("Total frames: %d (animation: %d, delay: %d)", len(frames), totalAnimFrames, totalDelayFrames)
	t.Logf("Total delay ticks: %d (%.1f seconds at 20fps)", totalDelayTicks, float64(totalDelayTicks)/20.0)

	// Log every non-delay frame in detail
	t.Logf("")
	t.Logf("=== Animation Frames ===")
	for _, f := range frames {
		if f.IsDelay {
			continue
		}
		t.Logf("Frame %d: %d animation(s)", f.FrameNum, len(f.Animations))
		for _, a := range f.Animations {
			animName := "Unknown"
			if a.AnimCode >= 0 && a.AnimCode < len(animationNames) {
				animName = animationNames[a.AnimCode]
			}
			t.Logf("  Dino %d | Anim %d (%s) | Encoded=0x%02X | RawArgs=%X",
				a.DinoIndex, a.AnimCode, animName, a.EncodedByte, a.RawArgs)

			if a.AnimCode == 2 {
				moveCount++
				t.Logf("    --- doMove argument breakdown (as contestReadFrame decodes) ---")
				t.Logf("    ax (heading change): %d (signed byte: 0x%02X)", a.MoveHeading, a.RawArgs[0])
				t.Logf("    bx/di (move mode):   0x%02X (%d)", a.MoveBX, a.MoveBX)
				t.Logf("    dx (stop/go):        %d", a.MoveDX)
				t.Logf("    stack arg_0 (isRunningSpeed):   %d", a.IsRunningSpeed)
				t.Logf("    ds:65b8 (moveCode):  0x%02X (%d)", a.MoveMoveCode, a.MoveMoveCode)
				t.Logf("    arg2 raw byte:       0x%02X (binary: %08b)", a.RawArgs[1], a.RawArgs[1])

				// Cross-reference with doc tables
				switch a.MoveBX {
				case 0x01:
					t.Logf("    >> Mode: CREEP (per doc: bx=1)")
				case 0x04:
					t.Logf("    >> Mode: RUN (per doc: bx=4)")
				case 0x0A:
					t.Logf("    >> Mode: WALK (per doc: bx=A)")
				default:
					t.Logf("    >> Mode: UNKNOWN (bx=0x%02X not in doc tables)", a.MoveBX)
				}

				// Flag potential issues
				if a.IsRunningSpeed == 1 {
					t.Logf("    >> IsRunningSpeed flag SET: doMove uses isRunningSpeed switch table (2nd table)")
					t.Logf("       Doc walk/creep tables show stack=0, run show stack=1, significant for switching from run to walk speed.")
					t.Logf("       If doc is correct, this should be 0 for all movement types.")
				} else {
					t.Logf("    >> IsRunningSpeed flag CLEAR: doMove uses simpler switch table (1st table)")
				}
			}
		}
	}

	t.Logf("")
	t.Logf("=== Move Command Count ===")
	t.Logf("Total Move Dino commands: %d", moveCount)

	if moveCount == 0 {
		t.Error("No Move Dino commands found in contest data")
	}
}

// TestMoveArg2Encoding verifies that the arg2 byte for Move Dino correctly encodes
// and decodes the three sub-fields: bx (mode), dx (stop/go), and stack (isRunningSpeed).
// This simulates both the Go server's encoding AND contestReadFrame's decoding.
func TestMoveArg2Encoding(t *testing.T) {
	cases := []struct {
		name       string
		isRunningSpeed    bool
		modeBits   byte // OR'd into arg2 (creep=0x01, walk=0x0A, run=0x04)
		wantBX     int
		wantDX     int
		wantStack  int
		docStackOK bool // whether doc tables agree with computed stack value
	}{
		// Walk with legs — Go code: arg2 = 0x90 | 0x0A = 0x9A
		{"walk_with_legs", true, 0x0A, 0x0A, 1, 1, false},
		// Walk without legs — Go code: arg2 = 0x10 | 0x0A = 0x1A
		{"walk_no_legs", false, 0x0A, 0x0A, 1, 0, true},
		// Run with legs — Go code: arg2 = 0x90 | 0x04 = 0x94
		{"run_with_legs", true, 0x04, 0x04, 1, 1, false},
		// Run without legs — Go code: arg2 = 0x10 | 0x04 = 0x14
		{"run_no_legs", false, 0x04, 0x04, 1, 0, true},
		// Creep with legs — Go code: arg2 = 0x90 | 0x01 = 0x91
		{"creep_with_legs", true, 0x01, 0x01, 1, 1, false},
		// Creep without legs — Go code: arg2 = 0x10 | 0x01 = 0x11
		{"creep_no_legs", false, 0x01, 0x01, 1, 0, true},
		// Dont move with legs — Go code: arg2 = 0x90 | 0x01 = 0x91 (same as creep)
		{"dontmove_with_legs", true, 0x01, 0x01, 1, 1, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate Go server encoding
			var arg2 byte
			if tc.isRunningSpeed {
				arg2 = 0x90
			} else {
				arg2 = 0x10
			}
			arg2 = arg2 | tc.modeBits

			// Simulate contestReadFrame decoding
			bx := int(arg2 & 0x0F)
			dx := int((arg2 & 0x70) >> 4)
			stack := int(arg2 >> 7)

			t.Logf("arg2=0x%02X (binary: %08b) → bx=0x%X dx=%d stack=%d", arg2, arg2, bx, dx, stack)

			if bx != tc.wantBX {
				t.Errorf("bx: got 0x%X, want 0x%X", bx, tc.wantBX)
			}
			if dx != tc.wantDX {
				t.Errorf("dx: got %d, want %d", dx, tc.wantDX)
			}
			if stack != tc.wantStack {
				t.Errorf("stack: got %d, want %d", stack, tc.wantStack)
			}

			if !tc.docStackOK {
				t.Logf("WARNING: Doc movement tables show stack=0 for all entries, but computed stack=%d", stack)
				t.Logf("  This causes doMove to use the isRunningSpeed switch table instead of no-legs table")
			}
		})
	}
}

// TestMoveArg2BitLayout documents and verifies the exact bit layout of arg2
// as contestReadFrame extracts it. This helps identify if bits are going to
// the wrong fields.
func TestMoveArg2BitLayout(t *testing.T) {
	// arg2 byte bit layout per contestReadFrame:
	//   Bit 7:     stack arg (isRunningSpeed flag) — arg2 >> 7
	//   Bits 6-4:  dx (stop/go flag)         — (arg2 & 0x70) >> 4
	//   Bits 3-0:  bx/di (movement mode)     — arg2 & 0x0F
	//
	// Go server encodes as:
	//   isRunningSpeed: arg2 = 0x90 (bit 7 = 1, bits 6-4 = 001)
	//   no-legs:  arg2 = 0x10 (bit 7 = 0, bits 6-4 = 001)
	//   Then OR in mode: creep=0x01, walk=0x0A, run=0x04

	t.Run("verify_base_values", func(t *testing.T) {
		// The base value 0x90 encodes TWO things: isRunningSpeed AND dx=1
		//   0x90 = 1001_0000
		//   bit 7 = 1 (isRunningSpeed)
		//   bits 6-4 = 001 (dx = 1, meaning "moving")
		//   bits 3-0 = 0000 (no mode set yet)
		base := byte(0x90)
		t.Logf("Base 0x%02X = %08b", base, base)
		t.Logf("  bit 7 (stack): %d", base>>7)
		t.Logf("  bits 6-4 (dx): %d", (base&0x70)>>4)
		t.Logf("  bits 3-0 (bx): 0x%X", base&0x0F)

		if (base >> 7) != 1 {
			t.Error("bit 7 should be 1 for isRunningSpeed base")
		}
		if ((base & 0x70) >> 4) != 1 {
			t.Error("dx should be 1 for 'moving' base")
		}

		// The base value 0x10 encodes only dx=1
		//   0x10 = 0001_0000
		//   bit 7 = 0 (no legs)
		//   bits 6-4 = 001 (dx = 1)
		//   bits 3-0 = 0000
		baseNoLegs := byte(0x10)
		t.Logf("Base 0x%02X = %08b", baseNoLegs, baseNoLegs)
		t.Logf("  bit 7 (stack): %d", baseNoLegs>>7)
		t.Logf("  bits 6-4 (dx): %d", (baseNoLegs&0x70)>>4)
		t.Logf("  bits 3-0 (bx): 0x%X", baseNoLegs&0x0F)
	})

	t.Run("walk_0x9A_full_decode", func(t *testing.T) {
		// arg2 = 0x9A is what CONT.018 has for both Move Dino commands
		arg2 := byte(0x9A)
		t.Logf("arg2 = 0x%02X = %08b", arg2, arg2)
		t.Logf("  bit 7 → stack (isRunningSpeed): %d", arg2>>7)
		t.Logf("  bits 6-4 → dx (stop/go): %d", (arg2&0x70)>>4)
		t.Logf("  bits 3-0 → bx (mode):    0x%X = %d decimal", arg2&0x0F, arg2&0x0F)
		t.Logf("")
		t.Logf("contestReadFrame passes to doMove:")
		t.Logf("  ax = heading (from arg1, sign-extended)")
		t.Logf("  bx = 0x%X → di in doMove", arg2&0x0F)
		t.Logf("  dx = %d", (arg2&0x70)>>4)
		t.Logf("  push %d (stack arg_0)", arg2>>7)
		t.Logf("")
		t.Logf("doMove behavior with these values:")
		if arg2>>7 == 1 {
			t.Logf("  arg_0=1 → TAKES isRunningSpeed branch (2nd switch table at loc_29B42)")
			t.Logf("  state 0: di=0x%X non-zero → checks leg type", arg2&0x0F)
			t.Logf("    leg type != 1 (sprawling) → state=9, schedules animation 14")
			t.Logf("    state 9 on next call → does NOTHING (no-op, state stays 9)")
		} else {
			t.Logf("  arg_0=0 → TAKES no-legs branch (1st switch table at loc_29A47)")
			t.Logf("  state 0: dx=%d non-zero → state=4", (arg2&0x70)>>4)
			t.Logf("    di=0x%X != 1 → no creep adjustment", arg2&0x0F)
			t.Logf("    schedules animation 14, then sub_2662A handles it")
		}
	})
}

// TestGenerateDelay validates the delay byte encoding.
func TestGenerateDelay(t *testing.T) {
	cases := []struct {
		count    int
		expected byte
		ticks    int
	}{
		{1, 0xFF, 1},
		{5, 0xFB, 5},
		{127, 0x81, 127},
		{128, 0x80, 128},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("delay_%d", tc.count), func(t *testing.T) {
			cr := NewContestResult()
			cr.GenerateDelay(tc.count)

			if len(cr.Actions) != 1 {
				t.Fatalf("Expected 1 byte, got %d", len(cr.Actions))
			}

			if cr.Actions[0] != tc.expected {
				t.Errorf("Expected byte 0x%02X, got 0x%02X", tc.expected, cr.Actions[0])
			}

			// Verify the game would interpret this correctly:
			// cbw sign-extends, neg makes positive (use -int() to avoid int8 overflow)
			signed := int8(cr.Actions[0])
			magnitude := -int(signed)
			if magnitude != tc.ticks {
				t.Errorf("Game would interpret as %d ticks, expected %d", magnitude, tc.ticks)
			}
		})
	}
}

// TestEndGame validates the terminator byte.
func TestEndGame(t *testing.T) {
	cr := NewContestResult()
	cr.EndGame()

	if len(cr.Actions) != 1 {
		t.Fatalf("Expected 1 byte, got %d", len(cr.Actions))
	}

	if cr.Actions[0] != 0x00 {
		t.Errorf("Expected terminator 0x00, got 0x%02X", cr.Actions[0])
	}
}
