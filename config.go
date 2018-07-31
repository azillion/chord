package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/azillion/whisper/internal/getconfig"
	"github.com/bwmarrin/discordgo"
)

var (
	ds discordgo.Session
)

const configHelp = `Configure whisper Discord settings.`

func (cmd *configCommand) Name() string      { return "config" }
func (cmd *configCommand) Args() string      { return "[OPTIONS]" }
func (cmd *configCommand) ShortHelp() string { return configHelp }
func (cmd *configCommand) LongHelp() string  { return configHelp }
func (cmd *configCommand) Hidden() bool      { return false }

func (cmd *configCommand) Register(fs *flag.FlagSet) {}

type configCommand struct{}

func (cmd *configCommand) Run(ctx context.Context, args []string) error {
	authConfig, err := getconfig.GetAuthConfig(email, password)
	if err != nil {
		return err
	}
	_, err = createDiscordSession(authConfig)
	if err != nil {
		return fmt.Errorf("You may need to login from a browser first or check your credentials\n%v", err)
	}
	fmt.Println("Created and saved a Discord auth token")

	return nil
}
