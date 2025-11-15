package commands

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

// ParseCmd implements the parse command
type ParseCmd struct {
	flags *Flags
}

// NewParseCmd creates a new parse command
func NewParseCmd(flags *Flags) *ParseCmd {
	return &ParseCmd{flags: flags}
}

// Register adds the parse command to the application
func (cmd *ParseCmd) Register(app *cli.Command) *cli.Command {
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "parse",
		Usage: "parse command",
		Flags: []cli.Flag{
			// Add command-specific flags here
		},
		Action: cmd.run,
	})

	return app
}

func (cmd *ParseCmd) run(ctx context.Context, c *cli.Command) error {
	log.Info().Msg("running parse command")

	fmt.Println("Hello World!")

	return nil
}
