package smee

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
)

// Command represents the CLI command for Smee
type Command struct {
	URL string `arg:"" optional:"" help:"The Smee.io URL to subscribe to. If not provided, checks SMEE_SOURCE env var, then creates a new channel if needed."`
}

// Run executes the Smee command
func (s *Command) Run(ctx context.Context) error {
	var (
		source string
		err    error
	)

	// Check command line argument first
	if s.URL != "" {
		source = s.URL
	} else {
		// Then check environment variable
		envSource := os.Getenv("SMEE_SOURCE")
		if envSource != "" {
			source = envSource
		} else {
			// Finally create a new channel if needed
			source, err = CreateChannel()
			if err != nil {
				return err
			}
		}
	}

	// Log the source of the Smee URL
	if s.URL != "" {
		fmt.Println("Subscribing to smee source (from command line): " + source)
	} else if os.Getenv("SMEE_SOURCE") != "" {
		fmt.Println("Subscribing to smee source (from SMEE_SOURCE env var): " + source)
	} else {
		fmt.Println("Subscribing to smee source (newly created): " + source)
	}

	events, err := OpenSSEUrl(ctx, source)
	if err != nil {
		return err
	}

	for ev := range events {
		// do what you want with the event
		fmt.Printf("Received event: id=%v, name=%v, payload=%v\n", ev.Id, ev.Name, string(ev.Data))
	}

	return nil
}

// Event represents a Server-Sent Event from Smee.io
type Event struct {
	Id   string
	Name string
	Data []byte
}

// CreateChannel creates a new Smee.io channel and returns its URL
func CreateChannel() (string, error) {
	httpClient := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := httpClient.Head("https://smee.io/new")
	if err != nil {
		return "", err
	}

	loc := resp.Header.Get("Location")
	return loc, nil
}

// OpenSSEUrl opens a connection to a Server-Sent Events endpoint
func OpenSSEUrl(ctx context.Context, url string) (<-chan Event, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "text/event-stream")
	req = req.WithContext(ctx)
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

		err := scanner.Err()
		if err == ctx.Err() {
			close(events)
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Error during resp.Body read:%s\n", err)
			close(events)
		}
	}()

	return events, nil
}
