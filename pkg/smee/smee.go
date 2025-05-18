package smee

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Command represents the CLI command for Smee
type Command struct {
	URL string `arg:"" optional:"" help:"The Smee.io URL to subscribe to. If not provided, a new channel will be created."`
}

// Run executes the Smee command
func (s *Command) Run() error {
	var source *string
	var err error

	if s.URL != "" {
		source = &s.URL
	} else {
		source, err = CreateChannel()
		if err != nil {
			return err
		}
	}

	fmt.Println("Subscribing to smee source: " + *source)

	logger := log.Logger{}

	target := make(chan Event)
	client := NewClient(source, target, &logger)

	fmt.Println("Client initialised")

	sub, err := client.Start()
	if err != nil {
		return err
	}

	fmt.Println("Client running")

	for ev := range target {
		// do what you want with the event
		fmt.Printf("Received event: id=%v, name=%v, payload=%v\n", ev.Id, ev.Name, string(ev.Data))
	}

	sub.Stop()
	return nil
}

// Client handles the connection to a Smee.io channel
type Client struct {
	source *string
	target chan<- Event
	logger *log.Logger
}

// Event represents a Server-Sent Event from Smee.io
type Event struct {
	Id   string
	Name string
	Data []byte
}

// CreateChannel creates a new Smee.io channel and returns its URL
func CreateChannel() (*string, error) {
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

// NewClient creates a new Client
func NewClient(source *string, target chan<- Event, logger *log.Logger) *Client {
	c := new(Client)
	c.source = source
	c.target = target
	c.logger = logger
	return c
}

// Start begins listening for events from the Smee.io channel
func (c *Client) Start() (*Subscription, error) {
	eventStream, err := OpenSSEUrl(*c.source)
	if err != nil {
		return nil, err
	}

	quit := make(chan interface{})
	go c.run(eventStream, quit)

	return &Subscription{terminator: quit}, nil
}

func (c *Client) run(sseEventStream <-chan Event, quit <-chan interface{}) {
	for {
		select {
		case event := <-sseEventStream:
			c.target <- event
		case <-quit:
			return
		}
	}
}

// Subscription represents an active subscription to a Smee.io channel
type Subscription struct {
	terminator chan<- interface{}
}

// Stop terminates the subscription to the Smee.io channel
func (c *Subscription) Stop() {
	c.terminator <- nil
}

// OpenSSEUrl opens a connection to a Server-Sent Events endpoint
func OpenSSEUrl(url string) (<-chan Event, error) {
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

	events := make(chan Event)

	var buf bytes.Buffer

	go func() {
		ev := Event{}
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
				ev = Event{}

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
