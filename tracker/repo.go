package tracker

type Repo interface {
	GetToday() []Record
	SaveAll(records []Record) error
}
