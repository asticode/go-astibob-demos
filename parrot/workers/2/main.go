package main

import (
	"flag"
	"fmt"

	"github.com/asticode/go-astibob"
	"github.com/asticode/go-astibob/abilities/audio_input"
	"github.com/asticode/go-astibob/abilities/audio_input/portaudio"
	"github.com/asticode/go-astibob/worker"
	"github.com/asticode/go-astilog"
)

func main() {
	// Parse flags
	flag.Parse()

	// Create logger
	l := astilog.NewFromFlags()
	defer l.Close()

	// Create worker
	w := worker.New("Worker #2", worker.Options{
		Index: astibob.ServerOptions{
			Addr:     "127.0.0.1:4000",
			Password: "admin",
			Username: "admin",
		},
		Server: astibob.ServerOptions{Addr: "127.0.0.1:4002"},
	}, l)
	defer w.Close()

	// Create portaudio
	p := portaudio.New(l)

	// Initialize portaudio
	if err := p.Initialize(); err != nil {
		l.Fatal(fmt.Errorf("main: initializing portaudio failed: %w", err))
	}
	defer p.Close()

	// Create default stream
	s, err := p.NewDefaultStream(portaudio.StreamOptions{
		BitDepth:         32,
		BufferLength:     5000,
		MaxSilenceLevel:  5 * 1e6,
		NumInputChannels: 2,
		SampleRate:       44100,
	})
	if err != nil {
		l.Fatal(fmt.Errorf("main: creating default stream failed: %w", err))
	}

	// Create runnable
	r := audio_input.NewRunnable("Audio input", s, l)

	// Register runnables
	w.RegisterRunnables(worker.Runnable{
		AutoStart: true,
		Runnable:  r,
	})

	// Register listenables
	// This is mandatory for the Web UI to work properly
	w.RegisterListenables(worker.Listenable{
		Listenable: r,
		Runnable:   "Audio input",
		Worker:     "Worker #2",
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
