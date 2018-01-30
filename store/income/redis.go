package income

import (
	"log"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
)

type score struct {
	id        int64
	score     int
	gameId    int64
	sessionId string
}

// zpop pops a value from the ZSET key using WATCH/MULTI/EXEC commands.
func (db *store) zpop(key string) (scores []Record, err error) {

	defer func() {
		// Return connection to normal state on error.
		if err != nil {
			db.redis.Do("DISCARD")
		}
	}()

	// Loop until transaction is successful.
	for {
		if _, err := db.redis.Do("WATCH", key); err != nil {
			return nil, err
		}

		members, err := redis.Strings(db.redis.Do("ZRANGE", key, 0, -1, "WITHSCORES"))
		if err != nil {
			return nil, err
		}
		db.redis.Send("MULTI")
		for i, _ := range members {
			if (i % 2) == 1 {
				continue
			}

			parts := strings.Split(members[i], "@")
			if len(parts) != 2 {
				log.Printf("Malformed redis key: %s", members[i])
				continue
			}
			id, _ := strconv.ParseInt(parts[1], 10, 64)

			parts = strings.Split(parts[0], "-")
			if len(parts) != 2 {
				log.Printf("Malformed redis key: %s", members[i])
				continue
			}
			gameId, _ := strconv.ParseInt(parts[0], 10, 64)
			sessionId := parts[1]

			v, err := strconv.Atoi(members[i+1])
			if err != nil {
				return nil, err
			}
			scores = append(scores, Record{
				UserId:        id,
				PartialAmount: v,
				GameId:        gameId,
				SessionId:     sessionId,
			})

			db.redis.Send("ZREM", key, members[i])
		}
		queued, err := db.redis.Do("EXEC")
		if err != nil {
			return nil, err
		}

		if queued != nil {
			break
		}
	}

	return
}

func (db *store) UpdateScores() error {
	records, err := db.zpop("shares")
	if err != nil {
		return err
	}

	tx, err := db.sqlx.Beginx()
	if err != nil {
		return err
	}

	for _, record := range records {
		err := db.CreateWithinTransaction(&record, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
