## Go 语言中的函数

[demo](https://github.com/lidonggg/alise-go/tree/master/function)

Go 语言中的函数是一等公民（first-class），函数类型也是一等的数据类型，这意味着函数不但可以用于封装代码、分割功能、解耦逻辑，还可以化身为普通的值，在其他函数间传递、赋予变量、做类型判断和转换等，就像切片和字典的值那样。更深层次的含义就是：函数值可以由此成为能够被随意传播的独立逻辑组件（功能模块）。

对于函数类型来说，它是一种对一组输入、输出进行模板化的重要工具，它比接口更加轻巧灵活，它的值也借此变成了可被热替换的逻辑组件。

如以下代码所示：
```go
package main

import "fmt"

// 函数声明，func 关键字在类型名称的右边
// func 关键字的右边就是参数列表和结果列表，参数列表必须用括号包起来，返回值列表多个的话也需要用括号包起来
type Printer func(contents string) (n int, err error)
// 函数签名
func printToStd(contents string) (bytesNum int, err error) {
    return fmt.Println(contents)
}

func main() {
    var p Printer
    p = printToStd
    p("something")
}
```

### 编写高阶函数
#### **什么是高阶函数**
只要满足以下的任意一个条件，我们即可以说这个函数是一个高阶函数：

- 参数列表：接收其他的函数作为参数传入；
- 结果列表：把其他的函数作为返回的结果。

接下来我们通过编写一个 calculate 函数来实现两整数之间的加减乘除操作，但是希望两个整数的具体操作都是由调用方给出的而不是由它自己本身来实现。

为了实现这一点，我们声明一个 operate 函数，用来执行具体的操作，代码如下：
```java
type operate func(x, y int) int
```
它接收两个 int 类型的参数，并且返回值也是一个 int 类型。

接下来编写 calculate 函数，这个函数除了需要接收两个 int 类型的参数之外，还应该接收一个 operate 类型的操作，用来执行真正的运算操作。calculate 函数应该返回两个结果，一个是代表操作结果的 int 类型，另外一个是 error 类型，因为如果 operate 类型的参数值为 nil 的话，那就应该直接返回一个错误。这里说明一下：**函数类型是引用类型，它的值是可以为 nil 的，而这种类型的 0 值恰好是 nil。**

calculate 函数的代码如下所示：
```go
func calculate(x int, y int, op operate) (int, error) {
	// 利用卫语句进行参数检查
	if op == nil {
		return 0, errors.New("invalid operation")
	}
	// 参数正确，进行相应的 operate 操作，没有错误发生，第二个返回值返回 nil
	return op(x, y), nil
}
```

那么如何编写 operate 类型的函数值？我们可以在调用之前先声明好一个函数，然后把它赋值给一个变量，也可以直接编写一个实现了 operate 类型的匿名函数：
```go
// 调用方法一：声明一个函数并赋值给一个变量
op := func(x, y int) int {
	return x + y
}
res, _ := calculate(1, 2, op)
fmt.Printf("res1: %d\n", res)

// 调用方法二：实现了 operate 类型的匿名函数
res2, _ := calculate(2, 3, func(x, y int) int {
	return x * y
})
fmt.Printf("res2: %d", res2)
```

完成了 calculate 函数的例子之后，我们可以发现其实它就是一个高阶函数，因为它接收了其他的函数作为参数的输入。

接下来看一下如何把其他的函数作为返回的结果，如<span id="genCalculator">下述代码</span>所示：
```go
type calculateFunc func(x int, y int) (int, error)

func genCalculator(op operate) calculateFunc {
	return func(x int, y int) (int, error) {
		if op == nil {
			return 0, errors.New("invalid operation")
		}
		return op(x, y), nil
	}
}

func main() {
	op := func(x, y int) int {
		return x + y
	}
	x, y := 56, 78
	add := genCalculator(op)
	res3, _ := add(x, y)
	fmt.Printf("res3: %d\n", res3)
}
```
### 实现闭包（Closure）
首先让我们理解一下外来标识符的含义。所谓的外来标识符，它既不代表函数的任何参数或者返回值，也不是函数内部声明的变量，它是直接从外面拿过来的。外来标识符也叫自由变量，可见它实际上就是一个变量。闭包就是一个持有外部变量的函数，它是一个由不确定变成确定的过程。闭包函数因为引用了自由变量而变成了一个“不确定”状态，也叫“开放”状态。也就是说闭包函数的内部结构是不完整的，有一部分逻辑需要这个自由变量参与完成，在闭包函数被定义的时候，这个自由变量的具体含义是未知的。即使对于像 Go 语言这种静态类型的编程语言而言，我们在定义闭包函数的时候也最多只能知道自由变量的类型，对于像 javascript 这种动态类型的编程语言，我们甚至连其类型也不一样。

**简单来说，闭包就是能够读取其他函数内部变量的函数**。Stuart 在他的 [ppt](https://app.box.com/s/elkumrpfng) 中对于闭包有这么一段描述：**In computer science, a closure is a function that is evaluated in an environment containing one or more bound variables. When called, the function can access these variables.**，大致意思就是在计算机科学领域中，闭包就是在一个包含一个或多个变量的环境中运行的函数，并且在调用的时候，它可以获取到这些变量。他还提到说：**closure:where a function remembers what happens around it**（闭包：一个能够记住它周围发生了什么的函数），以及 **one functiondefined inside another**（在另外一个函数中定义的函数）。

对于[上述代码](#genCalculator)来说，genCalculator 函数内部实际上就实现了一个闭包（return func），它里面使用的变量 op 既不代表它的任何参数或者结果也不是它自己声明的，而是定义它的 genCalculator 函数的参数，所以是一个自由变量。当运行到 **if op == nil** 这一行的时候，Go 语言编译器会试着去寻找 op 所代表的的含义，它会发现 op 代表的是 genCalculator 函数的参数，然后它会把这两者联系起来，这时可以说，自由变量 op 被捕获了。

利用闭包，我们可以在程序运行的过程中，根据需要生成不同的函数，继而影响后续的程序行为，这与设计模式中的“模板模式”有着异曲同工之妙。

### 函数参数传递
如下述代码所示：
```go
package main

import "fmt"

//
// 函数参数值相关
//
func main() {
	// 示例1
	array1 := [3]string{"a", "b", "c"}
	fmt.Printf("The array: %v\n", array1)
	// 数组是值类型，每一次复制都会重新拷贝一份，是一种深层复制
	array2 := modifyArray(array1)                  // [a b c]
	fmt.Printf("The modified array: %v\n", array2) // [a x c]
	fmt.Printf("The original array: %v\n", array1) // [a b c]
	fmt.Println()

	// 示例2
	slice1 := []string{"x", "y", "z"}
	fmt.Printf("The slice: %v\n", slice1) // [x y z]
	// 切片是引用类型，作为参数的时候复制的是它本身，是一种浅表复制
	slice2 := modifySlice(slice1)
	fmt.Printf("The modified slice: %v\n", slice2) // [x i z]
	fmt.Printf("The original slice: %v\n", slice1) // [x i z]
	fmt.Println()

	// 示例3
	complexArray1 := [3][]string{
		{"d", "e", "f"},
		{"g", "h", "i"},
		{"j", "k", "l"},
	}
	fmt.Printf("The complex array: %v\n", complexArray1) // [[d e f] [g h i] [j k l]]
	complexArray2 := modifyComplexArray(complexArray1)
	fmt.Printf("The modified complex array: %v\n", complexArray2) // [[d e f] [g s i] [o p q]]
	fmt.Printf("The original complex array: %v\n", complexArray1) // [[d e f] [g s i] [j k l]]
}

// 修改数组的元素值
func modifyArray(a [3]string) [3]string {
	a[1] = "x"
	return a
}

// 修改切片的元素值
func modifySlice(a []string) []string {
	a[1] = "i"
	return a
}

// 修改复杂数组，修改后产生的效果会随着修改的目标元素的类型的不同而不同
func modifyComplexArray(a [3][]string) [3][]string {
	a[1][1] = "s"
	a[2] = []string{"o", "p", "q"}
	return a
}
```
关于函数参数传递，我们有一点需要注意：所有传给函数的参数值都会被复制，函数在其内部使用的并不是参数值的原值，而是它的副本。

由于数组是值类型的，所以每一次复制都会拷贝它以及它的所有元素值，此时如果修改它的副本，并不会对它本身造成影响。

而对于引用类型，如切片、字典、通道等，复制它们的时候，其实是拷贝它们本身，而不会拷贝它们的底层数组，也就是说这这是一种浅层的复制，而不是深层复制。比如对于切片来说，复制的是它指向的底层数组中某一个元素的指针以及它的长度和容量等属性，因此此时修改新副本的切片，会对原来的切片产生同样的影响。

### 总结
Go 语言中，函数是一等公民（first class），它既可以被独立声明也可以被当成普通的值来传递或赋予变量。除此之外，我们还可以在其他函数的内部声明匿名函数并把它直接赋给变量。

Go 语言可以实现高阶函数，并且也可以通过高阶函数来实现闭包，从而可以做到逻辑的动态生成。

关于函数参数的传递，我们要注意针对不同的类型（值类型、引用类型），参数的复制过程。







