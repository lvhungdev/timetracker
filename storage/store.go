package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/lvhungdev/tt/tracker"
)

// TODO implement a logger to log errors to a file
type Store struct {
	path string
}

func NewStore(path string) (Store, error) {
	err := os.MkdirAll(path, 0666)
	if err != nil {
		return Store{}, err
	}

	return Store{
		path: path,
	}, nil
}

func (s *Store) GetToday() []tracker.Record {
	name := time.Now().Local().Format("2006-01-02")
	data, err := s.readFile(name)
	if err != nil {
		return []tracker.Record{}
	}

	records := []tracker.Record{}
	err = json.Unmarshal(data, &records)
	if err != nil {
		return []tracker.Record{}
	}

	return records
}

func (s *Store) GetFromDate(date time.Time) []tracker.Record {
	name := date.Local().Format("2006-01-02")

	data, err := s.readFile(name)
	if err != nil {
		return []tracker.Record{}
	}

	records := []tracker.Record{}
	err = json.Unmarshal(data, &records)
	if err != nil {
		return []tracker.Record{}
	}

	return records
}

// TODO improve performance for this
func (s *Store) SaveAll(records []tracker.Record) error {
	for _, r := range records {
		err := s.save(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) save(record tracker.Record) error {
	records := s.GetFromDate(record.Start)

	var exist *tracker.Record
	for i := 0; i < len(records); i++ {
		if records[i].Id == record.Id {
			exist = &records[i]
		}
	}

	if exist == nil {
		records = append(records, record)
	} else {
		*exist = record
	}

	data, err := json.Marshal(&records)
	if err != nil {
		return err
	}

	return s.writeFile(record.Start.Format("2006-01-02"), data)
}

func (s *Store) readFile(name string) ([]byte, error) {
	file, err := os.Open(filepath.Join(s.path, name))
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	data := []byte{}
	for {
		buffer := make([]byte, 512)
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				return data, nil
			} else {
				return nil, err
			}
		}

		data = append(data, buffer[:n]...)
	}
}

func (s *Store) writeFile(name string, data []byte) error {
	file, err := os.Create(filepath.Join(s.path, name))
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	return writer.Flush()
}
