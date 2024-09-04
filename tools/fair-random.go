package tools

import (
	"fmt"
	"math/rand"
	"time"
)

func FindNextRunnable(pCount int) {
	// 第1步：选择一个随机数并对pCount取模
	rand.Seed(time.Now().UnixNano())
	initial := rand.Intn(pCount)

	// 第2步：找到小于pCount且与pCount互质的数
	pNumber := findPrimeNumber(pCount)
	if len(pNumber) == 0 {
		fmt.Println("没有与pCount互质的数")
		return
	}

	// 使用互质数进行迭代计算
	fmt.Printf("Initial value: %d\n", initial)
	for i := 0; i < pCount; i++ {
		initial = (initial + pNumber[initial%len(pNumber)]) % pCount
		fmt.Printf("Next value: %d\n", initial)
	}
}

// 返回小于n且与n互质的所有数
func findPrimeNumber(n int) []int {
	var res []int
	for i := 1; i < n; i++ {
		if gcd(i, n) == 1 {
			res = append(res, i)
		}
	}
	return res
}

// 计算a和b的最大公约数
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
