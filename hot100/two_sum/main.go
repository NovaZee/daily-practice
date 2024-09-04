package main

var testList = []int{0, 6}
var target = 6

func main() {
	println(twoSum(testList, target))
}

func twoSum(list []int, target int) (int, int) {
	if len(list) == 1 {
		if list[0] != target {
			return 0, 0
		}
	}
	hash := make(map[int]int, len(list))
	for i := range list {
		if v, ok := hash[target-list[i]]; ok {
			return i, v
		}
		hash[list[i]] = i
	}
	return 0, 0
}
