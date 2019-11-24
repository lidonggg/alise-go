## Go 语言中的通道

[demo]()

---

通道（channel）是 Go 语言中最具有特色的数据类型，我们可以利用通道在多个 goroutine 之间传递数据。

### 啥是通道？
通道类型是天然并发安全的，它是 Go 语言自带的、唯一一个可以满足并发安全性的类型。

在多个 goroutine 并发中，我们可以通过原子函数和互斥锁保证对共享资源的安全性，消除竞争状态，但这样会比较影响性能。除此之外，我们还可以通过使用通道，在多个 goroutine 之间发送和接收共享的数据，达到数据同步的目的。这就像在两个 goroutine 之间架设的管道，一个 goroutine 可以往管道里面塞数据，另一个 goroutine 可以从这个管道里面取数据，就类似于我们所熟知的队列（从下文可以知道它是一个先进先出（FIFO）的队列）。

### 通道的使用
在声明并初始化一个通道的时候，我们需要用到 Go 语言内建的 make 函数，它需要接收的第一个参数代表了通道的具体类型的类型字面量，如下述代码所示：
```go
ch := make(chan int)
```
其中 chan 是表示通道类型的关键字，而 int 则说明了该通道的元素类型。

在初始化通道的时候，make 函数除了必须接收这样的类型字面量作为参数，还可以接收一个 int 类型的可选参数（值必须大于等于 0 ）用于表示通道的容量。所谓的通道的容量，就是指通道最多可以缓存多少个元素值。当不传第二个参数或第二个参数值为 0 的时候，此时的通道我们称为**无缓冲的通道**，否则称为**有缓冲的通道**。如下述代码所示：
```go
ch := make(chan int)	  // 无缓冲的通道
ch2 := make(chan int, 0)  // 无缓冲的通道
ch3 := make(chan int, 3)  // 有缓冲的通道
```

