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
	"io/ioutil"
	"os"
	"os/signal"
	"os/user"
	"syscall"

	"github.com/azillion/chord/internal/getconfig"
	"github.com/azillion/chord/version"
	"github.com/bwmarrin/discordgo"
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

const configFile string = ".chord.config"

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

	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
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

		// Set the log level.
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		return nil
	}
	// Run our program.
	p.Run()
}

func createDiscordSession(authConfig getconfig.AuthConfig) (*discordgo.Session, error) {
	var dToken string
	if email != "" && password != "" {
		ds, err := createDiscordSessionFromLogin(email, password)
		if err != nil {
			return new(discordgo.Session), err
		}
		return ds, nil
	}

	// everything besides configcmd will use this
	if authConfig.Token == "" && authConfig.Email == "" && authConfig.Password == "" {
		// read the token in from a file
		dTokenBytes, _ := ioutil.ReadFile(configPath)
		if len(dTokenBytes) > 0 {
			dToken = string(dTokenBytes)
		}

		// if file token is not empty
		if dToken != "" {
			ds, err := createDiscordSessionFromToken(dToken)
			if err != nil {
				return new(discordgo.Session), err
			}

			// create token file
			file, err := os.Create(configPath)
			if err != nil {
				return new(discordgo.Session), err
			}
			defer file.Close()

			// write token to token file
			if _, err := file.WriteString(ds.Token); err != nil {
				return new(discordgo.Session), err
			}
			return ds, nil
		}
		return new(discordgo.Session), fmt.Errorf("empty auth credentials provided, try 'chord config'")
	}

	// if a token is passed in
	if authConfig.Token != "" {
		ds, err := createDiscordSessionFromToken(authConfig.Token)
		if err != nil {
			return new(discordgo.Session), err
		}
		return ds, nil
	}

	// if an email and password are passed in, config will use this
	ds, err := createDiscordSessionFromLogin(authConfig.Email, authConfig.Password)
	if err != nil {
		return new(discordgo.Session), err
	}
	// save the token to the token file
	// create token file
	file, err := os.Create(configPath)
	if err != nil {
		return new(discordgo.Session), err
	}
	defer file.Close()

	// write token to token file
	if _, err := file.WriteString(ds.Token); err != nil {
		return new(discordgo.Session), err
	}
	return ds, nil
}

func createDiscordSessionFromToken(authToken string) (*discordgo.Session, error) {
	if token != "" {
		authToken = token
	}
	if authToken == "" {
		return new(discordgo.Session), fmt.Errorf("Empty auth token provided, try 'chord config'")
	}
	ds, err := discordgo.New(authToken)
	if err != nil {
		return new(discordgo.Session), err
	}
	dUser, err = ds.User("@me")
	if err != nil {
		return new(discordgo.Session), err
	}
	logrus.Debugf("Logged in as %s\n", dUser.String())
	return ds, nil
}

func createDiscordSessionFromLogin(emailIn, passwordIn string) (*discordgo.Session, error) {
	if email != "" {
		emailIn = email
	}
	if password != "" {
		passwordIn = password
	}
	if emailIn == "" || passwordIn == "" {
		return new(discordgo.Session), fmt.Errorf("No email or password entered, try 'chord config'")
	}

	// Create a new Discord session using the provided login information.
	ds, err := discordgo.New(emailIn, passwordIn)
	if err != nil {
		return new(discordgo.Session), err
	}
	dUser, err = ds.User("@me")
	if err != nil {
		return new(discordgo.Session), err
	}
	logrus.Debugf("Logged in as %s\n", dUser.String())
	return ds, nil
}
