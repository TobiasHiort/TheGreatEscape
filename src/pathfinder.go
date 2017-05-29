package main

import (
	"math"
	"sync"
	//	"fmt"
)

type Direction struct {
	xDir int //-1,0,1
	yDir int //-1,0,1
}

type jp struct {
	jp *tile   // jp
	fn []*tile // forced neighbors from jp
}

var ( // TODO: define this a weak ago...
	n = Direction{-1, 0}
	e = Direction{0, 1}
	s = Direction{1, 0}
	w = Direction{0, -1}

	nw = Direction{-1, -1}
	ne = Direction{-1, 1}
	se = Direction{1, 1}
	sw = Direction{1, -1}
)

func contains(tiles []*tile, t *tile) bool {
	for _, ti := range tiles {
		if ti == t {
			return true
		}
	}
	return false
}

func stepCost(t tile) float32 {
	cost := float32(1)
	cost += float32(t.heat) / 5 //TODO how much cost for fire etc??
	cost += float32(t.smoke)    //NEW, correct??
	if t.fireLevel > 0 {
		cost = float32(math.Inf(1))
	}
	return cost
}

func getNeighborsPruned(current *tile, dir Direction) []*tile {
	neighbors := []*tile{}

	north := validTile(current.neighborNorth)
	east := validTile(current.neighborEast)
	west := validTile(current.neighborWest)
	south := validTile(current.neighborSouth) // replace !?

	if dir.yDir == 0 {
		if dir.xDir == 0 {
			return getNeighbors(current)
		}
		if dir.xDir == -1 { // går rakt uppåt
			if north {
				neighbors = append(neighbors, current.neighborNorth)
			}
			if !validTile(current.neighborSW) && west {
				neighbors = append(neighbors, current.neighborWest)
				if north && validTile(current.neighborNW) {
					neighbors = append(neighbors, current.neighborNW)
				}
			}
			if !validTile(current.neighborSE) && east {
				neighbors = append(neighbors, current.neighborEast)
				if north && validTile(current.neighborNE) {
					neighbors = append(neighbors, current.neighborNE)
				}
			}

		} else { // går rakt neråt
			if validTile(current.neighborSouth) {
				neighbors = append(neighbors, current.neighborSouth)
			}
			if !validTile(current.neighborNW) && west {
				neighbors = append(neighbors, current.neighborWest)
				if south && validTile(current.neighborSW) {
					neighbors = append(neighbors, current.neighborSW)
				}
			}
			if !validTile(current.neighborNE) && east {
				neighbors = append(neighbors, current.neighborEast)
				if south && validTile(current.neighborSE) {
					neighbors = append(neighbors, current.neighborSE)
				}
			}
		}
	} else if dir.yDir == 1 {
		if dir.xDir == 1 { // går SE
			if validTile(current.neighborEast) {
				neighbors = append(neighbors, current.neighborEast)
			}
			if validTile(current.neighborSouth) {
				neighbors = append(neighbors, current.neighborSouth)
			}
			if east && south && validTile(current.neighborSE) {
				neighbors = append(neighbors, current.neighborSE)
			}
		} else if dir.xDir == -1 { // går NE
			if east {
				neighbors = append(neighbors, current.neighborEast)
			}
			if east && north && validTile(current.neighborSE) {
				neighbors = append(neighbors, current.neighborNE)
			}
			if north {
				neighbors = append(neighbors, current.neighborNorth)
			}
		} else { // ydir = 0,  går höger   //fixed!
			if east {
				neighbors = append(neighbors, current.neighborEast)
			}
			if !validTile(current.neighborNW) && north {
				neighbors = append(neighbors, current.neighborNorth)
				if east && validTile(current.neighborNE) {
					neighbors = append(neighbors, current.neighborNE)
				}
			}
			if !validTile(current.neighborSW) && south {
				neighbors = append(neighbors, current.neighborSouth)
				if east && validTile(current.neighborSE) {
					neighbors = append(neighbors, current.neighborSE)
				}
			}
		}
	} else { //xdir = -1
		if dir.xDir == -1 { // går NW
			if west {
				neighbors = append(neighbors, current.neighborWest)
			}
			if north {
				neighbors = append(neighbors, current.neighborNorth)
			}
			if west && north && validTile(current.neighborNW) {
				neighbors = append(neighbors, current.neighborNW)
			}
		} else if dir.xDir == 1 { // går SW
			if west {
				neighbors = append(neighbors, current.neighborWest)
			}
			if south {
				neighbors = append(neighbors, current.neighborSouth)
			}
			if west && south && validTile(current.neighborSW) {
				neighbors = append(neighbors, current.neighborSW)
			}
		} else { // ydir = 0,  går vänster
			if west {
				neighbors = append(neighbors, current.neighborWest)
			}
			if !validTile(current.neighborNE) && north {
				neighbors = append(neighbors, current.neighborNorth)
				if west && validTile(current.neighborNW) {
					neighbors = append(neighbors, current.neighborNW)
				}
			}
			if !validTile(current.neighborSE) && south {
				neighbors = append(neighbors, current.neighborSouth)
				if west && validTile(current.neighborSW) {
					neighbors = append(neighbors, current.neighborSW)
				}
			}
		}
	}
	return neighbors
}

