/* Copyright 2020 Kilobit Labs Inc. */

// High level command line processing utilities.
//
package cmd // import "kilobit.ca/go/args/cmd"

import "fmt"
import _ "errors"

import "strings"
import . "kilobit.ca/go/args"

// Check a given parameter to make sure it is of the correct type etc.
//
type Validator func(value string) error

// Returns valid for any string.
func Any(value string) error { return nil }

// Returns valid if the argument matches one of the provided opts.
//
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

// Parse a command line command.
//
// A nil args will default to os.Args.
//
func Parse(cmd []string, opts Opts, args *Args) (params ValueMap, rest []string, err error) {

	ap := NewArgParser(cmd)
	params = ValueMap{}

	for {
		o := ap.NextOpt()
		if o == "" {
			break
		}

		opt, ok := opts[o]
		if !ok {
			return nil, nil, fmt.Errorf("Unknown option, '%s'.", o)
		}

		arg := ""
		if opt.v != nil {
			arg = ap.NextArg()
			err := opt.v(arg)
			if err != nil {
				return nil, nil, fmt.Errorf("Option error, %s: %s", o, err)
			}
		}

		params[opt.name] = arg
	}

	// Check number of arguments.
	if len(ap.Args()) < len(*args) {
		return nil, nil, fmt.Errorf("Insufficient arguments, expected %d, got %d", len(*args), len(ap.Args()))
	}

	for _, arg := range *args {

		param := ap.NextArg()
		if arg.v != nil {
			err := arg.v(param)
			if err != nil {
				return nil, nil, fmt.Errorf("Argument error, %s: %s", arg.name, err)
			}
		}

		params[arg.name] = param
	}

	return params, ap.Args(), nil
}

// Usage
func Usage(name, desc string, opts Opts, args *Args) string {

	sb := &strings.Builder{}

	fmt.Fprintf(sb, "\n%s - %s\n\n", name, desc)

	fmt.Fprintf(sb, "Usage: %s", name)
	if len(opts) > 0 {
		fmt.Fprint(sb, " [OPTIONS]")
	}
	for _, arg := range *args {
		fmt.Fprintf(sb, " %s", arg.name)
	}
	fmt.Fprintf(sb, "\n\n")

	fmt.Fprintf(sb, "Where OPTIONS include:\n\n")

	for name, opt := range opts {

		if len(name) == 1 {
			fmt.Fprintf(sb, "\t-%-10s\t%s\n", name, opt.desc)
		} else {
			fmt.Fprintf(sb, "\t--%-10s\t%s\n", name, opt.desc)
		}

	}

	fmt.Fprintf(sb, "\nArguments:\n\n")

	for _, arg := range *args {
		fmt.Fprintf(sb, "\t%-10s\t%s\n", arg.name, arg.desc)
	}

	return sb.String()
}
