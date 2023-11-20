package entity

type Callback = func(p1, p2 *Player)

type Player struct {
	Id      uint32
	Name    string
	X       uint
	Y       uint
	Model   string // w、m、wm （Watcher、Marker）
	Players map[uint32]*Player
}
