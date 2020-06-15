/* Copyright 2019 Kilobit Labs Inc. */

package cmd // import "kilobit.ca/go/args/cmd"

import "testing"
import "kilobit.ca/go/tested/assert"

func TestCommandTest(t *testing.T) {

	assert.Expect(t, true, true)
}

func TestParseEmptyCmd(t *testing.T) {

	params, rest, err := Parse([]string{}, NewOpts(), NewArgs())
	if err != nil {
		t.Error(err)
	}

	assert.ExpectDeep(t, ValueMap{}, params)
	assert.ExpectDeep(t, []string{}, rest)
}

func TestParseCmd(t *testing.T) {

	opts := NewOpts().
		Add(NewParam("verbose", "Test verbose output.", nil), "v").
		Add(NewParam("foo", "a whole lot of foo", Any), "f")

	args := NewArgs(NewParam("hello", "A test argument", nil))

	params, rest, err := Parse(
		[]string{"--verbose", "-f", "bar", "Hello World!", "rest1", "--rest2", "rest3"},
		opts, args,
	)
	if err != nil {
		t.Error(err)
	}

	assert.Expect(t, "", params["verbose"])
	assert.Expect(t, "bar", params["foo"])
	assert.Expect(t, "Hello World!", params["hello"])

	assert.ExpectDeep(t, []string{"rest1", "--rest2", "rest3"}, rest)
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
