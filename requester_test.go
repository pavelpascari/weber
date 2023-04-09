package main

import (
	"context"
	"github.com/chromedp/cdproto/network"
	"reflect"
	"testing"
	"time"
)

func TestProcessor(t *testing.T) {

	reqs := []network.EventRequestWillBeSent{
		{
			RequestID: "1",
			Request: &network.Request{
				URL: "http://example.com/1",
			},
		},
		{
			RequestID: "2",
			Request: &network.Request{
				URL: "http://example.com/2",
			},
		},
		// repeat the request once more
		{
			RequestID: "1",
			Request: &network.Request{
				URL: "http://example.com/1",
			},
		},
	}

	resps := []network.EventResponseReceived{
		{
			RequestID: "1",
			Response: &network.Response{
				URL:    "http://example.com/1",
				Status: 200,
				Headers: network.Headers{
					"Cache-Control": "max-age=0",
				},
			},
		},
		{
			RequestID: "2",
			Response: &network.Response{
				URL:    "http://example.com/2",
				Status: 200,
				Headers: network.Headers{
					"Cache-Control": "max-age=0",
				},
			},
		},
	}

	p := NewProcessor()

	go func() {
		for _, req := range reqs {
			p.networkRequests <- evt{
				RequestID: string(req.RequestID),
				Request:   req.Request,
			}
		}

		for _, resp := range resps {
			p.networkResponse <- evt{
				RequestID: string(resp.RequestID),
				Response:  resp.Response,
			}
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// assertions go here

	for i := 0; i < 2; i++ {
		select {
		case e := <-p.out:
			if e.Request == nil || e.Response == nil {
				t.Fatalf("expected a request and response but got nil: %v", e)
			}
		case <-ctx.Done():
			t.Log(p.Close())
			t.Fatal("timed out")
		}
	}
}

func Test_getRecord(t *testing.T) {
	dummyLog := func(s string, a ...any) {}

	tests := []struct {
		name string
		e    evt
		cfg  config
		want []string
	}{
		{
			name: "basic",
			e: evt{
				Request: &network.Request{
					URL:    "http://example.com",
					Method: "GET",
				},
				Response: &network.Response{
					URL:    "http://example.com",
					Status: 200,
					Headers: network.Headers{
						"Cache-Control": "max-age=0",
					},
				},
			},
			cfg: config{
				methods:    []string{},
				domains:    []string{},
				outputCols: []string{"url", "method", "status", "cache-control"},
			},
			want: []string{"http://example.com", "GET", "200", "max-age=0"},
		},
		{
			name: "filter-methods",
			e: evt{
				Request: &network.Request{
					URL:    "http://example.com",
					Method: "GET",
				},
				Response: &network.Response{
					URL:    "http://example.com",
					Status: 200,
					Headers: network.Headers{
						"Cache-Control": "max-age=0",
					},
				},
			},
			cfg: config{
				methods: []string{"POST"},
			},
			want: nil,
		},
		{
			name: "filter-domains",
			e: evt{
				Request: &network.Request{
					URL:    "http://example.com",
					Method: "GET",
				},
				Response: &network.Response{
					URL:    "http://example.com",
					Status: 200,
					Headers: network.Headers{
						"Cache-Control": "max-age=0",
					},
				},
			},
			cfg: config{
				domains: []string{"foo.bar"},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRecord(tt.e, tt.cfg, dummyLog); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}
