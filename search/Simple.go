package search

func Common(sneak, food [][]int) []int {

	p1 := sneak[0]
	p2 := food[0]

	var step []int
	x, y := p1[0]-p2[0], p1[1]-p2[1]

	if x > 0 {
		for i := 0; i < x; i++ {
			step = append(step, MoveUp)
		}
	} else {
		for i := 0; i < -x; i++ {
			step = append(step, MoveDown)
		}
	}

	if y > 0 {
		for i := 0; i < y; i++ {
			step = append(step, MoveLeft)
		}
	} else {
		for i := 0; i < -y; i++ {
			step = append(step, MoveRight)
		}
	}

	return step
}
