package main

import "fmt"
import "sync"
import "math"


import(
  //  "os"
    //"bufio"
	"encoding/json"
	//"time"
    "io/ioutil"
)
//import "time"

const MINHEAT = 10
const MEDIUMHEAT = 20
const MAXHEAT = 30

type tile struct {
	xCoord int
	yCoord int

	heat      int //how hot a tile is before fire
	fireLevel int //strength of the fire
	smoke int
	
	wall bool
	door bool

	occupied *Person
	personID int

	outOfBounds bool

	neighborNorth *tile
	neighborEast  *tile
	neighborSouth *tile
	neighborWest  *tile

	neighborNW *tile
	neighborNE *tile
	neighborSE *tile
	neighborSW *tile
}

//Initializes the fire
func SetFire(thisTile *tile) {
	thisTile.heat = MINHEAT
	thisTile.fireLevel = 1
	thisTile.smoke = 1
}

func FireSpread(tileMap [][]tile) {
	for x := 0; x < len(tileMap); x++ {
		for y := 0; y < len(tileMap[0]); y++ {
			tl := &(tileMap[x][y])
			// tl.heat < 30 {fireSpreadTile(tl)}
			fireSpreadTile(tl)
		}
	}
}

func fireSpreadTile(thisTile *tile) { //TODO: fixa tiles på riktigt!!
	if thisTile.heat >= MINHEAT {
		thisTile.fireLevel = 1
	}
	if thisTile.heat >= MEDIUMHEAT {
		thisTile.fireLevel = 2
	}
	if thisTile.heat >= MAXHEAT {
		thisTile.fireLevel = 3
	}
	fire := thisTile.fireLevel
	if fire > 10 {fire = 10}

	if thisTile.neighborNorth != nil && !thisTile.neighborNorth.wall && thisTile.fireLevel != 0 {
		(thisTile.neighborNorth.heat) += fire//thisTile.fireLevel
	}
	if thisTile.neighborEast != nil && !thisTile.neighborEast.wall && thisTile.fireLevel != 0 {
		(thisTile.neighborEast.heat) += fire//thisTile.fireLevel
	}
	if thisTile.neighborWest != nil && !thisTile.neighborWest.wall && thisTile.fireLevel != 0 {
		(thisTile.neighborWest.heat) += fire//thisTile.fireLevel
	}
	if thisTile.neighborSouth != nil && !thisTile.neighborSouth.wall && thisTile.fireLevel != 0 {
		(thisTile.neighborSouth.heat) += fire//thisTile.fireLevel
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

	if x > 0 && y > 0 {
		thisTile.neighborNW = &tileMap[x-1][y-1]
	}
	if x > 0 && y < maxY-1 {
		thisTile.neighborNE = &tileMap[x-1][y+1]
	}
	if x < maxX-1 && y < maxY-1 {
		thisTile.neighborSE = &tileMap[x+1][y+1]
	}
	if x < maxX-1 && y > 0 {
		thisTile.neighborSW = &tileMap[x+1][y-1]
	}
}

func makeNewTile(thisPoint int, x int, y int) tile {

	//makes a basic floor tile with no nothin on it
	//and also no neighbors
	newTile := tile{x, y, 0, 0, 0, false, false, nil, 0, false, nil, nil, nil, nil, nil, nil, nil, nil}

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
	//	tile := GetTile(inMap, person[0], person[1])  // inverted!
		tile := GetTile(inMap, person[1], person[0])
		peopleArray[i] = makePerson(tile)
	}
	return peopleArray
}

func Run(m *[][]tile, ppl []*Person, statsList *[][]int) {
	//sList := []Stats{}

	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}
	wg.Add(len(ppl))
	*statsList = [][]int{}
	for i, pers := range ppl {
		go func(p *Person, ind int){//, ind int){
			//	ind := i
		//	p := pers
			defer wg.Done()

			//	sList :=  []int{}// append(sList, p.getStats())
			p.MovePerson(m)
			sList := &[]int{}
			p.getStats((sList))//(statsList[ind])
			mutex.Lock()
		//	(*statsList)[ind] = *sList
			*statsList = append(*statsList, *sList)
			mutex.Unlock()
			//	fmt.Println(len(statsList))

		}(pers, i)//, i)
	}
	step++
	wg.Wait()
	
	if math.Mod(float64(step), 2) == 0 {
		FireSpread(*m)
		SmokeSpread(*m)
		InitPlans(m)
	} // MOVED!
	
	
