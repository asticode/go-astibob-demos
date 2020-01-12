package main

import (
	"flag"
	"fmt"

	"github.com/asticode/go-astibob"
	"github.com/asticode/go-astibob/index"
	"github.com/asticode/go-astilog"
)

func main() {
	// Parse flags
	flag.Parse()

	// Create logger
	l := astilog.NewFromFlags()
	defer l.Close()

	// Create index
	i, err := index.New(index.Options{
		Server: astibob.ServerOptions{
			Addr:     "127.0.0.1:4000",
			Password: "admin",
			Username: "admin",
		},
	}, l)
	if err != nil {
		l.Fatal(fmt.Errorf("main: creating index failed: %w", err))
	}
	defer i.Close()

	// Handle signals
	i.HandleSignals()

	// Serve
	i.Serve()

	// Blocking pattern
	i.Wait()
}
