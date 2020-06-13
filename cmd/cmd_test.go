/* Copyright 2019 Kilobit Labs Inc. */

package cmd // import "kilobit.ca/go/args/cmd"

import "testing"
import "kilobit.ca/go/tested/assert"

func TestCommandTest(t *testing.T) {

	assert.Expect(t, true, true)
}

func cmdHandler(params ValueMap, rest []string) error {
	return nil
}

func TestCommandRun(t *testing.T) {

	cmd := New("test", "A test command.", cmdHandler)
	err := cmd.Run([]string{})
	if err != nil {
		t.Error(err)
	}
}

func TestCommandRunArgs(t *testing.T) {

	cmd := New("test", "Another test command.", func(params ValueMap, rest []string) error {

		return nil
	})

	cmd.AddArgs(NewArg("foo", "a whole lot of foo", nil))

}
