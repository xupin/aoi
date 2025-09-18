package grid

import (
	"log"
	"strings"

	"github.com/xupin/aoi/entity"
)

type Aoi struct {
	Players      map[uint32]*entity.Player
	PlayersX     map[uint]map[uint32]*entity.Player
	PlayersY     map[uint]map[uint32]*entity.Player
	VisibleRange uint
}

const (
	AOI_WATCHER = "w"
	AOI_MARKER  = "m"
)

const (
	// 地图尺寸
	MAP_ROWS = 20 // y
	MAP_COLS = 20 // x
)

func NewAoi() *Aoi {
	return &Aoi{
		Players:      make(map[uint32]*entity.Player, 0),
		PlayersX:     make(map[uint]map[uint32]*entity.Player, 0),
		PlayersY:     make(map[uint]map[uint32]*entity.Player, 0),
		VisibleRange: 5,
	}
}

func (r *Aoi) Enter(p *entity.Player, f entity.Callback) map[uint32]*entity.Player {
	r.Players[p.Id] = p
	if _, ok := r.PlayersX[p.X]; !ok {
		r.PlayersX[p.X] = make(map[uint32]*entity.Player)
	}
	r.PlayersX[p.X][p.Id] = p
	if _, ok := r.PlayersY[p.Y]; !ok {
		r.PlayersY[p.Y] = make(map[uint32]*entity.Player)
	}
	r.PlayersY[p.Y][p.Id] = p
	log.Printf("玩家[%s]进入地图 x%d,y%d \n", p.Name, p.X, p.Y)
	// 如果玩家是被观察者，广播消息给视野内所有观察者
	if r.IsMarker(p) {
		r.Broadcast(p, f)
	}
	// 如果玩家是观察者，广播消息给视野内所有被观察者
	if r.IsWatcher(p) {
		return r.findNeighbors(p, AOI_MARKER)
	}
	// log.Printf("内存[aoi] %d", unsafe.Sizeof(r))
	return nil
}

func (r *Aoi) Move(p *entity.Player, x, y uint, move, leave, enter entity.Callback) []*entity.Player {
	log.Printf("玩家[%s]移动坐标 x%d,y%d ->  x%d,y%d \n", p.Name, p.X, p.Y, x, y)
	// 获取当前坐标视野内的观察者、被观察者
	bWatchers := r.findNeighbors(p, AOI_WATCHER)
	bMarkers := r.findNeighbors(p, AOI_MARKER)
	// 移动
	p.X, p.Y = x, y
	// 获取移动后坐标视野内的观察者、被观察者
	aWatchers := r.findNeighbors(p, AOI_WATCHER)
	aMarkers := r.findNeighbors(p, AOI_MARKER)
	//
	if r.IsMarker(p) {
		// 离开对方视野
		for id, p1 := range bWatchers {
			if _, ok := aWatchers[id]; ok {
				continue
			}
			leave(p, p1)
		}
		// 进入对方视野
		for id, p1 := range aWatchers {
			if _, ok := bWatchers[id]; ok {
				move(p, p1)
			} else {
				enter(p, p1)
			}
		}
	}
	// 新的视野邻居
	players := []*entity.Player{}
	if r.IsWatcher(p) {
		for id := range aMarkers {
			if p1, ok := bMarkers[id]; !ok {
				players = append(players, p1)
			}
		}
	}
	return players
}

func (r *Aoi) Leave(p *entity.Player, f entity.Callback) {
	delete(r.Players, p.Id)
	delete(r.PlayersX[p.X], p.Id)
	delete(r.PlayersY[p.Y], p.Id)
	// 如果玩家是被观察者，广播消息给视野内所有观察者
	if r.IsMarker(p) {
		r.Broadcast(p, f)
	}
}

func (r *Aoi) Broadcast(p *entity.Player, f entity.Callback) {
	players := r.findNeighbors(p, AOI_MARKER)
	for _, p1 := range players {
		f(p, p1)
	}
}

func (r *Aoi) Get(id uint32) *entity.Player {
	return r.Players[id]
}

func (r *Aoi) IsMarker(p *entity.Player) bool {
	return strings.Contains(p.Model, "m")
}

func (r *Aoi) IsWatcher(p *entity.Player) bool {
	return strings.Contains(p.Model, "w")
}

func (r *Aoi) findNeighbors(p *entity.Player, model string) map[uint32]*entity.Player {
	// 地图边界
	xMin := int64(p.X - r.VisibleRange)
	if xMin < 0 {
		xMin = 0
	}
	xMax := p.X + r.VisibleRange
	if xMax > MAP_ROWS {
		xMax = MAP_ROWS
	}
	// 感兴趣的邻居
	neighbors := make(map[uint32]*entity.Player, 0)
	for x := uint(xMin); x < uint(xMax); x++ {
		players, ok := r.PlayersX[x]
		if !ok {
			continue
		}
		for _, p1 := range players {
			if p1.Id == p.Id {
				continue
			}
			// 判断玩家aoi模型
			if model == "w" {
				if !r.IsWatcher(p1) {
					continue
				}
			} else {
				if !r.IsMarker(p1) {
					continue
				}
			}
			// 判断Y轴是否在视野内
			if abs(int(p.Y-p1.Y)) > int(r.VisibleRange) {
				continue
			}
			neighbors[p1.Id] = p1
		}
	}
	return neighbors
}

func abs(n int) int {
	y := n >> 63
	return (n ^ y) - y
}
