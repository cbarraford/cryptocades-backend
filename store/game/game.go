package game

type Record struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
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
			Name:        "Goblin Stacks",
			Description: "Earn lotto tickets as your goblin builds your tower. Hit tower milestones to receive bonuses!",
		},
		{
			Id:          2,
			Name:        "Asteroid Tycoon",
			Description: "Mine asteroids to earn iron ore which can be used to buy lotto tickets or upgrade your mining rig.",
		},
	}
}
