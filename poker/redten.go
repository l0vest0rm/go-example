package main

type RedTen struct {
	playerNum int
}

func NewRedTen(playerNum int) IPlay {
	t := &RedTen{playerNum: playerNum}
	return t
}

// ModVals 修改值
func (t *RedTen) ModVals(vals []int) []int {
	l := len(vals)
	for i := 0; i < l; i++ {
		if vals[i] == 2 || vals[i] == 15 || vals[i] == 28 || vals[i] == 41 {
			if t.playerNum == 5 {
				vals[i] = vals[l-1]
				l--
			}
		}

	}

	return vals[:l]
}

func (t *RedTen) Dispacther() {

}
