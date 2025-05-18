package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/kklipsch/billy-bot/pkg/openrouter"
	"github.com/kklipsch/billy-bot/pkg/smee"
)

type CLI struct {
	Smee       smee.Command       `cmd:"smee" help:"Run the Smee client to receive webhook events."`
	OpenRouter openrouter.Command `cmd:"openrouter" help:"Send requests to OpenRouter AI models."`
}

func main() {
	ctx := signalContext()

	k := kong.Parse(&CLI{},
		kong.Name("billy-bot"),
		kong.Description("The _worst_ code bot."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.BindTo(ctx, (*context.Context)(nil)),
	)

	err := k.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func signalContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
		cancel()
	}()

	return ctx
}
