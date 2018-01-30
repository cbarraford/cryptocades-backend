package manager

import (
	"fmt"
	"log"
	"time"

	"github.com/cbarraford/cryptocades-backend/store"
	"github.com/cbarraford/cryptocades-backend/store/jackpot"
)

const MAX_JACKPOTS = 1

func Start(store store.Store) {
	// spawn jackpot(s)
	tickJack := time.NewTicker(5 * time.Second)
	tickScores := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-tickJack.C:
				if err := ManageJackpots(store.Jackpots); err != nil {
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

func ManageJackpots(store jackpot.Store) error {
	jackpots, err := store.GetActiveJackpots()
	if err != nil {
		return fmt.Errorf("Error getting active jackpots: %+v", err)
	}
	if len(jackpots) == 0 {
		jackpot := jackpot.Record{
			Jackpot: 100,
			// One week end time
			EndTime: time.Now().UTC().Add(168 * time.Hour),
		}
		err = store.Create(&jackpot)
		if err != nil {
			return fmt.Errorf("Unable to create new jackpot: %+v", err)
		}
	} else {
	}
	return nil
}
