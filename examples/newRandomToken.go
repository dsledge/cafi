package main

// Import the CAFI Module
import (
	"flag"

	cafi "github.com/dsledge/cafi"
	scribble "github.com/dsledge/scribble"
)

var (
	logfile  = flag.String("logfile", "console", "The default log file will log to system console")
	loglevel = flag.Int("loglevel", 0, "Sets the default log level to INFO messages and higher")
	bytes    = flag.Int("bytes", 16, "Sets the byte length used to generate the token")
)

// Generate a new random token
func main() {
	flag.Parse()

	// Configure the CAFI SDK
	cafi.Configure(logfile, loglevel)

	token, err := cafi.NewRandomToken(*bytes)
	if err != nil {
		scribble.Fatal("Unable to generate a new random token: %s", err)
	}
	scribble.Info("%d byte random token generated: %s", bytes, token)
}
