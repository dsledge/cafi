/*
Module: CAFI
Package: CAFI
Description: Cross Account Function Iterator. This package is used to iterate through mutliple cloud based accounts and or projects using the supported provider packages.
*/
package cafi

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/dsledge/scribble"
)

// Configure the CAFI sdk for logging
func Configure(logfile *string, loglevel *int) {
	// Configuring the log for console or file
	if *logfile == "console" {
		scribble.NewConsoleLogger(*loglevel)
	} else {
		scribble.NewFileLogger(*loglevel, *logfile)
	}
}

// Create a new random token as hex string and return to the caller
func NewRandomToken(byte_length int) (string, error) {
	bytes := make([]byte, byte_length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
