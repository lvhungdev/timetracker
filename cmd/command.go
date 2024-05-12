package main

import (
	"fmt"
	"strconv"
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
	from time.Time
	to   time.Time
}

func newCmdReport(from, to time.Time) command {
	return cmdReport{
		cmdBase: cmdBase{time.Now()},
		from:    from,
		to:      to,
	}
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
		return parseCmdReport(args[1:])

	default:
		return nil, fmt.Errorf("unknown command: %s", args[0])
	}
}

func parseCmdReport(args []string) (command, error) {
	now := time.Now()

	if len(args) == 0 {
		from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)

		return newCmdReport(from, to), nil
	}

	if args[0] == "yesterday" {
		yesterday := now.AddDate(0, 0, -1)
		from := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.Local)
		to := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, time.Local)
		return newCmdReport(from, to), nil
	}

	from, to, err := parseDuration(args)
	if err == nil {
		return newCmdReport(from, to), nil
	}

	if len(args) != 2 {
		return nil, err
	}

	from, err1 := parseTime(args[0])
	to, err2 := parseTime(args[1])
	if err1 != nil {
		return nil, fmt.Errorf("invalid option: %s", args[0])
	}
	if err2 != nil {
		return nil, fmt.Errorf("invalid option: %s", args[1])
	}

	return newCmdReport(from, to), nil
}

func parseDuration(args []string) (time.Time, time.Time, error) {
	numb := 1
	if len(args) == 2 {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid option: %s", args[0])
		}
		numb = n
	}

	now := time.Now()
	switch args[len(args)-1] {
	case "day":
		from := now.AddDate(0, 0, -numb+1)
		from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.Local)
		to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
		return from, to, nil

	case "week":
		monday := now.AddDate(0, 0, -int(now.Weekday())+1-7*(numb-1))
		sunday := now.AddDate(0, 0, -int(now.Weekday())+7)
		from := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.Local)
		to := time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 23, 59, 59, 0, time.Local)
		return from, to, nil

	default:
		return time.Time{}, time.Time{}, fmt.Errorf("unknown option: %s", args[len(args)-1])
	}
}

func parseTime(s string) (time.Time, error) {
	if len(s) == len("15:04:05") {
		s = time.Now().Format("2006-01-02 ") + s
	}

	return time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
}
