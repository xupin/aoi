package tower

import (
	"fmt"
	"math"
	"strconv"

	"github.com/xupin/aoi/entity"
)

type Tower struct {
	Id       int
	Name     string
	X        uint
	Y        uint
	Watchers map[uint32]*entity.Player
	Markers  map[uint32]*entity.Player
}

type Aoi struct {
	Towers       map[uint]map[uint]*Tower
	TowerWidth   uint
	TowerHeight  uint
	VisibleRange uint
}

const (
	// 地图尺寸
	MAP_ROWS = 50 // y
	MAP_COLS = 50 // x
)

func NewAoi() *Aoi {
	return &Aoi{
		Towers:       make(map[uint]map[uint]*Tower, 0),
		TowerWidth:   5,
		TowerHeight:  5,
		VisibleRange: 5,
	}
}

func (r *Aoi) Start() {
	// 计算灯塔数量
	maxX := uint(math.Ceil(float64(MAP_COLS) / float64(r.TowerWidth)))
	maxY := uint(math.Ceil(float64(MAP_ROWS) / float64(r.TowerHeight)))
	// 生成灯塔
	id := 1
	for x := uint(0); x < maxX; x++ {
		r.Towers[x] = make(map[uint]*Tower, maxY)
		for y := uint(0); y < maxY; y++ {
			name := "灯塔" + strconv.Itoa(id)
			r.Towers[x][y] = &Tower{
				Id:       id,
				Name:     name,
				X:        x,
				Y:        y,
				Watchers: make(map[uint32]*entity.Player, 0),
				Markers:  make(map[uint32]*entity.Player, 0),
			}
			// fmt.Printf("灯塔%d[%d,%d]加入 \n", id, x, y)
			id++
		}
	}
}

func (r *Aoi) Enter(p *entity.Player, f entity.Callback) {
	// 加入灯塔
	tower := r.getTower(p.X, p.Y)
	tower.Markers[p.Id] = p
	fmt.Printf("玩家[%s]进入地图 \n", p.Name)
	// 获取视野内的灯塔、被观察者
	towers := r.getWatchedTowers(p.X, p.Y)
	for _, tower := range towers {
		for id, p1 := range tower.Watchers {
			if p.Id == p1.Id {
				continue
			}
			if _, ok := p.Players[id]; !ok {
				// 互相加入可视列表
				p.Players[id] = p1
				p1.Players[p.Id] = p
				// 回调
				f(p, p1)
			}
		}
		tower.addWatcher(p)
	}
}

func (r *Aoi) Leave(p *entity.Player, f entity.Callback) {
	// 离开灯塔
	tower := r.getTower(p.X, p.Y)
	tower.removeMarker(p)
	// 获取视野内的灯塔、被观察者
	towers := r.getWatchedTowers(p.X, p.Y)
	for _, tower := range towers {
		// 离开灯塔
		tower.removeWatcher(p)
	}
	for id, p1 := range p.Players {
		// 互相离开可视列表
		delete(p.Players, id)
		delete(p1.Players, p.Id)
		// 回调
		f(p, p1)
	}
	fmt.Printf("玩家[%s]离开地图 \n", p.Name)
}

func (r *Aoi) Move(p *entity.Player, x, y uint, move, leave, enter entity.Callback) {
	fmt.Printf("玩家[%s]移动坐标 x%d,y%d -> x%d,y%d \n", p.Name, p.X, p.Y, x, y)
	// 离开、加入新的灯塔
	bTower := r.getTower(p.X, p.Y)
	aTower := r.getTower(x, y)
	if bTower.Id != aTower.Id {
		bTower.removeMarker(p)
		aTower.addMarker(p)
	}
	// 移动，判断视野内的灯塔有无变化
	bTowers := r.getWatchedTowers(p.X, p.Y)
	aTowers := r.getWatchedTowers(x, y)
	if r.TowersEqual(bTowers, aTowers) {
		for _, p1 := range p.Players {
			move(p, p1)
		}
		return
	}
	// 需要离开的灯塔
	for _, tower := range r.TowersDiff(bTowers, aTowers) {
		tower.removeWatcher(p)
	}
	// 离开玩家视野
	bPlayers := r.getWatchers(p.X, p.Y)
	aPlayers := r.getWatchers(x, y)
	for _, p1 := range r.PlayersDiff(bPlayers, aPlayers) {
		delete(p1.Players, p.Id)
		if _, ok := p.Players[p1.Id]; ok {
			delete(p.Players, p1.Id)
			leave(p, p1)
		}
	}
	// 新加入的灯塔
	for _, tower := range r.TowersDiff(aTowers, bTowers) {
		tower.addWatcher(p)
		// 新的观察者
		for _, p1 := range tower.Watchers {
			if p.Id == p1.Id {
				continue
			}
			if _, ok := p.Players[p1.Id]; !ok {
				p.Players[p1.Id] = p1
				p1.Players[p.Id] = p
				enter(p, p1)
			}
		}
	}
}

