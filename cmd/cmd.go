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

func Any(value string) error { return nil }

func OneOf(opts ...string) Validator {

	m := map[string]struct{}{}
	for _, opt := range opts {
		m[opt] = struct{}{}
	}

	return func(value string) error {
		if _, ok := m[value]; !ok {
			return fmt.Errorf("Invalid parameter, %s", value)
		}
		return nil
	}
}

// Type representing an Opt arg or and Arg.
//
// A nil Validator when used as an Opt indicates a boolean flag.
//
// A nil validator when used as an Arg indicates that no validation
// will be done on the parameter value.
//
type Param struct {
	name string
	desc string
	v    Validator
}

func NewParam(name, desc string, v Validator) *Param {
	return &Param{name, desc, v}
}

// Type Opts is a mapping of options to names and aliases.
//
// Returns the Opts map for chaining.  Duplicate names / aliases will
// overwrite prior entries.
//
type Opts map[string]*Param

func NewOpts() Opts {
	return Opts{}
}

func (opts Opts) Add(param *Param, aliases ...string) Opts {

	names := append([]string{param.name}, aliases...)
	for _, name := range names {
		opts[name] = param
	}

	return opts
}

func (opts Opts) Merge(o Opts) Opts {

	for n, p := range o {
		opts[n] = p
	}

	return opts
}

// Type Params contains a set of arguments.
//
// Returns Args for chaining.  Use destructuring to concatenate
// argument lists.
//
type Args []*Param

func NewArgs(param ...*Param) *Args {
	args := (Args)(param)
	return &args
}

func (args *Args) Add(param ...*Param) *Args {

	*args = append(*args, param...)

	return args
}

func (args *Args) Concat(as *Args) *Args {

	params := ([]*Param)(*as)
	*args = append(*args, params...)

	return args
}

// Mapping of a Param name to it's assigned value.
//
type ValueMap map[string]string

// Interface for running Command actions.  Can also be satisfied with
// a CmdHandler.
//
type CmdRunner interface {
	RunCommand(params ValueMap, rest []string) error
}

// A function type that performs the Command action.  This type
// satisfies the CmdRunner interface.
//
// The params is a map of the named options and args.
//
type CmdRunnerFunc func(params ValueMap, rest []string) error

func (h CmdRunnerFunc) RunCommand(params ValueMap, rest []string) error {
	return h(params, rest)
}

// Type representing a command, it's parameters and meta data.
//
type Command struct {
	name   string
	desc   string
	opts   Opts
	args   *Args
	runner CmdRunner
}

func New(name, desc string, opts Opts, args *Args, runner CmdRunner) *Command {

	cmd := Command{name, desc, opts, args, runner}

	return &cmd
}

// Run a command via it's handler.
//
// A nil args will default to os.Args.
//
func (cmd *Command) Run(args []string) error {

	ap := NewArgParser(args)
	params := ValueMap{}

	for n := range ap.NextOptC() {
		opt, ok := cmd.opts[n]
		if !ok {
			return fmt.Errorf("Unknown option, '%s'.", n)
		}

		arg := ""
		if opt.v != nil {
			arg = ap.NextArg()
			err := opt.v(arg)
			if err != nil {
				return fmt.Errorf("Option error, %s: %s", n, err)
			}
		}

		params[opt.name] = arg
	}

	// Check number of arguments.

	for _, arg := range *cmd.args {

		param := ap.NextArg()
		if arg.v != nil {
			err := arg.v(param)
			if err != nil {
				return fmt.Errorf("Argument error, %s: %s", arg.name, err)
			}
		}

		params[arg.name] = param
	}

	return cmd.runner.RunCommand(params, ap.Args())
}
