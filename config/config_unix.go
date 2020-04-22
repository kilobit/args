// +build !windows

/* Copyright 2020 Kilobit Labs Inc. */

package config // import "kilobit.ca/go/args/config"

import "golang.org/x/sys/unix"

const DEFAULT_SIGNAL = unix.SIGUSR1
