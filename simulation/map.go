package main

import "fmt"
import "sync"

const MINHEAT = 10
const MEDIUMHEAT = 20
const MAXHEAT = 30

type tile struct {
	xCoord int
	yCoord int

	heat      int //how hot a tile is before fire
	fireLevel int //strength of the fire

	wall bool
	door bool

	occupied *Person
	personID int

	outOfBounds bool

	neighborNorth *tile
	neighborEast  *tile
	neighborSouth *tile
	neighborWest  *tile
}

//Initializes the fire
func SetFire(thisTile *tile) {
	thisTile.heat = MINHEAT
	thisTile.fireLevel = 1
}

func FireSpread(tileMap [][]tile) {
	for x := 0; x < len(tileMap); x++ {
		for y := 0; y < len(tileMap[0]); y++ {
			fireSpreadTile(&(tileMap[x][y]))
		}
	}

}

func fireSpreadTile(thisTile *tile) {
	if thisTile.heat >= MINHEAT {
		thisTile.fireLevel = 1
	}
	if thisTile.heat >= MEDIUMHEAT {
		thisTile.fireLevel = 2
	}
	if thisTile.heat >= MAXHEAT {
		thisTile.fireLevel = 3
	}

	if thisTile.neighborNorth != nil && thisTile.fireLevel != 0 {
		(thisTile.neighborNorth.heat) += thisTile.fireLevel
	}
	if thisTile.neighborEast != nil && thisTile.fireLevel != 0 {
		(thisTile.neighborEast.heat) += thisTile.fireLevel
	}
	if thisTile.neighborWest != nil && thisTile.fireLevel != 0 {
		(thisTile.neighborWest.heat) += thisTile.fireLevel
	}
	if thisTile.neighborSouth != nil && thisTile.fireLevel != 0 {
		(thisTile.neighborSouth.heat) += thisTile.fireLevel
	}
}

func assignNeighbor(thisTile *tile, x int, y int, maxX int, maxY int, tileMap [][]tile) {
	if x > 0 {
		thisTile.neighborNorth = &tileMap[x-1][y]
	}

	if y > 0 {
		thisTile.neighborWest = &tileMap[x][y-1]
	}

	if x < maxX-1 {
		thisTile.neighborSouth = &tileMap[x+1][y]
	}

	if y < maxY-1 {
		thisTile.neighborEast = &tileMap[x][y+1]
	}
}

func makeNewTile(thisPoint int, x int, y int) tile {

	//makes a basic floor tile with no nothin on it
	//and also no neighbors
	newTile := tile{x, y, 0, 0, false, false, nil, 0, false, nil, nil, nil, nil}

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

func TileConvert(inMap [][]int) [][]tile {
	mapXSize := len(inMap)
	mapYSize := len(inMap[0])

	//Initiates a slice of tile slices (2D tile slice)
	tileMap := make([][]tile, mapXSize)

	for x := 0; x < mapXSize; x++ {
		//initiates slice of tiles
		tileMap[x] = make([]tile, mapYSize)

		for y := 0; y < mapYSize; y++ {
			//constructs a new tile
			newTile := makeNewTile(inMap[x][y], x, y)

			//inserts tile into 2d slice
			tileMap[x][y] = newTile

		}
	}

	//Assigns 4 neighbors to each tile
	for x := 0; x < mapXSize; x++ {
		for y := 0; y < mapYSize; y++ {
			assignNeighbor(&(tileMap[x][y]), x, y, mapXSize, mapYSize, tileMap)
		}
	}

	return tileMap

}

func GetTile(inMap [][]tile, x int, y int) *tile {
	for i := range inMap {
		for j := range inMap[i] {
			if inMap[i][j].xCoord == x && inMap[i][j].yCoord == y {
				return &inMap[i][j]
			}
		}
	}
	return nil
}

func PeopleInit(inMap [][]tile, peopleList [][]int) []*Person {
	size := len(peopleList)
	peopleArray := make([]*Person, size)
	for i, person := range peopleList {
		tile := GetTile(inMap, person[0], person[1])
		peopleArray[i] = makePerson(tile)
	}
	return peopleArray
}


func Run(inMap [][]tile, peopleArray []*Person) {
	// go run ruitnes for concurrency
	for _, person := range peopleArray {
		person.MovePerson(&inMap)
	}
}

func RunGo(inMap [][]tile, peopleArray []*Person) {
	var wg sync.WaitGroup
	wg.Add(len(peopleArray))
	for _, person := range peopleArray {
		go func(currentPerson *Person) {
			defer wg.Done()
			currentPerson.MovePerson(&inMap)
		}(person)
	}
	wg.Wait()
}

func printTileP(thisTile tile) {
	if thisTile.occupied != nil{
		fmt.Print("X")
	} else if thisTile.wall {
		fmt.Print("1")
	} else if thisTile.door {
		fmt.Print("2")
	} else if thisTile.outOfBounds {
		fmt.Print("3")
	} else {
		fmt.Print("0")
	}
  
}


func printTileMapP(inMap [][]tile) {
	mapXSize := len(inMap)
	mapYSize := len(inMap[0])

	for x:= 0; x < mapXSize; x++ {
		for y:= 0; y < mapYSize; y++{
			printTileP(inMap[x][y])
		}
		fmt.Print("\n")
	}
}

 func CheckFinish (peopleArray []*Person) bool {
 	for i := 0; i < len(peopleArray); i++ {
 		if (peopleArray[i].safe == false && peopleArray[i].alive == true) {
 			return false
 		}  
 	}
 	return true
 }
/*
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
*/

func main() {

	matrix := [][]int{
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{1, 0, 1, 1, 1, 1, 1},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 2, 0, 0, 0}}
	testmap := TileConvert(matrix)
	/*
		var tile = GetTile (testmap, 2, 0)
		printTile(*tile)
		fmt.Print("/n")
/*
		start1 := &testmap[1][0]
		start2 := &testmap[1][2]
		start3 := &testmap[0][1]
		start4 := &testmap[3][4]
		start5 := &testmap[2][0]
		start6 := &testmap[5][5]
/*
		var p1 = *makePerson(start1)
		var p2 = *makePerson(start2)
		var p3 = *makePerson(start1)
		var p4 = *makePerson(start2)
		var p5 = *makePerson(start1)
		var p6 = *makePerson(start2)
*/
		/*
		list := make([][]int, 0)
		list.append([1][2])
		list.append([0][2])
		list.append([2][3])*/
		list := [][]int{
			{1, 2},
			{0, 2},
			{3, 0}}

		 peopleArray := PeopleInit (testmap, list)
	for _, people := range peopleArray {
		if people != nil {
			fmt.Print("True")
			fmt.Print("\n")
		}
	}
/*
	printTileMapP(testmap)
	Run(testmap, peopleArray)
	fmt.Print("\n")
	printTileMapP(testmap)
	Run(testmap, peopleArray)
	fmt.Print("\n")
	printTileMapP(testmap) */

	for !CheckFinish(peopleArray) {
		printTileMapP(testmap)
		RunGo(testmap, peopleArray)
		fmt.Print("\n")
	}
	
	if CheckFinish (peopleArray) == false {
		fmt.Print("false")
		fmt.Print("\n")
	}



	//mainPath()
//	MainPeople()

}