func getNeighbors(current *tile) []*tile {
	neighbors := []*tile{}

	north := validTile(current.neighborNorth)
	east := validTile(current.neighborEast)
	west := validTile(current.neighborWest)
	south := validTile(current.neighborSouth)

	if north {
		neighbors = append(neighbors, current.neighborNorth)
		if west && validTile(current.neighborNW) {
			neighbors = append(neighbors, current.neighborNW)
		}
		if east && validTile(current.neighborNE) {
			neighbors = append(neighbors, current.neighborNE)
		}
	}
	if east {
		neighbors = append(neighbors, current.neighborEast)
	}
	if west {
		neighbors = append(neighbors, current.neighborWest)
	}
	if south {
		neighbors = append(neighbors, current.neighborSouth)
		if west && validTile(current.neighborSW) {
			neighbors = append(neighbors, current.neighborSW)
		}
		if east && validTile(current.neighborSE) {
			neighbors = append(neighbors, current.neighborSE)
		}
	}
	return neighbors
}

func validTile(t *tile) bool {
	if t == nil {
		return false
	}
	return !t.wall && !t.outOfBounds && t.heat < 1 //&& t.smoke < 100
}

func canGo(t *tile) bool {
	if validTile(t) {
		if t.occupied != nil {
			return !t.occupied.screwed
		}
		return true
	}
	return false
	/*	if t == nil {return false}
		if t.occupied != nil && t.occupied.screwed {return false}
		return !t.wall && !t.outOfBounds && t.heat < 2 */
}

func compactPath(parentOf map[*tile]*tile, from *tile, to *tile) ([]*tile, bool) {
	path := []*tile{to}

	current := to

	for current.xCoord != from.xCoord || current.yCoord != from.yCoord {
		path = append([]*tile{parentOf[current]}, path...)

		ok := true
		current, ok = parentOf[current]

		if !ok {
			return nil, false
		}
	}
	return path, true
}

func getDir(from *tile, to *tile) Direction {
	if from == nil {
		return Direction{1, 1}
	}
	//	if from == to {panic("same same")}

	x := to.xCoord - from.xCoord
	y := to.yCoord - from.yCoord
	if x > 1 {
		x = 1
	} else if x < 0 {
		x = -1
	}
	if y > 1 {
		y = 1
	} else if y < 0 {
		y = -1
	}
	return Direction{x, y}
}

func (t *tile) followDir(dir Direction) *tile { // diagonalt!
	if dir.xDir == 1 {
		if !validTile(t.neighborSouth) {
			return nil
		}
		if dir.yDir == 1 {
			if !validTile(t.neighborEast) {
				return nil
			}
			if !validTile(t.neighborSE) {
				return nil
			}
			return t.neighborSE
		}
		if dir.yDir == -1 {
			if !validTile(t.neighborWest) {
				return nil
			}
			if !validTile(t.neighborSW) {
				return nil
			}
			return t.neighborSW
		}
	}
	if dir.xDir == -1 {
		if !validTile(t.neighborNorth) {
			return nil
		}
		if dir.yDir == 1 {
			if !validTile(t.neighborEast) {
				return nil
			}
			if !validTile(t.neighborNE) {
				return nil
			}
			return t.neighborNE
		}
		if dir.yDir == -1 {
			if !validTile(t.neighborWest) {
				return nil
			}
			if !validTile(t.neighborNW) {
				return nil
			}
			return t.neighborNW
		}
	}
	return nil
}

func (t *tile) neighbor(dir Direction) *tile {
	if dir == n {
		return t.neighborNorth
	}
	if dir == e {
		return t.neighborEast
	}
	if dir == s {
		return t.neighborSouth
	}
	if dir == w {
		return t.neighborWest
	}

	if dir == nw {
		return t.neighborNW
	}
	if dir == ne {
		return t.neighborNE
	}
	if dir == se {
		return t.neighborSE
	}
	if dir == sw {
		return t.neighborSW
	}

	return nil
}

