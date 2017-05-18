package main

func MapInit(peopleList [][]int, newMap [][]int) [][]tile {
	//gets map data from GM and inits the map
	var currentMap ([][]tile)
	currentMap = TileConvert(newMap)

	return currentMap
}

func GameLoop(inMap [][]int, peopleList [][]int, fireStartPos []int) {
	//newMap := MapInit(foo, bar)
	//do all the Inits
	statsList := StatsInit(len(peopleList))
	fireList := StatsInit(1)

	currentMap := MapInit(peopleList, inMap)
	peopleArray := PeopleInit(currentMap, peopleList)
	statsList = StatsStart(statsList, peopleArray)
	SetFire(GetTile(currentMap, fireStartPos[0], fireStartPos[1]))
	for !CheckFinish(peopleArray) {
		//if *a == *b {
		Run(&currentMap, peopleArray, statsList)
		fireList = FireStats(&currentMap)
		SendToPipe(statsList, fireList)
		//} //PrintTileMapP(aMap)
	}

}
