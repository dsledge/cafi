package cafi

import (
	"testing"

	"github.com/dsledge/scribble"
)

func TestConfigure(t *testing.T) {
	console := "console"
	level := scribble.TRACE
	Configure(&console, &level)
	scribble.Trace("Testing that TRACE logging is working")
	scribble.Debug("Testing that DEBUG logging is working")
	scribble.Info("Testing that INFO logging is working")
	scribble.Warn("Testing that WARN logging is working")
	scribble.Error("Testing that ERROR logging is working")
}
