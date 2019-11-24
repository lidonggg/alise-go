package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	//timeout := time.After(5 * time.Second)
	//doSelect()
	//simulateTimeout()
	doSelectWithCycle()
}

func doSelect() {
	// 准备好几个通道
	intChannels := [3]chan int{
		make(chan int, 1),
		make(chan int, 1),
		make(chan int, 1),
	}
	// 随机选择一个通道，并向它发送元素值
	index := rand.Intn(3)
	fmt.Printf("The index: %d\n", index)
	intChannels[index] <- index

	// 哪一个通道中有可取的元素值，哪个对应的分支就会被执行
	select {
	case <-intChannels[0]:
		fmt.Println("first")
	case <-intChannels[1]:
		fmt.Println("second")
	case elem := <-intChannels[2]:
		fmt.Printf("The third candidate case is selected, the element is %d.\n", elem)
	default:
		fmt.Println("No candidate case is selected!")
	}
}

func doSelectWithCycle() {
	intChan := make(chan int, 1)
	// 一秒后关闭通道。
	time.AfterFunc(time.Second, func() {
		close(intChan)
	})
	select {
	case _, ok := <-intChan:
		if !ok {
			fmt.Println("The candidate case is closed.")
			break
		}
		fmt.Println("The candidate case is selected.")
	}
}

//
// 模拟超时
//
func simulateTimeout() {
	ch := make(chan int, 2)
	timeout := time.After(5 * time.Second)
	for {
		select {
		case elem := <-ch:
			fmt.Printf("Channel got value %v", elem)
		case <-timeout:
			fmt.Println("Timeout")
			return
		}
	}
}