func Jp(current *tile, dir Direction) []jp { //[]*tile {
	jps := []jp{}
	if (dir.xDir == 0 && dir.yDir != 0) || (dir.yDir == 0 && dir.xDir != 0) {
		tmpJP := getJumpPoint(current, dir)
		for tmpJP.jp != nil && tmpJP.jp.occupied != nil {
			jps = append(jps, tmpJP)
			tmpJP = getJumpPoint(tmpJP.jp.neighbor(dir), dir)
		}
		return append(jps, tmpJP)
	}
	for current != nil {
		jpX := getJumpPoint(current, Direction{dir.xDir, 0})
		if jpX.jp != nil {
			tmpJP := jp{current, []*tile{jpX.jp}}
			jps = append(jps, tmpJP)
			//	jps = append(jps, jp{})
		}
		jpY := getJumpPoint(current, Direction{0, dir.yDir})
		if jpY.jp != nil {
			tmpJP := jp{current, []*tile{jpY.jp}}
			jps = append(jps, tmpJP)
		}
		tempJP := sneJP(current, dir)
		if tempJP.jp != nil {
			//	fmt.Println("temp?:", tempJP.jp)
			jps = append(jps, tempJP)
		}
		current = current.followDir(dir)
	}
	return jps
}

func sneJP(current *tile, dir Direction) jp {
	curJP := jp{}
	if dir == nw {
		if !validTile(current.neighborSW) && validTile(current.neighborWest) {
			curJP.jp = current
			curJP.fn = append(curJP.fn, current.neighborWest)
		}
		if !validTile(current.neighborNE) && validTile(current.neighborNorth) {
			curJP.jp = current
			curJP.fn = append(curJP.fn, current.neighborNorth)
		}
	} else if dir == ne {
		if !validTile(current.neighborNW) && validTile(current.neighborNorth) {
			curJP.jp = current
			curJP.fn = append(curJP.fn, current.neighborNorth)
		}
		if !validTile(current.neighborSE) && validTile(current.neighborEast) {
			curJP.jp = current
			curJP.fn = append(curJP.fn, current.neighborEast)
		}
	} else if dir == se {
		if !validTile(current.neighborNE) && validTile(current.neighborEast) {
			curJP.jp = current
			curJP.fn = append(curJP.fn, current.neighborEast)
		}
		if !validTile(current.neighborSW) && validTile(current.neighborSouth) {
			curJP.jp = current
			curJP.fn = append(curJP.fn, current.neighborSouth)
		}
	} else if dir == sw {
		if !validTile(current.neighborNW) && validTile(current.neighborWest) {
			curJP.jp = current
			curJP.fn = append(curJP.fn, current.neighborWest)
		}
		if !validTile(current.neighborSE) && validTile(current.neighborSouth) {
			curJP.jp = current
			curJP.fn = append(curJP.fn, current.neighborSouth)
		}
	}
	return curJP
}

func getJumpPoint(current *tile, dir Direction) jp {
	curJP := jp{}
	if current == nil {
		return curJP
	}
	if current.door {
		return jp{current, nil}
	}
	if dir.xDir == 0 {
		if dir.yDir == 1 { // höger
			return eastJP(current)
		}
		if dir.yDir == -1 { // vänster
			return westJP(current)
		}
	}
	if dir.yDir == 0 {
		if dir.xDir == 1 { // neråt
			return (southJP(current))
		}
		if dir.xDir == -1 { // uppåt
			return northJP(current)

		}
	}
	return curJP
}

func eastJP(current *tile) jp {
	curJP := jp{}
	if (!validTile(current.neighborNW)) || !validTile(current.neighborSW) {
		if !validTile(current.neighborNW) {
			if validTile(current.neighborNorth) {
				curJP.fn = append(curJP.fn, current.neighborNorth)
				if validTile(current.neighborEast) && validTile(current.neighborNE) {
					curJP.fn = append(curJP.fn, current.neighborNE)
				}
			}
			if validTile(current.neighborSouth) && validTile(current.neighborEast) && validTile(current.neighborSE) && !validTile(current.neighborSW) {
				curJP.fn = append(curJP.fn, current.neighborSE)
			}
		}
		if !validTile(current.neighborSW) {
			if validTile(current.neighborSouth) {
				curJP.fn = append(curJP.fn, current.neighborSouth)
				if validTile(current.neighborEast) && validTile(current.neighborSE) {
					curJP.fn = append(curJP.fn, current.neighborSE)
				} //TODO: prettify!! implement for similair funcs
			}
			if validTile(current.neighborNorth) && validTile(current.neighborEast) && validTile(current.neighborNE) && !validTile(current.neighborNW) {
				curJP.fn = append(curJP.fn, current.neighborNE)
			}
		}
		if len(curJP.fn) > 0 {
			curJP.jp = current
			return curJP
		}
	}
	if validTile(current.neighborEast) {
		return getJumpPoint(current.neighborEast, e)
	} else {
		return curJP
	}
}

