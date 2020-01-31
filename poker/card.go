package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
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
	r := rand.New(rand.NewSource(time.Now().Unix())) //根据系统时间戳初始化Random
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
	if val == 0 {
		return red("King")
	} else if val < 14 {
		return fmt.Sprintf("%s%d", red("♥"), val)
	} else if val < 27 {
		return fmt.Sprintf("%s%d", blue("♠"), val-13)
	} else if val < 40 {
		return fmt.Sprintf("%s%d", green("♣"), val-26)
	} else if val < 53 {
		return fmt.Sprintf("%s%d", yellow("♦"), val-39)
	} else {
		return cyan("Queen")
	}
}

// ConvertVal2Str 将牌数字转换成字符串
func ConvertVal2Str(val int) string {
	if val == 0 {
		return "K0"
	} else if val < 14 {
		return fmt.Sprintf("A%d", val)
	} else if val < 27 {
		return fmt.Sprintf("B%d", val-13)
	} else if val < 40 {
		return fmt.Sprintf("C%d", val-26)
	} else if val < 53 {
		return fmt.Sprintf("D%d", val-39)
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
		return num, nil
	case "B":
		return 13 + num, nil
	case "C":
		return 26 + num, nil
	case "D":
		return 39 + num, nil
	case "Q":
		return 53, nil
	case "K":
		return 0, nil
	default:
		return -1, errors.New("wrong color")

	}
}
