package game

type Record struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
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
		{Id: 1, Name: "Goblin Stacks"},
		{Id: 2, Name: "Asteroid Tycoon"},
	}
}
