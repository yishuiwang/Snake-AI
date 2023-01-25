package search

func Circle(sneak [][]int) []int {
	var step []int
	sz := 0
loop:
	for j := 0; j < 11; j++ {
		// 向下走一步
		step = append(step, 1)
		// 向右走直到尽头减一
		for i := 0; i < height-2; i++ {
			step = append(step, 3)
		}
		// 向下走一步
		step = append(step, 1)
		// 向左走直到尽头
		for i := 0; i < height-2; i++ {
			step = append(step, 2)
		}
	}
	// 向下走一步
	step = append(step, 1)
	// 向右走直到尽头
	for i := 0; i < height-1; i++ {
		step = append(step, 3)
	}
	// 向上走直到尽头
	for i := 0; i < width-1; i++ {
		step = append(step, 0)
	}
	// 向左走直到尽头
	for i := 0; i < height-1; i++ {
		step = append(step, 2)
	}
	sz++
	if sz == 200 {
		return step
	}
	goto loop
}
