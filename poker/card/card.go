package card

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	HONG_1  = 0
	HONG_2  = 1
	HONG_3  = 2
	HONG_4  = 3
	HONG_5  = 4
	HONG_6  = 5
	HONG_7  = 6
	HONG_8  = 7
	HONG_9  = 8
	HONG_10 = 9
	HONG_11 = 10
	HONG_12 = 11
	HONG_13 = 12
	FANG_1  = 13
	FANG_2  = 14
	FANG_3  = 15
	FANG_4  = 16
	FANG_5  = 17
	FANG_6  = 18
	FANG_7  = 19
	FANG_8  = 20
	FANG_9  = 21
	FANG_10 = 22
	FANG_11 = 23
	FANG_12 = 24
	FANG_13 = 25
	HEI_1   = 26
	HEI_2   = 27
	HEI_3   = 28
	HEI_4   = 29
	HEI_5   = 30
	HEI_6   = 31
	HEI_7   = 32
	HEI_8   = 33
	HEI_9   = 34
	HEI_10  = 35
	HEI_11  = 36
	HEI_12  = 37
	HEI_13  = 38
	MEI_1   = 39
	MEI_2   = 40
	MEI_3   = 41
	MEI_4   = 42
	MEI_5   = 43
	MEI_6   = 44
	MEI_7   = 45
	MEI_8   = 46
	MEI_9   = 47
	MEI_10  = 48
	MEI_11  = 49
	MEI_12  = 50
	MEI_13  = 51
)

var (
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
)

// NewCards 构造一局新牌
func NewCards() []int {
	//构造一个大小为54的数组,按照红、黑、梅、方、小王、大王
	vals := make([]int, 54)
	for i := 0; i < 54; i++ {
		vals[i] = i
	}

	return vals
}

// Shuffle 洗牌
func Shuffle(vals []int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //根据系统时间戳初始化Random
	for n := len(vals); n > 0; n-- {
		randIndex := r.Intn(n)                                  //得到随机index
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1] //最后一张牌和第randIndex张牌互换
	}
}

func PrintCards(vals []int) {
	for i := 0; i < len(vals); i++ {
		fmt.Printf("%s,", ConvertVal2PrintChars(vals[i]))
	}
}

func ConvertVal2PrintChars(val int) string {
	if val < 13 {
		return fmt.Sprintf("%s%d", red("♥"), val+1)
	} else if val < 26 {
		return fmt.Sprintf("%s%d", blue("♦"), val-12)
	} else if val < 39 {
		return fmt.Sprintf("%s%d", green("♠"), val-25)
	} else if val < 52 {
		return fmt.Sprintf("%s%d", yellow("♣"), val-38)
	} else if val == 52 {
		return cyan("King")
	} else {
		return cyan("Queen")
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
		return "K0"
	} else {
		return "Q0"
	}
}

func ConvertVals2PrintChars(vals []int) []string {
	strs := make([]string, 0)
	for _, val := range vals {
		strs = append(strs, ConvertVal2PrintChars(val))
	}
	return strs
}

func ConvertStr2Val(str string) (int, error) {
	color := strings.ToUpper(str[:1])
	numStr := str[1:]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return -1, err
	}

	switch color {
	case "A":
		return num - 1, nil
	case "B":
		return 12 + num, nil
	case "C":
		return 25 + num, nil
	case "D":
		return 38 + num, nil
	case "Q":
		return 53, nil
	case "K":
		return 52, nil
	default:
		return -1, errors.New("wrong color")

	}
}

func RemoveCard(cards []int, card int) ([]int, error) {
	l := len(cards)
	found := false
	for i := 0; i < l; i++ {
		if card == cards[i] {
			found = true
			cards = append(cards[:i], cards[i+1:]...)
			break
		}
	}

	if found {
		return cards, nil
	}

	return cards, fmt.Errorf("card not found:%d", card)
}

func RemoveCards(remainCards []int, cards []int) ([]int, error) {
	var err error
	tmpCards := make([]int, len(remainCards))
	copy(tmpCards, remainCards)
	for _, card := range cards {
		tmpCards, err = RemoveCard(tmpCards, card)
		if err != nil {
			return remainCards, err
		}
	}

	return tmpCards, nil
}
