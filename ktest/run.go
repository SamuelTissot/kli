package ktest

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

const KT_ENV_FLAG = "KLI_TEST_EXEC"

type KT struct {
	Out    []byte
	ErrOut bytes.Buffer
	Err    error
	StdErr io.Writer
}

func NewKT() *KT {
	return &KT{}
}

func (kt *KT) Exec(t *testing.T, f func(t *testing.T)) {

	if os.Getenv(KT_ENV_FLAG) == "1" {
		f(t)
	}

	n := kt.CallerName()
	cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=%s", n))
	cmd.Stderr = bufio.NewWriter(&kt.ErrOut)
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=1", KT_ENV_FLAG))
	kt.Out, kt.Err = cmd.Output()
}

func (kt *KT) CallerName() string {
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
