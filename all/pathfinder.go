package main

import (
	"fmt"
	"math"
	"sync"
)

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

	for len(costQueue) != 0 && !current.tile.door {
		current = (&costQueue).Pop()
		neighbors := getNeighbors(current.tile, costQueue)
		var wg sync.WaitGroup
		wg.Add(len(neighbors))
		var mutex = &sync.Mutex{}
		for _, neighbor := range neighbors {		
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
