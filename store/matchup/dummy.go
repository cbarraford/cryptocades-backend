package matchup

import "errors"

var kaboom = errors.New("Not Implemented")

type Dummy struct{}

func (*Dummy) KeyName(s string, o int) string { return "" }
func (*Dummy) Get(s string, o int, i int64) (Record, error) {
	return Record{}, kaboom
}
func (*Dummy) GetTopPerformers(s string, o int, t int) ([]Record, error) {
	return nil, kaboom
}
func (*Dummy) ExpandRecords(r []Record) ([]Record, error) { return nil, kaboom }
