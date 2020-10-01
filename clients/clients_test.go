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
}

// TODO move this in
func (u *TestUtil) OpenURL(url string) error {
	return u.compareURLS(u.openURL, url)
}

func (u *TestUtil) ReadLine(instruction string, isPassword bool) (string, error) {
	if u.reader == nil {
		u.reader = bufio.NewReader(strings.NewReader(u.input))
	}

	line, err := u.reader.ReadString('\n')
	return strings.TrimSpace(line), err
}

func (u *TestUtil) compareURLS(expected, got string) error {
	expectedURL, err := url.Parse(expected)
	if err != nil {
		return fmt.Errorf("could not parse url %s: %w", expected, err)
	}

	gotURL, err := url.Parse(got)
	if err != nil {
		return fmt.Errorf("could not parse url %s: %w", got, err)
	}
	gotParams := gotURL.Query()

	for i, q := range expectedURL.Query() {
		val, ok := gotParams[i]
		if !ok {
			return fmt.Errorf("%s does not exist as part of query parameters for %s", i, got)
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
