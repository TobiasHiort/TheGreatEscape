package main

//import "fmt"

//toppen av stupröret, yttersta lagret av löken
/*
type personEvent struct {
	newPos *tile
	newAliveState *bool
	currentPerson *Person
}
*/

/*
type updateEvent {
	//each person has a spot in this array. if they dont move, the coords are blank
	//if they die, maybe a -1,-1???
	personMoveList		*([][]int)
	//sends a map of the coords that have fire that has increased
	fireIncreaseList	*([][]int)
	//timeStamp int
}
*/

//func get ppl movement from map

/*
func timeStamp(inMap [][]tile, currentTime int) Event{
	var ev Event
	ev.tileMap = inMap
	ev.timeStamp = currentTime
	return ev
}
*/

//this function will in the future preferably grab just the changes on the map, rather than it as a whole
//outEvent := Event
//outEvent = timeStamp(tileMap, time)
//grabs data from map
//returns an Event

func MapInit(peopleList [][]int, newMap [][]int) [][]tile {
	//gets map data from GM and inits the map
	var currentMap ([][]tile)
	currentMap = TileConvert(newMap)

	///peopleArray = PeopleInit(currentMap, peopleList)

	//gets ppl data from GM and calls the initpppl function in map
	return currentMap
}

func GameLoop(inMap [][]int, peopleList [][]int, fireStartPos []int) {
	//newMap := MapInit(foo, bar)
	//do all the Inits
	statsList := StatsInit(len(peopleList))
	fireList := StatsInit(1)
	/*
		inMap := [][]int{
			{0, 0, 0, 1, 0, 0, 0},
			{0, 0, 0, 1, 0, 0, 0},
			{1, 0, 1, 1, 1, 1, 1},
			{0, 0, 0, 1, 0, 0, 0},
			{0, 0, 0, 1, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 2, 0, 0, 0}}

			peopleList := [][]int{
				{1, 2},
				{0, 2},
				{3, 0}}
	*/

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

/*
func main() {
	GameLoop()
}
*/
