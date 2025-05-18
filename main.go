package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/kklipsch/billy-bot/pkg/smee"
)

type CLI struct {
	Smee smee.Command `cmd:"smee" help:"Run the Smee client to receive webhook events."`
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Name("billy-bot"),
		kong.Description("The _worst_ code bot."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)
	err := ctx.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
