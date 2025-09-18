package linkedlistfast

import (
	"fmt"
	"log"
	"math"

	"github.com/xupin/aoi/entity"
)

type Aoi struct {
	nodes        map[uint32]*node
	xList        *list
	yList        *list
	visibleRange uint
	step         int
}

type list struct {
	head *node
	tail *node
	size int
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
		nodes:        make(map[uint32]*node, 0),
		xList:        &list{},
		yList:        &list{},
		visibleRange: 5,
		step:         50,
	}
}

func (r *Aoi) Enter(p *entity.Player, f entity.Callback) {
	node := r.Add(p)
	log.Printf("玩家[%s]进入地图 \n", p.Name)
	players := r.findNeighbors(node, "wm")
	for _, p1 := range players {
		f(p, p1.player)
	}
}

func (r *Aoi) Move(p *entity.Player, x, y uint, move, leave, enter entity.Callback) {
	node, ok := r.nodes[p.Id]
	if !ok {
		return
	}
	p.X, p.Y = x, y
	log.Printf("玩家[%s]移动坐标 x%d,y%d -> x%d,y%d \n", p.Name, p.X, p.Y, x, y)
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
	node, ok := r.nodes[p.Id]
	if !ok {
		return
	}
	log.Printf("玩家[%s]离开地图 \n", p.Name)
	players := r.findNeighbors(node, "wm")
	for _, p1 := range players {
		f(p, p1.player)
	}
	r.Remove(node.Id)
}

func (r *Aoi) Add(player *entity.Player) *node {
	if v, ok := r.nodes[player.Id]; ok {
		return v
	}
	r.xList.size += 1
	r.yList.size += 1
	node := &node{
		Id:     player.Id,
		x:      player.X,
		y:      player.Y,
		player: player,
	}
	r.nodes[node.Id] = node
	if r.xList.head == nil || r.yList.head == nil {
		r.xList, r.yList = &list{
			head: node,
			tail: node,
		}, &list{
			head: node,
			tail: node,
		}
		return node
	}
	r.addX(node)
	r.addY(node)
	return node
}

func (r *Aoi) Remove(id uint32) {
	node, ok := r.nodes[id]
	if !ok {
		return
	}
	r.removeX(node)
	r.removeY(node)
	delete(r.nodes, id)
	r.xList.size -= 1
	r.yList.size -= 1
}

func (r *Aoi) removeX(delNode *node) {
	// 根节点
	if delNode.xPrev == nil {
		if delNode.xNext == nil {
			r.xList.head, r.xList.tail = nil, nil
		} else {
			r.xList.head = delNode.xNext
			r.xList.head.xPrev = nil
		}
	} else if delNode.xNext == nil { // 尾节点
		delNode.xPrev.xNext = nil
		r.xList.tail = delNode.xPrev
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
			r.yList.head, r.yList.tail = nil, nil
		} else {
			r.yList.head = delNode.yNext
			r.yList.head = nil
		}
	} else if delNode.yNext == nil { // 尾节点
		delNode.yPrev.yNext = nil
		r.yList.tail = delNode.yPrev
	} else if delNode.yNext != nil {
		delNode.yPrev.yNext = delNode.yNext
		delNode.yNext.yPrev = delNode.yPrev
	} else {
		delNode.yPrev.yNext = nil
	}
	delNode.yPrev, delNode.yNext = nil, nil
}

func (r *Aoi) addX(targetNode *node) {
	// 直接插入尾部
	if r.xList.tail.x <= targetNode.x {
		r.xList.tail.xNext = targetNode
		targetNode.xPrev = r.xList.tail
		r.xList.tail = targetNode
		return
	}
	var (
		prev *node
		slow = r.xList.head
		skip = int(math.Ceil(float64(r.xList.size) / float64(r.step)))
	)
	for i := 0; i < r.step; i++ {
		// 移动快指针
		fast := r.xFastMove(slow, skip)
		// 快指针的坐标小于插入坐标，移动慢指针至快指针位置
		if fast.x < targetNode.x {
			slow = fast
			continue
		}
		for ; slow != nil; slow = slow.xNext {
			// 插入cur之前
			if slow.x >= targetNode.x {
				// 根节点
				if slow.xPrev == nil {
					r.xList.head = targetNode
				} else {
					targetNode.xPrev = slow.xPrev // 当前节点的前置节点作为新节点的前置节点
					slow.xPrev.xNext = targetNode // 新节点成为当前节点的前置节点的后置节点
				}
				slow.xPrev = targetNode // 新节点作为当前节点的前置节点
				targetNode.xNext = slow // 当前节点作为新节点的后置节点
				break
			}
			prev = slow
		}
		// 插入尾部
		if prev != nil && slow == nil {
			prev.xNext = targetNode
			targetNode.xPrev = prev
			r.xList.tail = targetNode
		}
		break
	}
}

