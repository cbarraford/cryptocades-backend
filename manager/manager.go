package manager

import (
	"fmt"
	"log"
	"time"

	"github.com/CBarraford/lotto/store"
	"github.com/CBarraford/lotto/store/jackpot"
)

const MAX_JACKPOTS = 1

func Start(store store.Store) {
	// spawn jackpot(s)
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := ManageJackpots(store.Jackpots); err != nil {
					// TODO: we should alert on this error
					log.Printf("%+v", err)
				}
			case <-quit:
				ticker.Stop()
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
			EndTime: time.Now().UTC().AddDate(0, 0, 7),
		}
		err = store.Create(&jackpot)
		if err != nil {
			return fmt.Errorf("Unable to create new jackpot: %+v", err)
		}
	} else {
	}
	return nil
}
