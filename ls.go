package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"

	"github.com/azillion/chord/internal/getconfig"
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
		logrus.Debugf("Session Failed \n%v\nexiting.", spew.Sdump(ds))
		return fmt.Errorf("You may need to login from a browser first or check your credentials\n%v", err)
	}

	channels, err := ds.UserChannels()
	logrus.Debugf("Retrieved Channels\n%v\n", spew.Sdump(channels))
	if err != nil {
		return err
	}
	fmt.Println("Available Private Channels:")
	for i, channel := range channels {
		logrus.Debugf("Channel %d\n%v\n", i, spew.Sdump(channel))

		// Switch of supported channel types
		switch chanType := channel.Type; chanType {
		case 1: // Direct Messages
			var flatRecipients string
			recipients := channel.Recipients
			logrus.Debugf("Channel %d recipients\n%v\n", i, spew.Sdump(recipients))
			for _, recipient := range recipients {
				flatRecipients = fmt.Sprintf("%s %s", flatRecipients, recipient.Username)
			}
			fmt.Printf("\t%d) DM to %s\n", i, strings.TrimSpace(flatRecipients))
		default:
			fmt.Println("No available channels")
			return nil
		}
	}
	// reader := bufio.NewReader(os.Stdin)

	// fmt.Print("Select a channel to switch to: ")
	// channelSelS, err := reader.ReadString('\n')
	// if err != nil {
	// 	return err
	// }
	// channelSel, err := strconv.Atoi(strings.TrimSpace(channelSelS))
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(channelSel)

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