func (r *Tower) addWatcher(p *entity.Player) {
	r.Watchers[p.Id] = p
	fmt.Printf("玩家[%s]关注灯塔[%d,%d] \n", p.Name, r.X, r.Y)
}

func (r *Tower) removeWatcher(p *entity.Player) {
	delete(r.Watchers, p.Id)
	fmt.Printf("玩家[%s]不再关注灯塔[%d,%d] \n", p.Name, r.X, r.Y)
}

func (r *Tower) addMarker(p *entity.Player) {
	r.Markers[p.Id] = p
	fmt.Printf("玩家[%s]加入灯塔[%d,%d] \n", p.Name, r.X, r.Y)
}

func (r *Tower) removeMarker(p *entity.Player) {
	delete(r.Markers, p.Id)
	fmt.Printf("玩家[%s]离开灯塔[%d,%d] \n", p.Name, r.X, r.Y)
}

func (r *Aoi) TowersEqual(bTowers, aTowers []*Tower) bool {
	if len(bTowers) != len(aTowers) {
		return false
	}
	for k := range bTowers {
		if bTowers[k].Id != aTowers[k].Id {
			return false
		}
	}
	return true
}

func (r *Aoi) TowersDiff(bTowers, aTowers []*Tower) []*Tower {
	towers := make([]*Tower, 0)
	inTowers := make(map[int]bool, 0)
	for _, tower := range aTowers {
		inTowers[tower.Id] = true
	}
	for _, tower := range bTowers {
		if _, ok := inTowers[tower.Id]; ok {
			continue
		}
		towers = append(towers, tower)
	}
	return towers
}

func (r *Aoi) PlayersDiff(bPlayers, aPlayers []*entity.Player) []*entity.Player {
	players := make([]*entity.Player, 0)
	newPlayers := make(map[uint32]bool, 0)
	for _, player := range aPlayers {
		newPlayers[player.Id] = true
	}
	for _, player := range bPlayers {
		if _, ok := newPlayers[player.Id]; ok {
			continue
		}
		players = append(players, player)
	}
	return players
}

func (r *Aoi) getTower(x, y uint) *Tower {
	x, y = r.transPos(x, y)
	tower, ok := r.Towers[x][y]
	if !ok {
		fmt.Printf("灯塔[异常]不存在的灯塔: %d,%d \n", x, y)
		return nil
	}
	return tower
}

func (r *Aoi) getWatchedTowers(x, y uint) []*Tower {
	xMin, xMax := int64(x-r.VisibleRange), int64(x+r.VisibleRange)
	if xMin < 0 {
		xMin = 0
	}
	if xMax > MAP_COLS {
		xMax = MAP_COLS
	}
	yMin, yMax := int64(y-r.VisibleRange), int64(y+r.VisibleRange)
	if yMin < 0 {
		yMin = 0
	}
	if yMax > MAP_ROWS {
		yMax = MAP_ROWS
	}
	towers := make([]*Tower, 0)
	for x := uint(xMin); x < uint(xMax); x += r.TowerWidth {
		for y := uint(yMin); y < uint(yMax); y += r.TowerHeight {
			tower := r.getTower(x, y)
			if tower == nil {
				continue
			}
			towers = append(towers, tower)
		}
	}
	return towers
}

func (r *Aoi) getWatchers(x, y uint) []*entity.Player {
	towers := r.getWatchedTowers(x, y)
	players := make([]*entity.Player, 0)
	for _, tower := range towers {
		for _, player := range tower.Watchers {
			players = append(players, player)
		}
	}
	return players
}

func (r *Aoi) transPos(x, y uint) (uint, uint) {
	x = uint(math.Floor(float64(x) / float64(r.TowerWidth)))
	y = uint(math.Floor(float64(y) / float64(r.TowerHeight)))
	return x, y
}
