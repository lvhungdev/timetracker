package tracker

import "time"

type Repo interface {
	GetToday() []Record
	GetFromDate(date time.Time) []Record
	SaveAll(records []Record) error
}
