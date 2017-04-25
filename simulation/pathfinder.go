package main


import	"fmt"

func getPath(m *[][]tile, from tile, to tile) ([]tile, bool){

	// map to keep track of the final path
	var parentOf map[tile]tile
	parentOf = make(map[tile]tile)
	//initialise 'costqueue', start-0, other-infinite
	costQueue := queue{}

	for _,list := range *m {
		for _, t := range list {
			costQueue.Add(t, 100)   // 10~infinite	
		}
	}
	costQueue.Add(from, 0)


	//slice of yet to check-tiles
	checkedQueue := costQueue

	current := tileCost{}
	//essential loop
	for len(costQueue) != 0 && current.tile != to{
		current = (&costQueue).Pop()
		//loop through neighbours of current tile
		neighbors := getNeighbors(&(current.tile))
		for _, neighbor := range neighbors {		
			cost := current.cost + 1 // TODO: 1 default cost improve! depending on heat, smoke etc
			if cost < costQueue.costOf(neighbor) {
				parentOf[neighbor] = current.tile
				costQueue.Update(neighbor, cost)					
			}
		}
		checkedQueue.AddTC(current)
		costQueue.Remove(current.tile)	
	}
	return compactPath(parentOf, from, to)
}

func getNeighbors(current *tile) []tile{
	neighbors := []tile{}

	if validTile(current.neighborNorth) {neighbors = append(neighbors, *current.neighborNorth)}
	if validTile(current.neighborEast) {neighbors = append(neighbors, *current.neighborEast)}
	if validTile(current.neighborWest) {neighbors = append(neighbors, *current.neighborWest)}
	if validTile(current.neighborSouth) {neighbors = append(neighbors, *current.neighborSouth)}

	return neighbors
}

func validTile(t *tile) bool {  // TODO implement avoiding fire etc
	if t == nil {
		return false
	}
	return !t.wall 
}

func compactPath(parentOf map[tile]tile, from tile, to tile) ([]tile, bool) {
	path := []tile{to}

	current := to

	for current.xCoord != from.xCoord || current.yCoord != from.yCoord {		
		path = append(path, parentOf[current])
	//	current = parentOf[current]
	
	
		ok := true
		current, ok = parentOf[current]
		if  !ok{
			return nil, false//[]tile{}	
		}
	}

	return path, true	
}

func printPath(path []tile) {
	if path == nil {
		fmt.Println("No valid path exists")
	}
	for i, t := range path {
		fmt.Println(i , ":", t.xCoord ,"," ,t.yCoord)
	}
}

func mainPath() {

	workingPath()
	fmt.Println("--------------")
	blockedPath()

}

func workingPath() {
	matrix := [][]int {
		{0,1,0,0},
		{0,0,1,0},
		{0,0,0,0}, 
		{0,0,1,0}}
	testmap := TileConvert(matrix)
	printTileMap(testmap)
	
	path, _ := getPath(&testmap, testmap[0][0], testmap[0][2])

	fmt.Println("\nWorking path:")
	printPath(path)
}


func blockedPath(){
	matrix := [][]int {
		{0,1,0,0},
		{0,0,1,0},
		{0,0,1,0}, 
		{0,0,1,0}}
	testmap := TileConvert(matrix)
	printTileMap(testmap)
	
	path, _ := getPath(&testmap, testmap[0][0], testmap[3][3])

	fmt.Println("\nBlocked path:")
	printPath(path)

}