func westJP(current *tile) jp {
	curJP := jp{}
	if current.door {
		curJP.jp = current
		return curJP
	}
	if !validTile(current.neighborNE) || !validTile(current.neighborSE) {
		if !validTile(current.neighborNE) {
			if validTile(current.neighborNorth) {
				curJP.fn = append(curJP.fn, current.neighborNorth)
				if validTile(current.neighborWest) && validTile(current.neighborNW) {
					curJP.fn = append(curJP.fn, current.neighborNW)
				}
			}
			if validTile(current.neighborWest) && validTile(current.neighborSouth) && validTile(current.neighborSW) && !validTile(current.neighborSE) {
				curJP.fn = append(curJP.fn, current.neighborSW)
			}
		}
		if !validTile(current.neighborSE) {
			if validTile(current.neighborSouth) {
				curJP.fn = append(curJP.fn, current.neighborSouth)
				if validTile(current.neighborWest) && validTile(current.neighborSW) {
					curJP.fn = append(curJP.fn, current.neighborSW)
				}
			}
			if validTile(current.neighborNorth) && validTile(current.neighborWest) && validTile(current.neighborNW) && !validTile(current.neighborNE) {
				curJP.fn = append(curJP.fn, current.neighborNW)
			}
		}
		if len(curJP.fn) > 0 {
			curJP.jp = current
			if validTile(current.neighborWest) {
				curJP.fn = append(curJP.fn, current.neighborWest)
			}
			return curJP
		}
	}
	if validTile(current.neighborWest) {
		return getJumpPoint(current.neighborWest, w) // TODO: 'southJP' instead of getjp
	} else {
		return curJP
	} //lr nil? right..?
}

func southJP(current *tile) jp {
	curJP := jp{}

	if !validTile(current.neighborNW) || !validTile(current.neighborNE) {
		if !validTile(current.neighborNW) {
			if validTile(current.neighborWest) {
				curJP.fn = append(curJP.fn, current.neighborWest)
			}
			if validTile(current.neighborSouth) && validTile(current.neighborEast) && validTile(current.neighborSE) && !validTile(current.neighborNE) {
				curJP.fn = append(curJP.fn, current.neighborSE)
			}
		}
		if !validTile(current.neighborNE) {
			if validTile(current.neighborEast) {
				curJP.fn = append(curJP.fn, current.neighborEast)
			}
			if validTile(current.neighborSouth) && validTile(current.neighborWest) && validTile(current.neighborSW) && !validTile(current.neighborNW) {

				curJP.fn = append(curJP.fn, current.neighborSW)
			}
		}
		if len(curJP.fn) > 0 {
			curJP.jp = current
			if validTile(current.neighborSouth) {
				curJP.fn = append(curJP.fn, current.neighborSouth)
			}
			return curJP
		}
	}
	if validTile(current.neighborSouth) {

		return getJumpPoint(current.neighborSouth, s) // TODO: 'southJP' instead of getjp
	}
	return curJP
}

func northJP(current *tile) jp {
	curJP := jp{}
	if !validTile(current.neighborSE) || !validTile(current.neighborSW) {
		if !validTile(current.neighborSE) {
			if validTile(current.neighborEast) {
				curJP.fn = append(curJP.fn, current.neighborEast)
			}
			if validTile(current.neighborWest) && validTile(current.neighborNorth) && validTile(current.neighborNW) && !validTile(current.neighborSW) {
				curJP.fn = append(curJP.fn, current.neighborNW)
			}
		}
		if !validTile(current.neighborSW) {
			if validTile(current.neighborWest) {
				curJP.fn = append(curJP.fn, current.neighborWest)
			}
			if validTile(current.neighborNorth) && validTile(current.neighborEast) && validTile(current.neighborNE) && !validTile(current.neighborSE) {
				curJP.fn = append(curJP.fn, current.neighborNE)
			}
		}
		if len(curJP.fn) > 0 {
			curJP.jp = current
			if validTile(current.neighborNorth) {
				curJP.fn = append(curJP.fn, current.neighborNorth)
			}
			return curJP
		}
	}
	if validTile(current.neighborNorth) {
		return getJumpPoint(current.neighborNorth, n)
	} else {
		return curJP
	}
}

func smplCost(t1 *tile, t2 *tile) float32 {
	cost := smplDistance(t1, t2)
	if t1.occupied != nil {
		cost += 10
	}
	if t2.occupied != nil {
		cost += 10
	}
	return cost
}

