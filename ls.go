package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/azillion/whisper/internal/getconfig"
)

const lsHelp = `List available Discord channels.`

func (cmd *lsCommand) Name() string      { return "ls" }
func (cmd *lsCommand) Args() string      { return "[OPTIONS]" }
func (cmd *lsCommand) ShortHelp() string { return lsHelp }
func (cmd *lsCommand) LongHelp() string  { return lsHelp }
func (cmd *lsCommand) Hidden() bool      { return false }

func (cmd *lsCommand) Register(fs *flag.FlagSet) {}

type lsCommand struct{}

func (cmd *lsCommand) Run(ctx context.Context, args []string) error {
	ds, err := createDiscordSession(getconfig.AuthConfig{})
	if err != nil {
		return fmt.Errorf("You may need to login from a browser first or check your credentials\n%v", err)
	}

	channels, err := ds.UserChannels()
	if err != nil {
		return err
	}
	fmt.Println("Available Private Channels:")
	for i, channel := range channels {
		if channel.Type == 1 {
			var flatRecipients string
			recipients := channel.Recipients
			for _, recipient := range recipients {
				flatRecipients = fmt.Sprintf("%s %s", flatRecipients, recipient.Username)
			}

			fmt.Printf("\t%d) DM to %s\n", i, strings.TrimSpace(flatRecipients))
			// spew.Dump(channel)
		}
	}

	// TODO: select from guild channels
	// guilds, err := ds.UserGuilds(10, "", "")
	// if err != nil {
	// 	return err
	// }
	// fmt.Println()
	// fmt.Println("Available Guilds:")
	// for _, guild := range guilds {
	// 	spew.Dump(guild)
	// }

	return nil
}
