package boost

import "errors"

var kaboom = errors.New("Not Implemented")

type Dummy struct{}

func (*Dummy) Create(record *Record) error           { return kaboom }
func (*Dummy) Get(id int64) (Record, error)          { return Record{}, kaboom }
func (*Dummy) ListByUser(id int64) ([]Record, error) { return nil, kaboom }
func (*Dummy) Assign(i, d int64) error               { return kaboom }
