package entry

import "errors"

var kaboom = errors.New("Not Implemented")

type Dummy struct{}

func (*Dummy) Create(record *Record) error           { return kaboom }
func (*Dummy) Get(id int64) (Record, error)          { return Record{}, kaboom }
func (*Dummy) GetOdds(j, i int64) (Odds, error)      { return Odds{}, kaboom }
func (*Dummy) List() ([]Record, error)               { return nil, kaboom }
func (*Dummy) ListByUser(id int64) ([]Record, error) { return nil, kaboom }
func (*Dummy) UserSpent(id int64) (int, error)       { return 0, kaboom }
