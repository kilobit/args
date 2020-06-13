/* Copyright 2020 Kilobit Labs Inc. */

// Command processor and runner.
//
package cmd // import "kilobit.ca/go/args/cmd"

import "fmt"
import _ "errors"

import . "kilobit.ca/go/args"

// Check a given parameter to make sure it is of the correct type etc.
//
type Validator func(value string) error

// Type representing an Opt arg or and Arg
//
type Param struct {
	name string
	desc string
	v    Validator
}

func NewArg(name, desc string, v Validator) *Param {
	return &Param{name, desc, v}
}

// Type representing an Option.
//
type Option struct {
	hasArg bool
	Param
}

// Mapping of a Param name to it's assigned value.
//
type ValueMap map[string]string

// A function type that performs the Command action.
//
// The params is a map of the named options and args.
//
type CmdHandler func(params ValueMap, rest []string) error

// Type representing a command, it's parameters and meta data.
//
type Command struct {
	name    string
	desc    string
	opts    map[string]*Option
	args    []*Param
	handler CmdHandler
}

func (cmd *Command) HasOpt(name string) (*Option, bool) {
	p, ok := cmd.opts[name]
	return p, ok
}

// Add and option to the command.
//
// Note that the Option name will be added automatically.
//
func (cmd *Command) AddOpt(opt *Option, aliases ...string) error {

	names := append([]string{opt.name}, aliases...)

	for _, name := range names {
		if _, ok := cmd.opts[name]; ok {
			return fmt.Errorf("Duplicate option, %s", name)
		}

		cmd.opts[name] = opt
	}

	return nil
}

func (cmd *Command) AddArgs(args ...*Param) {

	cmd.args = append(cmd.args, args...)
}

func New(name, desc string, handler CmdHandler) *Command {

	cmd := Command{name, desc, map[string]*Option{}, []*Param{}, handler}

	return &cmd
}

func (cmd *Command) Run(args []string) error {

	ap := NewArgParser(args)
	params := map[string]string{}

	for n := range ap.NextOptC() {
		opt, ok := cmd.opts[n]
		if !ok {
			return fmt.Errorf("Unknown option, '%s'.", n)
		}

		arg := ""
		if opt.hasArg {
			arg = ap.NextArg()
		}

		err := opt.v(arg)
		if err != nil {
			return fmt.Errorf("Option error, %s: %s", n, err)
		}

		params[n] = arg
	}

	// Check number of arguments.

	for _, arg := range cmd.args {

		param := ap.NextArg()
		err := arg.v(param)
		if err != nil {
			return fmt.Errorf("Argument error, %s: %s", arg.name, err)
		}

		params[arg.name] = param
	}

	return cmd.handler(params, ap.Args())
}
