package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/lvhungdev/tt/command"
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
	cmd, err := command.Get(args)
	if err != nil {
		return err
	}

	switch cmd := cmd.(type) {
	case command.GetCurrent:
		r := t.GetCurrent()
		if r == nil {
			fmt.Println("no active time tracking")
			return nil
		}
		renderer.RenderRecord(os.Stdout, *r)

	case command.StartTracking:
		old, curr, err := t.StartTracking(cmd.Name, cmd.Time())
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

	case command.StopTracking:
		r, err := t.StopTracking(cmd.Time())
		if err != nil {
			return err
		}
		if err := t.Save(); err != nil {
			return err
		}
		renderer.RenderRecord(os.Stdout, *r)

	case command.Report:
		records := t.GetAll(cmd.From, cmd.To)

		renderer.RenderRecords(os.Stdout, records)
	}

	return nil
}