func (r *Aoi) addY(targetNode *node) {
	// 直接插入尾部
	if r.yList.tail.y <= targetNode.y {
		r.yList.tail.yNext = targetNode
		targetNode.yPrev = r.yList.tail
		r.yList.tail = targetNode
		return
	}
	var (
		prev *node
		slow = r.yList.head
		skip = int(math.Ceil(float64(r.yList.size) / float64(r.step)))
	)
	for i := 0; i < r.step; i++ {
		// 移动快指针
		fast := r.yFastMove(slow, skip)
		// 快指针的坐标小于插入坐标，移动慢指针至快指针位置
		if fast.y < targetNode.y {
			slow = fast
			continue
		}
		for ; slow != nil; slow = slow.yNext {
			// 插入cur之前
			if slow.y >= targetNode.y {
				// 根节点
				if slow.yPrev == nil {
					r.yList.head = targetNode
				} else {
					targetNode.yPrev = slow.yPrev // 当前节点的前置节点作为新节点的前置节点
					slow.yPrev.yNext = targetNode // 新节点成为当前节点的前置节点的后置节点
				}
				slow.yPrev = targetNode // 新节点作为当前节点的前置节点
				targetNode.yNext = slow // 当前节点作为新节点的后置节点
				break
			}
			prev = slow
		}
		// 插入尾部
		if prev != nil && slow == nil {
			prev.yNext = targetNode
			targetNode.yPrev = prev
			r.yList.tail = targetNode
		}
		break
	}
}

func (r *Aoi) PrintNode() {
	for cur := r.xList.head; cur != nil; cur = cur.xNext {
		if cur.xNext == nil {
			log.Print(cur.Id, "->", "nil\n")
		} else {
			log.Print(cur.Id, "->")
		}
	}
	for cur := r.yList.head; cur != nil; cur = cur.yNext {
		if cur.yNext == nil {
			log.Print(cur.Id, "->", "nil\n")
		} else {
			log.Print(cur.Id, "->")
		}
	}
}

// 感兴趣的邻居
func (r *Aoi) findNeighbors(targetNode *node, model string) map[uint32]*node {
	neighbors := make(map[uint32]*node, 0)
	// 向后找
	for cur := targetNode.xNext; cur != nil; cur = cur.xNext {
		// 当前节点已经超出范围
		if cur.x-targetNode.x > r.visibleRange {
			break
		}
		// y轴不符合
		if abs(int(cur.y-targetNode.y)) > int(r.visibleRange) {
			continue
		}
		neighbors[cur.Id] = cur
	}
	// 向前找
	for cur := targetNode.xPrev; cur != nil; cur = cur.xPrev {
		// 当前节点已经超出范围
		if targetNode.x-cur.x > r.visibleRange {
			fmt.Println(targetNode.x, cur.x)
			break
		}
		// y轴不符合
		if abs(int(targetNode.y-cur.y)) > int(r.visibleRange) {
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

func (r *Aoi) xFastMove(cur *node, skip int) *node {
	for i := 1; i <= skip; i++ {
		if cur.xNext == nil {
			break
		}
		cur = cur.xNext
	}
	return cur
}

func (r *Aoi) yFastMove(cur *node, skip int) *node {
	for i := 1; i <= skip; i++ {
		if cur.yNext == nil {
			break
		}
		cur = cur.yNext
	}
	return cur
}

func abs(n int) int {
	y := n >> 63
	return (n ^ y) - y
}
