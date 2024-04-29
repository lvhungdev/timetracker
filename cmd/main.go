package main

import (
	"log"

	"github.com/lvhungdev/tt/storage"
	"github.com/lvhungdev/tt/tracker"
)

func main() {
	s, err := storage.NewStore(".")
	if err != nil {
		log.Fatalf("unable to initialize store %v", err)
	}

	t := tracker.New(&s)

	t.StartTracking("task 1")

	err = t.Save()
	if err != nil {
		log.Fatalln(err)
	}
}
