package utils

import (
	"bufio"
	"fmt"
	"os"
)

var (
	stdinReader *bufio.Reader
)

func ReadLine(instruction string) (string, error) {
	if stdinReader == nil {
		stdinReader = bufio.NewReader(os.Stdin)
	}
	fmt.Printf("# %s:\n> ", instruction)
	// TODO refractor this
	line, err := stdinReader.ReadString('\n')
	if err != nil {
		return line, fmt.Errorf("error reading from standard input: %v", err)
	}
	return line, nil
}
