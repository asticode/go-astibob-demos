package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/asticode/go-astibob"
	"github.com/asticode/go-astibob/abilities/audio_input"
	"github.com/asticode/go-astibob/abilities/speech_to_text"
	"github.com/asticode/go-astibob/abilities/speech_to_text/deepspeech"
	"github.com/asticode/go-astibob/abilities/text_to_speech"
	"github.com/asticode/go-astibob/worker"
)

const wd = "tmp"

func main() {
	// Parse flags
	flag.Parse()

	// Create logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Create worker
	w := worker.New("Worker #3", worker.Options{
		Index: astibob.ServerOptions{
			Addr:     "127.0.0.1:4000",
			Password: "admin",
			Username: "admin",
		},
		Server: astibob.ServerOptions{Addr: "127.0.0.1:4003"},
	}, l)
	defer w.Close()

	// Create deepspeech
	mp := wd + "/model/en"
	d := deepspeech.New(deepspeech.Options{
		AlphabetPath:   mp + "/alphabet.txt",
		BeamWidth:      1024,
		ClientPath:     wd + "/DeepSpeech/DeepSpeech.py",
		LMPath:         mp + "/lm.binary",
		LMWeight:       0.75,
		ModelPath:      mp + "/output_graph.pb",
		PrepareDirPath: wd + "/prepare",
		TrainingArgs: map[string]string{
			"checkpoint_dir":   wd + "/model/custom/checkpoints",
			"dev_batch_size":   "1",
			"export_dir":       wd + "/model/custom",
			"noearly_stop":     "",
			"test_batch_size":  "1",
			"train_batch_size": "1",

			// Mozilla values
			"learning_rate": "0.0001",
			"dropout_rate":  "0.15",
			"lm_alpha":      "0.75",
			"lm_beta":       "1.85",
		},
		TriePath:             mp + "/trie",
		ValidWordCountWeight: 1.85,
	}, l)
	defer d.Close()

	// Initialize deepspeech
	if err := d.Init(); err != nil {
		l.Fatal(fmt.Errorf("main: initializing deepspeech failed: %w", err))
	}

	// Create runnable
	r := speech_to_text.NewRunnable("Speech to Text", d, l, speech_to_text.RunnableOptions{
		SpeechesDirPath: wd + "/speeches",
	})

	// Initialize runnable
	if err := r.Init(); err != nil {
		l.Fatal(fmt.Errorf("main: initializing runnable failed: %w", err))
	}
	defer r.Close()

	// Register runnables
	w.RegisterRunnables(worker.Runnable{
		AutoStart: true,
		Runnable:  r,
	})

	// Register listenables
	w.RegisterListenables(
		// Forwards audio samples to the "Speech to Text" ability
		worker.Listenable{
			Listenable: audio_input.NewListenable(audio_input.ListenableOptions{
				OnSamples: func(from astibob.Identifier, samples []int, bitDepth, numChannels, sampleRate int, maxSilenceLevel float64) (err error) {
					// Send message
					if err = w.SendMessage(worker.MessageOptions{
						Message: speech_to_text.NewSamplesMessage(
							from,
							samples,
							bitDepth,
							numChannels,
							sampleRate,
							maxSilenceLevel,
						),
						Runnable: "Speech to Text",
						Worker:   "Worker #3",
					}); err != nil {
						err = fmt.Errorf("main: sending message failed: %w", err)
						return
					}
					return
				},
			}),
			Runnable: "Audio input",
			Worker:   "Worker #2",
		},
		worker.Listenable{
			// Say detected words out loud
			Listenable: speech_to_text.NewListenable(speech_to_text.ListenableOptions{
				OnText: func(from astibob.Identifier, text string) (err error) {
					// Send message
					if err = w.SendMessage(worker.MessageOptions{
						Message:  text_to_speech.NewSayMessage(text),
						Runnable: "Text to Speech",
						Worker:   "Worker #1",
					}); err != nil {
						err = fmt.Errorf("main: sending message failed: %w", err)
						return
					}
					return
				},
			}),
			Runnable: "Speech to Text",
			Worker:   "Worker #3",
		},
	)

	// Handle signals
	w.HandleSignals()

	// Serve
	w.Serve()

	// Register to index
	w.RegisterToIndex()

	// Blocking pattern
	w.Wait()
}
