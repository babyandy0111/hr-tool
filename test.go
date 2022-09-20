package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	//var age int
	//var name string
	//var phone string
	//
	//fmt.Println("what's your name：")
	//fmt.Scanln(&name)
	//
	//fmt.Println("how old are you：")
	//fmt.Scanln(&age)
	//
	//fmt.Println("may i have your phone number：")
	//fmt.Scanln(&phone)
	//
	//fmt.Println("我認識你了！")
	//fmt.Println("你的名字是：", name)
	//fmt.Println("年齡：", age)
	//fmt.Println("電話：", phone)

	rand.Seed(time.Now().UnixNano())
	var box = rand.Intn(20)
	var num int
	fmt.Println(box)

	// fmt.Println("請猜數字：")
	//fmt.Scanln(&num)
	for {
		fmt.Println("請輸入一個數字，20以下：")
		fmt.Scan(&num)
		if num == box {
			fmt.Println("猜對了！")
			break
		} else {
			fmt.Println("猜錯了!，請重新輸入：")
		}
	}

}
