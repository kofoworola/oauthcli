package clients

import (
	"bufio"
	"fmt"
	"net/url"
	"strings"
)

type TestUtil struct {
	reader *bufio.Reader
	// openURL will be used to compare against whatever parameter is passed to OpenURL
	openURL string
	// input contains ReadLine input
	// each entry to ReadLine should be arranged in order and separated by a line break
	input string
	// password is what will be returned when read password is called
	password string
}

func (u *TestUtil) OpenURL(open string) error {
	expectedURL, err := url.Parse(u.openURL)
	if err != nil {
		return fmt.Errorf("could not parse url %s: %w", u.openURL, err)
	}

	gotURL, err := url.Parse(open)
	if err != nil {
		return fmt.Errorf("could not parse url %s: %w", open, err)
	}
	gotParams := gotURL.Query()

	expectedResolved := fmt.Sprintf("%s//%s/%s", expectedURL.Scheme, expectedURL.Host, expectedURL.Path)
	gotResolved := fmt.Sprintf("%s//%s/%s", gotURL.Scheme, gotURL.Host, gotURL.Path)
	if expectedResolved != gotResolved {
		return fmt.Errorf("expcted URI %s got %s", expectedResolved, gotResolved)
	}
	// check params
	for i, q := range expectedURL.Query() {
		val, ok := gotParams[i]
		if !ok {
			return fmt.Errorf("%s does not exist as part of query parameters for %s", i, open)
		}
		for _, x := range q {
			found := false
			for _, y := range val {
				if x == y {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("val %s missing in query paramaeters for %s", x, i)
			}
		}
	}
	return nil
}

func (u *TestUtil) ReadLine(instruction string, isPassword bool) (string, error) {
	if isPassword {
		return u.password, nil
	}
	if u.reader == nil {
		u.reader = bufio.NewReader(strings.NewReader(u.input))
	}

	line, err := u.reader.ReadString('\n')
	return strings.TrimSpace(line), err
}
