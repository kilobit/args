/* Copyright 2019 Kilobit Labs Inc. */

package args // import "kilobit.ca/go/args"

import "bufio"
import "strings"
import "io"

// Read-Eval-Print Loop for Args.
//
type REPL struct {
	br *bufio.Reader
}

// A handler function type for processing repl commands.
//
// End processing by returning false.
//
type REPLHandler func(ap *ArgParser) bool

// Create a new REPL by passing in a reader from which lines will be
// read.
//
func NewREPL(r io.Reader) *REPL {

	return &REPL{
		bufio.NewReader(r),
	}
}

func (repl *REPL) Run(f REPLHandler) error {

	scanner := bufio.NewScanner(repl.br)

	var line string
	for {
		r := scanner.Scan()
		if !r {
			return scanner.Err()
		}

		line = scanner.Text()
		args := strings.Fields(line)
		ap := NewArgParser(args)

		r = f(ap)
		if !r {
			break
		}
	}

	return nil
}
