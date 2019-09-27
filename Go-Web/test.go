package main

import (
	"fmt"
	"os"
)

type friend struct {
	left  int
	right int
}

func main() {
	n := 0
	m := 0
	n, _ = fmt.Scan(&n, &m)

	friends := make([]friend, m)

	var tmpFriend friend

	a := 0
	b := 0
	for i := 0; i < m; i++ {
		n, _ := fmt.Scan(&a, &b)
		if n == 0 {
			break
		} else {
			friends[i].left = a
			friends[i].right = b
		}
	}

	if m == 1 {
		fmt.Println(friends[0].right - friends[0].left + 1)
		os.Exit(0)
	}

	if m == 2 {
		for i := 0; i < 2; i++ {
			if i == 0 {
				fmt.Println(friends[i].right - friends[i].left + 1)
			}
			tmp := quantity(friends[i-1], friends[i])
			fmt.Println(tmp.right - tmp.left + 1)
			os.Exit(0)
		}
	}

	for i := 0; i < m; i++ {
		if i == 0 {
			fmt.Println(friends[0].right - friends[0].left + 1)
			continue
		}
		if i == 1 {
			tmpFriend = quantity(friends[i-1], friends[i])
			result := tmpFriend.right - tmpFriend.left + 1
			fmt.Println(result)
			continue
		}
		tmpFriend = quantity(tmpFriend, friends[i])
	}
	result := tmpFriend.right - tmpFriend.left + 1
	fmt.Println(result)
}

func quantity(a, b friend) friend {
	var tmp friend
	tmp.left = 0
	tmp.right = 0
	if a.left < b.left {
		tmp.left = a.left
	} else {
		tmp.left = b.left
	}
	if a.right > b.right {
		tmp.right = a.right
	} else {
		tmp.right = b.right
	}
	return tmp
}
