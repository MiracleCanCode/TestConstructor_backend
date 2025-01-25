package repository

import (
	"fmt"
	"log"
	"sort"
)

func r() {
	var nums = make([]int, 3)
	for i := 0; i < len(nums); i++ {
		_, err := fmt.Scan(&nums[i])
		if err != nil {
			log.Fatal(err)
		}
	}
	sort.Ints(nums)

	for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}

	fmt.Println(nums[1])
}
