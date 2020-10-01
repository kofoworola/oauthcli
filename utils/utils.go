package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	stdinReader *bufio.Reader
)

type Util interface {
	ReadLine(instruction string, isPassword bool) (string, error)
	OpenURL(url string) error
}

type MainUtil struct {
	reader *bufio.Reader
}

// ReadLine tries to read the next line from the reader (usually stdin except test)
// it should try to read 3 times before returning the error
func (u *MainUtil) ReadLine(instruction string, isPassword bool) (string, error) {
	if u.reader == nil {
		u.reader = bufio.NewReader(os.Stdin)
	}

	var (
		line []byte
		err  error
	)
	max := 3
	count := 0
	for (len(line) == 0 || err != nil) && count < max {
		count++
		fmt.Printf("# %s: ", instruction)
		if isPassword {
			line, err = terminal.ReadPassword(int(syscall.Stdin))
			fmt.Println("")
		} else {
			line, err = u.reader.ReadBytes('\n')
		}
		if err != nil {
			PrintError(fmt.Errorf("error reading input: %w", err))
		}
	}
	return strings.TrimSpace(string(line)), err
}

func PrintError(err error) {
	fmt.Printf("ERROR: %s", err.Error())
}

func (u *MainUtil) OpenURL(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}
