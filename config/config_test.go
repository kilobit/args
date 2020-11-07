/* Copyright 2020 Kilobit Labs Inc. */

package config // import "kilobit.ca/go/args/config"

import "os"
import "io/ioutil"
import "strings"
import "testing"
import "kilobit.ca/go/tested/assert"
import "encoding/json"

func TestConfigTest(t *testing.T) {

	assert.Expect(t, true, true)
}

var loadTests = []struct {
	s      string
	exc    *Config
	haserr bool
}{
	{"{\"foo\":\"bar\"}", &Config{"foo": "bar"}, false},
	{"{\"foo\":\"bar\"}\n", &Config{"foo": "bar"}, false},
	{"{\"foo\":{\"bing\":\"bong\"}}", &Config{"foo": map[string]interface{}{"bing": "bong"}}, false},
	{"52", &Config{}, true},
	{"[1, 2, 3, 4,]", &Config{}, true},
}

func TestConfigWrite(t *testing.T) {

	for _, data := range loadTests {

		if data.haserr {
			continue
		}

		var w strings.Builder
		err := data.exc.Write(&w)
		assert.Ok(t, err)

		// t.Log(w.String())

		var obj interface{}
		err = json.Unmarshal([]byte(w.String()), &obj)
		assert.Ok(t, err)

		// t.Log(data.s)

		var exp interface{}
		err = json.Unmarshal([]byte(data.s), &exp)
		assert.Ok(t, err)

		assert.ExpectDeep(t, exp, obj)
	}
}

func TestConfigWriteFile(t *testing.T) {

	for _, data := range loadTests {

		if data.haserr {
			continue
		}

		f, err := ioutil.TempFile("", "writefile_test")
		if err != nil {
			t.Fatal(err)
		}

		//defer os.Remove(f.Name()) // Cleanup

		f.Close()

		if err := data.exc.WriteFile(f.Name()); err != nil {
			t.Fatal(err)
		}

		c, err := FromFile(f.Name())
		if err != nil {
			t.Fatal(err)
		}

		assert.ExpectDeep(t, data.exc, c)
	}
}

func TestConfigLoad(t *testing.T) {

	for _, data := range loadTests {

		r := strings.NewReader(data.s)

		c, err := Load(r)

		assert.ExpectDeep(t, data.exc, c)
		assert.Expect(t, data.haserr, err != nil)
	}
}

func TestConfigFromFile(t *testing.T) {

	for _, data := range loadTests {

		f, err := ioutil.TempFile("", "fromfile_test")
		if err != nil {
			t.Fatal(err)
		}

		defer os.Remove(f.Name()) // Cleanup

		if _, err := f.Write([]byte(data.s)); err != nil {
			t.Fatal(err)
		}

		c, err := FromFile(f.Name())

		assert.ExpectDeep(t, data.exc, c)
		assert.Expect(t, data.haserr, err != nil)
	}
}

func TestConfigWatch(t *testing.T) {

	data := loadTests[0]

	f, err := ioutil.TempFile("", "fromfile_test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(f.Name()) // Cleanup

	if _, err := f.Write([]byte(data.s)); err != nil {
		t.Fatal(err)
	}

	confs := Watch(f.Name())

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	proc.Signal(DEFAULT_SIGNAL)

	c := <-confs

	assert.ExpectDeep(t, data.exc, c)
	assert.Expect(t, data.haserr, err != nil)
}
