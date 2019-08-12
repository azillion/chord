package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/azillion/chord/internal/getconfig"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

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
