package keystore

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func PromptMissing(opts *Options) error {
	reader := bufio.NewReader(os.Stdin)

	if strings.TrimSpace(opts.KeystoreName) == "" {
		fmt.Print("Keystore name [upload-keystore]: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		opts.KeystoreName = strings.TrimSpace(line)
	}

	if strings.TrimSpace(opts.StorePassword) == "" {
		password, err := readPassword("Store password: ")
		if err != nil {
			return err
		}
		opts.StorePassword = password
	}

	if strings.TrimSpace(opts.KeyPassword) == "" {
		password, err := readPassword("Key password (blank = store password): ")
		if err != nil {
			return err
		}
		opts.KeyPassword = password
	}

	if strings.TrimSpace(opts.KeyAlias) == "" {
		fmt.Print("Key alias [upload]: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		opts.KeyAlias = strings.TrimSpace(line)
	}

	if strings.TrimSpace(opts.CommonName) == "" {
		fmt.Print("Certificate common name [Android Upload]: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		opts.CommonName = strings.TrimSpace(line)
	}

	return nil
}

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	if term.IsTerminal(int(os.Stdin.Fd())) {
		bytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
