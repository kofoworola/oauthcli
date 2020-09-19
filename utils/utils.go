package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	stdinReader *bufio.Reader
)

func ReadLine(instruction string) (string, error) {
	if stdinReader == nil {
		stdinReader = bufio.NewReader(os.Stdin)
	}
	fmt.Printf("# %s: ", instruction)
	// TODO refractor this
	line, err := stdinReader.ReadString('\n')
	if err != nil {
		return strings.TrimSpace(line), fmt.Errorf("error reading from standard input: %v", err)
	}
	return strings.TrimSpace(line), nil
}

func ReadPassType(instruction string) (string, error) {
	fmt.Printf("# %s:", instruction)
	pass, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Printf("\n")
	return string(pass), err
}

func PrintError(err error) {
	fmt.Printf("ERROR: %s", err.Error())
}
