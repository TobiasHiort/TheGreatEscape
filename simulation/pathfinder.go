package main

import (
	"fmt"
	"math"
	"sync"
)

type Direction struct {
	xDir int   //-1,0,1
	yDir int   //-1,0,1
}

func getPath(m *[][]tile, from *tile) ([]*tile, bool) {

	// map to keep track of the final path
	var parentOf map[*tile]*tile
	parentOf = make(map[*tile]*tile)

	//initialise 'costqueue', start-0, other-infinite
	costQueue := queue{}

	for i, list := range *m {
		for j, _ := range list {
			costQueue.Add(&(*m)[i][j], float32(math.Inf(1)))	
		}
	}

	costQueue.Update(from, 0)

//	checkedQueue := queue{}  // TODO: implement this later for a more efficient algorithm: dummer! går inte fortare..

	v := float32(0)
	current := tileCost{&tile{}, &v}
	currentDir := Direction{1,1}
	
	for len(costQueue) != 0 && !current.tile.door {
	//	fmt.Println("----")
		current = (&costQueue).Pop()
		currentDir = getDir(parentOf[current.tile], current.tile)  // for reference!
	//	fmt.Println(currentDir)
		//neighbors := getNeighbors(current.tile, costQueue)
		neighbors := getNeighborsPruned(current.tile, currentDir)
		var wg sync.WaitGroup
		wg.Add(len(neighbors))
		var mutex = &sync.Mutex{}
		for _, neighbor := range neighbors {
		//	fmt.Println(neighbor)
			go func(n *tile) {			
				defer wg.Done()			
				cost := *current.cost + stepCost(*n)
				if Diagonal(current.tile, n) {cost += float32(math.Sqrt(2)) - 1}
				if n.occupied.IsWaiting() {cost += 1}

				// TODO: 1 default cost improve!? depending on heat, smoke etc
				mutex.Lock()
				if cost < costQueue.costOf(n) {
					
					parentOf[n] = current.tile
					costQueue.Update(n, cost)				
				}
				mutex.Unlock()
			}(neighbor)		
		}
		wg.Wait()	
		//	checkedQueue.AddTC(current)	
	}	
	return compactPath(parentOf, from, current.tile)
}

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
	if t.fireLevel > 0 {
		cost = float32(math.Inf(1))
	}
	return cost
}


/*
func getJumpPoint(m *[][]tile, current *tile, dir Direction, from *tile, to *tile) *tile {
	//from+to onödig(?)
	nextX := current.xCoord + dir.xDir
	nextY := current.yCoord + dir.yDir
	nextTile := GetTile(m, nextX, nextY)

	if nextTile == nil {return nil}
	
	if nextTile.door {return nextTile}

	//	if 
}*/

