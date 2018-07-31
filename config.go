package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/azillion/whisper/internal/getconfig"
	"github.com/sirupsen/logrus"
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
	logrus.Debugf("email: %s\n", authConfig.Email)
	if err != nil {
		return err
	}
	ds, err := createDiscordSession(authConfig)
	if err != nil {
		return fmt.Errorf("You may need to login from a browser first or check your credentials\n%v", err)
	}
	logrus.Debugf("Session %v\n", spew.Sdump(ds))
	fmt.Println("Created and saved a Discord auth token")

	return nil
}
