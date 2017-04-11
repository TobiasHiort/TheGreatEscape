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

func SetFire(thisTile *tile) {
	thisTile.fireLevel = 1
}

func FireSpread(tileMap *([][]tile)) {
	for x:= 0; x < len(&tileMap); x++{
		for y:= 0; y < len(&tileMap[0]); y++{
			fireSpreadTile(&(tileMap[x][y]))
		}
	}

}

func fireSpreadTile(thisTile *tile){
	if thisTile.heat > 9 {
		thisTile.fireLevel = 1
	}
	if thisTile.heat > 19 {
		thisTile.fireLevel = 2
	}
	if thisTile.heat > 29 {
		thisTile.fireLevel = 3
	}

	if thisTile.neighborNorth != nil {
		(thisTile.neighborNorth.heat) += thisTile.fireLevel
	}
	if thisTile.neighborEast != nil {
		(thisTile.neighborEast.heat)	+= thisTile.fireLevel
	}
	if thisTile.neighborWest != nil {
		(thisTile.neighborWest.heat)	+= thisTile.fireLevel
	}
	if thisTile.neighborSouth != nil {
		(thisTile.neighborSouth.heat) += thisTile.fireLevel
	}
}

func assignNeighbor(thisTile *tile, x int, y int, maxX int, maxY int, tileMap [][]tile) {

	if x > 0 {
		//thisTile.neighborWest = &tileMap[x-1][y]
		thisTile.neighborNorth = &tileMap[x-1][y]
	}

	if y > 0 {
		//thisTile.neighborNorth = &tileMap[x][y-1]
		thisTile.neighborWest = &tileMap[x][y-1]
	}

	if x < maxX-1 {
		//thisTile.neighborEast = &tileMap[x+1][y]
		thisTile.neighborSouth = &tileMap[x+1][y]
	}

	if y < maxY-1 {
		//thisTile.neighborSouth = &tileMap[x][y+1]
		thisTile.neighborEast = &tileMap[x][y+1]
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
			assignNeighbor(&(tileMap[x][y]), x, y, mapXSize, mapYSize, tileMap)
		}
	}

	return tileMap
}

func printTile(thisTile tile) {
	if thisTile.wall {
		fmt.Print("[vägg(")
	} else if thisTile.door {
		fmt.Print("[dörr(")
	} else if thisTile.outOfBounds {
		fmt.Print("[ute(")
	} else {
		fmt.Print("[golv(")
	}
	fmt.Print(thisTile.fireLevel)

	fmt.Print(")] ")
}

func printTileMap(inMap [][]tile) {
	mapXSize := len(inMap)
	mapYSize := len(inMap[0])

	for x:= 0; x < mapXSize; x++{
		for y:= 0; y < mapYSize; y++{
			printTile(inMap[x][y])
		}
		fmt.Print("\n")
	}
}
func printNeighbors(atile tile) {
	if atile.neighborNorth != nil {
		fmt.Print("North: ")
		printTile(*(atile.neighborNorth))
		fmt.Print("\n")
	} else {
		fmt.Print("North: nil\n")
	}
	if atile.neighborWest != nil {
		fmt.Print("West: ")
		printTile(*(atile.neighborWest))
		fmt.Print("\n")
	} else {
		fmt.Print("West: nil\n")
	}
	if atile.neighborEast != nil {
		fmt.Print("East: ")
		printTile(*(atile.neighborEast))
		fmt.Print("\n")
	} else {
		fmt.Print("East: nil\n")
	}
	if atile.neighborSouth != nil {
		fmt.Print("South: ")
		printTile(*(atile.neighborSouth))
		fmt.Print("\n")
	} else {
		fmt.Print("South: nil\n")
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
		fmt.Print("\n")
		printNeighbors(amap[4][4])

		//fire testing
		SetFire(&(amap[2][2]))
		for i := 0; i < 100; i++{
			FireSpread(&amap)
			if i%10 == 0{
				printTileMap(amap)
			}
		}
	}