//	if math.Mod(float64(step), 40) == 0 {InitPlans(m)}
	// TODO: takes up tiiime!!1 fixy-changy

	//if math.Mod(float64(step), 5) == 0 {InitPlans(m)}
}

 func CheckFinish (peopleArray []*Person) bool {
 	for i := 0; i < len(peopleArray); i++ {
 		if (peopleArray[i].safe == false && peopleArray[i].alive == true) {
 			return false
 		}
 	}
 	return true
 }

func mainMap() {
	if math.Mod(float64(step), 20) == 0 {}
//	mainPath()
//	MainPeople()
//      testRedirect()
//	testDiagPpl()
//	testDiag()
	debugging()
//	testSame()
//	testDiagonally()
//	testMovePeople()
//	Whut()
//	testJP()
	//	GLoop()

//	testDiagonally()
//	testMovePeople()

}

func testRedirect() {
	matrix := [][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{1, 0, 1, 1, 1, 0, 1},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 2, 0, 0, 0, 0}}

	list := [][]int{
		{1, 1},
		{1, 2},
		{0, 2},
		{0, 0},
		{2, 1},
		{3, 1},
		{0, 3},
		{0, 1},
		{1, 0},
		{1, 3}}

	tryThis(matrix, list, -1, -1)
}

func testDiag() {
	matrix := [][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 2}}

	list := [][]int{
		{0, 0},
		{0, 6}}
	tryThis(matrix, list, -1, -1)
}



func testDiagPpl() {
	matrix := [][]int{
		{0, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 2}}

	list := [][]int{
		{0, 0},
		{0, 1},
		{1, 0},
		{0, 6},
		{2, 4}}

	tryThis(matrix, list, -1, -1)
}

func testDiagonally() {
	matrix := [][]int{
		{0, 1, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0, 0, 0},
		{0, 1, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 2}}

	list := [][]int{
		{0, 0},
		{0, 6},
		{2, 4}}

	tryThis(matrix, list, 1, 2)
}

func testMovePeople() {
	matrix := [][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0},
		{1, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 2}}

	list := [][]int{
		{0, 0}}
	//	{0, 6},
	//	{2, 4}}

	tryThis(matrix, list, -1, -1)
	// Note: it takes 1 timeunit to take a step from the door and away
}


func testJP() {
	matrix := [][]int{
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{2, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0, 1}}

	list := [][]int{
		{0, 2}}
	//	{0, 6},
	//	{2, 4}}


	testmap := TileConvert(matrix)
	pplArray := PeopleInit(testmap, list)

	MovePeople(&testmap, pplArray)

	for i, p := range pplArray {
		fmt.Println("Person", i, "time:  ", p.time, "\n         health:", p.hp)
	}

}



func tryThis(matrix [][]int, ppl [][]int, x, y int) {
	testmap := TileConvert(matrix)
	pplArray := PeopleInit(testmap, ppl)

	InitPlans(&testmap)
	fmt.Println("init")
	if x >= 0 && y >= 0 {SetFire(&testmap[x][y])}
	MovePeople(&testmap, pplArray)

	for i, p := range pplArray {
		fmt.Println("Person", i, "time:  ", p.time, "\n         health:", p.hp)
	}
}


func testSame() {
	matrix := [][]int {
		{1,1,1,1,1,1,1,1,1,1,1,1,1},
		{1,0,0,0,0,0,1,0,0,0,0,0,1},
		{1,0,0,0,0,0,1,0,0,1,0,0,1},
		{1,0,0,0,0,0,1,1,1,1,0,0,1},
		{1,1,1,0,0,0,1,0,0,0,0,0,1},
		{2,0,0,0,0,0,0,0,0,0,0,0,1},
		{1,1,1,1,1,1,1,1,1,1,1,1,1}}

	list := [][]int{{3,2}, {5,5}}

	tryThis(matrix, list, -1, -1)
}



func debugging() {

	if true {b, err3 := ioutil.ReadFile("../src/mapfile.txt")
	if err3 != nil{
		panic(err3)
	}

	var m = [][]int{}
	err := json.Unmarshal(b, &m)
	if err != nil{
		panic(err)
	}
	//m[8][1] = 2
	//m[13][0] = 2

	//	testmap := TileConvert(m)


	c, err4 := ioutil.ReadFile("../src/playerfile.txt")
	if err4 != nil{
		panic(err4)
	}

	var mm = [][]int{}
	err5 := json.Unmarshal(c, &mm)
	if err5 != nil{
		panic(err5)
	}
		tryThis(m, mm, 20, 20)// 31, 31)
	}
/*	m := [][]int{
		{1,1,1,1,1,1},
		{2,0,0,0,0,1},
		{1,0,0,0,0,1},
		{1,0,0,0,0,1},
		{1,0,0,0,0,1},
		{1,0,0,0,0,1},
		{1,0,0,0,0,1},
		{1,0,0,0,0,1},		
		{1,0,0,0,0,2},
		{1,1,1,1,1,1}}
	mm := [][]int{{1,4}}
*/
//	ppl := PeopleInit(testmap, mm)

//	list := [][]int{{89,33}}//{104, 28}, {105, 29}}  // lr tvärtom?

}



