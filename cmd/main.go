package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/lvhungdev/tt/renderer"
	"github.com/lvhungdev/tt/storage"
	"github.com/lvhungdev/tt/tracker"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("unable to get home directory, use current directory instead: %v", err)
		homeDir = "."
	}

	s, err := storage.NewStore(path.Join(homeDir, ".timetracker"))
	if err != nil {
		log.Fatalf("unable to initialize store %v", err)
	}

	t := tracker.New(&s)

	if err := handle(t, os.Args[1:]); err != nil {
		fmt.Println(err)
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
		records := t.GetAll(cmd.from, cmd.to)

		renderer.RenderRecords(os.Stdout, records)
	}

	return nil
}
