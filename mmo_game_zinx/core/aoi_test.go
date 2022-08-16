package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)
	fmt.Println(aoiMgr)
}

func TestAOIManagerSurroundGridsByGrid(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)
	fmt.Println(aoiMgr)
	for gid := range aoiMgr.grids {
		//得到gid周围九宫格的信息
		grids := aoiMgr.GetSurroundGridsByGid(gid)
		fmt.Println("gid:", gid, "grids len:", len(grids))
		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.GID)
		}
		fmt.Println("surrounding grid:", gIDs)
	}
}
