/* Copyright 2020 Kilobit Labs Inc. */

// Simple argument parsing for Golang.
//
// Configuration parser
//
package config // import "kilobit.ca/go/args/config"

import "io"
import "os"
import "encoding/json"
import "log"

var logger *log.Logger = log

type Config map[string]interface{}

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
