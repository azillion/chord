package getconfig

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

type AuthConfig struct {
	Email, Password string
}

// GetAuthConfig returns the Discord AuthConfig.
// Optionally takes in the authentication values, otherwise pulls them from a
// config file.
func GetAuthConfig(email, password string) (AuthConfig, error) {
	if email != "" && password != "" {
		return AuthConfig{
			Email:    email,
			Password: password,
		}, nil
	}

	logrus.Debugf("TODO: Handle usage of config file")

	reader := bufio.NewReader(os.Stdin)

	if email == "" {
		fmt.Print("Enter Discord Email: ")
		email_in, err := reader.ReadString('\n')
		if err != nil {
			return AuthConfig{}, err
		}
		email = email_in
	}

	if password == "" {
		fmt.Print("Enter Discord Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return AuthConfig{}, err
		}
		password = string(bytePassword)
	}

	email, password = strings.TrimSpace(email), strings.TrimSpace(password)
	return AuthConfig{
		Email:    email,
		Password: password,
	}, nil
}
