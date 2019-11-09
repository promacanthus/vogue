// 在Go中，数组是定长的元素编号序列

package main

import "fmt"

func main() {
	var a [5]int // 创建一个存储5个int的数组a,元素类型和长度都是数组类型的一部分,默认情况下数组中存储零值
	fmt.Println("emp:", a)

	a[4] = 100 //使用array[index]=value语法,在索引处设置值
	fmt.Println("set:", a)
	fmt.Println("get:", a[4])   // 使用array[index]来获取值
	fmt.Println("len:", len(a)) // 内置函数len返回数组的长度

	b := [5]int{1, 2, 3, 4, 5} // 使用这种语法在一行中声明并初始化数组
	fmt.Println("dcl:", b)     // 使用fmt.Println将会以[v1 v2 v3 ...]形式输出数组

	var twoD [2][3]int // 数组类型是一维的,但是可以通过组合类型来创建多维数据结果
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d: ", twoD)
}

// 切片比数组在Go中更常用
