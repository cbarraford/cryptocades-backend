package jackpot

import "errors"

var kaboom = errors.New("Not Implemented")

type Dummy struct{}

func (*Dummy) Create(record *Record) error          { return kaboom }
func (*Dummy) Get(id int64) (Record, error)         { return Record{}, kaboom }
func (*Dummy) Update(r *Record) error               { return kaboom }
func (*Dummy) List() ([]Record, error)              { return nil, kaboom }
func (*Dummy) GetActiveJackpots() ([]Record, error) { return nil, kaboom }
