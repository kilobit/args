ARGS
====

A super simple argument processing library for Golang.

Status: In-Development

```{.go}
import "kilobit.ca/args"

func main() {

	p := args.NewArgParser(nil)
	
	for opt := range p.NextOptC() {
		handle(opt)
	}
	
	cmd := p.NextArg()
	
	for opt := p.NextOpt(); opt != ""; opt = p.NextOpt() {
		handle_more_opts(opt)
	}
	
	docmd(cmd)
}
```

Features
--------

- Simple,  unopinionated module that fits your code structure.
- GNU style option handling.
- Channel or iterator style interface.

Installation
------------

```{.bash}
go get kilobit.ca/args
```

Building
--------

```{.bash}
cd kilobit.ca/args
go test -v
go build
```

Contribute
----------

Please submit a pull request with any bug fixes or feature requests
that you have.  All submissions imply consent to use / distribute
under the terms of the LICENSE.

Support
-------

Submit tickets through [github](https://github.com/kilobit/args).

License
-------

See LICENSE.

--
Created: Sept 23, 2019
By: Christian Saunders <cps@kilobit.ca>
Copyright 2019 Kilobit Labs Inc.
