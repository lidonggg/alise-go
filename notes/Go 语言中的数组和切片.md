# Go 语言中的数组和切片

go 语言中的数组和切片都属于集合类的类型，它们之间最主要的区别是：数组的长度是固定的，而切片是可变长度的。

数组的长度在声明它的时候就必须给出，并且在整个生命周期中都是不会再变的，因此可以说数组的长度是其类型的一部分。而切片的类型字面量中只有元素的类型，而没有长度，切片的长度是可以自动地随着其中元素数量的增长而增长，但不会随着元素数量的减少而减少。

我们可以把切片当成是对数组的一层封装，因为在每个切片的底层数据结构中，一定会包含一个数组。数组可以看成是切片的底层数据结构，而切片则可以看成是对数组的某个连续片段的引用。正因为如此，go 语言中的切片是引用类型，而数组则是值类型。

通过以下代码我们可以看到数组和切片在类型上的不同：
```go
// 数组是值类型的
numbers1 := [...]int{1, 2, 3, 4, 5, 6}
maxIndex1 := len(numbers1) - 1
// range 表达式的求值结果会被复制，因此被迭代的对象是 range 表达式结果值的副本而不是原值
for i, e := range numbers1 {
	if i == maxIndex1 {
		numbers1[0] += e
	} else {
		numbers1[i+1] += e
	}
}
fmt.Println(numbers1) // [7 3 5 7 9 11]

// 切片是引用类型
numbers2 := []int{1, 2, 3, 4, 5, 6}
maxIndex2 := len(numbers2) - 1
for i, e := range numbers2 {
	if i == maxIndex2 {
		numbers2[0] += e
	} else {
		numbers2[i+1] += e
	}
}
	fmt.Println(numbers2) // [22 3 6 10 15 21]
```

我们可以通过内建函数 len() 得到数组或切片的长度，通过内建函数 cap() 可以得到它们的容量，需要注意的是，数组的容量和长度永远是相等且不可变的。

## 切片的扩容规则
切片由于其长度可变的特性，一旦一个切片无法容纳更多地元素，Go 语言就会想办法扩容。Go 语言的扩容机制是不改变原有的切片，而是生成一个容量更大的切片，然后把原有切片的元素一并拷贝到新的切片中。

一般情况下，新切片的容量是老切片容量的两倍，但是当原切片的长度大于或等于 1024 的时候，Go 语言将会以原容量的 1.25 倍作为基准进行扩容。新容量基准会被调整直到结果不小于原长度与所要追加的长度之和。另外如果一次追加的元素过多，以至于新长度比原来长度的 2 倍还要大，那么新容量就会直接以新的长度作为基准。通过以下示例我们可以验证上述的扩容规则：
```go
package main
import "fmt"

func main() {
	// 示例1
	s1 := make([]int, 0)
	fmt.Printf("The capacity of s1: %d\n", cap(s1))  // 1，1
	for i := 1; i <= 5; i++ {
		s1 = append(s1, i)
		fmt.Printf("s1(%d): len: %d, cap: %d\n", i, len(s1), cap(s1))
	}
	fmt.Println()

	// 示例2
	s2 := make([]int, 1024)
	fmt.Printf("The capacity of s2: %d\n", cap(s2))  // 1024
	s2e1 := append(s2, make([]int, 200)...)
	fmt.Printf("s7e1: len: %d, cap: %d\n", len(s2e1), cap(s2e1))  // 1224（1024+200），1280（1024*1.25）
	s2e2 := append(s2, make([]int, 400)...)
	fmt.Printf("s2e2: len: %d, cap: %d\n", len(s2e2), cap(s2e2))  // 1424，1696
	s2e3 := append(s2, make([]int, 600)...)
	fmt.Printf("s2e3: len: %d, cap: %d\n", len(s2e3), cap(s2e3))  // 1624,2048
	fmt.Println()

	// 示例3
	s3 := make([]int, 10)
	fmt.Printf("The capacity of s3: %d\n", cap(s3))  // 10
	s3a := append(s3, make([]int, 11)...)
	fmt.Printf("s3a: len: %d, cap: %d\n", len(s3a), cap(s3a))  // 21,22
	s3b := append(s3a, make([]int, 23)...)
	fmt.Printf("s3b: len: %d, cap: %d\n", len(s3b), cap(s3b))  // 44,44
	s3c := append(s3b, make([]int, 45)...)
	fmt.Printf("s3c: len: %d, cap: %d\n", len(s3c), cap(s3c))  // 89,96
}
```

上述代码中用到了 append() 函数，在无需扩容的时候，它返回的是指向原底层数组的心切片，而在需要扩容的时候，它返回的是指向新底层数组的新切片。同时需要注意的是，只要新长度不会超过切片的原容量，那么使用 append 函数对其追加元素的时候就不会引起扩容。