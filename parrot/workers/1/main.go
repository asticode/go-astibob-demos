package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/asticode/go-astibob"
	"github.com/asticode/go-astibob/abilities/text_to_speech"
	"github.com/asticode/go-astibob/abilities/text_to_speech/speak"
	"github.com/asticode/go-astibob/worker"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilog"
)

func main() {
	// Parse flags
	flag.Parse()

	// Create logger
	l := astilog.NewFromFlags()
	defer l.Close()

	// Create worker
	w := worker.New("Worker #1", worker.Options{
		Index: astibob.ServerOptions{
			Addr:     "127.0.0.1:4000",
			Password: "admin",
			Username: "admin",
		},
		Server: astibob.ServerOptions{Addr: "127.0.0.1:4001"},
	}, l)
	defer w.Close()

	// Say "Hello world" when the runnable has started
	o := &sync.Once{}
	w.On(astibob.DispatchConditions{
		From: astibob.NewRunnableIdentifier("Text to Speech", "Worker #1"),
		Name: astikit.StrPtr(astibob.RunnableStartedMessage),
	}, func(m *astibob.Message) (err error) {
		o.Do(func() {
			// Send message
			if err = w.SendMessage(worker.MessageOptions{
				Message:  text_to_speech.NewSayMessage("Hello world"),
				Runnable: "Text to Speech",
				Worker:   "Worker #1",
			}); err != nil {
				err = fmt.Errorf("main: sending message failed: %w", err)
				return
			}
		})
		return
	})

	// Create speaker
	s := speak.New(speak.Options{}, l)

	// Initialize speaker
	if err := s.Initialize(); err != nil {
		l.Fatal(fmt.Errorf("main: initializing speaker failed: %w", err))
	}
	defer s.Close()

	// Register runnables
	w.RegisterRunnables(worker.Runnable{
		Runnable: text_to_speech.NewRunnable("Text to Speech", s, l),
	})

	// Handle signals
	w.HandleSignals()

	// Serve
	w.Serve()

	// Register to index
	w.RegisterToIndex()

	// Blocking pattern
	w.Wait()
}
