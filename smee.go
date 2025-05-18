package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

// SmeeClient handles the connection to a Smee.io channel
type SmeeClient struct {
	source *string
	target chan<- SSEvent
	logger *log.Logger
}

// SSEvent represents a Server-Sent Event from Smee.io
type SSEvent struct {
	Id   string
	Name string
	Data []byte
}

// CreateSmeeChannel creates a new Smee.io channel and returns its URL
func CreateSmeeChannel() (*string, error) {
	httpClient := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := httpClient.Head("https://smee.io/new")
	if err != nil {
		return nil, err
	}

	loc := resp.Header.Get("Location")
	return &loc, nil
}

// NewSmeeClient creates a new SmeeClient
func NewSmeeClient(source *string, target chan<- SSEvent, logger *log.Logger) *SmeeClient {
	c := new(SmeeClient)
	c.source = source
	c.target = target
	c.logger = logger
	return c
}

// Start begins listening for events from the Smee.io channel
func (c *SmeeClient) Start() (*SmeeClientSubscription, error) {
	eventStream, err := OpenSSEUrl(*c.source)
	if err != nil {
		return nil, err
	}

	quit := make(chan interface{})
	go c.run(eventStream, quit)

	return &SmeeClientSubscription{terminator: quit}, nil
}

func (c *SmeeClient) run(sseEventStream <-chan SSEvent, quit <-chan interface{}) {
	for {
		select {
		case event := <-sseEventStream:
			c.target <- event
		case <-quit:
			return
		}
	}
}

// SmeeClientSubscription represents an active subscription to a Smee.io channel
type SmeeClientSubscription struct {
	terminator chan<- interface{}
}

// Stop terminates the subscription to the Smee.io channel
func (c *SmeeClientSubscription) Stop() {
	c.terminator <- nil
}

// OpenSSEUrl opens a connection to a Server-Sent Events endpoint
func OpenSSEUrl(url string) (<-chan SSEvent, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "text/event-stream")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error: resp.StatusCode == %d\n", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "text/event-stream" {
		return nil, fmt.Errorf("Error: invalid Content-Type == %s\n", resp.Header.Get("Content-Type"))
	}

	events := make(chan SSEvent)

	var buf bytes.Buffer

	go func() {
		ev := SSEvent{}
		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			line := scanner.Bytes()

			switch {

			// start of event
			case bytes.HasPrefix(line, []byte("id:")):
				ev.Id = string(line[4:])

				// event name
			case bytes.HasPrefix(line, []byte("event:")):
				ev.Name = string(line[7:])

				// event data
			case bytes.HasPrefix(line, []byte("data:")):
				buf.Write(line[6:])

				// end of event
			case len(line) == 0:
				ev.Data = buf.Bytes()
				buf.Reset()
				events <- ev
				ev = SSEvent{}

			default:
				fmt.Fprintf(os.Stderr, "Error during EventReadLoop - Default triggerd! len:%d\n%s", len(line), line)
				close(events)

			}
		}

		if err = scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error during resp.Body read:%s\n", err)
			close(events)

		}
	}()

	return events, nil
}
