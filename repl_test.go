/* Copyright 2019 Kilobit Labs Inc. */

package args // import "kilobit.ca/go/args"

import "testing"
import "strings"

func TestREPLTest(t *testing.T) {

	if true != true {
		t.Errorf("Sanity check failed.")
	}
}

const replscript string = `This is the first line
this is the second line
`

func TestREPL(t *testing.T) {

	r := strings.NewReader(replscript)
	repl := NewREPL(r)

	err := repl.Run(func(ap *ArgParser) bool {

		var args []string
		for arg := range ap.NextArgC() {
			args = append(args, arg)
		}

		t.Logf("%#v", args)

		if len(args) != 5 {
			t.Error("Wrong number of arguments in repl.")
		}

		return true
	})

	if err != nil {
		t.Error(err)
	}
}
