package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/profx5/jordi/internal/app"
	"github.com/profx5/jordi/internal/config"
	"github.com/profx5/jordi/internal/version"
)

var (
	exit = os.Exit

	flags        = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	help         = flags.Bool("help", false, "Print usage instructions and exit.")
	printVersion = flags.Bool("version", false, "Print version and exit.")
	insecure     = flags.Bool("insecure", false, `Skip TLS certificate verification. (NOT SECURE!)`)
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
%s [flags] [address] [method]

The 'address' will typically be in the form "host:port" where host can be an IP
address or a hostname and port is a numeric port or service name.

The optional 'method' is the fully qualified name of the method to invoke, in the form
"package.Service/Method".
If the method is not specified, the user will be prompted to select a service and method.

Available flags:
`, os.Args[0])
	flags.PrintDefaults()
}

func fail(err error, msg string, args ...interface{}) {
	if err != nil {
		msg += ": %v"
		args = append(args, err)
	}
	fmt.Fprintf(os.Stderr, msg, args...)
	fmt.Fprintln(os.Stderr)
	if err != nil {
		exit(1)
	} else {
		// nil error means it was CLI usage issue
		fmt.Fprintf(os.Stderr, "Try '%s -help' for more details.\n", os.Args[0])
		exit(2)
	}
}

func main() {
	flags.Usage = usage
	err := flags.Parse(os.Args[1:])
	if err != nil {
		fail(err, "Failed to parse flags")
	}

	if *help {
		usage()
		os.Exit(0)
	}
	if *printVersion {
		fmt.Fprintf(os.Stderr, "%s %s\n", filepath.Base(os.Args[0]), version.Version)
		os.Exit(0)
	}

	args := flags.Args()

	var target, method string
	switch len(args) {
	case 0:
		fail(nil, "Too few arguments.")
	case 1:
		target = args[0]
	case 2:
		target = args[0]
		method = args[1]
	default:
		fail(nil, "Too many arguments.")
	}

	config := config.New(target, method, *insecure)
	app := app.New(config)
	if err := app.Run(context.Background()); err != nil {
		fail(err, "Failed")
	}
}
