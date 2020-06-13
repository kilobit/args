/* Copyright 2019 Kilobit Labs Inc. */

package cmd // import "kilobit.ca/go/args/cmd"

import "testing"
import "kilobit.ca/go/tested/assert"
import "errors"

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

		//t.Logf("params: %s", params)
		//t.Logf("rest:   %s", rest)

		v, ok := params["verbose"]
		if !ok {
			t.Errorf("Missing expected option, verbose.")
		}
		assert.Expect(t, "", v)

		f, ok := params["foo"]
		if !ok {
			t.Errorf("Missing expected option, foo.")
		}
		assert.Expect(t, "bar", f)

		hw, ok := params["hello"]
		if !ok {
			t.Error("Missing expected argument, hw.")
		}

		assert.Expect(t, "Hello World!", hw)

		assert.ExpectDeep(t, []string{"rest1", "rest2"}, rest)

		return errors.New("")
	})

	cmd.AddOpt(NewParam("verbose", "Test verbose output.", nil), "v")
	cmd.AddOpt(NewParam("foo", "a whole lot of foo", Any), "f")
	cmd.AddArgs(NewParam("hello", "A test argument", nil))

	err := cmd.Run([]string{"--verbose", "-f", "bar", "Hello World!", "rest1", "rest2"})
	assert.ExpectDeep(t, errors.New(""), err)
}
