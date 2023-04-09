package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	x = flag.String("X", "", "")
	h = flag.String("H", "", "")
	o = flag.String("o", "", "")
	c = flag.String("c", "", "")
	v = flag.Bool("v", false, "")
	q = flag.Bool("q", false, "")

	supportedMethods = []string{
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
	}
)

const usage = `Usage: weber -o output.csv [OPTIONS] <url>

Options:
  -X <method>   Comma-separated list of HTTP methods to watch for (GET, POST, OPTIONS, PUT, DELETE). Default behavior is to consider all methods.
  -H <string>   Comma-separated list of hostname or IP address to watch for. Default behavior is to consider all hosts.
  -o <file>     Write the response to a file. CSV is the default and only supported format.
  -c <string>   Comma-separated list of columns to write to the output file. Default is "url,method,status". Available columns are:
                    %s
  -v            Enable verbose logging to observe all browser events.
  -q            Disable all logging.

Examples:
      # Watch for all requests to https://example.com
      weber -o output.csv https://example.com

      # Watch for all requests to https://example.com> and https://example.org
      weber -H example.com,example.org -o output.csv https://example.com

      # Watch for GET requests on example.org and output the URL, request method, the status code, and cache-control header
      weber -X GET -H example.org -o output.csv -c "url,method,status,Cache-Control" https://example.com
`

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

	// url is the URL we observe being loaded
	url string

	// outputPath is the path to the file where the output will be written
	outputPath string
	// outputCols is the list of columns to write to the output file
	outputCols []string
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, mergedSupportedCols())
	}

	cfg, err := flagsToConfig()
	if err != nil {
		usageAndExit(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	log := logger(cfg)

	if err := WatchNetworkFor(ctx, cfg.url, cfg, log); err != nil {
		errAndExit(err.Error())
	}

	log("Done.")
}

type logF func(string, ...any)

func logger(cfg config) logF {
	return func(msg string, args ...any) {
		if cfg.quiet {
			return
		}

		fmt.Fprintf(os.Stderr, msg, args...)
		fmt.Fprint(os.Stderr, "\n")
	}
}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprint(os.Stderr, msg)
		fmt.Fprint(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}

func errAndExit(msg string) {
	fmt.Fprint(os.Stderr, msg)
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}

func flagsToConfig() (config, error) {
	flag.Parse()
	if flag.NArg() < 1 {
		return config{}, fmt.Errorf("missing required argument: url")
	}

	// Make sure the host is provided
	url := flag.Args()[0]
	if url == "" {
		return config{}, fmt.Errorf("provided url argument is empty")
	}

	cfg := config{
		url:         url,
		giveUpAfter: 5 * time.Second,
		outputCols:  []string{"url", "method", "status"},
	}

	if *x != "" {
		methods := strings.Split(*x, ",")
		if len(methods) > 0 {
			for _, m := range methods {
				if !contains(supportedMethods, m) {
					return config{}, fmt.Errorf("unsupported HTTP method: %s", m)
				}
				cfg.methods = append(cfg.methods, m)
			}
		}
	}

	if *h != "" {
		cfg.domains = strings.Split(*h, ",")
	}

	if *c != "" {
		cols := strings.Split(*c, ",")
		if len(cols) > 0 {
			cfg.outputCols = cols
		}
	}

	if *o == "" {
		return config{}, fmt.Errorf("missing required argument: -o")
	} else {
		cfg.outputPath = *o
	}

	if *v {
		cfg.verbose = true
	}

	if *q {
		cfg.quiet = true
		cfg.verbose = false
	}

	return cfg, nil
}
