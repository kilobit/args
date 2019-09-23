ARGS
====

A super simple argument processing library for Golang.  This library
was built out of curiosity and a desire for a simple argument library
that doesn't impose structure on the calling application.

With compliments to [spf13](https://spf13.com/) for
[Cobra](https://github.com/spf13/cobra), a really great library. 

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

In-Progress:
 - Get the currently unprocessed arguments.
 - A simple REPL.

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
