package matchup

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/garyburd/redigo/redis"
	"github.com/jmoiron/sqlx"
)

type Store interface {
	KeyName(matchup string, offset int) string
	Get(matchup string, offset int, userId int64) (record Record, err error)
	GetTopPerformers(matchup string, offset int, top int) (records []Record, err error)
	ExpandRecords(records []Record) ([]Record, error)
}

type store struct {
	Store
	sqlx  *sqlx.DB
	redis redis.Conn
}

func NewStore(db *sqlx.DB, redis redis.Conn) Store {
	return &store{sqlx: db, redis: redis}
}

type Record struct {
	UserId   int64  `json:"-"`
	Username string `json:"username"`
	Rank     int    `json:"rank"`
	Score    int    `json:"score"`
}

// Keyname generates the redis key for a specific matchup
// "matchup" should be a string (ie 'daily') while offset which one to retrieve
// (ie current aka, previous aka -1, etc)
func (db *store) KeyName(matchup string, offset int) string {
	if offset > 0 {
		offset = -offset
	}
	now := time.Now().UTC()
	if matchup == "daily" {
		now = now.Add(time.Duration(offset) * 24 * time.Hour)
		return fmt.Sprintf(
			"%d-%02d-%02d", now.Year(), now.Month(), now.Day(),
		)
	} else if matchup == "hourly" {
		now = now.Add(time.Duration(offset) * time.Hour)
		return fmt.Sprintf(
			"%d-%02d-%02d %d", now.Year(), now.Month(), now.Day(), now.Hour(),
		)
	}
	return ""
}

func (db *store) Get(matchup string, offset int, userId int64) (record Record, err error) {
	keyname := db.KeyName(matchup, offset)
	if keyname == "" {
		return Record{}, fmt.Errorf("Invalid matchup interval: %s", matchup)
	}

	record.UserId = userId
	record.Rank, err = redis.Int(db.redis.Do("ZRANK", keyname, userId))
	if err != nil && err != redis.ErrNil {
		return record, err
	}
	// rank starts with zero being top rank. Increment by one to offset
	record.Rank = record.Rank + 1

	record.Score, err = redis.Int(db.redis.Do("ZSCORE", keyname, userId))
	if err != nil && err != redis.ErrNil {
		return record, err
	}

	records, err := db.ExpandRecords([]Record{record})
	return records[0], err
}

func (db *store) GetTopPerformers(matchup string, offset int, top int) ([]Record, error) {
	var err error
	records := []Record{}
	keyname := db.KeyName(matchup, offset)
	if keyname == "" {
		return nil, fmt.Errorf("Invalid matchup interval: %s", matchup)
	}
	scores, err := redis.Strings(db.redis.Do("ZREVRANGE", keyname, 0, top, "WITHSCORES"))
	if err != nil {
		return nil, err
	}

	// the data we get from redis is an array alternating between userId and
	// score. Because of that, we do some work to convery an array of ints, to
	// an array of Records
	for i, v := range scores {
		if i%2 == 0 {
			value, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				log.Printf("Error converting value to int: %s", v)
			}
			records = append(records, Record{UserId: value, Rank: i/2 + 1})
		} else {
			value, err := strconv.Atoi(v)
			if err != nil {
				log.Printf("Error converting value to int: %s", v)
			}
			records[len(records)-1].Score = value
		}
	}

	return db.ExpandRecords(records)
}

func (db *store) ExpandRecords(records []Record) ([]Record, error) {
	var err error

	ids := make([]int64, len(records))
	questionMarks := make([]string, len(records))
	for i, _ := range records {
		ids[i] = records[i].UserId
		questionMarks[i] = "?"
	}

	var userRecords []user.Record
	query, args, err := sqlx.In("SELECT * FROM users WHERE id IN (?);", ids)
	if err != nil {
		return nil, err
	}
	query = db.sqlx.Rebind(query)

	err = db.sqlx.Select(&userRecords, query, args...)
	if err != nil {
		return nil, err
	}

	for i, r := range records {
		records[i].Username = "Unknown"
		for _, user := range userRecords {
			if r.UserId == user.Id {
				records[i].Username = user.Username
			}
		}
	}

	return records, nil
}
