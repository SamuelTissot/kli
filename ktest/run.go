package ktest

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

const KT_ENV_FLAG = "KLI_TEST_EXEC"

// KT
// helps to test command line tools with
// os.exit() calls
type KT struct {
	// Out is the stdout
	Out []byte
	// ErrOut is the stderr
	ErrOut bytes.Buffer
	// the execution error
	Err error
}

func NewKT() *KT {
	return &KT{}
}

// Exec run the test function and capture it's output
func (kt *KT) Exec(t *testing.T, f func(t *testing.T)) {
	if os.Getenv(KT_ENV_FLAG) == "1" {
		f(t)
		return
	}

	n := kt.callerName()
	cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=%s", n))
	cmd.Stderr = bufio.NewWriter(&kt.ErrOut)
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=1", KT_ENV_FLAG))
	kt.Out, kt.Err = cmd.Output()
}

// callerName returns the calling function name
func (kt *KT) callerName() string {
	n := kt.getFrame(2).Function
	pos := strings.LastIndex(n, ".")
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(".")
	if adjustedPos >= len(n) {
		return ""
	}
	return n[adjustedPos:]
}

// getFrame
// got it from stackoverflow https://stackoverflow.com/questions/35212985/is-it-possible-get-information-about-caller-function-in-golang
func (kt *KT) getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}
