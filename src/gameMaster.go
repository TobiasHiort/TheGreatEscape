package main

import "fmt"

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

	for i, st := range peopleList {
	//	if st[0] == peopleList[0] && st[1] == (*sList)[1] && st[0] != 0 && st[1] != 0 {
		for j, st2 := range peopleList {
			if j > i {
				if st[0] == st2[0] && st[1] == st2[1] {panic(fmt.Sprintf("%v", st))}
			}
		}
	}
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

		Run(&currentMap, peopleArray, &statsList)
		fireList, smokeList = FireStats(&currentMap)		
		//smokeList = SmokeStats(&currentMap)
		SendToPipe(&statsList, &fireList, &smokeList)
	}
	readStats(peopleArray, currentMap)
}
