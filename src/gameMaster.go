package main

func MapInit(peopleList [][]int, newMap [][]int) [][]tile {
	//gets map data from GM and inits the map
	var currentMap ([][]tile)
	currentMap = TileConvert(newMap)

	return currentMap
}

func GameLoop(inMap [][]int, peopleList [][]int, fireStartPos [][]int) {
	//newMap := MapInit(foo, bar)
	//do all the Inits
	//statsList := StatsInit(len(peopleList))

//	toPipe(&[][]int{})
	
	currentMap := MapInit(peopleList, inMap)
	peopleArray := PeopleInit(currentMap, peopleList)
	statsList := StatsStart(peopleArray)
	fireList := FireInit(currentMap, fireStartPos)
	smokeList := SmokeStats(&currentMap)
	InitPlans(&currentMap)
	SendToPipe(&statsList, &fireList, &smokeList)
	for !CheckFinish(peopleArray) {
		//if *a == *b {
		/*		peopleArray = *///Run(&currentMap, peopleArray, &statsList)
	//	fmt.Print((currentMap))
	//	statsList = StatsStart(peopleArray)
	//	fmt.Println(statsList)	
		Run(&currentMap, peopleArray, &statsList)
		fireList, smokeList = FireStats(&currentMap)		
		//smokeList = SmokeStats(&currentMap)
		SendToPipe(&statsList, &fireList, &smokeList)
		//} //PrintTileMapP(aMap)
	}

}
