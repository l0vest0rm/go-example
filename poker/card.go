package main

import (
	"fmt"
	"math/rand"
	"time"
)

// ICards 牌接口
type ICards interface {
	Shuffle()
	PrintCards()
}

// Cards 一局牌结构
type cards struct {
	vals []int
}

// CreateNewCards 构造一局新牌
func NewCards() ICards {
	//构造一个大小为54的数组,按照红、黑、梅、方、小王、大王
	vals := make([]int, 54)
	for i := 0; i < 54; i++ {
		vals[i] = i
	}

	cards := &cards{vals: vals}

	return cards
}

// Shuffle 洗牌
func (t *cards) Shuffle1() {
	vals := t.vals
	r := rand.New(rand.NewSource(time.Now().Unix())) //根据系统时间戳初始化Random
	for len(vals) > 0 {                              //根据牌面数组长度遍历
		n := len(vals)                                          //数组长度
		randIndex := r.Intn(n)                                  //得到随机index
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1] //最后一张牌和第randIndex张牌互换
		vals = vals[:n-1]
	}
}

// Shuffle 洗牌
func (t *cards) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().Unix())) //根据系统时间戳初始化Random
	for n := len(t.vals); n > 0; n-- {
		randIndex := r.Intn(n)                                          //得到随机index
		t.vals[n-1], t.vals[randIndex] = t.vals[randIndex], t.vals[n-1] //最后一张牌和第randIndex张牌互换
	}
}

func (t *cards) PrintCards() {
	for i := 0; i < 54; i++ {
		fmt.Printf("%s,", ConvertVal2Str(t.vals[i]))
	}
}

// ConvertVal2Str 将牌数字转换成字符串
func ConvertVal2Str(val int) string {
	if val < 13 {
		return fmt.Sprintf("A%d", val+1)
	} else if val < 26 {
		return fmt.Sprintf("B%d", val-12)
	} else if val < 39 {
		return fmt.Sprintf("C%d", val-25)
	} else if val < 52 {
		return fmt.Sprintf("D%d", val-38)
	} else if val == 52 {
		return "Q0"
	} else {
		return "K0"
	}
}
