package manager

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	coinApi "github.com/miguelmota/go-coinmarketcap"

	"github.com/cbarraford/cryptocades-backend/store"
	"github.com/cbarraford/cryptocades-backend/store/entry"
	"github.com/cbarraford/cryptocades-backend/store/jackpot"
	"github.com/cbarraford/cryptocades-backend/store/user"
	"github.com/cbarraford/cryptocades-backend/util"
)

const MAX_JACKPOTS = 1

var target_price int64

func init() {
	var err error
	target_price, err = strconv.ParseInt(os.Getenv("TARGET_PRICE"), 10, 32)
	if err != nil {
		// TODO: we should alert on this error
		log.Fatalf("Failed to read TARGET_PRICE")
	}
}

func Start(store store.Store) {
	// spawn jackpot(s)
	tickJack := time.NewTicker(5 * time.Second)
	tickScores := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-tickJack.C:
				if err := ManageJackpots(store.Jackpots, store.Entries, store.Users); err != nil {
					// TODO: we should alert on this error
					log.Printf("Manage Jackpot Error: %+v", err)
				}
			case <-tickScores.C:
				if err := store.Incomes.UpdateScores(); err != nil {
					// TODO: we should alert on this error
					log.Printf("Update Scores Error: %+v", err)
				}
			case <-quit:
				tickJack.Stop()
				tickScores.Stop()
				return
			}
		}
	}()
}

func ManageJackpots(store jackpot.Store, entryStore entry.Store, userStore user.Store) error {
	jackpots, err := store.GetActiveJackpots()
	if err != nil {
		return fmt.Errorf("Error getting active jackpots: %+v", err)
	}
	if len(jackpots) == 0 {
		coinInfo, err := coinApi.GetCoinData("bitcoin")
		jackpot := jackpot.Record{
			Jackpot: util.ToFixed(float64(target_price)/coinInfo.PriceUsd, 5),
			// One week end time
			EndTime: time.Now().UTC().Add(168 * time.Hour),
		}
		err = store.Create(&jackpot)
		if err != nil {
			return fmt.Errorf("Unable to create new jackpot: %+v", err)
		}
	}

	jackpots, err = store.GetIncompleteJackpots()
	if err != nil {
		return fmt.Errorf("Error getting incomplete jackpots: %+v", err)
	}
	for _, jackpot := range jackpots {
		jackpot.WinnerId, err = PickWinner(entryStore, jackpot.Id)
		if err != nil {
			return fmt.Errorf("Error picking jackpot winner: %+v", err)
		}

		if jackpot.WinnerId > 0 {
			user, err := userStore.Get(jackpot.WinnerId)
			if err != nil {
				return fmt.Errorf("Error getting jackpot winner: %+v", err)
			}
			jackpot.WinnerBTCAddr = user.BTCAddr
			err = store.Update(&jackpot)
			if err != nil {
				return fmt.Errorf("Error updating jackpot winner: %+v", err)
			}
		}
	}
	return nil
}

func PickWinner(store entry.Store, jackpotId int64) (int64, error) {
	records, err := store.ListByJackpot(jackpotId)
	if err != nil {
		return 0, err
	}

	// if no entries, no winners
	if len(records) == 0 {
		return 0, nil
	}

	// shuffle our records returned
	for i := len(records) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		records[i], records[j] = records[j], records[i]
	}

	// count total entries in this jackpot
	var totalEntries int
	for _, record := range records {
		totalEntries = totalEntries + record.Amount
	}

	// rand.Intn includes 0, so increment by one
	winnerNum := rand.Intn(totalEntries) + 1

	winner := findWinner(records, winnerNum)

	if winner.UserId > 0 {
		return winner.UserId, nil
	}

	return 0, fmt.Errorf("Unable to find a jackpot winner: %d", jackpotId)
}

func findWinner(records []entry.Record, winnerNum int) entry.Record {
	if winnerNum == 0 {
		return entry.Record{}
	}

	var counter int
	for _, record := range records {
		counter = counter + record.Amount
		if counter >= winnerNum {
			return record
		}
	}

	return entry.Record{}
}
