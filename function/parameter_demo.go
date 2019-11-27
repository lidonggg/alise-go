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
