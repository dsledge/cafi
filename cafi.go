/*
Module: CAFI
Package: CAFI
Description: Cross Account Function Iterator. This package is used to iterate through mutliple cloud based accounts and or projects using the supported provider packages.
*/
package cafi

import (
	"github.com/dsledge/scribble"
)

func Configure(logfile *string, loglevel *int) {
	// Configuring the log for console or file
	if *logfile == "console" {
		scribble.NewConsoleLogger(*loglevel)
	} else {
		scribble.NewFileLogger(*loglevel, *logfile)
	}
}