func getNeighborsPruned(current *tile, dir Direction) []*tile{
	neighbors := []*tile{}

	north := validTile(current.neighborNorth) 
	east := validTile(current.neighborEast)
	west := validTile(current.neighborWest)
	south := validTile(current.neighborSouth)   // replace !?
	
	if dir.yDir == 0 {  // horisontal/vertical? hur vare med coordsen..
		if dir.xDir == 0 {return getNeighbors(current, queue{})}
		if dir.xDir == -1 {  // går rakt uppåt
			if north {
				neighbors = append(neighbors, current.neighborNorth)
				//	if !west {neighbors = append(neighbors, current.neighborNW )}
				//	if !east {neighbors = append(neighbors, current.neighborNE )}
			
			}
			if !validTile(current.neighborSW) && west {
				neighbors = append(neighbors, current.neighborWest)
				if north && validTile(current.neighborNW) {neighbors = append(neighbors, current.neighborNW)}
			}
			if !validTile(current.neighborSE) && east {
				neighbors = append(neighbors, current.neighborEast)
				if north && validTile(current.neighborNE) {neighbors = append(neighbors, current.neighborNE)}
			}
			
		} else { // går rakt neråt
			if validTile(current.neighborSouth) {
				neighbors = append(neighbors, current.neighborSouth)
				//	if !validTile(current.neighborWest) {neighbors = append(neighbors, current.neighborSW )}
				//	if !validTile(current.neighborEast) {neighbors = append(neighbors, current.neighborSE )}				
			}
			if !validTile(current.neighborNW) && west {
				neighbors = append(neighbors, current.neighborWest)
				if south && validTile(current.neighborSW) {neighbors = append(neighbors, current.neighborSW)}	
			}
			if !validTile(current.neighborNE) && east {
				neighbors = append(neighbors, current.neighborEast)
				if south && validTile(current.neighborSE) {neighbors = append(neighbors, current.neighborSE)}
			}  // done so far!			
		}		
	} else if dir.yDir == 1 { 
		if dir.xDir == 1 { // går SE
			if validTile(current.neighborEast) {
				neighbors = append(neighbors, current.neighborEast)			
			}
			if validTile(current.neighborSouth) {
				neighbors = append(neighbors, current.neighborSouth)
			//	if !validTile(current.neighborNW) && validTile(current.neighborNW) {
			//		neighbors = append(neighbors, current.neighborNW)}
			}
			if east && south && validTile(current.neighborSE) {neighbors = append(neighbors, current.neighborSE)}

			
		} else if dir.xDir == -1 { // går NE    
			if east {neighbors = append(neighbors, current.neighborEast)}
			if east && north && validTile(current.neighborSE) {neighbors = append(neighbors, current.neighborNE)}
			if north {
				neighbors = append(neighbors, current.neighborNorth)
				//	if !west && validTile(current.neighborSW) {neighbors = append(neighbors, current.neighborSW)}
			}
			
		} else {  // ydir = 0,  går höger   //fixed!
			if east {
				neighbors = append(neighbors, current.neighborEast)
			}
			if !validTile(current.neighborNW) && north {
				neighbors = append(neighbors, current.neighborNorth)
				if east && validTile(current.neighborNE) {neighbors = append(neighbors, current.neighborNE)}
			}
			if !validTile(current.neighborSW) && south {
				neighbors = append(neighbors, current.neighborSouth)
				if east && validTile(current.neighborSE) {neighbors = append(neighbors, current.neighborSE)}
			}
		}

	} else { //xdir = -1
		if dir.xDir == 1 { // går NW
			if west {neighbors = append(neighbors, current.neighborWest)}
			if north {neighbors = append(neighbors, current.neighborNorth)}
			if west && north && validTile(current.neighborNW) {neighbors = append(neighbors, current.neighborNW)}	
		} else if dir.xDir == -1 { // går SW
			if west {neighbors = append(neighbors, current.neighborWest)}
			if south {neighbors = append(neighbors, current.neighborSouth)}
			if west && south && validTile(current.neighborSW) {neighbors = append(neighbors, current.neighborSW)}
		} else {  // ydir = 0,  går vänster
			if west {neighbors = append(neighbors, current.neighborWest)}
			if !validTile(current.neighborNE) && north {
				neighbors = append(neighbors, current.neighborNorth)
				if west && validTile(current.neighborNW) {neighbors = append(neighbors, current.neighborNW)}
			}
			if !validTile(current.neighborSE) && south {
				neighbors = append(neighbors, current.neighborSouth)
				if west && validTile(current.neighborSW) {neighbors = append(neighbors, current.neighborSW)}
			}
			
		}
	}

	return neighbors
}


