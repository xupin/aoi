package linkedlist

import (
	"fmt"

	"github.com/xupin/aoi/entity"
)

type Aoi struct {
	list         map[uint32]*node
	xList        *node
	yList        *node
	VisibleRange uint
}

type node struct {
	Id     uint32
	xPrev  *node
	xNext  *node
	yPrev  *node
	yNext  *node
	x, y   uint
	player *entity.Player
}

func NewAoi() *Aoi {
	return &Aoi{
		list:         make(map[uint32]*node, 0),
		VisibleRange: 5,
	}
}

func (r *Aoi) Enter(p *entity.Player, f entity.Callback) {
	node := r.Add(p)
	fmt.Printf("玩家[%s]进入地图 \n", p.Name)
	players := r.findNeighbors(node, "wm")
	for _, p1 := range players {
		f(p, p1.player)
	}
}

func (r *Aoi) Move(p *entity.Player, x, y uint, move, leave, enter entity.Callback) {
	node, ok := r.list[p.Id]
	if !ok {
		return
	}
	p.X, p.Y = x, y
	fmt.Printf("玩家[%s]移动坐标 x%d,y%d -> x%d,y%d \n", p.Name, p.X, p.Y, x, y)
	// 离开玩家视野
	bPlayers := r.findNeighbors(node, "wm")
	r.Remove(node.Id)
	node = r.Add(p)
	aPlayers := r.findNeighbors(node, "wm")
	for _, p1 := range bPlayers {
		if _, ok := aPlayers[p1.Id]; ok {
			continue
		}
		leave(p, p1.player)
	}
	for _, p1 := range aPlayers {
		if _, ok := bPlayers[p1.Id]; ok {
			move(p, p1.player)
		} else {
			enter(p, p1.player)
		}
	}
}

func (r *Aoi) Leave(p *entity.Player, f entity.Callback) {
	node, ok := r.list[p.Id]
	if !ok {
		return
	}
	fmt.Printf("玩家[%s]离开地图 \n", p.Name)
	players := r.findNeighbors(node, "wm")
	for _, p1 := range players {
		f(p, p1.player)
	}
	r.Remove(node.Id)
}

func (r *Aoi) Add(player *entity.Player) *node {
	if v, ok := r.list[player.Id]; ok {
		return v
	}
	node := &node{
		Id:     player.Id,
		x:      player.X,
		y:      player.Y,
		player: player,
	}
	r.list[node.Id] = node
	if r.xList == nil || r.yList == nil {
		r.xList, r.yList = node, node
		return node
	}
	r.addX(node)
	r.addY(node)
	return node
}

func (r *Aoi) Remove(id uint32) {
	node, ok := r.list[id]
	if !ok {
		return
	}
	r.removeX(node)
	r.removeY(node)
	delete(r.list, id)
}

func (r *Aoi) removeX(delNode *node) {
	// 根节点
	if delNode.xPrev == nil {
		if delNode.xNext == nil {
			r.xList = nil
		} else {
			r.xList = delNode.xNext
			r.xList.xPrev = nil
		}
	} else if delNode.xNext != nil {
		delNode.xPrev.xNext = delNode.xNext
		delNode.xNext.xPrev = delNode.xPrev
	} else {
		delNode.xPrev.xNext = nil
	}
	delNode.xPrev, delNode.xNext = nil, nil
}

func (r *Aoi) removeY(delNode *node) {
	// 根节点
	if delNode.yPrev == nil {
		if delNode.yNext == nil {
			r.yList = nil
		} else {
			r.yList = delNode.yNext
			r.yList.yPrev = nil
		}
	} else if delNode.yNext != nil {
		delNode.yPrev.yNext = delNode.yNext
		delNode.yNext.yPrev = delNode.yPrev
	} else {
		delNode.yPrev.yNext = nil
	}
	delNode.yPrev, delNode.yNext = nil, nil
}

func (r *Aoi) addX(newNode *node) {
	var (
		prev *node
		cur  *node
	)
	for cur = r.xList; cur != nil; cur = cur.xNext {
		// 插入cur之前
		if cur.x > newNode.x {
			// 根节点
			if cur.xPrev == nil {
				r.xList = newNode
			} else {
				newNode.xPrev = cur.xPrev // 当前节点的前置节点作为新节点的前置节点
				cur.xPrev.xNext = newNode // 新节点成为当前节点的前置节点的后置节点
			}
			cur.xPrev = newNode // 新节点作为当前节点的前置节点
			newNode.xNext = cur // 当前节点作为新节点的后置节点
			break
		}
		prev = cur
	}
	// 插入尾部
	if prev != nil && cur == nil {
		prev.xNext = newNode
		newNode.xPrev = prev
	}
}

func (r *Aoi) addY(newNode *node) {
	var (
		prev *node
		cur  *node
	)
	for cur = r.yList; cur != nil; cur = cur.yNext {
		// 插入cur之前
		if cur.y > newNode.y {
			// 根节点
			if cur.yPrev == nil {
				r.yList = newNode
			} else {
				newNode.yPrev = cur.yPrev // 当前节点的前置节点作为新节点的前置节点
				cur.yPrev.yNext = newNode // 新节点成为当前节点的前置节点的后置节点
			}
			cur.yPrev = newNode // 新节点作为当前节点的前置节点
			newNode.yNext = cur // 当前节点作为新节点的后置节点
			break
		}
		prev = cur
	}
	// 插入尾部
	if prev != nil && cur == nil {
		prev.yNext = newNode
		newNode.yPrev = prev
	}
}

func (r *Aoi) PrintNode() {
	for cur := r.xList; cur != nil; cur = cur.xNext {
		if cur.xNext == nil {
			fmt.Print(cur.Id, "->", "nil\n")
		} else {
			fmt.Print(cur.Id, "->")
		}
	}
	for cur := r.yList; cur != nil; cur = cur.yNext {
		if cur.yNext == nil {
			fmt.Print(cur.Id, "->", "nil\n")
		} else {
			fmt.Print(cur.Id, "->")
		}
	}
}

func (r *Aoi) findNeighbors(targetNode *node, model string) map[uint32]*node {
	// 感兴趣的邻居
	neighbors := make(map[uint32]*node, 0)
	// 向后找
	for cur := targetNode.xNext; cur != nil; cur = cur.xNext {
		// 当前节点已经超出范围
		if cur.x-targetNode.x > r.VisibleRange {
			break
		}
		// y轴不符合
		if abs(int(cur.y-targetNode.y)) > int(r.VisibleRange) {
			continue
		}
		neighbors[cur.Id] = cur
	}
	// 向前找
	for cur := targetNode.xPrev; cur != nil; cur = cur.xPrev {
		// 当前节点已经超出范围
		if targetNode.x-cur.x > r.VisibleRange {
			fmt.Println(targetNode.x, cur.x)
			break
		}
		// y轴不符合
		if abs(int(targetNode.y-cur.y)) > int(r.VisibleRange) {
			continue
		}
		// 已存在
		if _, ok := neighbors[cur.Id]; ok {
			continue
		}
		neighbors[cur.Id] = cur
	}
	return neighbors
}

func abs(n int) int {
	y := n >> 63
	return (n ^ y) - y
}