func smplDistance(t1 *tile, t2 *tile) float32 {
	xDif := math.Max(float64(t1.xCoord), float64(t2.xCoord)) - math.Min(float64(t1.xCoord), float64(t2.xCoord))
	yDif := math.Max(float64(t1.yCoord), float64(t2.yCoord)) - math.Min(float64(t1.yCoord), float64(t2.yCoord))
	if xDif == 0 {
		return float32(yDif)
	}
	if yDif == 0 {
		return float32(xDif)
	}

	return float32(math.Sqrt(xDif*xDif + yDif*yDif))
}

func notSoSmplCost(t1 *tile, t2 *tile) float32 {
	dir := getDir(t1, t2)
	next := nextTile(t1, dir)
	cost := stepCost(*next) - 1
	for next != t2 {
		next := nextTile(next, dir)
		cost += stepCost(*next) - 1
	}
	return cost
}

func nextTile(t1 *tile, dir Direction) *tile {
	if dir == n {
		return t1.neighborNorth
	}
	if dir == e {
		return t1.neighborEast
	}
	if dir == s {
		return t1.neighborSouth
	}
	if dir == w {
		return t1.neighborWest
	}

	if dir == nw {
		return t1.neighborNW
	}
	if dir == ne {
		return t1.neighborNE
	}
	if dir == se {
		return t1.neighborSE
	}
	if dir == sw {
		return t1.neighborSW
	}

	return nil // default?
}

func getPath3(m *[][]tile, from []*tile) { //INIT!

	// map över jp
	var parentOf map[*tile]*tile
	parentOf = make(map[*tile]*tile)

	// map över costOf, in beta!
	var costOf map[*tile]float32
	costOf = make(map[*tile]float32)
	//

	cq := queue{}

	/*for i, list := range *m {
		for j, _ := range list {
			val := float32(math.Inf(1))
			cq = append(cq, tileCost{&(*m)[i][j], &val})
			costOf[&(*m)[i][j]] = val  //beta
		}
	}*/

	for _, f := range from {
		cq.Update(f, 0)
		costOf[f] = 0 //beta
	}
	v := float32(0)
	current := tileCost{&tile{}, &v}
	currentDir := Direction{0, 0}

	for len(cq) != 0 {
		current = (&cq).Pop()
		if *current.cost == float32(math.Inf(1)) {
			return
		}
		_, ok := parentOf[current.tile]

		//	if !ok || (ok && (current.tile.occupied == nil || (current.tile.occupied != nil && len(current.tile.occupied.plan) == 0))) {
		if ok {
			currentDir = getDir(parentOf[current.tile], current.tile)

		} else {
			currentDir = Direction{0, 0}
		}
		neighbors := /*getNeighbors(current.tile, cq)*/ getNeighborsPruned(current.tile, currentDir)

		var wg sync.WaitGroup
		wg.Add(len(neighbors))
		var mutex = &sync.Mutex{}
		for _, neighbor := range neighbors {

			go func(n *tile) {
				defer wg.Done()
				//n := neighbor
				jps := JpInit(n, getDir(current.tile, n))
				for _, jp := range jps {
					if jp.jp != nil {
						mutex.Lock()
						cost := *current.cost + smplCost(current.tile, jp.jp) + 100*float32(jp.jp.smoke)
						p, ok := parentOf[jp.jp]
						//if !ok || (ok && cost < cq.costOf(jp.jp) && p.smoke >= current.tile.smoke) {
						if !ok || (ok && cost < costOf[jp.jp] && p.smoke >= current.tile.smoke) { //beta
							parentOf[jp.jp] = current.tile
							cq.Update(jp.jp, cost)
							costOf[jp.jp] = cost //beta

							if jp.jp.occupied != nil {
								setPlan(parentOf, jp.jp)
							}
						}
						for _, n := range jp.fn {
							fnCost := cost + smplCost(jp.jp, n) + 100*float32(n.smoke) //float32(math.Mod(float64(n.smoke), 50))
							p, ok := parentOf[n]
							//							if jp.jp != n && (!ok || (n != nil && fnCost < cq.costOf(n) && p.smoke >= n.smoke)) {
							if jp.jp != n && (!ok || (n != nil && fnCost < costOf[n] && p.smoke >= n.smoke)) { //beta
								parentOf[n] = jp.jp
								cq.Update(n, fnCost)
								costOf[n] = fnCost //beta
								if n.occupied != nil {
									setPlan(parentOf, n)
								}
							}
						}
						mutex.Unlock()
					}
				}
			}(neighbor)
		}
		wg.Wait()
		//}
	}
}

