package user

import "errors"

var kaboom = errors.New("Not Implemented")

type Dummy struct{}

func (*Dummy) Create(record *Record) error              { return kaboom }
func (*Dummy) Update(record *Record) error              { return kaboom }
func (*Dummy) MarkAsConfirmed(record *Record) error     { return kaboom }
func (*Dummy) PasswordSet(record *Record) error         { return kaboom }
func (*Dummy) Get(i int64) (Record, error)              { return Record{}, kaboom }
func (*Dummy) GetByBTCAddress(u string) (Record, error) { return Record{}, kaboom }
func (*Dummy) GetByUsername(u string) (Record, error)   { return Record{}, kaboom }
func (*Dummy) GetByEmail(e string) (Record, error)      { return Record{}, kaboom }
func (*Dummy) List() ([]Record, error)                  { return nil, kaboom }
func (*Dummy) Authenticate(u, p string) (Record, error) { return Record{}, kaboom }
func (*Dummy) AppendScore(scores []score) error         { return kaboom }
func (*Dummy) Delete(id int64) error                    { return kaboom }
