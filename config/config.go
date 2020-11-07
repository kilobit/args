/* Copyright 2020 Kilobit Labs Inc. */

// Simple argument parsing for Golang.
//
// Configuration parser
//
package config // import "kilobit.ca/go/args/config"

import "io"
import "os"
import "os/signal"
import "encoding/json"

type Config map[string]interface{}

func (c *Config) Write(w io.Writer) error {

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(c)
}

func (c *Config) WriteFile(filename string) error {

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {

		return err
	}

	defer f.Close()

	return c.Write(f)
}

func Load(r io.Reader) (*Config, error) {

	c := Config{}
	dec := json.NewDecoder(r)
	err := dec.Decode(&c)

	return &c, err
}

func FromFile(filename string) (*Config, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return Load(f)
}

// Re-parse a config when a signal is received.
//
// Default signal is SIGUSR1.  Returns a channel of Config objects, on
// error a nil pointer will be passed on the channel and no other
// signals will be processed.
//
// On windows SIGHUP will be used instead.
//
func Watch(filename string, sig ...os.Signal) chan *Config {

	confs := make(chan *Config, 1)
	if len(sig) == 0 {
		sig = []os.Signal{DEFAULT_SIGNAL}
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, sig...)

	go func() {

		<-sigs

		c, err := FromFile(filename)
		if err != nil {
			signal.Reset(sig...)
			return
		}

		confs <- c
	}()

	return confs
}
