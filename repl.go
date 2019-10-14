/* Copyright 2019 Kilobit Labs Inc. */

package args // import "kilobit.ca/go/args"

import "bufio"
import "strings"
import "io"
import "os"

const default_prompt string = "> "

// Read-Eval-Print Loop for Args.
//
type REPL struct {
	br     *bufio.Reader
	w      io.Writer
	prompt string
}

// Options to modify the functioning of the REPL.
//
type REPLOpt func(repl *REPL)

// Set the prompt at REPL creation time.
//
func REPLOptPrompt(prompt string) REPLOpt {
	return func(repl *REPL) {
		repl.prompt = prompt
	}
}

// Interface for a REPL Handler.
//
type REPLHandler interface {
	HandleCmd(ap *ArgParser) bool
}

// A handler function type for processing repl commands.
//
// End processing by returning false.
//
type REPLHandlerFunc func(ap *ArgParser) bool

func (f REPLHandlerFunc) HandleCmd(ap *ArgParser) bool {
	return f(ap)
}

// Create a new REPL.
//
// Reader and writer will default to os.Stdin, os.Stdout if nil.
//
// Additional REPLOpt functions modify the behaviour of the REPL.
//
func NewREPL(r io.Reader, w io.Writer, opts ...REPLOpt) *REPL {

	if r == nil {
		r = os.Stdin
	}
	if w == nil {
		w = os.Stdout
	}

	repl := &REPL{
		bufio.NewReader(r),
		w,
		default_prompt,
	}

	for _, opt := range opts {
		opt(repl)
	}

	return repl
}

// Get the current prompt value.
//
func (repl *REPL) Prompt() string {
	return repl.prompt
}

// Set the prompt value.
//
// Set prompt to the empty string to supress the prompt.
//
func (repl *REPL) SetPrompt(prompt string) {
	repl.prompt = prompt
}

// Run the REPL.
//
func (repl *REPL) Run(handler REPLHandler) error {

	scanner := bufio.NewScanner(repl.br)

	var line string
	for {
		if repl.prompt != "" {
			repl.w.Write([]byte(repl.prompt))
		}

		r := scanner.Scan()
		if !r {
			return scanner.Err()
		}

		line = scanner.Text()
		args := strings.Fields(line)
		ap := NewArgParser(args)

		r = handler.HandleCmd(ap)
		if !r {
			break
		}
	}

	return nil
}