func setPlan(parentOf map[*tile]*tile, pers *tile) {
	path := []*tile{pers}

	current := pers
	ok := true

	curDir := Direction{}
	lastDir := Direction{0, 0}
	for !current.door {
		curDir = getDir(current, parentOf[current])
		if curDir == lastDir {
			path = path[:len(path)-1]
		}
		path = append(path, parentOf[current])
		if parentOf[current].xCoord == 28 && parentOf[current].yCoord == 49 {
		}

		current, ok = parentOf[current]
		lastDir = curDir

		if !ok {
			//fmt.Println("nope?")
			//pers.occupied.plan = []*tile{pers.safestTile()}
			//pers.occupied.plan = []*tile{}
			return
		}
	}
	if ok {
		pers.occupied.plan = path
	} else {
		//fmt.Println("nope?2")
		/*	pers.occupied.plan = []*tile{pers.safestTile()}*/
	}

	/*else if len(pers.occupied.plan) > 0 && pers.occupied.dir != getDir(pers, pers.occupied.plan[0]){
	d := getDir(pers.occupied.currentTile(), pers.occupied.plan[0])
	if d.xDir == 0 && d.yDir == 0 {
		pers.occupied.plan = pers.occupied.plan[1:]
		if len(pers.occupied.plan) > 0 {pers.occupied.dir = getDir(pers, pers.occupied.plan[0])}
	} */ // PBS above tmp borttaget..

	/*else {pers.occupied.dir = getDir(pers.occupied.currentTile(), pers.occupied.plan[0])}// Maybe onödig??*/
	//	}
}

func getJumpPointInit(current *tile, dir Direction) jp {
	curJP := jp{}
	if current == nil || !validTile(current) {
		return curJP
	}
	if current.occupied != nil {
		return jp{current, nil}
	}

	if dir.xDir == 0 {
		if dir.yDir == 1 { // höger
			return eastJPInit(current)
		}
		if dir.yDir == -1 { // vänster
			return westJPInit(current)
		}
	}
	if dir.yDir == 0 {
		if dir.xDir == 1 { // neråt
			return (southJPInit(current))
		}
		if dir.xDir == -1 { // uppåt
			return northJPInit(current)
		}
	}
	return curJP
}

func JpInit(current *tile, dir Direction) []jp {
	jps := []jp{}
	if current.occupied != nil {
		jps = append(jps, jp{current, []*tile{}})
	}
	if (dir.xDir == 0 && dir.yDir != 0) || (dir.yDir == 0 && dir.xDir != 0) {
		tmpJP := getJumpPointInit(current, dir)
		for tmpJP.jp != nil && tmpJP.jp.occupied != nil {
			jps = append(jps, tmpJP)
			tmpJP = getJumpPointInit(tmpJP.jp.neighbor(dir), dir)
		}
		return append(jps, tmpJP)
	}

	for current != nil {
		jpX := getJumpPointInit(current, Direction{dir.xDir, 0})
		if jpX.jp != nil {
			tmpJP := jp{current, []*tile{jpX.jp}}
			jps = append(jps, tmpJP)
		}
		jpY := getJumpPointInit(current, Direction{0, dir.yDir})
		if jpY.jp != nil {
			tmpJP := jp{current, []*tile{jpY.jp}}
			jps = append(jps, tmpJP)
		}
		tempJP := sneJP(current, dir)
		if tempJP.jp != nil {
			jps = append(jps, tempJP)
		}
		current = current.followDir(dir)
		if current != nil && current.occupied != nil {
			jps = append(jps, jp{current, []*tile{}})
		}
	}
	return jps
}

func eastJPInit(current *tile) jp {
	curJP := jp{}
	if (!validTile(current.neighborNW)) || !validTile(current.neighborSW) {
		if !validTile(current.neighborNW) {
			if validTile(current.neighborNorth) {
				curJP.fn = append(curJP.fn, current.neighborNorth)
				if validTile(current.neighborEast) && validTile(current.neighborNE) {
					curJP.fn = append(curJP.fn, current.neighborNE)
				}
			}
			if validTile(current.neighborSouth) && validTile(current.neighborEast) && validTile(current.neighborSE) && !validTile(current.neighborSW) {
				curJP.fn = append(curJP.fn, current.neighborSE)
			}
		}
		if !validTile(current.neighborSW) {
			if validTile(current.neighborSouth) {
				curJP.fn = append(curJP.fn, current.neighborSouth)
				if validTile(current.neighborEast) && validTile(current.neighborSE) {
					curJP.fn = append(curJP.fn, current.neighborSE)
				} //TODO: prettify!! implement for similair funcs
			}
			if validTile(current.neighborNorth) && validTile(current.neighborEast) && validTile(current.neighborNE) && !validTile(current.neighborNW) {
				curJP.fn = append(curJP.fn, current.neighborNE)
			}

		}
		if len(curJP.fn) > 0 {
			curJP.jp = current
			return curJP
		}
	}
	if validTile(current.neighborEast) {
		return getJumpPointInit(current.neighborEast, e)
	} else {
		return curJP
	}
}

