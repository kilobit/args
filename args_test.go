/* Copyright 2019 Kilobit Labs Inc. */

package args // import "kilobit.ca/go/args"

import "testing"
import "strings"

func Expect(t *testing.T, expected interface{}, actual interface{}) {

	if expected != actual {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestArgsTest(t *testing.T) {

	if true != true {
		t.Errorf("Sanity check failed.")
	}
}

// Test data for Option / Command / Option style args.
var ocoargs map[string][]string = map[string][]string{

	// Golden path test.
	"/foo/bar -1 -2 -3 -456 command -7 -89": toSlice("1", "2", "3", "4", "5", "6", "command", "7", "8", "9"),

	// - stops option parsing.
	"/foo/bar -123 foo -": toSlice("1", "2", "3", "foo"),
}

func toSlice(ss ...string) []string {
	return ss
}

func check(t *testing.T, ex []string, act []string) {

	for i, arg := range act {
		if arg != ex[i] {
			t.Errorf("Expected %s, Got %s.", ex[i], arg)
		}
	}
}

func TestOptCmdOpt(t *testing.T) {

	for sargs, argset := range ocoargs {
		args := strings.Split(sargs, " ")
		p := NewArgParser(args[1:])
		r := []string{}

		for arg := range p.NextOptC() {
			r = append(r, arg)
		}

		cmd := p.NextArg()
		r = append(r, cmd)

		for arg := range p.NextOptC() {
			r = append(r, arg)
		}

		check(t, argset, r)
	}
}

func TestArgParserArgs(t *testing.T) {

	args := strings.Split("This is a set of 7 arguments", " ")
	p := NewArgParser(args)

	if len(p.Args()) != 7 {
		t.Error("Expected 7 arguments remaining.")
	}

	p.NextArg()

	if len(p.Args()) != 6 {
		t.Error("Expected 6 arguments remaining.")
	}

	for range p.NextArgC() {
	}

	if len(p.Args()) != 0 {
		t.Error("Expected 0 arguments remaining.")
	}
}

func TestOptWithArgs(t *testing.T) {

	args := strings.Split("--fmt json -z mdt", " ")
	p := NewArgParser(args)

	opt1 := p.NextOpt()
	val1 := p.NextArg()

	Expect(t, "fmt", opt1)
	Expect(t, "json", val1)

	opt2 := p.NextOpt()
	val2 := p.NextArg()

	Expect(t, "z", opt2)
	Expect(t, "mdt", val2)
}

func TestOptWithArg(t *testing.T) {

	args := strings.Split("-nfoo remain", " ")
	p := NewArgParser(args)

	opt := p.NextOpt()
	arg := p.NextArg()

	Expect(t, "n", opt)
	Expect(t, "foo", arg)
	Expect(t, "remain", p.Args()[0])
}
