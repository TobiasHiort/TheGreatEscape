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
	smoke     int

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
	for _, pers := range ppl {
		//	p := pers
		go func(p *Person){//, ind int){
			defer wg.Done()
			p.MovePerson(m)
			mutex.Lock()
			sList := &[]int{}
			p.getStats((sList))//(statsList[ind])
			
			// Below: check if there's more then 1 person on the same tile
			/*	for _, st := range *statsList {
			if st[0] == (*sList)[0] && st[1] == (*sList)[1] && st[0] != 0 && st[1] != 0 {
				panic(fmt.Sprintf("Multiple occupants at: %v", p.currentTile()))}
				//panic(fmt.Println(p.safe))}
		}*/
			*statsList = append(*statsList, *sList) 
			mutex.Unlock()
			
		}(pers)//, i)
	}
	step++
	wg.Wait()

	if math.Mod(float64(step), 2) == 0 {
		FireSpread(*m)
		if math.Mod(float64(step), 4) == 0 {
			SmokeSpread(*m)
			//InitPlans(m)
		}
		if math.Mod(float64(step), 2) == 0 {InitPlans(m)}   // change it up?
//		SmokeSpread(*m)
//		InitPlans(m)
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

func mainMap() {
	matrix := [][]int {
		{1,1,1,1,1,1,1,1,1,1,1,1,1},
		{1,0,0,0,0,0,1,0,0,0,0,0,1},
		{1,0,0,0,0,0,1,0,0,1,0,0,1},
		{1,0,0,0,0,0,1,1,1,1,0,0,1},
		{1,1,1,0,0,0,1,0,0,0,0,0,1},
		{2,0,0,0,0,0,0,0,0,0,0,0,1},
		{1,1,1,1,1,1,1,1,1,1,1,1,1}}

//	list := [][]int{{3,2}, {5,5}}
//	testmap := TileConvert(matrix)
//	SetFire(&testmap[2][2])
//	SmokeSpread(testmap)

//	fmt.Println((&testmap[3][3]).safestTile())

	ppl := [][]int{{1,1}, {1,2}, {2,1}, {1,3}}
	tryThis(matrix, ppl, -1, -1)
	
	//debugging()
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
}

//send heat
func FireStats(m *[][]tile) ([][]int, [][]int){
	fire := [][]int{}
	smoke := [][]int{}

	for i, list := range *m {
		for j, _ := range list {
			//tl := GetTile(*m, i, j)
			tl := &(*m)[i][j]
			if tl.heat  > 0 {fire = append(fire, []int{tl.yCoord, tl.xCoord, tl.heat})}
			if tl.smoke  > 0 {smoke = append(smoke, []int{tl.yCoord, tl.xCoord, tl.smoke})}
		}
	}
	return fire, smoke
}

// send smoke
func SmokeStats(m *[][]tile) [][]int{
	smoke := [][]int{}

	for i, list := range *m {
		for j, _ := range list {
			tl := &(*m)[i][j]
			if tl.smoke  > 0 {smoke = append(smoke, []int{tl.yCoord, tl.xCoord, tl.smoke})}
		}
	}
	return smoke
}

func SmokeSpread(tileMap [][]tile) {
	smokeTiles := []*tile{}
	for x := 0; x < len(tileMap); x++ {
		for y := 0; y < len(tileMap[0]); y++ {
			tl := &(tileMap[x][y])
			if tl.smoke > 0 {smokeTiles = append(smokeTiles, tl)}
		}
	}
	for _, s := range smokeTiles {
		if s.smoke > 0 {SmokeSpreadTile(s)}
	}
}

func SmokeSpreadTile(thisTile *tile) {
	smoke := thisTile.smoke
	if smoke >= 1 {smoke = 1}

	if thisTile.neighborNorth != nil && !thisTile.neighborNorth.wall && !thisTile.neighborNorth.outOfBounds {
		(thisTile.neighborNorth.smoke) += smoke //thisTile.smoke/30
	}
	if thisTile.neighborEast != nil && !thisTile.neighborEast.wall && !thisTile.neighborEast.outOfBounds {
		(thisTile.neighborEast.smoke) += smoke //thisTile.smoke/30
	}
	if thisTile.neighborWest != nil && !thisTile.neighborWest.wall && !thisTile.neighborWest.outOfBounds {
		(thisTile.neighborWest.smoke) += smoke //thisTile.smoke/30
	}
	if thisTile.neighborSouth != nil && !thisTile.neighborSouth.wall && !thisTile.neighborSouth.outOfBounds {
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
		SetFire(GetTile(currentMap, fire[1], fire[0]))
		tempList := []int{}
		tempList = append(tempList, fire[0])
		tempList = append(tempList, fire[1])
		tempList = append(tempList, MINHEAT)
		fireStats = append(fireStats, tempList)
	}
	return fireStats
}

func PeopleStats(peopleArray []*Person) []float32 {

	aliveAmount := 0
	deadAmount := 0
	injuredAmount := 0

	for i := 0; i < len(peopleArray); i++ {
		if ((peopleArray[i]).alive) && ((peopleArray[i]).hp < 70) {
			injuredAmount++
		} else if (peopleArray[i].alive) {
			aliveAmount++
		} else {
			deadAmount++
		}
	}
	return []float32{float32(aliveAmount), float32(deadAmount), float32(injuredAmount)}
}

func MapStats(inMap [][]tile) []int{
	fireTiles := 0
	for i := range inMap {
		for j := range inMap[i] {
			if inMap[i][j].heat >= MINHEAT {
					fireTiles ++
				}
			}
		}
	return []int{fireTiles}
}

func DoorCoord(inMap [][]tile) [][]int {

  var door [][]int
	for i := range inMap {
		for j := range inMap[i] {
			if inMap[i][j].door {
					door = append(door, []int{i, j})
				}
			}
		}
	return door
}

func (t *tile) safestTile() *tile {
	nbrs := getNeighbors(t)

	if len(nbrs) == 0 {return nil}
	
	safest := nbrs[0]

	for _, nbr := range nbrs {
		if nbr.smoke < safest.smoke {safest = nbr}
	}
	return safest 
	/*
	tmp1 := safest(t.neighborNorth, t.neighborWest)
	tmp2 := safest(t.neighborSouth, t.neighborEast)

	tmp3 := safest(t.neighborNW, t.neighborNE)
	tmp4 := safest(t.neighborSW, t.neighborSE)

	tmp1 = safest(tmp1, tmp2)
	tmp2 = safest(tmp3, tmp4)
	
	return safest(tmp1, tmp2)*/
}

func safest(t1, t2 *tile) *tile {
//	if !validTile(t1) {//|| t.occupied != nil {
	if !canGo(t1) {//|| t.occupied != nil {
//		if validTile(t2) {// && t.occupied == nil {
		if canGo(t2) {// && t.occupied == nil {
		return t2} else {return nil}}
//	if !validTile(t2) /*&& t1.occupied == nil*/ {return t1}
	if !canGo(t2) /*&& t1.occupied == nil*/ {return t1}	
	cst1 := math.Min(float64(t1.smoke), float64(t2.smoke))
	tmp1 := t1
	if int(cst1) < tmp1.smoke {tmp1 = t2}
	//if !validTile(tmp1) {return nil}
	return tmp1
}
