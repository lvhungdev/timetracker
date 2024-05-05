package main

import (
	"log"
	"os"

	"github.com/lvhungdev/tt/renderer"
	"github.com/lvhungdev/tt/storage"
	"github.com/lvhungdev/tt/tracker"
)

func main() {
	s, err := storage.NewStore(".")
	if err != nil {
		log.Fatalf("unable to initialize store %v", err)
	}

	t := tracker.New(&s)

	curr, new, err := t.StartTracking("task 1")
	if err != nil {
		log.Fatalln(err)
	}

	err = t.Save()
	if err != nil {
		log.Fatalln(err)
	}

	if curr != nil {
		renderer.RenderRecord(os.Stdout, *curr)
	}
	if new != nil {
		renderer.RenderRecord(os.Stdout, *new)
	}
}

func handle(args []string) error {
	if len(args) == 0 {

	} else if args[0] == "start" {

	} else if args[0] == "stop" {

	}

	return nil
}
