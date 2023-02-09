package main

import (
	"Snake-go/game"
	"Snake-go/search"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

var (
	playGround *game.PlayGround
)

func init() {
	rand.Seed(3141592653589793238)
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	playGround = game.NewPlayGround()
	playGround.Render(playGround.Background)
}

func main() {
	t := search.Train(100, 50)
	nn := search.NewNeuralNetwork()
	nn.UpdateWeights(t.Genes)
	for {
		state := playGround.GetState()
		steps := nn.Predict(state)
		if playGround.Move(steps) {
			playGround.Render(playGround.Background)
		} else {
			game.FinallyScore(playGround.Score)
			time.Sleep(3 * time.Second)
			termbox.Close()
			return
		}
		time.Sleep(time.Second / 10)
	}
}

func testGame() {
	for {
		if ev := termbox.PollEvent(); ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
				termbox.Close()
				return
			}
			s := game.MoveEvent(ev)
			if s != -1 {
				if playGround.Move(s) {
					playGround.Render(playGround.Background)
				} else {
					game.FinallyScore(playGround.Score)
					time.Sleep(3 * time.Second)
					termbox.Close()
					return
				}
			}
		}
	}
}

func testAStart() {
	for {
		key := termbox.PollEvent().Key
		if key == termbox.KeyEsc || key == termbox.KeyCtrlC {
			termbox.Close()
			return
		}
		steps := search.AStart(playGround.Snake, playGround.Food)
		for _, v := range steps {
			if playGround.Move(v) {
				playGround.Render(playGround.Background)
			} else {
				game.FinallyScore(playGround.Score)
				time.Sleep(3 * time.Second)
				termbox.Close()
				return
			}
			time.Sleep(time.Second / 10)
		}
	}
}

func testSimple() {
	for {
		steps := search.Simple(playGround.Snake, playGround.Food)
		for _, v := range steps {
			if playGround.Move(v) {
				playGround.Render(playGround.Background)
			} else {
				game.FinallyScore(playGround.Score)
				time.Sleep(3 * time.Second)
				termbox.Close()
				return
			}
			time.Sleep(time.Second / 10)
		}
	}
}

func testCircle() {
	for {
		key := termbox.PollEvent().Key
		if key == termbox.KeyEsc || key == termbox.KeyCtrlC {
			termbox.Close()
			return
		}
		steps := search.Circle(playGround.Snake)
		for _, v := range steps {
			if playGround.Move(v) {
				playGround.Render(playGround.Background)
			} else {
				game.FinallyScore(playGround.Score)
				time.Sleep(3 * time.Second)
				termbox.Close()
				return
			}
			time.Sleep(time.Second / 10)
		}
	}
}

func testGetState() {
	for {
		ev := termbox.PollEvent()
		if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
			return
		}
		playGround.GetState()
		t := game.MoveEvent(ev)
		if t != -1 {
			playGround.Move(t)
			playGround.Render(playGround.Background)
		}
	}
}
