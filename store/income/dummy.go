package income

import "errors"

var kaboom = errors.New("Not Implemented")

type Dummy struct{}

func (*Dummy) Create(record *Record) error                 { return kaboom }
func (*Dummy) Get(id int64) (Record, error)                { return Record{}, kaboom }
func (*Dummy) ListByUser(id int64) ([]Record, error)       { return nil, kaboom }
func (*Dummy) UserIncome(id int64) (int, error)            { return 0, kaboom }
func (*Dummy) UserIncomeRank(id int64) (int, error)        { return 0, kaboom }
func (*Dummy) UpdateScores() error                         { return kaboom }
func (*Dummy) CountBonuses(i int64, p string) (int, error) { return 0, kaboom }