//send heat
func FireStats(m *[][]tile) [][]int{
	//func FireStats(start []*tile, dir Direction) [][]int{
	fire := [][]int{}

	for i, list := range *m {
		for j, _ := range list {
			//tl := GetTile(*m, i, j)
			tl := &(*m)[i][j]
			if tl.heat  > 0 {fire = append(fire, []int{tl.yCoord, tl.xCoord, tl.heat})}
		}
	}
	return fire
}

// send smoke

func SmokeStats(m *[][]tile) [][]int{
	//func FireStats(start []*tile, dir Direction) [][]int{
	smoke := [][]int{}

	for i, list := range *m {
		for j, _ := range list {
			//tl := GetTile(*m, i, j)
			tl := &(*m)[i][j]
			if tl.smoke  > 0 {smoke = append(smoke, []int{tl.yCoord, tl.xCoord, tl.smoke})}
		}
	}
	return smoke
}



// testing smoke

func SmokeSpread(tileMap [][]tile) {
	smokeTiles := []*tile{}
	for x := 0; x < len(tileMap); x++ {
		for y := 0; y < len(tileMap[0]); y++ {
			tl := &(tileMap[x][y])
			// tl.heat < 30 {fireSpreadTile(tl)}
			if tl.smoke > 0 {smokeTiles = append(smokeTiles, tl)}
		}
	}
	for _, s := range smokeTiles {
		if s.smoke > 0 {SmokeSpreadTile(s)}
	}

}

func SmokeSpreadTile(thisTile *tile) { //TODO: fixa tiles på riktigt!!
/*	if thisTile.heat >= MINHEAT {
		thisTile.fireLevel = 1
	}
	if thisTile.heat >= MEDIUMHEAT {
		thisTile.fireLevel = 2
	}
	if thisTile.heat >= MAXHEAT {
		thisTile.fireLevel = 3
	}*/
	//	if thisTile.smoke < 1 {return}
	smoke := thisTile.smoke
	if smoke >= 1 {smoke = 1}
	
	if thisTile.neighborNorth != nil && !thisTile.neighborNorth.wall {
		(thisTile.neighborNorth.smoke) += smoke //thisTile.smoke/30
	}
	if thisTile.neighborEast != nil && !thisTile.neighborEast.wall {
		(thisTile.neighborEast.smoke) += smoke //thisTile.smoke/30
	}
	if thisTile.neighborWest != nil && !thisTile.neighborWest.wall {
		(thisTile.neighborWest.smoke) += smoke //thisTile.smoke/30
	}
	if thisTile.neighborSouth != nil && !thisTile.neighborSouth.wall {
		(thisTile.neighborSouth.smoke) += smoke //thisTile.smoke/30
	}
/*	
	if thisTile.neighborNW != nil && !thisTile.neighborNW.wall {
		(thisTile.neighborNW.smoke) += 1//thisTile.smoke
	}
	if thisTile.neighborNE != nil && !thisTile.neighborNE.wall {
		(thisTile.neighborNE.smoke) += 1//thisTile.smoke
	}
	if thisTile.neighborSE != nil && !thisTile.neighborSE.wall {
		(thisTile.neighborSE.smoke) += 1//thisTile.smoke
	}
	if thisTile.neighborSW != nil && !thisTile.neighborSW.wall {
		(thisTile.neighborSW.smoke) += 1//thisTile.smoke
	}*/
}




// merging

func StatsStart(peopleArray []*Person) [][]int {
	statsList := [][]int{}
	for i := 0; i < len(peopleArray); i++ {
		statsList = append(statsList, peopleArray[i].GetStats())
	}
	return statsList
}

func FireInit(currentMap [][]tile, fireList [][]int) [][]int {
	fireStats := [][]int{}
	for _, fire := range fireList {
		SetFire(GetTile(currentMap, fire[0], fire[1]))
		tempList := []int{}
		tempList = append(tempList, fire[0])
		tempList = append(tempList, fire[1])
		tempList = append(tempList, MINHEAT)
		fireStats = append(fireStats, tempList)
	}
	return fireStats
}
