package models

type interaction struct {
	userID int64
	at     int64
}

type Interactions struct {
	arr []interaction
}

func NewInteractions() Interactions {
	return Interactions{}
}

func (i Interactions) Add(userId int64, at int64) {
	i.arr = append(i.arr, interaction{userID: userId, at: at})
}

func (i Interactions) Search(userId int64) bool {
	//TODO: add Logic
	return false
}