func getNeighbors(current *tile, costQueue queue) []*tile {
	neighbors := []*tile{}

	north := validTile(current.neighborNorth) 
	east := validTile(current.neighborEast)
	west := validTile(current.neighborWest)
	south := validTile(current.neighborSouth)

	if north {
		neighbors = append(neighbors, current.neighborNorth)
		if west && validTile(current.neighborNW) {
			neighbors = append(neighbors, current.neighborNW)}
		if east && validTile(current.neighborNE) {
			neighbors = append(neighbors, current.neighborNE)}
	}
	if east {neighbors = append(neighbors, current.neighborEast)}
	if west {neighbors = append(neighbors, current.neighborWest)}
	if south {
		neighbors = append(neighbors, current.neighborSouth)
		if west && validTile(current.neighborSW) {
			neighbors = append(neighbors, current.neighborSW)}
		if east && validTile(current.neighborSE) {
			neighbors = append(neighbors, current.neighborSE)}	
	}

	//  --this--
	// nedanför kollar om värdet finns i costQueue också.. tar längre tid för 100*100 och 100*200 iaf
	
/*	north := validTile(current.neighborNorth) && costQueue.inQueue(current.neighborNorth)
	east := validTile(current.neighborEast) && costQueue.inQueue(current.neighborEast)
	west := validTile(current.neighborWest) && costQueue.inQueue(current.neighborWest)
	south := validTile(current.neighborSouth) && costQueue.inQueue(current.neighborSouth)

	if north {
		neighbors = append(neighbors, current.neighborNorth)
		if west && validTile(current.neighborNW) && costQueue.inQueue(current.neighborNW) {
			neighbors = append(neighbors, current.neighborNW)}
		if east && validTile(current.neighborNE) && costQueue.inQueue(current.neighborNE){
			neighbors = append(neighbors, current.neighborNE)}
	}
	if east {neighbors = append(neighbors, current.neighborEast)}
	if west {neighbors = append(neighbors, current.neighborWest)}
	if south {
		neighbors = append(neighbors, current.neighborSouth)
		if west && validTile(current.neighborSW) && costQueue.inQueue(current.neighborSW){
			neighbors = append(neighbors, current.neighborSW)}
		if east && validTile(current.neighborSE) && costQueue.inQueue(current.neighborSE){
			neighbors = append(neighbors, current.neighborSE)}	
	} */
//  --this--
	
	/*
	
	if validTile(current.neighborNorth) && costQueue.inQueue(current.neighborNorth){
		neighbors = append(neighbors, current.neighborNorth)
	}
	if validTile(current.neighborEast) && costQueue.inQueue(current.neighborEast){
		neighbors = append(neighbors, current.neighborEast)
	}
	if validTile(current.neighborWest) && costQueue.inQueue(current.neighborWest){
		neighbors = append(neighbors, current.neighborWest)
	}
	if validTile(current.neighborSouth) && costQueue.inQueue(current.neighborSouth){
		neighbors = append(neighbors, current.neighborSouth)
	}
	//
	if validTile(current.neighborNW) && costQueue.inQueue(current.neighborNW){
		neighbors = append(neighbors, current.neighborNW)
	}
	if validTile(current.neighborNE) && costQueue.inQueue(current.neighborNE){
		neighbors = append(neighbors, current.neighborNE)
	}
	if validTile(current.neighborSE) && costQueue.inQueue(current.neighborSE){
		neighbors = append(neighbors, current.neighborSE)
	}
	if validTile(current.neighborSW) && costQueue.inQueue(current.neighborSW){
		neighbors = append(neighbors, current.neighborSW)
	}
	*/
	//

	return neighbors
}



func validTile(t *tile) bool {
	if t == nil {
		return false
	}
	return !t.wall && !t.outOfBounds
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

func printPath(path []*tile) {
	if path == nil {
		fmt.Println("No valid path exists")
	}
	for i, t := range path {
		if (t == nil) {
			fmt.Println("End")
		} else {fmt.Println(i, ":", t.xCoord, ",", t.yCoord)}
	}
}

func mainPath() {

	workingPath()
	fmt.Println("--------------")
/*	blockedPath()
	fmt.Println("--------------")
	firePath()*/
	fmt.Println("--------------")
	doorsPath()
}

func workingPath() {
	matrix := [][]int{
		{0, 1, 2, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 1, 0}}
	testmap := TileConvert(matrix)

	path, _ := getPath(&testmap, &testmap[0][0])

	fmt.Println("\nWorking path:")
	printPath(path)
}

func blockedPath() {
	matrix := [][]int{
		{0, 1, 2, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0}}
	testmap := TileConvert(matrix)

	path, _ := getPath(&testmap, &testmap[0][0])

	fmt.Println("\nBlocked path:")
	printPath(path)

}

func firePath() {
	matrix := [][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 2, 0, 0, 0}}
	testmap := TileConvert(matrix)
	SetFire(&(testmap[3][2]))
	for i := 0; i < 10; i++ {
		FireSpread(testmap)
	}

	path, _ := getPath(&testmap, &testmap[0][3])
	fmt.Println("\nFire path:")
	printPath(path)
}

