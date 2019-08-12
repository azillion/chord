// Copyright © 2018 Alexander Zillion <alex@alexzillion.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"syscall"

	"github.com/azillion/chord/internal/getconfig"
	"github.com/azillion/chord/version"
	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	"github.com/genuinetools/pkg/cli"
	"github.com/sirupsen/logrus"
)

var (
	email    string
	password string
	token    string

	debug      bool
	configPath string

	// // DSession global Discord session
	// DSession discordgo.Session
	dUser *discordgo.User
)

const (
	configFile        string = ".chord.config"
	discordSessionKey string = "discordSession"
)

func init() {
	currentUser, err := user.Current()
	if err != nil {
		os.Exit(1)
	}
	configPath = currentUser.HomeDir + "/" + configFile
}

func main() {
	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "chord"
	p.Description = "A Discord TUI for direct messaging."
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Build list of available commands
	p.Commands = []cli.Command{
		&configCommand{},
		&lsCommand{},
		&tuiCommand{},
	}

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("global", flag.ExitOnError)

	p.FlagSet.StringVar(&email, "email", "", "email for Discord account")
	p.FlagSet.StringVar(&email, "e", "", "email for Discord account")

	p.FlagSet.StringVar(&password, "password", "", "password for Discord account")
	p.FlagSet.StringVar(&password, "p", "", "password for Discord account")

	p.FlagSet.StringVar(&token, "token", "", "token for Discord account")
	p.FlagSet.StringVar(&token, "t", "", "token for Discord account")

	p.FlagSet.BoolVar(&debug, "d", true, "enable debug logging")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		ds, err := createDiscordSession(getconfig.AuthConfig{})
		if err != nil {
			logrus.Debugf("Session Failed \n%v\nexiting.", spew.Sdump(ds))
			err = fmt.Errorf("You may need to login from a browser first or check your credentials\n%v", err)
			panic(err, ds)
		}
		ctx = setContextValue(ctx, discordSessionKey, ds)

		// On ^C, or SIGTERM handle exit.
		signals := make(chan os.Signal, 0)
		signal.Notify(signals, os.Interrupt)
		signal.Notify(signals, syscall.SIGTERM)
		_, cancel := context.WithCancel(ctx)
		go func() {
			for sig := range signals {
				cancel()
				logrus.Infof("Received %s, exiting.", sig.String())
				os.Exit(0)
			}
		}()

		return nil
	}
	// Run our program.
	p.Run()
}
