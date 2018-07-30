package main

import (
	"context"
	"flag"
	"fmt"
	"bufio"
	"io/ioutil"
	"os"
	"os/user"
)

var (
	dToken string
	user User
	path string
	ds 
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
	user, err := user.Current()
	if err != nil {
		return err
	}

	path := user.HomeDir + ".whisper.config"
	if _, err := os.Stat("/path/to/whatever"); err == nil {
	token, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
	}

	ds, err := createDiscordSession()
	if err != nil {
		return err
	}

	err := ioutil.WriteFile(path, ds.Token, 0644)

	return nil
}
