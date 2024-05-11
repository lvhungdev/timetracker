package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

	if err := handle(t, os.Args[1:]); err != nil {
		fmt.Printf("[ERROR] %v", err)
	}
}

func handle(t *tracker.Tracker, args []string) error {
	cmd, err := getCommand(args)
	if err != nil {
		return err
	}

	switch cmd := cmd.(type) {
	case cmdGetCurrent:
		r := t.GetCurrent()
		if r == nil {
			fmt.Println("no active time tracking")
			return nil
		}
		renderer.RenderRecord(os.Stdout, *r)

	case cmdStartTracking:
		old, curr, err := t.StartTracking(cmd.name)
		if err != nil {
			return err
		}
		if err := t.Save(); err != nil {
			return err
		}
		if old != nil {
			renderer.RenderRecord(os.Stdout, *old)
		}
		renderer.RenderRecord(os.Stdout, *curr)

	case cmdStopTracking:
		r, err := t.StopTracking()
		if err != nil {
			return err
		}
		if err := t.Save(); err != nil {
			return err
		}
		renderer.RenderRecord(os.Stdout, *r)

	case cmdReport:
		now := time.Now()
		from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
		records := t.GetAll(from, to)

		renderer.RenderRecords(os.Stdout, records)
	}

	return nil
}
