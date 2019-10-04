/* Copyright 2019 Kilobit Labs Inc. */

package args // import "kilobit.ca/go/args"

import "testing"
import "strings"

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
