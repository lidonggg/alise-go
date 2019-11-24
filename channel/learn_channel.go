package main

import (
	"fmt"
)

func main() {

	readChan()
}

func deadlock1() {
	//ch := make(chan int) // 无缓冲的通道
	//ch <- 2
	//fmt.Print(<-ch)

	ch := make(chan int) // 无缓冲的通道
	go func() {
		ch <- 2
	}()
	fmt.Print(<-ch)
}

func deadlock2() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		ch2 <- 2
		ch1 <- 1
	}()

	fmt.Print(<-ch2)
}

func deadlock3() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		ch2 <- 2
		ch1 <- 1
	}()

	<-ch1
	<-ch2
}

func deadlock4() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		ch2 <- 2
		ch1 <- 1
	}()
	go func() {
		<-ch1
		<-ch2
	}()

}

func deadlock5() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	for c := range ch {
		fmt.Println(c)
	}
}

func getIntChan() <-chan int {
	num := 5
	ch := make(chan int, num)
	for i := 0; i < num; i++ {
		ch <- i
	}
	close(ch)
	return ch
}

func printInChan() {
	intChan := getIntChan()
	for elem := range intChan {
		fmt.Printf("The element in intChan: %v\n", elem)
	}
}

//func mirroredQuery() string {
//	responses := make(chan string, 3)
//	go func() { responses <- request("asia.gopl.io") }()
//	go func() { responses <- request("europe.gopl.io") }()
//	go func() { responses <- request("americas.gopl.io") }()
//	return <-responses
//}
//func request(hostname string) (response string) { /* ... */ }

func readChan() {
	ch := make(chan int, 3)
	ch <- 1   // 发送数值 1 给这个通道
	ch <- 2   // 发送数值 2 给这个通道
	ch <- 3   // 发送数值 3 给这个通道
	x := <-ch // 从通道里读取值，并把读取的值赋值给 x 变量
	fmt.Printf("The first value received from channel ch is %d\n", x)
	y := <-ch // 从通道里读取值，并把读取的值赋值给 y 变量
	fmt.Printf("The second value received from channel ch is %d\n", y)
	fmt.Println(cap(ch))
	fmt.Println(len(ch))
}
