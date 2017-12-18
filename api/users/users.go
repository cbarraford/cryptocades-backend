package users

import (
	"time"
)

type input struct {
	BTCAddr     string    `json:"btc_address"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	MinedHashes int       `json:"mined_hashes"`
	BonusHashes int       `json:"bonus_hashes"`
	CreatedTime time.Time `json:"created_time"`
	UpdatedTime time.Time `json:"updated_time"`
}
