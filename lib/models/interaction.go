package models

type interaction struct {
	IIN string
	at  int64
}

type Interactions struct {
	arr []interaction
}

func NewInteractions() *Interactions {
	return &Interactions{}
}

func (i *Interactions) Add(IIN string, at int64) {
	i.arr = append(i.arr, interaction{IIN: IIN, at: at})
}

func (i *Interactions) Search(IIN string) bool {
	for _, val := range i.arr {
		if IIN == val.IIN {
			return true
		}
	}
	return false
}
