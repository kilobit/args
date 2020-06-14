/* Copyright 2019 Kilobit Labs Inc. */

package cmd // import "kilobit.ca/go/args/cmd"

import "testing"
import "kilobit.ca/go/tested/assert"
import "errors"

func TestCommandTest(t *testing.T) {

	assert.Expect(t, true, true)
}

var noprunner CmdRunnerFunc = func(params ValueMap, rest []string) error {
	return nil
}

func TestCommandRun(t *testing.T) {

	cmd := New("test", "A test command.", noprunner)
	err := cmd.Run([]string{})
	if err != nil {
		t.Error(err)
	}
}

func TestCommandRunArgs(t *testing.T) {

	var runner CmdRunnerFunc = func(params ValueMap, rest []string) error {

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
	}

	cmd := New("test", "Another test command.", runner)

	cmd.AddOpt(NewParam("verbose", "Test verbose output.", nil), "v")
	cmd.AddOpt(NewParam("foo", "a whole lot of foo", Any), "f")
	cmd.AddArgs(NewParam("hello", "A test argument", nil))

	err := cmd.Run([]string{"--verbose", "-f", "bar", "Hello World!", "rest1", "rest2"})
	assert.ExpectDeep(t, errors.New(""), err)
}

func TestOneOfValidator(t *testing.T) {

	v := OneOf("one", "two", "three")

	if err := v("one"); err != nil {
		t.Error(err)
	}

	if err := v("two"); err != nil {
		t.Error(err)
	}

	if err := v("three"); err != nil {
		t.Error(err)
	}

	if err := v("four"); err == nil {
		t.Error("Expected error for invalid input, 'four'.")
	}

}

func TestOneOfEmptyValidator(t *testing.T) {

	v := OneOf()

	if err := v("any"); err == nil {
		t.Error("Expected error for invalid input, 'any'.")
	}
}

func TestOpts(t *testing.T) {

	opt1 := NewParam("opt1", "Test Option 1", nil)
	opt2 := NewParam("opt2", "Test Option 2", nil)
	opt3 := NewParam("opt3", "Test Option 3", nil)

	opts := NewOpts().Add(opt1, "1", "2")

	opts.Add(opt2, "2", "3")

	opts2 := NewOpts().Add(opt3, "3", "4")

	opts.Merge(opts2)

	//t.Logf("%#v", opts)

	assert.Expect(t, opts["opt1"], opt1)
	assert.Expect(t, opts["opt2"], opt2)
	assert.Expect(t, opts["1"], opt1)
	assert.Expect(t, opts["2"], opt2)
	assert.Expect(t, opts["3"], opt3)
	assert.Expect(t, opts["4"], opt3)
}

func TestArgs(t *testing.T) {

	arg1 := NewParam("arg1", "First Argument.", nil)
	arg2 := NewParam("arg2", "Second Argument.", nil)
	arg3 := NewParam("arg3", "Third Argument.", nil)
	arg4 := NewParam("arg4", "Fourth Argument.", nil)

	args := NewArgs(arg1).Add(arg2, arg3)

	args.Add(arg4)

	args.Concat(args)

	//t.Logf("%#v", args)

	assert.Expect(t, (*args)[0], arg1)
	assert.Expect(t, (*args)[1], arg2)
	assert.Expect(t, (*args)[2], arg3)
	assert.Expect(t, (*args)[3], arg4)
	assert.Expect(t, (*args)[4], arg1)
	assert.Expect(t, (*args)[5], arg2)
	assert.Expect(t, (*args)[6], arg3)
	assert.Expect(t, (*args)[7], arg4)
}
