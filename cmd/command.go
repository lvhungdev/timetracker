package main

import (
	"fmt"
	"time"
)

type command interface {
	time() time.Time
}

type cmdBase struct {
	t time.Time
}

func (c cmdBase) time() time.Time {
	return c.t
}

type cmdGetCurrent struct {
	cmdBase
}

func newCmdGetCurrent() command {
	return cmdGetCurrent{cmdBase{time.Now()}}
}

type cmdStartTracking struct {
	cmdBase
	name string
}

func newCmdStartTracking(name string) command {
	return cmdStartTracking{cmdBase{time.Now()}, name}
}

type cmdStopTracking struct {
	cmdBase
}

func newCmdStopTracking() command {
	return cmdStopTracking{cmdBase{time.Now()}}
}

type cmdReport struct {
	cmdBase
}

func newCmdReport() command {
	return cmdReport{cmdBase{time.Now()}}
}

func getCommand(args []string) (command, error) {
	if len(args) == 0 {
		return newCmdGetCurrent(), nil
	}

	switch args[0] {
	case "start":
		if len(args) == 1 {
			return nil, fmt.Errorf("name is required")
		}
		return newCmdStartTracking(args[1]), nil

	case "stop":
		return newCmdStopTracking(), nil

	case "report":
		return newCmdReport(), nil

	default:
		return nil, fmt.Errorf("unknown command: %s", args[0])
	}
}
