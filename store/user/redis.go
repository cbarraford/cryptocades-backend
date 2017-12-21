package user

import (
	"strconv"

	"github.com/garyburd/redigo/redis"
)

type score struct {
	addr  string
	score int
}

// zpop pops a value from the ZSET key using WATCH/MULTI/EXEC commands.
func (db *store) zpop(key string) (scores []score, err error) {

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
			v, err := strconv.Atoi(members[i+1])
			if err != nil {
				return nil, err
			}
			scores = append(scores, score{members[i], v})

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
	scores, err := db.zpop("hashes")
	if err != nil {
		return err
	}
	return db.AppendScore(scores)
}
