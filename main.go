package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

type config struct {
	// verbose enables verbose logging to observe all browser events
	verbose bool
	// quiet disables all logging
	quiet bool
	// giveUpAfter is the time after which the program will give up waiting for another request
	giveUpAfter time.Duration

	// methods is a list of HTTP methods to watch for
	methods []string
	// domains is a list of domains to watch for
	domains []string

	// outputPath is the path to the file where the output will be written
	outputPath string
	// outputCols is the list of columns to write to the output file
	outputCols []string
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	cfg := config{
		quiet:       true,
		giveUpAfter: 5 * time.Second,

		methods: []string{"GET"},
		//domains:     []string{"compass.pressekompass.net"},
		outputPath: "output.csv",
		outputCols: []string{"url", "method", "status", "Content-Type", "Cache-Control"},
	}

	url := "https://compass.pressekompass.net/compasses/braunschweigerzeitung/was-halten-sie-vom-leihrollerverbot-in-p-xH4vV5"

	if err := WatchNetworkFor(ctx, url, cfg); err != nil {
		if !cfg.quiet {
			log.Fatalf("error watching network for %s: %v", url, err)
		}
		os.Exit(1)
	}

	if !cfg.quiet {
		fmt.Println("Done.")
	}
}
