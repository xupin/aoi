package main

import (
	"fmt"

	"github.com/xupin/aoi/grid"
	"github.com/xupin/aoi/tower"
)

func main() {
	gridTest()
	// towerTest()
}

func gridTest() {
	aoi := &grid.Aoi{
		Players:      make(map[uint64]*grid.Player),
		PlayersX:     make(map[uint]map[uint64]*grid.Player),
		PlayersY:     make(map[uint]map[uint64]*grid.Player),
		VisibleRange: 5,
	}
	p1 := &grid.Player{
		Id:    1,
		Name:  "pp",
		X:     0,
		Y:     0,
		Model: "wm",
	}
	p2 := &grid.Player{
		Id:    2,
		Name:  "wl",
		X:     2,
		Y:     20,
		Model: "wm",
	}
	p3 := &grid.Player{
		Id:    3,
		Name:  "sd",
		X:     0,
		Y:     0,
		Model: "wm",
	}
	aoi.Enter(p1, func(p1, p2 *grid.Player) {
		fmt.Printf("玩家[%s]遇见玩家[%s] \n", p1.Name, p2.Name)
	})
	aoi.Enter(p2, func(p1, p2 *grid.Player) {
		fmt.Printf("玩家[%s]遇见玩家[%s] \n", p1.Name, p2.Name)
	})
	aoi.Enter(p3, func(p1, p2 *grid.Player) {
		fmt.Printf("玩家[%s]遇见玩家[%s] \n", p1.Name, p2.Name)
	})
	aoi.Move(p2, 2, 10,
		func(p1, p2 *grid.Player) {
			fmt.Printf("玩家[%s]移动坐标，通知玩家[%s] \n", p1.Name, p2.Name)
		},
		func(p1, p2 *grid.Player) {
			fmt.Printf("玩家[%s]离开视野，通知玩家[%s] \n", p1.Name, p2.Name)
		},
		func(p1, p2 *grid.Player) {
			fmt.Printf("玩家[%s]进入视野，通知玩家[%s] \n", p1.Name, p2.Name)
		},
	)
	aoi.Leave(p3, func(p1, p2 *grid.Player) {
		fmt.Printf("玩家[%s]离开视野，通知玩家[%s] \n", p1.Name, p2.Name)
	})
	aoi.Move(p2, 20, 10,
		func(p1, p2 *grid.Player) {
			fmt.Printf("玩家[%s]移动坐标，通知玩家[%s] \n", p1.Name, p2.Name)
		},
		func(p1, p2 *grid.Player) {
			fmt.Printf("玩家[%s]离开视野，通知玩家[%s] \n", p1.Name, p2.Name)
		},
		func(p1, p2 *grid.Player) {
			fmt.Printf("玩家[%s]进入视野，通知玩家[%s] \n", p1.Name, p2.Name)
		},
	)
	aoi.Leave(p2, func(p1, p2 *grid.Player) {
		fmt.Printf("玩家[%s]离开地图，通知玩家[%s] \n", p1.Name, p2.Name)
	})
}

func towerTest() {
	aoi := &tower.Aoi{
		Towers:       make(map[uint]map[uint]*tower.Tower, 0),
		TowerWidth:   5,
		TowerHeight:  5,
		VisibleRange: 5,
	}
	aoi.Init()
	p1 := &tower.Player{
		Id:      1,
		Name:    "pp",
		X:       49,
		Y:       49,
		Players: make(map[uint64]*tower.Player, 0),
	}
	p2 := &tower.Player{
		Id:      2,
		Name:    "wl",
		X:       8,
		Y:       8,
		Players: make(map[uint64]*tower.Player, 0),
	}
	p3 := &tower.Player{
		Id:      3,
		Name:    "sd",
		X:       0,
		Y:       0,
		Players: make(map[uint64]*tower.Player, 0),
	}
	aoi.Enter(p1, func(p1, p2 *tower.Player) {
		fmt.Printf("玩家[%s]遇见玩家[%s] \n", p1.Name, p2.Name)
	})
	aoi.Enter(p2, func(p1, p2 *tower.Player) {
		fmt.Printf("玩家[%s]遇见玩家[%s] \n", p1.Name, p2.Name)
	})
	aoi.Enter(p3, func(p1, p2 *tower.Player) {
		fmt.Printf("玩家[%s]遇见玩家[%s] \n", p1.Name, p2.Name)
	})
	aoi.Leave(p1, func(p1, p2 *tower.Player) {
		fmt.Printf("玩家[%s]离开，通知玩家[%s] \n", p1.Name, p2.Name)
	})
	aoi.Move(p2, 9, 9, func(p1, p2 *tower.Player) {
		fmt.Printf("玩家[%s]移动视野，通知玩家[%s] \n", p1.Name, p2.Name)
	}, func(p1, p2 *tower.Player) {
		fmt.Printf("玩家[%s]离开视野，通知玩家[%s] \n", p1.Name, p2.Name)
	}, func(p1, p2 *tower.Player) {
		fmt.Printf("玩家[%s]进入视野，通知玩家[%s] \n", p1.Name, p2.Name)
	})
	aoi.Leave(p3, func(p1, p2 *tower.Player) {
		fmt.Printf("玩家[%s]离开，通知玩家[%s] \n", p1.Name, p2.Name)
	})
}