func westJPInit(current *tile) jp {
	curJP := jp{}
	if current.door {
		curJP.jp = current
		return curJP
	}
	if !validTile(current.neighborNE) || !validTile(current.neighborSE) {

		if !validTile(current.neighborNE) {
			if validTile(current.neighborNorth) {
				curJP.fn = append(curJP.fn, current.neighborNorth)
				if validTile(current.neighborWest) && validTile(current.neighborNW) {
					curJP.fn = append(curJP.fn, current.neighborNW)
				}
			}
			if validTile(current.neighborWest) && validTile(current.neighborSouth) && validTile(current.neighborSW) && !validTile(current.neighborSE) {
				curJP.fn = append(curJP.fn, current.neighborSW)
			}
		}
		if !validTile(current.neighborSE) {
			if validTile(current.neighborSouth) {
				curJP.fn = append(curJP.fn, current.neighborSouth)
				if validTile(current.neighborWest) && validTile(current.neighborSW) {
					curJP.fn = append(curJP.fn, current.neighborSW)
				}
			}
			if validTile(current.neighborNorth) && validTile(current.neighborWest) && validTile(current.neighborNW) && !validTile(current.neighborNE) {
				curJP.fn = append(curJP.fn, current.neighborNW)
			}
		}
		if len(curJP.fn) > 0 {
			curJP.jp = current
			if validTile(current.neighborWest) {
				curJP.fn = append(curJP.fn, current.neighborWest)
			}
			return curJP
		}
	}
	if validTile(current.neighborWest) {
		return getJumpPointInit(current.neighborWest, w)
	} else {
		return curJP
	}
}

func southJPInit(current *tile) jp {
	curJP := jp{}

	if !validTile(current.neighborNW) || !validTile(current.neighborNE) {
		if !validTile(current.neighborNW) {
			if validTile(current.neighborWest) {
				curJP.fn = append(curJP.fn, current.neighborWest)
			}

			if validTile(current.neighborSouth) && validTile(current.neighborEast) && validTile(current.neighborSE) && !validTile(current.neighborNE) {
				curJP.fn = append(curJP.fn, current.neighborSE)
			}
		}
		if !validTile(current.neighborNE) {
			if validTile(current.neighborEast) {
				curJP.fn = append(curJP.fn, current.neighborEast)
			}
			if validTile(current.neighborSouth) && validTile(current.neighborWest) && validTile(current.neighborSW) && !validTile(current.neighborNW) {

				curJP.fn = append(curJP.fn, current.neighborSW)
			}
		}
		if len(curJP.fn) > 0 {
			curJP.jp = current
			if validTile(current.neighborSouth) {
				curJP.fn = append(curJP.fn, current.neighborSouth)
			}
			return curJP
		}
	}
	if validTile(current.neighborSouth) {
		return getJumpPointInit(current.neighborSouth, s)
	}
	return curJP
}

func northJPInit(current *tile) jp {

	curJP := jp{}
	if !validTile(current.neighborSE) || !validTile(current.neighborSW) {
		if !validTile(current.neighborSE) {
			if validTile(current.neighborEast) {
				curJP.fn = append(curJP.fn, current.neighborEast)
			}
			if validTile(current.neighborWest) && validTile(current.neighborNorth) && validTile(current.neighborNW) && !validTile(current.neighborSW) {
				curJP.fn = append(curJP.fn, current.neighborNW)
			}
		}
		if !validTile(current.neighborSW) {
			if validTile(current.neighborWest) {
				curJP.fn = append(curJP.fn, current.neighborWest)
			}
			if validTile(current.neighborNorth) && validTile(current.neighborEast) && validTile(current.neighborNE) && !validTile(current.neighborSE) {
				curJP.fn = append(curJP.fn, current.neighborNE)
			}
		}
		if len(curJP.fn) > 0 {
			curJP.jp = current
			if validTile(current.neighborNorth) {
				curJP.fn = append(curJP.fn, current.neighborNorth)
			}
			return curJP
		}
	}
	if validTile(current.neighborNorth) {
		return getJumpPointInit(current.neighborNorth, n)
	} else {
		return curJP
	}
}

