/* Copyright 2019 Kilobit Labs Inc. */

// Simple argument parsing for Golang.
//
package args // import "kilobit.ca/go/args"

import "os"
import "strings"
import "unicode/utf8"

type ArgParser struct {
	shopt bool
	si    uint
	args  []string
}

func NewArgParser(args []string) *ArgParser {
	if args == nil {
		args = os.Args[1:]
	}

	return &ArgParser{false, 0, args}
}

// Get the remaining arguments.
//
func (ap *ArgParser) Args() []string {

	args := make([]string, len(ap.args))
	copy(args, ap.args)

	return args
}

// Peek at the next argument.
//
func (ap *ArgParser) PeekArg() string {
	if len(ap.args) == 0 {
		return ""
	}

	return ap.args[0][ap.si:]
}

// Returns the whole next argument or the empty string.
//
// Note that if NextArg is called after a short option, it will
// terminate the short option parsing and treat the remaining
// characters as the following argument.
//
// e.g. -nfoo would be interpreted as 'foo' if NextArg is called
// immediately after the NextOpt return 'n'.
//
func (ap *ArgParser) NextArg() string {
	ap.shopt = false

	if len(ap.args) == 0 {
		return ""
	}

	if len(ap.args[0]) <= (int)(ap.si) {
		ap.args = ap.args[1:]
		ap.si = 0
	}

	arg := ap.args[0][ap.si:]
	ap.args = ap.args[1:]

	ap.si = 0

	return arg
}

// Channel version of the NextArg method, useful for looping.
//
func (ap *ArgParser) NextArgC() chan string {
	c := make(chan string)

	go func() {
		for {
			arg := ap.NextArg()
			if arg == "" {
				break
			}
			c <- arg
		}

		close(c)
	}()

	return c
}

// Returns the next option or the empty string.
//
// Note: Supports gnu style long and short opts.
//
func (ap *ArgParser) NextOpt() string {
	if len(ap.args) == 0 {
		return ""
	}

	if ap.shopt {
		return ap.nextShopt()
	}

	arg := ap.args[0]
	switch {

	case strings.HasPrefix(arg, "--"):
		opt := arg[2:]
		ap.args = ap.args[1:]
		return opt

	case strings.HasPrefix(arg, "-"):
		ap.shopt = true
		ap.si = 1
		return ap.nextShopt()
	default:
		return ""
	}
}

// Helper function for handling short opts.
//
func (ap *ArgParser) nextShopt() string {

	if len(ap.args) == 0 {
		return ""
	}

	arg := ap.args[0]

	// Test if there are more shopts.
	if uint(len(arg)) <= ap.si {
		ap.si = 0
		ap.shopt = false
		ap.args = ap.args[1:]
		return ap.NextOpt()
	}

	r, size := utf8.DecodeRuneInString(arg[ap.si:])
	ap.si += uint(size)

	opt := string(r)
	return opt
}

// Channel version of the NextOpt method, useful for looping.
//
func (ap *ArgParser) NextOptC() chan string {
	c := make(chan string)

	go func() {
		for {
			opt := ap.NextOpt()
			if opt == "" {
				break
			}
			c <- opt
		}

		close(c)
	}()

	return c
}
