package main

import (
	"context"
	"flag"
	"fmt"
)

var dToken string

const configHelp = `Configure whisper Discord settings.`

func (cmd *configCommand) Name() string      { return "config" }
func (cmd *configCommand) Args() string      { return "[OPTIONS]" }
func (cmd *configCommand) ShortHelp() string { return configHelp }
func (cmd *configCommand) LongHelp() string  { return configHelp }
func (cmd *configCommand) Hidden() bool      { return false }

func (cmd *configCommand) Register(fs *flag.FlagSet) {}

type configCommand struct{}

func (cmd *configCommand) Run(ctx context.Context, args []string) error {
	dToken, err := createDiscordSession()
	if err != nil {
		return err
	}
	fmt.Printf("\nToken: %s\n", dToken)
	return nil
}
