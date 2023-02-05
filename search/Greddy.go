package search

func Greedy(sneak, food [][]int) []int {

	p1 := sneak[0]
	p2 := food[0]

	var step []int
	x, y := p1[0]-p2[0], p1[1]-p2[1]

	if x > 0 {
		for i := 0; i < x; i++ {
			step = append(step, 0)
		}
	} else {
		for i := 0; i < -x; i++ {
			step = append(step, 1)
		}
	}

	if y > 0 {
		for i := 0; i < y; i++ {
			step = append(step, 2)
		}
	} else {
		for i := 0; i < -y; i++ {
			step = append(step, 3)
		}
	}

	return step
}
