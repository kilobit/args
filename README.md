ARGS
====

A super simple argument processing library for Golang.

This library was built out of curiosity and a desire for a simple
argument library that doesn't impose structure on the calling
application.

With compliments to [spf13](https://spf13.com/) for
[Cobra](https://github.com/spf13/cobra), a really great library. 

Status: In-Development

```{.go}
import "kilobit.ca/go/args"

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

```{.go}
import "kilobit.ca/go/args"

func echo(ap *args.ArgParser) bool {

	var args []string
	for arg := range ap.NextArgC {
		args = append(args, arg)
	}
	
	fmt.Println(args)
	
	return true
}

func main() {
	
	repl := args.NewREPL(nil, nil, REPLOptPrompt("$ "))
	
	repl.Run(echo)
}
```

```{.go}
import "kilobit.ca/go/args/config"

func main() {

	c, _ := config.FromFile("myconfig.json")

	confs := config.Watch("myconfig.json")

	for {
		conf := <- confs

		// update running config
	}

}
```

```{.go}
import "kilobit.ca/go/args/config"

func main() {

	c, _ := config.FromFile("myconfig.json")

	c['foo'] = "bar"
	
	c.WriteFile("myconfig.json")
}
```


Features
--------

 - Simple,  unopinionated module that fits your code structure.
 - GNU style option handling.
 - Channel or iterator style interface.
 - A simple REPL.
 - Get the currently unprocessed arguments.
 - Load json based configuration files.
 - Signal based configuration watcher.
 - Write configs back to json files.
 - [ ] Use non-json files for configs.

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
