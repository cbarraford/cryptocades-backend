package game

type Record struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Store interface {
	List() []Record
}

type store struct {
	Store
}

func NewStore() Store {
	return &store{}
}

func (db *store) List() []Record {
	return []Record{
		{
			Id:          1,
			Name:        "Tallest Tower",
			Type:        "Passive",
			Description: "Build your tower taller and taller. Hit tower milestones to get jackpot plays",
		},
		{
			Id:          2,
			Name:        "Asteroid Tycoon",
			Type:        "Active",
			Description: "Mine asteroids to gather iron ore and upgrade your ship or trade for jackpot plays",
		},
	}
}
