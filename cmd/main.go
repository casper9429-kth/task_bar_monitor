package main

import (
	"flag"
	"log"
	"os"

	"github.com/casper9429-kth/task_bar_monitor/internal/app"
)

func main() {
	// Parse command line flags
	debugFlag := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// Set up logging
	if *debugFlag {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.Println("Debug logging enabled")
	} else {
		// Only errors in normal mode
		log.SetOutput(os.Stderr)
		log.SetFlags(log.Ldate | log.Ltime)
	}

	log.Println("Starting Ubuntu System Monitor")

	// Initialize and run the application
	application := app.NewApplication()
	err := application.Run()
	if err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