func doorsPath() {
	matrix := [][]int{
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{1, 1, 1, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{2, 0, 0, 1, 0, 0, 0}}

	testmap := TileConvert(matrix)

	path, _ := getPath(&testmap, &testmap[0][0])
	fmt.Println("\nDoors path:")
	printPath(path)
}



// new funcs

func getDir(from *tile, to *tile) Direction{
	if from == nil{ 
		return Direction{1,1}}
	return Direction {to.xCoord - from.xCoord, to.yCoord - from.yCoord}
}

func (t *tile)followDir(dir Direction) *tile{  // diagonalt!
	if dir.xDir == 1 {
		if !validTile(t.neighborSouth) {return nil}
		if dir.yDir == 1 {
			if !validTile(t.neighborEast) {return nil}
			if !validTile(t.neighborSE) {return nil}
			return t.neighborSE
		}
		if dir.yDir == -1 {
			if !validTile(t.neighborWest) {return nil}
			if !validTile(t.neighborSW) {return nil}
			return t.neighborSW
		}
	}
	if dir.xDir == -1 {
		if !validTile(t.neighborNorth) {return nil}
		if dir.yDir == 1 {
			if !validTile(t.neighborEast) {return nil}
			if !validTile(t.neighborNE) {return nil}
			return t.neighborNE
		}
		if dir.yDir == -1 {
			if !validTile(t.neighborWest) {return nil}
			if !validTile(t.neighborNW) {return nil}
			return t.neighborNW
		}
	}
	return nil
}



func Jp(current *tile, dir Direction) []jp {//[]*tile {
	jps := []jp{}
	if current.door {return []jp{jp{current, []*tile{}}}}
	if dir.xDir == 0 {
		if dir.yDir == 1 {return []jp{getJumpPoint(current, dir)}} // höger
		if dir.yDir == -1 {return []jp{getJumpPoint(current, dir)}} // vänster		
	}
	if dir.yDir == 0 {
		if dir.xDir == 1 {return []jp{getJumpPoint(current, dir)}} // neråt
		if dir.xDir == -1 {return []jp{getJumpPoint(current, dir)}} // uppåt
	}

	for current != nil {
		jpX := getJumpPoint(current, Direction{dir.xDir, 0})
		if jpX.jp != nil {
			jps = append(jps, jpX)
			jps = append(jps, jp{})
		}
		jpY := getJumpPoint(current, Direction{0, dir.yDir})
		if jpY.jp != nil {jps = append(jps, jpY)}
		
		tempJP := sneJP(current, dir)
		if tempJP.jp != nil {
			fmt.Println("temp?:", tempJP.jp)
			jps = append(jps, tempJP)
		}
		current = current.followDir(dir) 
	/*	getJPs(current, dir, &jps)
		if current.door {return jps}
		current = current.followDir(dir) */
		
		//	fmt.Println(jps[0])
		//	fmt.Println(jps[1])
	}
	return jps
}

func sneJP(current *tile, dir Direction) jp{
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
/*
func getJPs(current *tile, dir Direction, jps *[]jp) { //*[]*tile) {
	//	current := currentTC.tile
	if current.door {
		*jps = append(*jps, jp{current, nil})
		return
	}// found jp!
	if dir.xDir == 1 {
		sPath := getJumpPoint(current, Direction{1,0})  // jp söderut?    // TODO: set 'parentof' properly!
		if sPath.jp != nil {*jps = append(*jps, sPath)}  // CONTINUE HERE!
		if dir.yDir == 1 { // sydöst
			ePath := getJumpPoint(current, Direction{0,1})  // jp österut?
			if ePath.jp != nil {*jps = append(*jps, ePath)}
			cPath := getJumpPoint(current, dir)
			if cPath.jp != nil {*jps = append(*jps, cPath)}


			if (!validTile(current.neighborSW) && validTile(current.neighborSouth)) || (!validTile(current.neighborNE) && validTile(current.neighborEast)) {
				if (!validTile(current.neighborSW) && validTile(current.neighborSouth)) {}
				if (!validTile(current.neighborNE) && validTile(current.neighborEast)) {}
			}
			
			
			//if !validTile(current.neighborSW) || !validTile(current.neighborNE) {
			//	*jps = append(*jps, current) // found jp!
			//}

				
			//if validTile(current.neighborSouth) && validTile(current.neighborEast) && validTile(current.neighborSE) {
			//	return getJumpPoint(current.neighborSE, dir)
			//} else {return current}  //lr nil? right..?
			return
		}
		if dir.yDir == -1 { // sydväst
			wPath := getJumpPoint(current, Direction{-1,0})  // jp västerut?
			if wPath.jp != nil {*jps = append(*jps, wPath)}
			if !validTile(current.neighborNW) || !validTile(current.neighborSE) {
				*jps = append(*jps, current) // found jp!
			}
			//if validTile(current.neighborWest) && validTile(current.neighborSouth) && validTile(current.neighborSW) {
		//		return getJumpPoint(current.neighborSW, dir)
		//	} else {return current}  //lr nil? right..?
			return
		}
	}
	if dir.xDir == -1 {  
		nPath := getJumpPoint(current, Direction{-1,0})  // jp norrut?
		if nPath != nil {*jps = append(*jps, nPath)}
		if dir.yDir == 1 { // nordöst
			ePath := getJumpPoint(current, Direction{0,1})  // jp österut?
			if ePath != nil {*jps = append(*jps, ePath)}
			if !validTile(current.neighborNW) || !validTile(current.neighborSE) {
				*jps = append(*jps, current) // found jp!
			} 
			//if validTile(current.neighborNorth) && validTile(current.neighborEast) && validTile(current.neighborNE) {
		//		return getJumpPoint(current.neighborNE, dir)
	//		} else {return current}  //lr nil? right..? 
//			return
		}
		if dir.yDir == -1 { // nordväst
			wPath := getJumpPoint(current, Direction{-1,0})  // jp västerut?
			if wPath != nil {*jps = append(*jps, wPath)}
			if !validTile(current.neighborSW) || !validTile(current.neighborNE) {
				*jps = append(*jps, current) // found jp!
			}
			//if validTile(current.neighborWest) && validTile(current.neighborNorth) && validTile(current.neighborNW) {
			//	return getJumpPoint(current.neighborNW, dir)
		//	} else {return current}  //lr nil? right..?
			return
		}	
	}
}
*/



func getJumpPoint(current *tile, dir Direction) jp{
	curJP := jp{}
	if current.door {return jp{current, nil}}
	if dir.xDir == 0 {
		if dir.yDir == 1 { // höger
			if (!validTile(current.neighborNW) /*&& validTile(current.neighborNorth)*/) || !validTile(current.neighborSW) /*&& validTile(current.neighborSouth)*/{
				curJP.jp = current
				
				//if !validTile(current.neighborNW) && validTile(current.neighborNorth) {curJP.fn = append(curJP.fn, current.neighborNorth)}
				if !validTile(current.neighborNW) {
					if validTile(current.neighborNorth) {curJP.fn = append(curJP.fn, current.neighborNorth)}
					if validTile(current.neighborSouth) && validTile(current.neighborEast) && validTile(current.neighborSE){
						curJP.fn = append(curJP.fn, current.neighborSE)}
				}
				//if !validTile(current.neighborSW) && validTile(current.neighborSouth) {curJP.fn = append(curJP.fn, current.neighborSouth)}
				if !validTile(current.neighborSW){
					if validTile(current.neighborSouth) {curJP.fn = append(curJP.fn, current.neighborSouth)}
					if validTile(current.neighborNorth) && validTile(current.neighborEast) && validTile(current.neighborNE){
						curJP.fn = append(curJP.fn, current.neighborNE)}
				}
				return curJP  // found jp!
				
			}
			if validTile(current.neighborEast) {
				return getJumpPoint(current.neighborEast, dir)
			} else {return curJP}			
		}       // OBS: above is modified differently so far!
		if dir.yDir == -1 {  // vänster
			if (!validTile(current.neighborNE) && validTile(current.neighborNorth)) || (!validTile(current.neighborSE) && validTile(current.neighborSouth)){
				//return current  // found jp!
				curJP.jp = current
				if !validTile(current.neighborNE) && validTile(current.neighborNorth) {curJP.fn = append(curJP.fn, current.neighborNorth)}
				if !validTile(current.neighborSE) && validTile(current.neighborSouth) {curJP.fn = append(curJP.fn, current.neighborSouth)}
				return curJP  // found jp!
			}
			if validTile(current.neighborWest) {
				return getJumpPoint(current.neighborWest, dir)
			} else {return curJP}		
		}
	}
	if dir.yDir == 0 {
		if dir.xDir == 1 { // neråt
			if (!validTile(current.neighborNW) && validTile(current.neighborWest)) || (!validTile(current.neighborNE) && validTile(current.neighborEast)) {
				//return current  // found jp!
				curJP.jp = current
				if !validTile(current.neighborNW) && validTile(current.neighborWest) {curJP.fn = append(curJP.fn, current.neighborWest)}
				if !validTile(current.neighborNE) && validTile(current.neighborEast) {curJP.fn = append(curJP.fn, current.neighborEast)}
				return curJP  // found jp!
			}
			if validTile(current.neighborSouth) {
				return getJumpPoint(current.neighborSouth, dir)
			} else {return curJP}  //lr nil? right..?		
		} /*grunkagrunkagrunkagrunka
		bra kaffe
		steg 1
		kaffe bönor
		2 kaffekvarn
		3 köksvåg
                vad tycker jenny? o.O majbi om man ändå dricker ofta kan ma ju göra de godare, kan ju vara värt de, majbi kuul mhmm*/
		if dir.xDir == -1 { // uppåt
			if (!validTile(current.neighborSW) && validTile(current.neighborWest)) || (!validTile(current.neighborSE) && validTile(current.neighborEast)) {
				//return current  // found jp!
				curJP.jp = current
				if !validTile(current.neighborSW) && validTile(current.neighborWest) {curJP.fn = append(curJP.fn, current.neighborWest)}
				if !validTile(current.neighborSE) && validTile(current.neighborEast) {curJP.fn = append(curJP.fn, current.neighborEast)}
				return curJP  // found jp!
			}
			if validTile(current.neighborNorth) {
				return getJumpPoint(current.neighborNorth, dir)
//				nextJP := getJumpPoint(current.neighborNorth, dir)
//				if nextJP.jp != nil {}
			} else {return curJP}  //lr nil? right..?		
		}
	}
/* härifrån	if dir.xDir == 1 {
		if dir.yDir == 1 { // sydöst
			if (!validTile(current.neighborSW) && validTile(current.neighborSouth)) || (!validTile(current.neighborNE) && validTile(current.neighborEast)) {
				fmt.Println("found this...?")
				return current // found jp!
			} */
			/*if validTile(current.neighborSE) {
				if validTile(current.neighborSouth) {	
					if validTile(current.neighborEast) {
						if validTile(current.neighborSE) {return getJumpPoint(current.neighborSE, dir)}
					} else {return current} // found jp!
				} else if validTile(current.neighborEast)
			}*/
	/*		
			if validTile(current.neighborSouth) && validTile(current.neighborEast) && validTile(current.neighborSE) {
				return getJumpPoint(current.neighborSE, dir)
			}
			if !validTile(current.neighborSouth) || !validTile(current.neighborEast){  // TODO: strict or
				return current  // found jp!    
			} else {return nil}  //lr nil? right..?
		}


		if dir.yDir == -1 { // sydväst
			if !validTile(current.neighborNW) || !validTile(current.neighborSE) {
				return current // found jp!
			}
			if validTile(current.neighborWest) && validTile(current.neighborSouth) && validTile(current.neighborSW) {
				return getJumpPoint(current.neighborSW, dir)
			} else {return nil}  //lr nil? right..?
		}
	}
	if dir.xDir == -1 {  //här
		if dir.yDir == 1 { // nordöst
			if !validTile(current.neighborNW) || !validTile(current.neighborSE) {
				return current // found jp!
			}
			if validTile(current.neighborNorth) && validTile(current.neighborEast) && validTile(current.neighborNE) {
				return getJumpPoint(current.neighborNE, dir)
			} else {return nil}  //lr nil? right..?
		}
		if dir.yDir == -1 { // nordväst
			if !validTile(current.neighborSW) || !validTile(current.neighborNE) {
				return current // found jp!
			}
			if validTile(current.neighborWest) && validTile(current.neighborNorth) && validTile(current.neighborNW) {
				return getJumpPoint(current.neighborNW, dir)
			} else {return nil}  //lr nil? right..?
		}	
	}*/  //hit
	return curJP
}


func Whut() {
	matrix := [][]int {
		{0,0,0,0,0,0,0},
		{0,0,0,0,0,0,0},
		{0,0,0,0,0,0,2},
		{0,0,0,1,0,0,0},
		{0,0,0,1,0,0,0},
		{0,0,0,1,0,0,0}}
	/*	matrix := [][]int{
		{0, 1, 0, 1, 0, 0, 0},
		{0, 1, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0, 2}}
*/
	testmap := TileConvert(matrix)	
//	jp := getJumpPoint(&testmap[1][1], Direction{1,1})

//	jps := Jp(&testmap[1][1], Direction{1,1})
	//	fmt.Println(jp)

/*	for _, jp := range jps {
		fmt.Println(jp.xCoord, jp.yCoord)
	}
*/
	pt, ok := getPath2(&testmap, &testmap[5][1])
	if ok {
		fmt.Println("jp path:")
		for _,t := range pt {
			fmt.Println(t.xCoord, t.yCoord)
		}
	}
}



func getPath2(m *[][]tile, from *tile) ([]*tile, bool) {

	// map över jps
	var parentOf map[*tile]*tile
	parentOf = make(map[*tile]*tile)

	cq := queue{}

	for i, list := range *m {
		for j, _ := range list {
			cq.Add(&(*m)[i][j], float32(math.Inf(1)))	
		}
	}

	cq.Update(from, 0)
	
	v := float32(0)
	current := tileCost{&tile{}, &v}
	currentDir := Direction{0,0}
	
	for len(cq) != 0 && !current.tile.door {
		current = (&cq).Pop()
		fmt.Println(current.tile, "\n--")
		if current.tile.door {compactPath(parentOf, from, current.tile)}
		_, ok := parentOf[current.tile]
		if ok {
			currentDir = getDir(parentOf[current.tile], current.tile)
		} else {currentDir = Direction{0,0}}
		neighbors := /*getNeighbors(current.tile, cq)*/ getNeighborsPruned(current.tile, currentDir)
	//.	var wg sync.WaitGroup
	//.	wg.Add(len(neighbors))
	//.	var mutex = &sync.Mutex{}
		for _, neighbor := range neighbors {
		//	fmt.Println(neighbor)
	//.		go func(n *tile) {
			//.			defer wg.Done()
			n := neighbor
				
				jps := Jp(n, getDir(current.tile, n))
				for _, jp := range jps {
				//	fmt.Println(jp.jp)
					
					if jp.jp == nil {
						//	fmt.Println("No jp")
					} else {
						fmt.Println("jp", jp.jp)
	//.					mutex.Lock()
						cost := cq.costOf(current.tile) + smplCost(current.tile, jp.jp) //TODO:!
						if cost < cq.costOf(jp.jp) {
							parentOf[jp.jp] = current.tile
							cq.Update(jp.jp, cost)
							for _, n := range jp.fn {
								if n != nil {
									parentOf[n] = jp.jp
									cq.Update(n, cost + 1)  // TODO +1 stämmer nog iinte!
								}
							}
							//	fmt.Println("current: ", current.tile.xCoord, current.tile.yCoord)
							//	fmt.Println("jp: ", jp.xCoord, jp.yCoord)
							//	fmt.Println("updated", cost)}
							
						}
	//.					mutex.Unlock()	
					}
				}
	//.		}(neighbor)						

		
		}
	//.	wg.Wait()		
	}
	fmt.Println(len(parentOf))
	return compactPath(parentOf, from, current.tile)
}

func smplCost(t1 *tile, t2 *tile) float32{
	xDif := math.Max(float64(t1.xCoord), float64(t2.xCoord)) - math.Min(float64(t1.xCoord), float64(t2.xCoord))
	yDif := math.Max(float64(t1.yCoord), float64(t2.yCoord)) - math.Min(float64(t1.yCoord), float64(t2.yCoord))
	if xDif == 0 {
		//	fmt.Println("no xdif")
		return float32(yDif)}
	if yDif == 0 {
	//	fmt.Println("no ydif")
		return float32(xDif)}
	return float32(math.Sqrt(xDif*xDif + yDif*yDif))
}


// get partial path
func getPPath(m *[][]tile, from *tile, to *tile) ([]*tile, bool) {

	// map to keep track of the final path
	var parentOf map[*tile]*tile
	parentOf = make(map[*tile]*tile)

	//initialise 'costqueue', start-0, other-infinite
	costQueue := queue{}

	for i, list := range *m {
		for j, _ := range list {
			costQueue.Add(&(*m)[i][j], float32(math.Inf(1)))	
		}
	}

	costQueue.Update(from, 0)

	v := float32(0)
	current := tileCost{&tile{}, &v}
	currentDir := Direction{1,1}
	
	for len(costQueue) != 0 && !(*current.tile == *to) {//!current.tile.door {

		current = (&costQueue).Pop()
		currentDir = getDir(parentOf[current.tile], current.tile)  // for reference!


		neighbors := getNeighborsPruned(current.tile, currentDir)
		var wg sync.WaitGroup
		wg.Add(len(neighbors))
		var mutex = &sync.Mutex{}
		for _, neighbor := range neighbors {
			go func(n *tile) {			
				defer wg.Done()			
				cost := *current.cost + stepCost(*n)
				if Diagonal(current.tile, n) {cost += float32(math.Sqrt(2)) - 1}
				if n.occupied.IsWaiting() {cost += 1}
			
				mutex.Lock()
				if cost < costQueue.costOf(n) {
					
					parentOf[n] = current.tile
					costQueue.Update(n, cost)				
				}
				mutex.Unlock()
			}(neighbor)		
		}
		wg.Wait()		
	}	
	return compactPath(parentOf, from, current.tile)
}

var (  // TODO: define this a weak ago...
	n = Direction{-1,0}
	e = Direction{1,0}
	s = Direction{1,0}
	w = Direction{0,-1}

	nw = Direction{-1,-1}
	ne = Direction{-1,1}
	se = Direction{1,1}
	sw = Direction{1,-1}
	
)

type jp struct {
	jp *tile   // jp
	fn []*tile // forced neighbors from jp
}
/*
func getForcedNeighbor(current *tile, dir Direction) *tile{
	// TODO: do this automagically in other functions!!
	if dir == n {}
}*/


