package game

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/nsf/termbox-go"
)

const (
	MoveUp    = 0
	MoveDown  = 1
	MoveLeft  = 2
	MoveRight = 3
	width     = 24
	height    = 32
)

type PlayGround struct {
	Score      int
	Steps      int
	Width      int
	Height     int
	HeadDir    int
	TailDir    int
	Snake      [][]int
	Food       [][]int
	Background [][]int
}

func NewPlayGround() *PlayGround {
	//Snake := [][]int{{0, 2}, {0, 1}, {0, 0}}
	snake := [][]int{{0, 0}}
	food := [][]int{{1, 2}}
	var background [][]int
	for i := 0; i < width; i++ {
		background = append(background, []int{})
		for j := 0; j < height; j++ {
			background[i] = append(background[i], 0)
		}
	}

	p := &PlayGround{
		Snake:      snake,
		Food:       food,
		Width:      width,
		Height:     height,
		HeadDir:    MoveRight,
		TailDir:    MoveRight,
		Background: background,
	}

	p.UpdatePlayGround()
	return p
}

func (p *PlayGround) UpdatePlayGround() {
	for i := 0; i < p.Width; i++ {
		for j := 0; j < p.Height; j++ {
			p.Background[i][j] = 0
		}
	}

	for _, v := range p.Snake {
		p.Background[v[0]][v[1]] = 1
	}

	for _, v := range p.Food {
		p.Background[v[0]][v[1]] = 2
	}
}

func (p *PlayGround) Move(direction int) bool {
	p.Steps++
	//避免反向移动
	p.UpdateDirection(direction)

	x := p.Snake[len(p.Snake)-1][0]
	y := p.Snake[len(p.Snake)-1][1]
	tail := []int{x, y}

	//移动身体
	if len(p.Snake) > 1 {
		for i := len(p.Snake) - 1; i > 0; i-- {
			p.Snake[i][0] = p.Snake[i-1][0]
			p.Snake[i][1] = p.Snake[i-1][1]
		}
	}

	switch p.HeadDir {
	case MoveUp:
		p.Snake[0][0]--
	case MoveDown:
		p.Snake[0][0]++
	case MoveLeft:
		p.Snake[0][1]--
	case MoveRight:
		p.Snake[0][1]++
	}

	//是否碰撞
	if p.isCollision() {
		return false
	}
	//是否吃到食物
	if p.eatFood() {
		p.Score++
		p.Snake = append(p.Snake, tail)
		if p.Score == p.Width*p.Height-1 {
			p.UpdatePlayGround()
			return false
		}
		p.Food = [][]int{}
		p.randomFood()
	}
	//更新背景
	p.UpdatePlayGround()
	return true
}

func (p *PlayGround) UpdateDirection(d int) {
	result := map[int]int{
		MoveUp:    MoveDown,
		MoveDown:  MoveUp,
		MoveLeft:  MoveRight,
		MoveRight: MoveLeft,
	}
	hash := map[string]int{
		"-10": MoveUp,
		"10":  MoveDown,
		"01":  MoveLeft,
		"0-1": MoveRight,
	}
	if p.HeadDir != result[d] {
		p.HeadDir = d
	}
	if len(p.Snake) == 1 {
		p.TailDir = d
	} else {
		//根据最后两个节点的位置判断尾部移动方向
		n1 := p.Snake[len(p.Snake)-1]
		n2 := p.Snake[len(p.Snake)-2]
		x, y := n1[0]-n2[0], n1[1]-n2[1]
		p.TailDir = hash[strconv.Itoa(x)+strconv.Itoa(y)]
	}
}

//碰撞检测
func (p *PlayGround) isCollision() bool {
	//边界碰撞
	if p.Snake[0][0] < 0 || p.Snake[0][0] >= p.Width || p.Snake[0][1] < 0 || p.Snake[0][1] >= p.Height {
		return true
	}
	//自身碰撞
	for i := 1; i < len(p.Snake); i++ {
		if p.Snake[0][0] == p.Snake[i][0] && p.Snake[0][1] == p.Snake[i][1] {
			return true
		}
	}
	return false
}

//吃食物
func (p *PlayGround) eatFood() bool {
	if p.Snake[0][0] == p.Food[0][0] && p.Snake[0][1] == p.Food[0][1] {
		return true
	}
	return false
}

//随机生成食物
func (p *PlayGround) randomFood() {
	var blank [][]int
	for i := 0; i < p.Width; i++ {
		for j := 0; j < p.Height; j++ {
			if p.Background[i][j] == 0 {
				blank = append(blank, []int{i, j})
			}
		}
	}
	n := rand.Intn(len(blank))
	p.Food = append(p.Food, [][]int{blank[n]}...)
}

func show(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += 1
	}
}

func (p *PlayGround) Render(background [][]int) {
	w, h := termbox.Size()
	l := w/2 - p.Height
	t := h/2 - p.Width/2
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
	show(l, t-2, termbox.ColorBlue, termbox.ColorDefault, "Score: "+strconv.Itoa(p.Score))
	//显示提示
	show(l, t+p.Width+1, termbox.ColorBlue, termbox.ColorDefault, "Press ESC or Ctrl+C to exit")

	termbox.Flush()
}

func FinallyScore(score int) {
	w, h := termbox.Size()
	l := w/2 - 10
	t := h/2 - 2
	if score == width*height-1 {
		show(l, t, termbox.ColorBlue, termbox.ColorDefault, "You Win!")
		return
	} else {
		show(l, t, termbox.ColorBlue, termbox.ColorDefault, "Game Over!")
	}
	show(l, t+1, termbox.ColorWhite, termbox.ColorDefault, "Your score is: "+strconv.Itoa(score))
	termbox.Flush()
}

func MoveEvent(ev termbox.Event) int {
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

func (p *PlayGround) GetState() []float64 {
	state := make([]float64, 0)
	var headDir []float64
	var tailDir []float64
	//转成四位二进制
	for i := 0; i < 4; i++ {
		if i == p.HeadDir {
			headDir = append(headDir, 1)
		} else {
			headDir = append(headDir, 0)
		}
		if i == p.TailDir {
			tailDir = append(tailDir, 1)
		} else {
			tailDir = append(tailDir, 0)
		}
	}
	dirs := [][]int{{0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}}

	for dir := range dirs {
		x, y := p.Snake[0][0]+dirs[dir][0], p.Snake[0][1]+dirs[dir][1]
		distance := 0.0
		hasFood := 0.0
		hasBody := 0.0
		for x >= 0 && x < p.Width && y >= 0 && y < p.Height {
			distance++
			if p.Background[x][y] == 2 {
				hasFood = 1
			}
			if p.Background[x][y] == 1 {
				hasBody = 1
			}
			x += dirs[dir][0]
			y += dirs[dir][1]
		}
		//state = append(state, 1.0/(distance+1), hasFood, hasBody)
		if distance == 0 {
			state = append(state, 0, hasFood, hasBody)
		} else {
			state = append(state, 1, hasFood, hasBody)
		}
	}
	state = append(state, headDir...)
	state = append(state, tailDir...)
	return state
}

func (p *PlayGround) Print() {
	for i := 0; i < p.Width; i++ {
		for j := 0; j < p.Height; j++ {
			fmt.Print(p.Background[i][j])
		}
		fmt.Println()
	}
}
