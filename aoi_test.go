package aoi_test

import (
	"testing"

	"github.com/xupin/aoi/entity"
	"github.com/xupin/aoi/grid"
	"github.com/xupin/aoi/linkedlist"
	"github.com/xupin/aoi/tower"
)

func TestLinkedListAOI(t *testing.T) {
	aoi := linkedlist.NewAoi()

	p1 := &entity.Player{Id: 1, Name: "pp", X: 0, Y: 0, Model: "wm"}
	p2 := &entity.Player{Id: 2, Name: "wl", X: 2, Y: 20, Model: "wm"}
	p3 := &entity.Player{Id: 3, Name: "sd", X: 0, Y: 0, Model: "wm"}
	p4 := &entity.Player{Id: 4, Name: "pp4", X: 1, Y: 21, Model: "wm"}

	aoi.Enter(p3, func(p1, p2 *entity.Player) {
		t.Logf("玩家[%s]遇见玩家[%s]", p1.Name, p2.Name)
	})
	aoi.Enter(p4, func(p1, p2 *entity.Player) {
		t.Logf("玩家[%s]遇见玩家[%s]", p1.Name, p2.Name)
	})
	aoi.Enter(p2, func(p1, p2 *entity.Player) {
		t.Logf("玩家[%s]遇见玩家[%s]", p1.Name, p2.Name)
	})
	aoi.Enter(p1, func(p1, p2 *entity.Player) {
		t.Logf("玩家[%s]遇见玩家[%s]", p1.Name, p2.Name)
	})

	aoi.Move(p4, 4, 4,
		func(p1, p2 *entity.Player) {
			t.Logf("玩家[%s]移动坐标，通知玩家[%s]", p1.Name, p2.Name)
		},
		func(p1, p2 *entity.Player) {
			t.Logf("玩家[%s]离开视野，通知玩家[%s]", p1.Name, p2.Name)
		},
		func(p1, p2 *entity.Player) {
			t.Logf("玩家[%s]进入视野，通知玩家[%s]", p1.Name, p2.Name)
		},
	)

	aoi.PrintNode()
}

func TestGridAOI(t *testing.T) {
	aoi := grid.NewAoi()

	p1 := &entity.Player{Id: 1, Name: "pp", X: 0, Y: 0, Model: "wm"}
	p2 := &entity.Player{Id: 2, Name: "wl", X: 2, Y: 20, Model: "wm"}
	p3 := &entity.Player{Id: 3, Name: "sd", X: 0, Y: 0, Model: "wm"}

	aoi.Enter(p1, func(p1, p2 *entity.Player) { t.Logf("玩家[%s]遇见玩家[%s]", p1.Name, p2.Name) })
	aoi.Enter(p2, func(p1, p2 *entity.Player) { t.Logf("玩家[%s]遇见玩家[%s]", p1.Name, p2.Name) })
	aoi.Enter(p3, func(p1, p2 *entity.Player) { t.Logf("玩家[%s]遇见玩家[%s]", p1.Name, p2.Name) })

	aoi.Move(p2, 2, 10,
		func(p1, p2 *entity.Player) { t.Logf("玩家[%s]移动坐标，通知玩家[%s]", p1.Name, p2.Name) },
		func(p1, p2 *entity.Player) { t.Logf("玩家[%s]离开视野，通知玩家[%s]", p1.Name, p2.Name) },
		func(p1, p2 *entity.Player) { t.Logf("玩家[%s]进入视野，通知玩家[%s]", p1.Name, p2.Name) },
	)

	aoi.Leave(p3, func(p1, p2 *entity.Player) { t.Logf("玩家[%s]离开视野，通知玩家[%s]", p1.Name, p2.Name) })

	aoi.Move(p2, 20, 10,
		func(p1, p2 *entity.Player) { t.Logf("玩家[%s]移动坐标，通知玩家[%s]", p1.Name, p2.Name) },
		func(p1, p2 *entity.Player) { t.Logf("玩家[%s]离开视野，通知玩家[%s]", p1.Name, p2.Name) },
		func(p1, p2 *entity.Player) { t.Logf("玩家[%s]进入视野，通知玩家[%s]", p1.Name, p2.Name) },
	)

	aoi.Leave(p2, func(p1, p2 *entity.Player) { t.Logf("玩家[%s]离开地图，通知玩家[%s]", p1.Name, p2.Name) })
}

func TestTowerAOI(t *testing.T) {
	aoi := tower.NewAoi()
	aoi.Start()

	p1 := &entity.Player{Id: 1, Name: "pp", X: 49, Y: 49, Players: make(map[uint32]*entity.Player)}
	p2 := &entity.Player{Id: 2, Name: "wl", X: 8, Y: 8, Players: make(map[uint32]*entity.Player)}
	p3 := &entity.Player{Id: 3, Name: "sd", X: 0, Y: 0, Players: make(map[uint32]*entity.Player)}

	aoi.Enter(p1, func(p1, p2 *entity.Player) { t.Logf("玩家[%s]遇见玩家[%s]", p1.Name, p2.Name) })
	aoi.Enter(p2, func(p1, p2 *entity.Player) { t.Logf("玩家[%s]遇见玩家[%s]", p1.Name, p2.Name) })
	aoi.Enter(p3, func(p1, p2 *entity.Player) { t.Logf("玩家[%s]遇见玩家[%s]", p1.Name, p2.Name) })

	aoi.Leave(p1, func(p1, p2 *entity.Player) { t.Logf("玩家[%s]离开，通知玩家[%s]", p1.Name, p2.Name) })

	aoi.Move(p2, 9, 9,
		func(p1, p2 *entity.Player) { t.Logf("玩家[%s]移动视野，通知玩家[%s]", p1.Name, p2.Name) },
		func(p1, p2 *entity.Player) { t.Logf("玩家[%s]离开视野，通知玩家[%s]", p1.Name, p2.Name) },
		func(p1, p2 *entity.Player) { t.Logf("玩家[%s]进入视野，通知玩家[%s]", p1.Name, p2.Name) },
	)

	aoi.Leave(p3, func(p1, p2 *entity.Player) { t.Logf("玩家[%s]离开，通知玩家[%s]", p1.Name, p2.Name) })
}
