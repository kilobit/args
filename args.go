/* Copyright 2019 Kilobit Labs Inc. */

package args

import "os"
import "strings"
import "unicode/utf8"

type ArgParser struct {
	shopt bool
	si uint
	args []string
}

func NewArgParser(args []string) *ArgParser {
	if args == nil {
		args = os.Args[1:]
	}
	
	return &ArgParser{false, 0, args}
}

// Peek at the next argument.
//
func (ap *ArgParser) PeekArg() string {
	if len(ap.args) == 0 {
		return ""
	}

	return ap.args[0]
}

// Returns the whole next argument or nil.
//
func (ap *ArgParser) NextArg() string {
	if len(ap.args) == 0 {
		return ""
	}

	arg := ap.args[0]
	ap.args = ap.args[1:]

	return arg
}

// Channel version of the NextArg method, useful for looping.
//
func(ap *ArgParser) NextArgC() chan string {
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

// Returns the next option or nil.
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
func(ap *ArgParser) NextOptC() chan string {
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