一个通道相当于一个先进先出（FIFO）的队列，通道中的各个元素都是严格地按照发送的顺序排列的，先被发送的元素值一定也是先被接收的。元素的接收和发送这个个操作的运算符都是 **<-** 。如下述代码所示：
```go
ch := make(chan int, 3)
ch <- 1   // 发送数值 1 给这个通道
ch <- 2   // 发送数值 2 给这个通道
ch <- 3   // 发送数值 3 给这个通道
x := <-ch // 从通道里读取值，并把读取的值赋值给 x 变量
fmt.Printf("The first value received from channel ch is %d\n", x)    // The first value received from channel ch is 1
y := <-ch // 从通道里读取值，并把读取的值赋值给 y 变量
fmt.Printf("The first value received from channel ch is %d\n", y)    // The second value received from channel ch is 2
```
以上代码只适合有缓冲通道，对于无缓冲通道，类似的做法会产生死锁（deadlock），如<span id="deadlock1">以下代码</span>所示：
```go
ch := make(chan int) // 无缓冲的通道
ch <- 2
fmt.Printf(<-ch)     // 产生 deadlock
```
原因就是无缓冲通道和有缓冲通道有着不同的数据传递方式。
### 无缓冲通道
无缓冲通道本身不存储信息，它只负责转手，有人传给它，它就必须要传给别人，如果只进行接收或者发送其中某一个操作，都会造成阻塞。对于[上述代码](#deadlock1)来说，只有一个 goroutine 即主 goroutine，第二行代码阻塞在传值，第三行代码阻塞在取值，因此主线程会一直卡主，系统一直在等待，所以会被判定为产生 deadlock 然后结束程序。

因此我们可以发现，无缓冲通道要求发送 goroutine 和接收 goroutine 必须同时准备好且是两个不同的协程。因此<span id="deadlockExtend1">以下代码</span>是不会发生死锁错误的：
```go
ch := make(chan int) // 无缓冲的通道
go func() {          // 新开辟一个协程
	ch <- 2
}()
fmt.Print(<-ch)      // 成功打印出 2
```

#### **死锁延伸1**
考虑下述代码：
```go
ch1 := make(chan int)
ch2 := make(chan int)
go func() {
	ch2 <- 2        // a
	ch1 <- 1        // b
}()

<-ch1
```
它会产生死锁吗？答案是肯定的，也就是说，依然会产生死锁。这段代码并不能保证是主线程的 <-ch1 先执行，还是子线程先执行。如果是主线程先执行，那么它会阻塞直到有其他的线程往 ch1 传值，然而如果子线程开始执行了，会首先执行 ch2 <- 2 这段代码，它同样会等待有其他协程去接收值，然而这里并没有，因此会阻塞在这里而发生 deadlock。如果是子线程先执行，那么会直接阻塞在 ch2 <- 2 语句而产生死锁。如果把上述代码 a,b 行调换位置，那么程序将会成功发送变量 ch1，不产生死锁但是 ch2 会一直阻塞下去然后程序终止运行

#### **死锁延伸2**
考虑如下代码：
```go
ch1 := make(chan int)
ch2 := make(chan int)
go func() {
	ch2 <- 2        // a
	ch1 <- 1        // b
}()

<-ch1
<-ch2
```
它的执行结果依然是死锁，还是会产生主协程和子协程相互等待的情况。但如果调换 a,b 处代码的位置，那么程序将会成功执行。

#### **死锁延伸3**
考虑如下代码：
```go
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
```
此时两个子协程之间会因为相互等待而发生死锁，但是不会影响主协程，所以程序不会报死锁错误。

#### **死锁延伸4**
考虑如下代码：
```go
ch := make(chan int, 2)
ch <- 1
ch <- 2
for c := range ch {
	fmt.Println(c)
}
```

输出结果为：
```go
1
2
fatal error: all goroutines are asleep - deadlock!
```
虽然这里的ch是带有缓冲的信道，但是容量只有两个，当两个输出完之后，可以简单的将此时的信道等价于无缓冲的信道。显然对于无缓冲的信道只是单纯的读取元素是会造成阻塞的，而且是在主协程，所以和最开始的死锁现场等价，故而会死锁。
### 有缓冲通道
考虑如下的应用场景：想获取服务端的一个数据，不过这个数据在三个镜像站点上都存在，这三个镜像分散在不同的地理位置，而我们的目的又是想最快的获取到数据，那么需要怎么做？

在这里我们可以定义一个容量为 3 的通道，然后同时发起 3 个并发的 goroutine 向这三个镜像获取数据，获取到的数据发送到通道中，然后直接返回接收到的第一条数据即可，代码如下所示：
```go
func mirroredQuery() string {
	responses := make(chan string, 3)
	go func() { responses <- request("asia.gopl.io") }()
	go func() { responses <- request("europe.gopl.io") }()
	go func() { responses <- request("americas.gopl.io") }()
	return <-responses
}
func request(hostname string) (response string) { /* ... */ }
```

### 通道的发送和接收操作的一些基本特性
- 对于同一个通道，发送操作之间是互斥的，接收操作之间是互斥的，同时对于通道中的同一个元素值来说，发送操作和接收操作之间也是互斥的。有一个细节需要注意，那就是元素值从外界进入通道时会被复制，也就是说进入通道的并不是接收操作符右边的那个元素值，而是它的副本。

- 发送操作和接收操作中对元素值的处理都是不可分割的。

- 发送操作在完全完成之前会被阻塞，接收操作也是如此。

以上三个特性都是为了保证通道的并发安全而存在的。

### 单向通道
有时候我们可能会限制某个通道只能接收或者只能发送，这种情况我们称之为“单向通道”。定义单向通道我们只要在定义的时候带上 <- 即可：
```go
ch := make(chan int, 1)
var send chan<- int = ch            //只能发送，<- 操作符在 chan 后面
var receive <-chan int = ch         //只能接收，<- 操作符在 chan 前面
```
与发送和接收操作相对应，单向通道里的发送和接收都是站在操作通道的角度上说的。

其实我们可以发现，单向通道是没有用的，通道就是为了传递数据而存在的，一个只有一端能用的通道就失去了意义，那么为什么会有单向通道的存在？

#### **单向通道的作用**
**简单来说，单向通道最主要的用途就是约束其他代码的行为。**这需要从两方面来讲，但都跟函数声明有关。来看下面的一段代码：
```go
func SendInt(ch chan<- int) {
    ch <- rand.Intn(1000)
}
```
这里的 SendInt 函数只接收一个 chan<- int 类型的参数，在这个函数的代码中只能向参数 ch 发送元素值，而不能从它那里接收元素值，这就起到了约束函数行为的作用。然而在这里可能意义不是很大，因为我们自己写的函数自己就能确定操作通道的方式，在实际场景中，这种约束一般会出现在接口类型声明中的某个方法定义上。请看下面的接口类型声明：
```go
type Notifier interface {
    SendInt(ch chan<- int)
}
```
在这里，接口中的 SendInt 方法只会接收一个发送通道作为参数，所以在该接口的所有实现类型中的 SendInt 方法都会受到限制，这种约束在我们编写模板代码或者可扩展的程序库的时候特别有用。

在调用 SendInt 函数的时候，我们只需要把一个元素类型匹配的双向通道传给它就行了，Go 语言会自动地把双向通道改成所需的单向通道：
```go
ch := make(chan int, 3)
SendInt(ch)
```

当然除了参数，单向通道还可以在函数声明的结果列表中被使用，如下述代码所示：
```go
func getIntChan() <-chan int {
	num := 5
	ch := make(chan int, num)
	for i := 0; i < num; i++ {
		ch <- i
	}
	close(ch)
	return ch
}
```
上述代码意味着得到该通道的程序，只能从通道中接收元素值。

在 Go 语言中，我们还可以声明函数类型，如果我们在函数类型中使用了双向通道，那么就相当于在约束所有实现了这个函数类型的函数，如下述代码所示：
```go
intChan := getIntChan()
for elem := range intChan {
	fmt.Printf("The element in intChan: %v\n", elem)
}
```
这条 for 语句会不断地从 intChan 中取出元素值，即使 intChan 被关闭，它也会在取出所有剩余元素值之后再结束。当 intChan 中没有元素的时候，for 循环会被阻塞在 for 关键字那行。直到有新的元素值可取；如果 intChan 的值为 nil，那么它同样会被永远阻塞在有 for 关键字的那一行。

### 管道
我们在操作 bash（shell）的时候有个管道操作符 “|”，它的意思是把上一个操作的输出当成下一个操作的输入，然后做一连串的操作。利用 Go 语言中的通道，我们也可以达到类似的效果，代码如下所示：
```go
ch1 := make(chan int)
ch2 := make(chan int)
go func() {
	ch1 <- 100
}()

go func() {
	v := <-ch1
	ch2 <- v
}()

fmt.Println(<-ch2)
```

### select 语句与通道的联用
Go 语言中的 select 语句是专门为了操作通道而存在的，也就是说 select 语句只能与通道联用。先看下面的一段代码：
```go
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
```
从上面的代码可以发现，select 语法的使用很类似于 switch，同样包含多个分支。由于 select 是专门为通道设计的，所以它的 case 子句只能包含操作通道类型的表达式。select 在使用过程中有如下一些约束：

- 代码执行到 select 时，case 语句会按照源代码的顺序被评估，且只评估一次；
- 除 default 外，如果只有一个 case 通过评估，那么就执行这个 case 里的语句；
- 除 default 外，如果有多个 case 通过评估，那么通过随机的方式去选取一个执行；
- 如果除 default 外的所有 case 都没有通过评估，那么执行 default 中的语句；
- 如果没有 default 语句，那么代码块会被阻塞直到有一个 case 通过评估，否则一直阻塞下去；
- 一条 select 语句中，default 分支最多只能有一条；
- select 语句的每次执行，包括 case 表达式求值和分支选择，都是独立的。不过它的执行是否是并发安全的，就要看 case 子句中的代码是否是并发安全的。

如果想要连续或者定时地操作其中的通道的那，那么就需要把 select 语句放入到一个 for 循环中，如下代码所示：
```go
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
```

### 管道
我们在操作 bash（shell）的时候有个管道操作符 “|”，它的意思是把上一个操作的输出当成下一个操作的输入，然后做一连串的操作。利用 Go 语言中的通道，我们也可以达到类似的效果，代码如下所示：
```go
ch1 := make(chan int)
ch2 := make(chan int)
go func() {
	ch1 <- 100
}()

go func() {
	v := <-ch1
	ch2 <- v
}()

fmt.Println(<-ch2)
```

### 获取通道的容量和长度
```go
cap(ch)     // 容量
len(ch)     // 长度
```
### 通道的关闭
```go
close(ch)
```