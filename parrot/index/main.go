package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/asticode/go-astibob"
	"github.com/asticode/go-astibob/index"
)

func main() {
	// Parse flags
	flag.Parse()

	// Create logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

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
