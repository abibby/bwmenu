package dmenu

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Option func(*Options) *Options

type Options struct {
	ReturnIndex bool // -i           return index of selected element
	// -v           show choose version
	NumRows         int    // -n [10]      set number of rows
	Width           int    // -w [50]      set width of choose window
	Font            string // -f [Menlo]   set font used by choose
	FontSize        int    // -s [26]      set font size used by choose
	HighlightColor  string // -c [0000FF]  highlight color for matched string
	BackgroundColor string // -b [222222]  background color of selected element
	// -u           disable underline and use background for matched string
	ReturnNonMatches bool // -m           return the query string in case it doesn't match any item
	// -p           defines a prompt to be displayed when query field is empty
	// -o           given a query, outputs results to standard output
}

func (o *Options) Command() (string, []string) {
	args := []string{}
	if o.ReturnNonMatches {
		args = append(args, "-m")
	}
	return "choose", args
}

func Open(items []string, options ...Option) (string, error) {
	opts := &Options{}
	for _, o := range options {
		opts = o(opts)
	}
	name, args := opts.Command()
	cmd := exec.Command(name, args...)
	cmd.Stdin = bytes.NewBufferString(strings.Join(items, "\n"))
	b, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("dmenu: %w", err)
	}
	return string(b), nil
}

func ReturnNonMatches() Option {
	return func(o *Options) *Options {
		o.ReturnNonMatches = true
		return o
	}
}
