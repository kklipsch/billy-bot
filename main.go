package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
	"github.com/kklipsch/billy-bot/pkg/openrouter"
	"github.com/kklipsch/billy-bot/pkg/smee"
)

type CLI struct {
	Smee       smee.Command        `cmd:"smee" help:"Run the Smee client to receive webhook events."`
	OpenRouter openrouter.Frinkiac `cmd:"openrouter" help:"Send requests to OpenRouter AI models."`

	EnvFile  string `default:".env" name:"env-file" short:"e" help:"Path to the .env file to load. Defaults to .env in the current directory. Set explicitly to empty to skip loading."`
	LogLevel string `default:"warn" name:"log-level" short:"l" help:"Set the log level. Options: debug, info, warn, error, fatal, panic. Defaults to info."`
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	ctx := signalContext()

	cli := CLI{}

	k := kong.Parse(&cli,
		kong.Name("billy-bot"),
		kong.Description("The _worst_ code bot."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.BindTo(ctx, (*context.Context)(nil)),
	)

	if cli.EnvFile != "" {
		if err := godotenv.Load(cli.EnvFile); err != nil {
			fmt.Printf("Error loading .env file: %v\n", err)
			os.Exit(1)
		}
	}

	loggLevel, err := zerolog.ParseLevel(cli.LogLevel)
	if err != nil {
		fmt.Println("Invalid log level:", err)
		os.Exit(1)
	}
	zerolog.SetGlobalLevel(loggLevel)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err = k.Run()
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
