package tracing

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func Trace(skip int) string {
	frame := GetCallerInfo(skip)
	return fmt.Sprintf("\t %s:%d | %s", filepath.Base(frame.File), frame.Line, frame.Function)
}

const maxTraceback = 10

func GetCallerInfo(skip int) runtime.Frame {
	pc := make([]uintptr, maxTraceback)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame
}
