package main

import "fmt"

type tile struct {
	xCoord int
	yCoord int

	heat			int
	fireLevel int

	wall bool
	door bool

	occupied bool
	personID int

	outOfBounds bool

	neighborNorth *tile
	neighborEast  *tile
	neighborSouth *tile
	neighborWest  *tile
}

/*
func FireSpread(thisTile tile){
	if thisTile.heat > 9 {
		thisTile.fireLevel = 1
	}
	if thisTile.heat > 19 {
		thisTile.fireLevel = 2
	}
	if thisTile.heat > 29 {
		thisTile.fireLevel = 3
	}

	tile.neighborNorth.heat += fireLevel
	tile.neighborEast.heat	+= fireLevel
	tile.neighborWest.heat	+= fireLevel
	tile.neighborSouth.heat += fireLevel

	/*
	tile.neighborNorth.heat += heat/10
	tile.neighborEast.heat += heat/10
	tile.neighborWest.heat += heat/10
	tile.neighborSouth.heat += heat/10
	*/
//}

func assignNeighbor(thisTile tile, x int, y int, tileMap [][]tile) {
	if x > 0 {
		thisTile.neighborWest = &tileMap[x-1][y]
	}

	if y > 0 {
		thisTile.neighborNorth = &tileMap[x][y-1]
	}

	if x < len(tileMap[x]) {
		thisTile.neighborEast = &tileMap[x+1][y]
	}

	if x < len(tileMap) {
		thisTile.neighborNorth = &tileMap[x][y+1]
	}
}

func makeNewTile(thisPoint int, x int, y int) tile{

	//makes a basic floor tile with no nothin on it
	newTile := tile{x, y, 0, 0, false, false, false, 0, false, nil, nil, nil, nil}

	if thisPoint == 0 {
		//make normal floor
		//helt normalt flour

		//append to tilemap
	} else if thisPoint == 1 {
		//wall
		newTile.wall = true
	} else if thisPoint == 2 {
		//door
		newTile.door = true
	} else if thisPoint == 3 {
		//out of bounds
		newTile.outOfBounds = true
	}

	return newTile
}

func TileConvert(inMap [][]int) [][]tile{
	mapXSize := len(inMap)
	mapYSize := len(inMap[0])

	//Initiates a slice of tile slices (2D tile slice)
	tileMap := make([][]tile, mapXSize)

	for x:= 0; x < mapXSize; x++{
		//initiates slice of tiles
		tileMap[x] = make([]tile, mapYSize)

		for y:= 0; y < mapYSize; y++{
			//constructs a new tile
			newTile := makeNewTile(inMap[x][y], x, y)

			//inserts tile into 2d slice
			tileMap[x][y] = newTile

		}
	}

	//Assigns 4 neighbors to each tile
	for x:= 0; x < mapXSize; x++{
		for y:= 0; y < mapYSize; y++{
			assignNeighbor(tileMap[x][y], x, y, tileMap)
		}
	}

	return tileMap
}

func printTile(thisTile tile) {
	if thisTile.wall {
		fmt.Print("[vägg] ")
	} else if thisTile.door {
		fmt.Print("[dörr] ")
	} else if thisTile.outOfBounds {
		fmt.Print("[ute] ")
	} else {
		fmt.Print("[golv] ")
	}
}

func printTileMap(inMap [][]tile) {
	mapXSize := len(inMap[0])
	mapYSize := len(inMap)

	for x:= 0; x < mapXSize; x++{
		for y:= 0; y < mapYSize; y++{
			printTile(inMap[x][y])
		}
    fmt.Print("\n")
	}
}

func main() {
	testMatrix := [][]int{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{1, 1, 0, 1, 1},
		{0, 0, 0, 3, 3},
		{2, 0, 0, 3, 3}}

		amap := TileConvert(testMatrix)
		//tileConvert(testMatrix)
		printTileMap(amap)
}
