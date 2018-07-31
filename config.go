package main

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/bwmarrin/discordgo"
)

var (
	dToken       string
	current_user user.User
	ds           discordgo.Session
	configPath   string
)

const configFile string = ".whisper.config"

func init() {
	current_user, err := user.Current()
	if err != nil {
		os.Exit(1)
	}
	configPath = current_user.HomeDir + "/" + configFile
}

const configHelp = `Configure whisper Discord settings.`

func (cmd *configCommand) Name() string      { return "config" }
func (cmd *configCommand) Args() string      { return "[OPTIONS]" }
func (cmd *configCommand) ShortHelp() string { return configHelp }
func (cmd *configCommand) LongHelp() string  { return configHelp }
func (cmd *configCommand) Hidden() bool      { return false }

func (cmd *configCommand) Register(fs *flag.FlagSet) {}

type configCommand struct{}

func (cmd *configCommand) Run(ctx context.Context, args []string) error {
	dTokenBytes, _ := ioutil.ReadFile(configPath)
	if len(dTokenBytes) > 0 {
		dToken = string(dTokenBytes)
	}

	ds, err := createDiscordSession(dToken)
	if err != nil {
		return err
	}

	if dToken == "" {
		file, err := os.Create(configPath)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(ds.Token); err != nil {
			return err
		}
	}

	return nil
}
