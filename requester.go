package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type evt struct {
	RequestID string
	Request   *network.Request
	Response  *network.Response
}

func NewProcessor() *Processor {
	p := &Processor{
		networkRequests: make(chan evt),
		networkResponse: make(chan evt),
		out:             make(chan evt),
		closing:         make(chan chan error),
	}

	go p.loop()

	return p
}

type Processor struct {
	networkRequests chan evt
	networkResponse chan evt

	closing chan chan error

	out chan evt
}

func (p *Processor) Close() error {
	errChan := make(chan error)
	p.closing <- errChan
	return <-errChan
}

func (p *Processor) loop() {
	requests := make(map[string]evt)
	var pending = []evt{}
	for {
		var first evt
		var updates chan evt

		if len(pending) > 0 {
			first = pending[0]
			updates = p.out
		} else {
			updates = nil
		}

		select {
		case errChan := <-p.closing:
			p.networkRequests = nil
			p.networkResponse = nil
			close(p.out)
			errChan <- nil
			return

		case e := <-p.networkRequests:
			req, ok := requests[e.RequestID]
			if !ok {
				requests[e.RequestID] = evt{RequestID: e.RequestID, Request: e.Request}
				continue
			}

			if req.Response != nil {
				pending = append(pending, evt{RequestID: e.RequestID, Request: e.Request, Response: req.Response})
				delete(requests, e.RequestID)
			}
		case e := <-p.networkResponse:
			req, ok := requests[e.RequestID]
			if !ok {
				requests[e.RequestID] = evt{RequestID: e.RequestID, Response: e.Response}
				continue
			}

			if req.Request != nil {
				pending = append(pending, evt{RequestID: e.RequestID, Request: req.Request, Response: e.Response})
				delete(requests, e.RequestID)
			}
		case updates <- first:
			pending = pending[1:]
		}
	}
}

func WatchNetworkFor(ctx context.Context, url string, cfg config, log logF) error {
	// Create a new ChromeDP context
	var opts []chromedp.ContextOption

	if cfg.verbose {
		opts = append(opts, chromedp.WithDebugf(log))
	}
	ctx, cancel := chromedp.NewContext(ctx, opts...)
	defer cancel()

	proc := NewProcessor()
	defer proc.Close()

	// Set up network event listeners
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *network.EventRequestWillBeSent:
			proc.networkRequests <- evt{RequestID: string(e.RequestID), Request: e.Request}
		case *network.EventResponseReceived:
			proc.networkResponse <- evt{RequestID: string(e.RequestID), Response: e.Response}
		}
	})

	// Run tasks
	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(url),
	)
	if err != nil {
		return fmt.Errorf("failed to run chromedp: %v", err)
	}

	file, err := os.OpenFile(cfg.outputPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer func() {
		log("Flushing writer...")
		writer.Flush()
	}()

	if err := writer.Write(cfg.outputCols); err != nil {
		return fmt.Errorf("failed to write header row: %v", err)
	}

	for {
		timeout := time.After(cfg.giveUpAfter)

		select {
		case e := <-proc.out:
			timeout = nil
			if !cfg.quiet {
				fmt.Print(".") // progress indicator
			}
			if record := getRecord(e, cfg, log); record != nil {
				if err := writer.Write(record); err != nil {
					return fmt.Errorf("failed to write record: %v", err)
				}
			}
		case <-timeout:
			log("\nGiving up waiting...")
			return nil
		case <-ctx.Done():
			log("context done...")
			return ctx.Err()
		}
	}
}

func getHeaderValue(h network.Headers, key string) string {
	for k, v := range h {
		if http.CanonicalHeaderKey(k) == http.CanonicalHeaderKey(key) {
			return v.(string)
		}
	}
	return ""
}

func getRecord(e evt, cfg config, log logF) []string {
	if len(cfg.methods) > 0 {
		if !contains(cfg.methods, e.Request.Method) {
			return nil
		}
	}

	if len(cfg.domains) > 0 {
		u, err := url.Parse(e.Request.URL)
		if err != nil {
			log("failed to parse url:", e.Request.URL)
			return nil
		}

		if !contains(cfg.domains, u.Hostname()) {
			return nil
		}
	}

	record := []string{}
	for _, col := range cfg.outputCols {
		record = append(record, getColValue(e, col))
	}

	return record
}

func contains(arr []string, target string) bool {
	for _, el := range arr {
		if el == target {
			return true
		}
	}
	return false
}

var (
	supportedGetters = map[string]func(evt) string{
		"url": func(e evt) string {
			return e.Request.URL
		},
		"method": func(e evt) string {
			return e.Request.Method
		},
		"Content-Type": func(e evt) string {
			return getHeaderValue(e.Response.Headers, "Content-Type")
		},
		"Cache-Control": func(e evt) string {
			return getHeaderValue(e.Response.Headers, "Cache-Control")
		},
		"Content-Length": func(e evt) string {
			return getHeaderValue(e.Response.Headers, "Content-Length")
		},

		"status": func(e evt) string {
			return strconv.FormatInt(e.Response.Status, 10)
		},
	}
)

func getColValue(e evt, col string) string {
	if getter, ok := supportedGetters[col]; ok {
		res := getter(e)
		return res
	}

	return "---"
}

func mergedSupportedCols() string {
	cols := []string{}
	for k := range supportedGetters {
		cols = append(cols, k)
	}
	return strings.Join(cols, ", ")
}
