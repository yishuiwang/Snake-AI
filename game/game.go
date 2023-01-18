package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	MoveUp = iota
	MoveDown
	MoveLeft
	MoveRight
)

var (
	score = 0
)

type PlayGround struct {
	direction  int
	snake      [][]int
	food       [][]int
	background [][]int
	width      int
	height     int
}

func newPlayGround() *PlayGround {
	snake := [][]int{{0, 2}, {0, 1}, {0, 0}}
	food := [][]int{{1, 2}}
	width := 25
	height := 32
	var background [][]int
	for i := 0; i < width; i++ {
		background = append(background, []int{})
		for j := 0; j < height; j++ {
			background[i] = append(background[i], 0)
		}
	}

	p := &PlayGround{
		direction:  MoveRight,
		snake:      snake,
		food:       food,
		background: background,
		width:      width,
		height:     height,
	}

	p.updatePlayGround()
	return p
}

func (p *PlayGround) updatePlayGround() {
	for i := 0; i < p.width; i++ {
		for j := 0; j < p.height; j++ {
			p.background[i][j] = 0
		}
	}

	for _, v := range p.snake {
		p.background[v[0]][v[1]] = 1
	}

	for _, v := range p.food {
		p.background[v[0]][v[1]] = 2
	}
}

func (p *PlayGround) move(direction int) bool {
	//避免反向移动
	p.changeDirection(direction)
	direction = p.direction

	x := p.snake[len(p.snake)-1][0]
	y := p.snake[len(p.snake)-1][1]
	tail := []int{x, y}

	//移动身体
	if len(p.snake) > 1 {
		for i := len(p.snake) - 1; i > 0; i-- {
			p.snake[i][0] = p.snake[i-1][0]
			p.snake[i][1] = p.snake[i-1][1]
		}
	}

	switch direction {
	case MoveUp:
		p.snake[0][0]--
	case MoveDown:
		p.snake[0][0]++
	case MoveLeft:
		p.snake[0][1]--
	case MoveRight:
		p.snake[0][1]++
	}

	//是否碰撞
	if p.isCollision() {
		return false
	}
	//是否吃到食物
	if p.eatFood() {
		score++
		p.snake = append(p.snake, tail)
		p.food = [][]int{}
		p.randomFood()
	}
	//更新背景
	p.updatePlayGround()
	return true
}

func (p *PlayGround) changeDirection(d int) bool {
	result := map[int]int{
		MoveUp:    MoveDown,
		MoveDown:  MoveUp,
		MoveLeft:  MoveRight,
		MoveRight: MoveLeft,
	}
	if p.direction != result[d] {
		p.direction = d
		return true
	}
	return false
}

//碰撞检测
func (p *PlayGround) isCollision() bool {
	//边界碰撞
	if p.snake[0][0] < 0 || p.snake[0][0] >= p.width || p.snake[0][1] < 0 || p.snake[0][1] >= p.height {
		return true
	}
	//自身碰撞
	for i := 1; i < len(p.snake); i++ {
		if p.snake[0][0] == p.snake[i][0] && p.snake[0][1] == p.snake[i][1] {
			return true
		}
	}
	return false
}

//吃食物
func (p *PlayGround) eatFood() bool {
	if p.snake[0][0] == p.food[0][0] && p.snake[0][1] == p.food[0][1] {
		return true
	}
	return false
}

//随机生成食物
func (p *PlayGround) randomFood() {
	//TODO
	//0-24随机数
	x := rand.Intn(p.width)
	y := rand.Intn(p.height)
	//判断是否与蛇重合
	for _, v := range p.snake {
		if v[0] == x && v[1] == y {
			p.randomFood()
			return
		}
	}
	p.food = append(p.food, []int{x, y})
}

func show(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += 1
	}
}

func (p *PlayGround) render(background [][]int) {
	w, h := termbox.Size()
	l := w/2 - p.height
	t := h/2 - p.width/2
	//在正中央渲染
	for i := 0; i < len(background); i++ {
		for j := 0; j < len(background[i]); j++ {
			if background[i][j] == 0 {
				show(l+j*2, i+t, termbox.ColorWhite, termbox.ColorWhite, "  ")
			}
			if background[i][j] == 1 {
				show(l+j*2, i+t, termbox.ColorGreen, termbox.ColorGreen, "██")
			}
			if background[i][j] == 2 {
				show(l+j*2, i+t, termbox.ColorRed, termbox.ColorRed, "██")
			}
		}
	}
	//显示分数
	show(l, t-2, termbox.ColorBlue, termbox.ColorDefault, "Score: "+strconv.Itoa(score))
	//显示提示
	show(l, t+p.width+1, termbox.ColorBlue, termbox.ColorDefault, "Press ESC or Ctrl+C to exit")

	termbox.Flush()
}

func finallyScore() {
	w, h := termbox.Size()
	l := w/2 - 10
	t := h/2 - 2
	show(l, t, termbox.ColorWhite, termbox.ColorDefault, "Game Over!")
	show(l, t+1, termbox.ColorWhite, termbox.ColorDefault, "Your score is: "+strconv.Itoa(score))
	termbox.Flush()
}

func moveEventKey(ev termbox.Event) int {
	var direction int
	if ev.Key == termbox.KeyArrowUp || ev.Ch == 'w' {
		direction = MoveUp
	} else if ev.Key == termbox.KeyArrowDown || ev.Ch == 's' {
		direction = MoveDown
	} else if ev.Key == termbox.KeyArrowLeft || ev.Ch == 'a' {
		direction = MoveLeft
	} else if ev.Key == termbox.KeyArrowRight || ev.Ch == 'd' {
		direction = MoveRight
	} else {
		direction = -1
	}
	return direction
}

func newGame() {
	rand.Seed(time.Now().UnixNano())
	playGround := newPlayGround()
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	playGround.render(playGround.background)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
				return
			}
			s := moveEventKey(ev)
			if s != -1 {
				if playGround.move(s) {
					playGround.render(playGround.background)
				} else {
					finallyScore()
					time.Sleep(3 * time.Second)
					return
				}
			}
			//time.Sleep(time.Second / 10)
		}
	}
}
