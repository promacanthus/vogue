# 53-random-number

```go
//  Go的math/rand包提供伪随机数生成。

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Print(rand.Intn(100), ",") // rand.Intn(100)返回0<=n<100的随机整数
	fmt.Print(rand.Intn(100))
	fmt.Println()

	fmt.Println(rand.Float64()) // rand.Float64()返回0.0<=f<1.0的64位浮点数

	fmt.Print((rand.Float64()*5)+5, ",")
	fmt.Print((rand.Float64() * 5) + 5)
	fmt.Println()

	// 默认的数字生成器具有确定性，因此默认情况下每次都会产生相同的数字序列
	// 为了生成不同的序列，提供一个变化的种子
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	// 对于需要保密的随机数，上述方式是不安全的，可以使用crypto/rand包

	// 调用rand.Rand的结果就像调用rand包的函数那样
	fmt.Print(r1.Intn(100), ",")
	fmt.Print(r1.Intn(100))
	fmt.Println()

	// 如果提供相同的种子，那么会生成一样是随机数序列
	s2 := rand.NewSource(42)
	r2 := rand.New(s2)
	fmt.Print(r2.Intn(100), ",")
	fmt.Print(r2.Intn(100))
	fmt.Println()

	s3 := rand.NewSource(42)
	r3 := rand.New(s3)
	fmt.Print(r3.Intn(100), ",")
	fmt.Print(r3.Intn(100))
}
```