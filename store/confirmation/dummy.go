package confirmation

import "errors"

var kaboom = errors.New("Not Implemented")

type Dummy struct{}

func (*Dummy) Create(record *Record) error        { return kaboom }
func (*Dummy) Get(i int64) (Record, error)        { return Record{}, kaboom }
func (*Dummy) GetByCode(c string) (Record, error) { return Record{}, kaboom }
func (*Dummy) Delete(id int64) error              { return kaboom }
