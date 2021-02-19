package main

const (
	PLAYER_BLACK = 'B'
	PLAYER_WHITE = 'W'
)

type Chess struct {
	Turn   byte
	Pieces map[int]byte
}

func (c *Chess) Init() {
	c.Turn = PLAYER_BLACK
	c.Pieces = make(map[int]byte)
}

func (c *Chess) GetPos(pos int) (x, y int) {
	x = pos/19 + 1
	y = pos%19 + 1
	return
}

func (c *Chess) posToInt(x, y int) int {
	return (x-1)*19 + y - 1
}

func (c *Chess) PlacePiece(x, y int) string {
	pos := c.posToInt(x, y)
	if _, ok := c.Pieces[pos]; ok {
		return "Already Exist"
	}
	c.Pieces[pos] = c.Turn
	if c.Win(pos) {
		return "Win"
	}
	if c.Turn == PLAYER_BLACK {
		c.Turn = PLAYER_WHITE
	} else if c.Turn == PLAYER_WHITE {
		c.Turn = PLAYER_BLACK
	}
	return "OK"
}

func (c *Chess) Reset() {
	c.Init()
}

func (c *Chess) Win(pos int) bool {
	offsets := [4]int{1, 19, 19 - 1, 19 + 1}
	for index := 0; index < 4; index++ {
		offset := offsets[index]
		backStop := false
		frontStop := false
		backNum := 0
		frontNum := 0
		for num := 1; num < 5; num++ {
			if !backStop {
				turn, ok := c.Pieces[pos-num*offset]
				if ok && turn == c.Turn {
					backNum++
				} else {
					backStop = true
				}
			}
			if !frontStop {
				turn, ok := c.Pieces[pos+num*offset]
				if ok && turn == c.Turn {
					frontNum++
				} else {
					frontStop = true
				}
			}
			if frontStop && backStop {
				break
			}
		}
		if backNum+frontNum+1 >= 5 {
			return true
		}
	}
	return false
}

func (c *Chess) Handle(msg string) {
	if msg == "Start" {
		c.Init()
		return
	} else if msg == "End" {
		return
	}

}