func (p *Person) redirect() bool {
	if p.dir == e {
		return p.redE()
	}
	if p.dir == s {
		return p.redS()
	}
	if p.dir == w {
		return p.redW()
	}
	if p.dir == n {
		return p.redN()
	}

	if p.dir == se {
		return p.redSE()
	}
	if p.dir == sw {
		return p.redSW()
	}
	if p.dir == nw {
		return p.redNW()
	}
	if p.dir == ne {
		return p.redNE()
	}

	return false
}

func (p *Person) redSE() bool {
	current := p.currentTile()
	if current.neighborSouth.smoke <= current.neighborEast.smoke {
		if p.moveTo(current.neighborSouth) {
			return true
		}
	}
	if p.moveTo(current.neighborEast) {
		return true
	}
	if p.moveTo(current.neighborSouth) {
		return true
	}
	//return p.moveTo(current.safestTile())
	return false
}

func (p *Person) redSW() bool {
	current := p.currentTile()
	if current.neighborSouth.smoke <= current.neighborWest.smoke {
		if p.moveTo(current.neighborSouth) {
			return true
		}
	}
	if p.moveTo(current.neighborWest) {
		return true
	}
	if p.moveTo(current.neighborSouth) {
		return true
	}
	//	return p.moveTo(current.safestTile())
	return false
}

func (p *Person) redNW() bool {
	current := p.currentTile()
	if current.neighborNorth.smoke <= current.neighborWest.smoke {
		if p.moveTo(current.neighborNorth) {
			return true
		}
	}
	if p.moveTo(current.neighborWest) {
		return true
	}
	if p.moveTo(current.neighborNorth) {
		return true
	}
	//return p.moveTo(current.safestTile())
	return false
}

func (p *Person) redNE() bool {
	current := p.currentTile()
	if current.neighborNorth.smoke <= current.neighborEast.smoke {
		if p.moveTo(current.neighborNorth) {
			return true
		}
	}
	if p.moveTo(current.neighborEast) {
		return true
	}
	if p.moveTo(current.neighborNorth) {
		return true
	}
	//return p.moveTo(current.safestTile())
	return false
}

func (p *Person) redE() bool {
	current := p.currentTile()
	if current.neighborSouth == nil {
		return p.moveTo(current.neighborNorth)
	}
	if current.neighborNorth == nil {
		return p.moveTo(current.neighborSouth)
	}

	if current.neighborSouth.smoke <= current.neighborNorth.smoke {
		if p.moveTo(current.neighborSouth) {
			return true
		}
	}
	if p.moveTo(current.neighborNorth) {
		return true
	}
	if p.moveTo(current.neighborSouth) {
		return true
	}
	//return p.moveTo(current.safestTile())
	return false
}

func (p *Person) redS() bool {
	current := p.currentTile()
	if current.neighborWest == nil {
		return p.moveTo(current.neighborEast)
	}
	if current.neighborEast == nil {
		return p.moveTo(current.neighborWest)
	}

	if current.neighborWest.smoke <= current.neighborEast.smoke {
		if p.moveTo(current.neighborWest) {
			return true
		}
	}
	if p.moveTo(current.neighborEast) {
		return true
	}
	if p.moveTo(current.neighborWest) {
		return true
	}
	//	return p.moveTo(current.safestTile())
	return false
}

func (p *Person) redW() bool {
	current := p.currentTile()
	if current.neighborSouth == nil {
		return p.moveTo(current.neighborNorth)
	}
	if current.neighborNorth == nil {
		return p.moveTo(current.neighborSouth)
	}

	if current.neighborSouth.smoke <= current.neighborNorth.smoke {
		if p.moveTo(current.neighborSouth) {
			return true
		}
	}
	if p.moveTo(current.neighborNorth) {
		return true
	}
	if p.moveTo(current.neighborSouth) {
		return true
	}
	//	return p.moveTo(current.safestTile())
	return false
}

func (p *Person) redN() bool {
	current := p.currentTile()
	if current.neighborWest == nil {
		return p.moveTo(current.neighborEast)
	}
	if current.neighborEast == nil {
		return p.moveTo(current.neighborWest)
	}

	if current.neighborWest.smoke <= current.neighborEast.smoke {
		if p.moveTo(current.neighborWest) {
			return true
		}
	}
	if p.moveTo(current.neighborEast) {
		return true
	}
	if p.moveTo(current.neighborWest) {
		return true
	}
	//	return p.moveTo(current.safestTile())   //TODO: obs unsure if 'safestfile' makes for a wierd looking redirect?
	return false
}

func (t *tile) endOfLine(dir Direction) *tile {
	next := t.followDir(dir)
	for next != nil {
		tmp := next.followDir(dir)
		if tmp == nil {
			return next
		}
		next = tmp
	}
	return nil
}
