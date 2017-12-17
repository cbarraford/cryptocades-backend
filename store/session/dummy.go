package session

import "errors"

var kaboom = errors.New("Not Implemented")

type Dummy struct{}

func (*Dummy) Create(record *Record, l int) error   { return kaboom }
func (*Dummy) GetByToken(t string) (Record, error)  { return Record{}, kaboom }
func (*Dummy) Authenticate(t string) (int64, error) { return 0, kaboom }
func (*Dummy) Delete(t string) error                { return kaboom }
