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
		if dir.xDir == -1 { // går NW
			if west {neighbors = append(neighbors, current.neighborWest)}
			if north {neighbors = append(neighbors, current.neighborNorth)}
			if west && north && validTile(current.neighborNW) {neighbors = append(neighbors, current.neighborNW)}	
		} else if dir.xDir == 1 { // går SW
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
	
	x := to.xCoord - from.xCoord
	y := to.yCoord - from.yCoord
	if x > 1 {x = 1} else if x < 0 {x = -1}
	if y > 1 {y = 1} else if y < 0  {y = -1}
	return Direction {x,y}//{math.Mod(to.xCoord - from.xCoord, ), to.yCoord - from.yCoord}
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

func (t *tile) neighbor(dir Direction) *tile{
	if dir == n {return t.neighborNorth}
	if dir == e {return t.neighborEast}
	if dir == s {return t.neighborSouth}
	if dir == w {return t.neighborWest}

	if dir == nw {return t.neighborNW}
	if dir == ne {return t.neighborNE}
	if dir == se {return t.neighborSE}
	if dir == sw {return t.neighborSW}

	return nil
}

func Jp(current *tile, dir Direction) []jp {//[]*tile {
//func Jp(current *tile, dir Direction, parentOf map[*tile]*tile) []jp {//[]*tile {
	jps := []jp{}
	//----
//	if current.occupied != nil{
//		jps = append(jps, jp{current, []*tile{}})
		//current.occupied.plan = compactPath2(parentOf, current)
		//current.occupied.dir = getDir(parent[current], current) //OBS: osäkert..maybe tvärtom?
//	}
	//----
	
//----	if current.door {return []jp{jp{current, []*tile{}}}}

	if (dir.xDir == 0 && dir.yDir != 0) ||(dir.yDir == 0 && dir.xDir != 0) {
		tmpJP := getJumpPoint(current, dir)
		for tmpJP.jp != nil && tmpJP.jp.occupied != nil {
			jps = append(jps, tmpJP)
			tmpJP = getJumpPoint(tmpJP.jp.neighbor(dir), dir)
	
		}
		
		return append(jps, tmpJP)
	}
	/*----
	if dir.xDir == 0 {
		if dir.yDir == 1 {return []jp{getJumpPoint(current, dir)}} // höger
		if dir.yDir == -1 {return []jp{getJumpPoint(current, dir)}} // vänster		
	}
	if dir.yDir == 0 {
		if dir.xDir == 1 {return []jp{getJumpPoint(current, dir)}} // neråt
		if dir.xDir == -1 {return []jp{getJumpPoint(current, dir)}} // uppåt
	}----*/

	for current != nil {
	//	fmt.Println("check", dir.xDir, 0)
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
		//----
	//	if current != nil && current.occupied != nil{
	//		jps = append(jps, jp{current, []*tile{}})
			//current.occupied.plan = compactPath2(parentOf, current)
			//current.occupied.dir = getDir(parent[current], current) //OBS: osäkert..maybe tvärtom?
	//	}
		//----
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
	if current == nil {return curJP}
	if current.door {
	//	fmt.Println("\n\n\n\n", current.xCoord, current.yCoord, "\n\n ")
		return jp{current, nil}}
//	if current.occupied != nil {return jp{current, nil}}
	
	if dir.xDir == 0 {
		if dir.yDir == 1 { // höger
			return eastJP(current)
		}       // OBS: above is modified differently so far!
		if dir.yDir == -1 {  // vänster
			return westJP(current)		
		}
	}
	if dir.yDir == 0 {
		if dir.xDir == 1 { // neråt
			return(southJP(current))		
		} /*grunkagrunkagrunkagrunka
		bra kaffe
		steg 1
		kaffe bönor
		2 kaffekvarn
		3 köksvåg
                vad tycker jenny? o.O majbi om man ändå dricker ofta kan ma ju göra de godare, kan ju vara värt de, majbi kuul mhmm*/
		if dir.xDir == -1 { // uppåt
			return northJP(current)
		/*	if (!validTile(current.neighborSW) && validTile(current.neighborWest)) || (!validTile(current.neighborSE) && validTile(current.neighborEast)) {
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
			} else {return curJP}  //lr nil? right..?	 */	
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

func eastJP(current *tile) jp{  //fns fixed(i think..)
//	fmt.Println("I'm here!")
//	fmt.Println(current.xCoord, current.yCoord)
	curJP := jp{}
	if (!validTile(current.neighborNW) /*&& validTile(current.neighborNorth)*/) || !validTile(current.neighborSW) /*&& validTile(current.neighborSouth)*/{
	//	fmt.Println("jp?")
	//	curJP.jp = current	
		//if !validTile(current.neighborNW) && validTile(current.neighborNorth) {curJP.fn = append(curJP.fn, current.neighborNorth)}
		if !validTile(current.neighborNW) {
			if validTile(current.neighborNorth) {
				curJP.fn = append(curJP.fn, current.neighborNorth)
				if validTile(current.neighborEast) && validTile(current.neighborNE) {
					curJP.fn = append(curJP.fn, current.neighborNE)} //check! (2/8)
			}
			if validTile(current.neighborSouth) && validTile(current.neighborEast) && validTile(current.neighborSE) && !validTile(current.neighborSW){
				curJP.fn = append(curJP.fn, current.neighborSE)}			
		}
		//if !validTile(current.neighborSW) && validTile(current.neighborSouth) {curJP.fn = append(curJP.fn, current.neighborSouth)}
		if !validTile(current.neighborSW){
			if validTile(current.neighborSouth) {
				curJP.fn = append(curJP.fn, current.neighborSouth)
				if validTile(current.neighborEast) && validTile(current.neighborSE) {
					curJP.fn = append(curJP.fn, current.neighborSE)} //TODO: prettify!! implement for similair funcs
			}
			if validTile(current.neighborNorth) && validTile(current.neighborEast) && validTile(current.neighborNE) && !validTile(current.neighborNW){
				curJP.fn = append(curJP.fn, current.neighborNE)}

		}
	//	return curJP  // found jp!

		if len(curJP.fn) > 0 {
			curJP.jp = current
			//if validTile(current.neighborEast) {curJP.fn = append(curJP.fn, current.neighborEast)}
			return curJP}
	}
	if validTile(current.neighborEast) {
	//	fmt.Println("continue")
	//	if current.neighborEast.door {return jp{current, nil}}
		return getJumpPoint(current.neighborEast, e)
	} else {return curJP}	
}

func westJP(current *tile) jp{ // check 4/8 done!
	curJP := jp{}
//	fmt.Println(current)
	if current.door {
		curJP.jp = current
		return curJP}
	if !validTile(current.neighborNE) || !validTile(current.neighborSE) {


		if !validTile(current.neighborNE) {
			if validTile(current.neighborNorth) {
				curJP.fn = append(curJP.fn, current.neighborNorth)
				if validTile(current.neighborWest) && validTile(current.neighborNW) {
					curJP.fn = append(curJP.fn, current.neighborNW)}
			}
			if validTile(current.neighborWest) && validTile(current.neighborSouth) && validTile(current.neighborSW) && !validTile(current.neighborSE){
				curJP.fn = append(curJP.fn, current.neighborSW)}
		}
		if !validTile(current.neighborSE) {
			if validTile(current.neighborSouth) {
				//fmt.Println("??")
				curJP.fn = append(curJP.fn, current.neighborSouth)
				if validTile(current.neighborWest) && validTile(current.neighborSW) {
					curJP.fn = append(curJP.fn, current.neighborSW)}
			}
			if validTile(current.neighborNorth) && validTile(current.neighborWest) && validTile(current.neighborNW) && !validTile(current.neighborNE){
				curJP.fn = append(curJP.fn, current.neighborNW)}
		}
		//if (!validTile(current.neighborNW) && validTile(current.neighborWest)) || (!validTile(current.neighborNE) && validTile(current.neighborEast)) {
		//return current  // found jp!
		//	if !validTile(current.neighborNW) && validTile(current.neighborWest) {curJP.fn = append(curJP.fn, current.neighborWest)}
	//	if !validTile(current.neighborNE) && validTile(current.neighborEast) {curJP.fn = append(curJP.fn, current.neighborEast)}
		if len(curJP.fn) > 0 {
			curJP.jp = current
			// TODO?: if validTile(current.neighborNorth) {curJP.append(current.neighborNorth)}
			if validTile(current.neighborWest) {curJP.fn = append(curJP.fn, current.neighborWest)}
			return curJP}
		//	return curJP  // found jp!
	}
	if validTile(current.neighborWest) {
		return getJumpPoint(current.neighborWest, w) // TODO: 'southJP' instead of getjp
	} else {
		return curJP}  //lr nil? right..?
}

func southJP(current *tile) jp{	
	curJP := jp{}

	if !validTile(current.neighborNW) || !validTile(current.neighborNE) {
		//curJP.jp = current
		if !validTile(current.neighborNW) {
			if validTile(current.neighborWest) {
				curJP.fn = append(curJP.fn, current.neighborWest)
			}

			if validTile(current.neighborSouth) && validTile(current.neighborEast) && validTile(current.neighborSE) && !validTile(current.neighborNE){
				curJP.fn = append(curJP.fn, current.neighborSE)}
		}
		if !validTile(current.neighborNE) {
			if validTile(current.neighborEast) {curJP.fn = append(curJP.fn, current.neighborEast)}
			if validTile(current.neighborSouth) && validTile(current.neighborWest) && validTile(current.neighborSW) && !validTile(current.neighborNW){
			
				curJP.fn = append(curJP.fn, current.neighborSW)}
		}
		//if (!validTile(current.neighborNW) && validTile(current.neighborWest)) || (!validTile(current.neighborNE) && validTile(current.neighborEast)) {
		//return current  // found jp!
		//	if !validTile(current.neighborNW) && validTile(current.neighborWest) {curJP.fn = append(curJP.fn, current.neighborWest)}
		//	if !validTile(current.neighborNE) && validTile(current.neighborEast) {curJP.fn = append(curJP.fn, current.neighborEast)}

		if len(curJP.fn) > 0 {
			curJP.jp = current
			// TODO?: if validTile(current.neighborNorth) {curJP.append(current.neighborNorth)}
			if validTile(current.neighborSouth) {curJP.fn = append(curJP.fn, current.neighborSouth)}
			return curJP}
		//	return curJP  // found jp!
	}
	if validTile(current.neighborSouth) {

		return getJumpPoint(current.neighborSouth, s) // TODO: 'southJP' instead of getjp
	}
/*	if current.neighborWest.door {
		curJP.fn = append(curJP.fn, current.neighborWest)
	}
	if current.neighborEast.door { curJP.fn = append(curJP.fn, current.neighborEast)
	} //else {return curJP}  //lr nil? right..?
	if len(curJP.fn) > 0 {
		curJP.jp = current
	} */
	return curJP
}


func northJP(current *tile) jp{

	curJP := jp{}
	if !validTile(current.neighborSE) || !validTile(current.neighborSW) {
	//	curJP.jp = current
		if !validTile(current.neighborSE) {
			if validTile(current.neighborEast) {curJP.fn = append(curJP.fn, current.neighborEast)}
			if validTile(current.neighborWest) && validTile(current.neighborNorth) && validTile(current.neighborNW) && !validTile(current.neighborSW){
				curJP.fn = append(curJP.fn, current.neighborNW)}
		}
		if !validTile(current.neighborSW) {
			if validTile(current.neighborWest) {curJP.fn = append(curJP.fn, current.neighborWest)}
			if validTile(current.neighborNorth) && validTile(current.neighborEast) && validTile(current.neighborNE) && !validTile(current.neighborSE){
				curJP.fn = append(curJP.fn, current.neighborNE)}
		}
		//if (!validTile(current.neighborNW) && validTile(current.neighborWest)) || (!validTile(current.neighborNE) && validTile(current.neighborEast)) {
		//return current  // found jp!
		//	if !validTile(current.neighborNW) && validTile(current.neighborWest) {curJP.fn = append(curJP.fn, current.neighborWest)}
	//	if !validTile(current.neighborNE) && validTile(current.neighborEast) {curJP.fn = append(curJP.fn, current.neighborEast)}

		if len(curJP.fn) > 0 {
			curJP.jp = current
			// TODO?: if validTile(current.neighborNorth) {curJP.append(current.neighborNorth)}
			if validTile(current.neighborNorth) {curJP.fn = append(curJP.fn, current.neighborNorth)}
			return curJP}
		//return curJP  // found jp!
	}
	if validTile(current.neighborNorth) {
		return getJumpPoint(current.neighborNorth, n) // TODO: 'southJP' instead of getjp
	} else {return curJP}  //lr nil? right..?
}


func Whut() {
	matrix := [][]int {
		{0,0,0,0,0,0,0},
		{0,0,1,0,0,0,0},
		{1,1,1,1,0,0,0},
		{0,0,0,1,0,0,0},
		{0,0,0,0,0,0,0},
		{2,0,0,1,1,0,0}}

	testmap := TileConvert(matrix)	

	path, _ := getPath2(&testmap, &testmap[0][6])
	printPath(path)
}



func getPath2(m *[][]tile, from *tile) ([]*tile, bool) {
	
	// map över jp
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
		_, ok := parentOf[current.tile]
		if ok {
			currentDir = getDir(parentOf[current.tile], current.tile)
		} else {currentDir = Direction{0,0}}
		neighbors := /*getNeighbors(current.tile, cq)*/ getNeighborsPruned(current.tile, currentDir)

		var wg sync.WaitGroup
		wg.Add(len(neighbors))
		var mutex = &sync.Mutex{}
		for _, neighbor := range neighbors {
			//	fmt.Println("neighbor",neighbor.xCoord, neighbor.yCoord)
			go func(n *tile) {
				defer wg.Done()
			//	n := neighbor

				jps := Jp(n, getDir(current.tile, n))
				for _, jp := range jps {

					if jp.jp != nil {
						
						mutex.Lock()
						//	if cq.costOf(current.tile) < 0 {fmt.Println("wtf?", current.tile, cq.costOf(current.tile))}
						//cost := cq.costOf(current.tile) + smplCost(current.tile, jp.jp) //TODO:!
						cost := *current.cost + smplCost(current.tile, jp.jp) //TODO:!
						//	if jp.jp == GetTile(*m, 5, 0) {fmt.Println("\nCOST: ",cost)}
						//	if cost < 0 {fmt.Println("neg cost?:",cost)}
					//	fmt.Println("jp", jp.jp)
					//	fmt.Println("jpcost?", cost)
					//	fmt.Println("whut?", cq.costOf(jp.jp))
						if cost < cq.costOf(jp.jp) {
							parentOf[jp.jp] = current.tile
							cq.Update(jp.jp, cost)
						//	fmt.Println("whut?", cq.costOf(jp.jp))
							for _, n := range jp.fn {
							//	fmt.Println("fn", n)
								fnCost := cost + smplCost(jp.jp, n)
							//	if fnCost < 1 {
								//	fmt.Println("neg cost?:",fnCost)
								//	fmt.Println(jp.jp)
								//	fmt.Println(n)
								//	fmt.Println(parentOf[jp.jp])
							//	}
								//fmt.Println("fn?", fnCost)
								if n != nil && fnCost < cq.costOf(n)  {
									parentOf[n] = jp.jp
									cq.Update(n, fnCost) 
								}
							}
						}
						mutex.Unlock()	
					}
				}
			}(neighbor)						
		
		}
		wg.Wait()		
	}
	return compactPath(parentOf, from, current.tile)
}

func smplCost(t1 *tile, t2 *tile) float32{
//	fmt.Println(t1)
//	fmt.Println(t2)
	xDif := math.Max(float64(t1.xCoord), float64(t2.xCoord)) - math.Min(float64(t1.xCoord), float64(t2.xCoord))
	yDif := math.Max(float64(t1.yCoord), float64(t2.yCoord)) - math.Min(float64(t1.yCoord), float64(t2.yCoord))
	if xDif == 0 {

		return float32(yDif)}
	if yDif == 0 {

		return float32(xDif)}

//	fmt.Println("other", float32(math.Sqrt(xDif*xDif + yDif*yDif)))
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
	e = Direction{0,1}
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



func getPath3(m *[][]tile, from []*tile) {
	
	// map över jp
	var parentOf map[*tile]*tile
	parentOf = make(map[*tile]*tile)

	cq := queue{}

	for i, list := range *m {
		for j, _ := range list {
			cq.Add(&(*m)[i][j], float32(math.Inf(1)))		
		}
	}

	for _, f := range from {
	//	fmt.Println("wtf f?:", f)
		cq.Update(f, 0)
	}
	v := float32(0)
	current := tileCost{&tile{}, &v}
	currentDir := Direction{0,0}	
	for len(cq) != 0 {//!current.tile.door {
		current = (&cq).Pop()
	//	fmt.Println("current:", current.tile.xCoord, current.tile.yCoord)
		if *current.cost == float32(math.Inf(1)) {break}
		_, ok := parentOf[current.tile]
		if ok {
			currentDir = getDir(parentOf[current.tile], current.tile)
		} else {currentDir = Direction{0,0}}
		neighbors := /*getNeighbors(current.tile, cq)*/ getNeighborsPruned(current.tile, currentDir)

		var wg sync.WaitGroup
		wg.Add(len(neighbors))
		var mutex = &sync.Mutex{}
		for _, neighbor := range neighbors {
			
			go func(n *tile) {
				defer wg.Done()
				//	n := neighbor

				jps := JpInit(n, getDir(current.tile, n))
				for _, jp := range jps {
				

					if jp.jp != nil {						
					//	fmt.Println("jp:", jp.jp.xCoord, jp.jp.yCoord)	
						mutex.Lock()
						cost := *current.cost + smplCost(current.tile, jp.jp) //TODO:!
						if cost < cq.costOf(jp.jp) {
							parentOf[jp.jp] = current.tile
							cq.Update(jp.jp, cost)

							if jp.jp.occupied != nil {
								setPlan(parentOf, jp.jp)  // + cq maybe??
							}
						}
						for _, n := range jp.fn {

							fnCost := cost + smplCost(jp.jp, n)
							if n != nil && fnCost < cq.costOf(n)  {							
								parentOf[n] = jp.jp
								cq.Update(n, fnCost)

								if n.occupied != nil {
									setPlan(parentOf, n)  // + cq maybe??
								}
							}
						}
						mutex.Unlock()	
					}
				
				}
				//}
			}(neighbor)						

			
		}
		wg.Wait()		
	}
	//	return compactPath(parentOf, from, current.tile)
}


//TODO: func compactpath2
func setPlan(parentOf map[*tile]*tile, pers *tile) {
	path := []*tile{pers}

	current := pers
	ok := true
	curDir := Direction{} // getDir(current, parentOf[current])
	lastDir := Direction{0,0}
	for !current.door {//current.xCoord != from.xCoord || current.yCoord != from.yCoord {
		//	path = append([]*tile{parentOf[current]}, path...)
		curDir = getDir(current, parentOf[current])
		if curDir == lastDir {
			path = path[:len(path)-1]
		}
		path = append(path, parentOf[current])	
			//	fmt.Println(current)
		current, ok = parentOf[current]
		lastDir = curDir
	
		if !ok {
			//return nil, false
			pers.occupied.plan = []*tile{}
			break
		}
	}
	if !ok {pers.occupied.plan = []*tile{}
	}else {

	//	printPath(pers.occupied.plan)
		pers.occupied.plan = path}
//	return path, true
}



func getJumpPointInit(current *tile, dir Direction) jp{
	curJP := jp{}
	if current == nil || !validTile(current) {return curJP}
	//if current.door {
	//	fmt.Println("\n\n\n\n", current.xCoord, current.yCoord, "\n\n ")
	//	return jp{current, nil}}
	if current.occupied != nil {return jp{current, nil}}
	
	if dir.xDir == 0 {
		if dir.yDir == 1 { // höger
			return eastJPInit(current)
		}       // OBS: above is modified differently so far!
		if dir.yDir == -1 {  // vänster
			return westJPInit(current)		
		}
	}
	if dir.yDir == 0 {
		if dir.xDir == 1 { // neråt
			return(southJPInit(current))		
		}
		if dir.xDir == -1 { // uppåt
			return northJPInit(current)	
		}
	}
	return curJP
}


func JpInit(current *tile, dir Direction) []jp {//[]*tile {
//func Jp(current *tile, dir Direction, parentOf map[*tile]*tile) []jp {//[]*tile {
	jps := []jp{}
	//----
	if current.occupied != nil{
		jps = append(jps, jp{current, []*tile{}})
		//current.occupied.plan = compactPath2(parentOf, current)
		//current.occupied.dir = getDir(parent[current], current) //OBS: osäkert..maybe tvärtom?
	}
	//----
	
//----	if current.door {return []jp{jp{current, []*tile{}}}}

	if (dir.xDir == 0 && dir.yDir != 0) ||(dir.yDir == 0 && dir.xDir != 0) {
		tmpJP := getJumpPointInit(current, dir)
		for tmpJP.jp != nil && tmpJP.jp.occupied != nil {
			jps = append(jps, tmpJP)
			tmpJP = getJumpPointInit(tmpJP.jp.neighbor(dir), dir)
	
		}
		
		return append(jps, tmpJP)
	}

	for current != nil {
	//	fmt.Println("check", dir.xDir, 0)
		jpX := getJumpPointInit(current, Direction{dir.xDir, 0})
		if jpX.jp != nil {
			tmpJP := jp{current, []*tile{jpX.jp}}
			jps = append(jps, tmpJP)
		//	jps = append(jps, jp{})
		}
		jpY := getJumpPointInit(current, Direction{0, dir.yDir})
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
		//----
		if current != nil && current.occupied != nil{
			jps = append(jps, jp{current, []*tile{}})
			//current.occupied.plan = compactPath2(parentOf, current)
			//current.occupied.dir = getDir(parent[current], current) //OBS: osäkert..maybe tvärtom?
		}
		//----
	}
	return jps
}




func eastJPInit(current *tile) jp{  //fns fixed(i think..)
//	fmt.Println("I'm here!")
//	fmt.Println(current.xCoord, current.yCoord)
	curJP := jp{}
	if (!validTile(current.neighborNW) /*&& validTile(current.neighborNorth)*/) || !validTile(current.neighborSW) /*&& validTile(current.neighborSouth)*/{
	//	fmt.Println("jp?")
	//	curJP.jp = current	
		//if !validTile(current.neighborNW) && validTile(current.neighborNorth) {curJP.fn = append(curJP.fn, current.neighborNorth)}
		if !validTile(current.neighborNW) {
			if validTile(current.neighborNorth) {
				curJP.fn = append(curJP.fn, current.neighborNorth)
				if validTile(current.neighborEast) && validTile(current.neighborNE) {
					curJP.fn = append(curJP.fn, current.neighborNE)} //check! (2/8)
			}
			if validTile(current.neighborSouth) && validTile(current.neighborEast) && validTile(current.neighborSE) && !validTile(current.neighborSW){
				curJP.fn = append(curJP.fn, current.neighborSE)}			
		}
		//if !validTile(current.neighborSW) && validTile(current.neighborSouth) {curJP.fn = append(curJP.fn, current.neighborSouth)}
		if !validTile(current.neighborSW){
			if validTile(current.neighborSouth) {
				curJP.fn = append(curJP.fn, current.neighborSouth)
				if validTile(current.neighborEast) && validTile(current.neighborSE) {
					curJP.fn = append(curJP.fn, current.neighborSE)} //TODO: prettify!! implement for similair funcs
			}
			if validTile(current.neighborNorth) && validTile(current.neighborEast) && validTile(current.neighborNE) && !validTile(current.neighborNW){
				curJP.fn = append(curJP.fn, current.neighborNE)}

		}
	//	return curJP  // found jp!

		if len(curJP.fn) > 0 {
			curJP.jp = current
			//if validTile(current.neighborEast) {curJP.fn = append(curJP.fn, current.neighborEast)}
			return curJP}
	}
	if validTile(current.neighborEast) {
	//	fmt.Println("continue")
	//	if current.neighborEast.door {return jp{current, nil}}
		return getJumpPointInit(current.neighborEast, e)
	} else {return curJP}	
}

func westJPInit(current *tile) jp{ // check 4/8 done!
	curJP := jp{}
//	fmt.Println(current)
	if current.door {
		curJP.jp = current
		return curJP}
	if !validTile(current.neighborNE) || !validTile(current.neighborSE) {


		if !validTile(current.neighborNE) {
			if validTile(current.neighborNorth) {
				curJP.fn = append(curJP.fn, current.neighborNorth)
				if validTile(current.neighborWest) && validTile(current.neighborNW) {
					curJP.fn = append(curJP.fn, current.neighborNW)}
			}
			if validTile(current.neighborWest) && validTile(current.neighborSouth) && validTile(current.neighborSW) && !validTile(current.neighborSE){
				curJP.fn = append(curJP.fn, current.neighborSW)}
		}
		if !validTile(current.neighborSE) {
			if validTile(current.neighborSouth) {
				//fmt.Println("??")
				curJP.fn = append(curJP.fn, current.neighborSouth)
				if validTile(current.neighborWest) && validTile(current.neighborSW) {
					curJP.fn = append(curJP.fn, current.neighborSW)}
			}
			if validTile(current.neighborNorth) && validTile(current.neighborWest) && validTile(current.neighborNW) && !validTile(current.neighborNE){
				curJP.fn = append(curJP.fn, current.neighborNW)}
		}
		//if (!validTile(current.neighborNW) && validTile(current.neighborWest)) || (!validTile(current.neighborNE) && validTile(current.neighborEast)) {
		//return current  // found jp!
		//	if !validTile(current.neighborNW) && validTile(current.neighborWest) {curJP.fn = append(curJP.fn, current.neighborWest)}
	//	if !validTile(current.neighborNE) && validTile(current.neighborEast) {curJP.fn = append(curJP.fn, current.neighborEast)}
		if len(curJP.fn) > 0 {
			curJP.jp = current
			// TODO?: if validTile(current.neighborNorth) {curJP.append(current.neighborNorth)}
			if validTile(current.neighborWest) {curJP.fn = append(curJP.fn, current.neighborWest)}
			return curJP}
		//	return curJP  // found jp!
	}
	if validTile(current.neighborWest) {
		return getJumpPointInit(current.neighborWest, w) // TODO: 'southJP' instead of getjp
	} else {
		return curJP}  //lr nil? right..?
}

func southJPInit(current *tile) jp{	
	curJP := jp{}

	if !validTile(current.neighborNW) || !validTile(current.neighborNE) {
		//curJP.jp = current
		if !validTile(current.neighborNW) {
			if validTile(current.neighborWest) {
				curJP.fn = append(curJP.fn, current.neighborWest)
			}

			if validTile(current.neighborSouth) && validTile(current.neighborEast) && validTile(current.neighborSE) && !validTile(current.neighborNE){
				curJP.fn = append(curJP.fn, current.neighborSE)}
		}
		if !validTile(current.neighborNE) {
			if validTile(current.neighborEast) {curJP.fn = append(curJP.fn, current.neighborEast)}
			if validTile(current.neighborSouth) && validTile(current.neighborWest) && validTile(current.neighborSW) && !validTile(current.neighborNW){
			
				curJP.fn = append(curJP.fn, current.neighborSW)}
		}
		if len(curJP.fn) > 0 {
			curJP.jp = current
			// TODO?: if validTile(current.neighborNorth) {curJP.append(current.neighborNorth)}
			if validTile(current.neighborSouth) {curJP.fn = append(curJP.fn, current.neighborSouth)}
			return curJP}
		//	return curJP  // found jp!
	}
	if validTile(current.neighborSouth) {

		return getJumpPointInit(current.neighborSouth, s) // TODO: 'southJP' instead of getjp
	}
	return curJP
}


func northJPInit(current *tile) jp{

	curJP := jp{}
	if !validTile(current.neighborSE) || !validTile(current.neighborSW) {
	//	curJP.jp = current
		if !validTile(current.neighborSE) {
			if validTile(current.neighborEast) {curJP.fn = append(curJP.fn, current.neighborEast)}
			if validTile(current.neighborWest) && validTile(current.neighborNorth) && validTile(current.neighborNW) && !validTile(current.neighborSW){
				curJP.fn = append(curJP.fn, current.neighborNW)}
		}
		if !validTile(current.neighborSW) {
			if validTile(current.neighborWest) {curJP.fn = append(curJP.fn, current.neighborWest)}
			if validTile(current.neighborNorth) && validTile(current.neighborEast) && validTile(current.neighborNE) && !validTile(current.neighborSE){
				curJP.fn = append(curJP.fn, current.neighborNE)}
		}
		if len(curJP.fn) > 0 {
			curJP.jp = current
			// TODO?: if validTile(current.neighborNorth) {curJP.append(current.neighborNorth)}
			if validTile(current.neighborNorth) {curJP.fn = append(curJP.fn, current.neighborNorth)}
			return curJP}
	}
	if validTile(current.neighborNorth) {
		return getJumpPointInit(current.neighborNorth, n) // TODO: 'southJP' instead of getjp
	} else {return curJP}  //lr nil? right..?
}



func (p *Person) redirect() bool {
	if p.dir == se {return p.reDirSE()}
	if p.dir == sw {return p.reDirSW()}
	if p.dir == nw {return p.reDirNW()}
	if p.dir == ne {return p.reDirNE()}
	return false
	
}

func (p *Person) reDirSE() bool{
	//TODO: continuieieie!!
//	if p.dir != se {return false}
	current := p.currentTile()
	//	if current.neighborNorth.occupied == nil { }
	if p.moveTo(current.neighborSouth) {
		next := current.followDir(p.dir)
		if next != nil {
			p.plan = append([]*tile{next}, p.plan...)
			for next != nil && next.xCoord != p.plan[1].xCoord{
				p.plan[0] = next
				next = next.followDir(p.dir)
			}
			tmp := p.plan[0]
		//	p.plan = append(p.plan[1:0])
			p.plan[0] = tmp.neighborEast
//			p.plan = append([]*tile{tmp}, p.plan...)
			p.plan = append([]*tile{p.plan[0].neighborWest}, p.plan[1:]...) //west??


			//([]*tile{p.plan[0], p.plan[0].neighborWest, p.plan[:]})
			
		} else {
			p.plan = append([]*tile{p.currentTile().neighborEast}, p.plan...)
		}
		
		return true
	} else 	if p.moveTo(current.neighborEast) {
		next := current.followDir(p.dir)
		if next != nil {
			p.plan = append([]*tile{next}, p.plan...)
			for next != nil && next.xCoord != p.plan[1].xCoord {
				p.plan[0] = next
				next = next.followDir(p.dir)
			}
			tmp := p.plan[0]
			//	p.plan = append(p.plan[1:0])
			p.plan[0] = tmp.neighborSouth
			//			p.plan = append([]*tile{tmp}, p.plan...)
			p.plan = append([]*tile{p.plan[0].neighborNorth}, p.plan[1:]...)


			//([]*tile{p.plan[0], p.plan[0].neighborWest, p.plan[:]})
			
		} else {
			p.plan = append([]*tile{p.currentTile().neighborSouth}, p.plan...)
		}
		return true
	}
	return false
}


func (p *Person) reDirSW() bool{
	//TODO: continuieieie!!
//	if p.dir != sw {return false}
	current := p.currentTile()
	//	if current.neighborNorth.occupied == nil { }
	if p.moveTo(current.neighborSouth) {
		next := current.followDir(p.dir)
		if next != nil {
			p.plan = append([]*tile{next}, p.plan...)
			for next != nil && next.xCoord != p.plan[1].xCoord {
				p.plan[0] = next
				next = next.followDir(p.dir)
			}
			tmp := p.plan[0]
		//	p.plan = append(p.plan[1:0])
			p.plan[0] = tmp.neighborWest
//			p.plan = append([]*tile{tmp}, p.plan...)
			p.plan = append([]*tile{p.plan[0].neighborEast}, p.plan[1:]...) 


			//([]*tile{p.plan[0], p.plan[0].neighborWest, p.plan[:]})
			
		} else {
			p.plan = append([]*tile{p.currentTile().neighborWest}, p.plan...)
		}
		
		return true
	} else 	if p.moveTo(current.neighborWest) {
		next := current.followDir(p.dir)
		if next != nil {
			p.plan = append([]*tile{next}, p.plan...)
			for next != nil && next.xCoord != p.plan[1].xCoord {
				p.plan[0] = next
				next = next.followDir(p.dir)
			}
			tmp := p.plan[0]
			//	p.plan = append(p.plan[1:0])
			p.plan[0] = tmp.neighborSouth
			//			p.plan = append([]*tile{tmp}, p.plan...)
			p.plan = append([]*tile{p.plan[0].neighborNorth}, p.plan[1:]...)


			//([]*tile{p.plan[0], p.plan[0].neighborWest, p.plan[:]})
			
		} else {
			p.plan = append([]*tile{p.currentTile().neighborSouth}, p.plan...)
		}
		return true
	}
	return false
}


func (p *Person) reDirNW() bool{
	//TODO: continuieieie!!
//	if p.dir != se {return false}
	current := p.currentTile()
	//	if current.neighborNorth.occupied == nil { }
	if p.moveTo(current.neighborNorth) {
		next := current.followDir(p.dir)
		if next != nil {
			p.plan = append([]*tile{next}, p.plan...)
			for next != nil && next.xCoord != p.plan[1].xCoord{
				p.plan[0] = next
				next = next.followDir(p.dir)
			}
			tmp := p.plan[0]
		//	p.plan = append(p.plan[1:0])
			p.plan[0] = tmp.neighborWest
//			p.plan = append([]*tile{tmp}, p.plan...)
			p.plan = append([]*tile{p.plan[0].neighborEast}, p.plan[1:]...) //west??


			//([]*tile{p.plan[0], p.plan[0].neighborWest, p.plan[:]})
			
		} else {
			p.plan = append([]*tile{p.currentTile().neighborWest}, p.plan...)
		}
		
		return true
	} else 	if p.moveTo(current.neighborWest) {
		next := current.followDir(p.dir)
		if next != nil {
			p.plan = append([]*tile{next}, p.plan...)
			for next != nil && next.xCoord != p.plan[1].xCoord {
				p.plan[0] = next
				next = next.followDir(p.dir)
			}
			tmp := p.plan[0]
			//	p.plan = append(p.plan[1:0])
			p.plan[0] = tmp.neighborNorth
			//			p.plan = append([]*tile{tmp}, p.plan...)
			p.plan = append([]*tile{p.plan[0].neighborSouth}, p.plan[1:]...)


			//([]*tile{p.plan[0], p.plan[0].neighborWest, p.plan[:]})
			
		} else {
			p.plan = append([]*tile{p.currentTile().neighborNorth}, p.plan...)
		}
		return true
	}
	return false
}

func (p *Person) reDirNE() bool{
	//TODO: continuieieie!!
//	if p.dir != se {return false}
	current := p.currentTile()
	//	if current.neighborNorth.occupied == nil { }
	if p.moveTo(current.neighborNorth) {
		next := current.followDir(p.dir)
		if next != nil {
			p.plan = append([]*tile{next}, p.plan...)
			for next != nil && next.xCoord != p.plan[1].xCoord {
				p.plan[0] = next
				next = next.followDir(p.dir)
			}
			tmp := p.plan[0]
		//	p.plan = append(p.plan[1:0])
			p.plan[0] = tmp.neighborEast
//			p.plan = append([]*tile{tmp}, p.plan...)
			p.plan = append([]*tile{p.plan[0].neighborWest}, p.plan[1:]...) //west??


			//([]*tile{p.plan[0], p.plan[0].neighborWest, p.plan[:]})
			
		} else {
			p.plan = append([]*tile{p.currentTile().neighborEast}, p.plan...)
		}
		
		return true
	} else 	if p.moveTo(current.neighborEast) {
		next := current.followDir(p.dir)
		if next != nil {
			p.plan = append([]*tile{next}, p.plan...)
			for next != nil && next.xCoord != p.plan[1].xCoord {
				p.plan[0] = next
				next = next.followDir(p.dir)
			}
			tmp := p.plan[0]
			//	p.plan = append(p.plan[1:0])
			p.plan[0] = tmp.neighborNorth
			//			p.plan = append([]*tile{tmp}, p.plan...)
			p.plan = append([]*tile{p.plan[0].neighborSouth}, p.plan[1:]...)


			//([]*tile{p.plan[0], p.plan[0].neighborWest, p.plan[:]})
			
		} else {
			p.plan = append([]*tile{p.currentTile().neighborNorth}, p.plan...)
		}
		return true
	}
	return false
}



// In Progress!!!!
func UpdateParentOf(m *[][]*tile, parentOf map[*tile]*tile, newBad []*tile) {
	for k, v := range parentOf {
		for _, bad := range newBad {
			if inBetween(k, v, bad) {
				
			}
		}
	}
}

func inBetween(t1 *tile, t2 *tile, between *tile) bool{
	dir1 := getDir(t1, between)
	dir2 := getDir(between, t2)
//	dir3 := getDir(t1, t2)
	return dir1 == dir2 
}

